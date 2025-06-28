package planner

import (
	"fmt"
	"matching-engine/internal/adapter/ortool"
	"matching-engine/internal/enums"
	"time"

	"github.com/rs/zerolog/log"
	"matching-engine/internal/model"
	"matching-engine/internal/service/pickupdropoffservice"
	"matching-engine/internal/service/timematrix"
)

type ORToolPlanner struct {
	pickupDropoffSelector pickupdropoffservice.PickupDropoffSelectorInterface
	timeMatrixSelector    timematrix.Selector
	orToolClient          *ortool.ORToolClient
	cfg                   *ORToolConfig
}

func NewORToolPlanner(
	pickupDropoffSelector pickupdropoffservice.PickupDropoffSelectorInterface,
	timeMatrixSelector timematrix.Selector,
	orToolClient *ortool.ORToolClient,

) PathPlanner {
	return &ORToolPlanner{
		pickupDropoffSelector: pickupDropoffSelector,
		timeMatrixSelector:    timeMatrixSelector,
		orToolClient:          orToolClient,
		cfg:                   NewORToolConfig(),
	}
}

func (p *ORToolPlanner) FindFirstFeasiblePath(
	offerNode *model.OfferNode,
	requestNode *model.RequestNode,
) ([]model.PathPoint, bool, error) {
	// Step 1: Get pickup & dropoff
	pickupDropoff, err := p.pickupDropoffSelector.
		GetPickupDropoffPointsAndDurations(requestNode.Request(), offerNode.Offer())
	if err != nil {
		return nil, false, fmt.Errorf("pickup/dropoff: %w", err)
	}

	// Step 2: Get time matrix and index map
	timeMatrixData, err := p.timeMatrixSelector.GetTimeMatrix(offerNode, requestNode)
	if err != nil {
		return nil, false, fmt.Errorf("time matrix: %w", err)
	}
	fullMatrix := timeMatrixData.TimeMatrix()
	pointIndex := timeMatrixData.PointIdToIndex()

	// Step 3: Construct full path (original path + pickup + dropoff)
	originalPath := offerNode.Offer().Path()
	path := make([]model.PathPoint, 0, len(originalPath)+2)

	// Add offer source and destination
	path = append(path, originalPath[0])

	// Add any intermediate points (2nd to second-last)
	if len(originalPath) > 2 {
		path = append(path, originalPath[1:len(originalPath)-1]...)
	}

	// Add pickup and dropoff points
	path = append(path, *pickupDropoff.Pickup())
	path = append(path, *pickupDropoff.Dropoff())

	path = append(path, originalPath[len(originalPath)-1])

	// Step 4: Build sub-matrix and time windows
	pathLen := len(path)
	subMatrix := make([][]int, pathLen)
	timeWindows := make([][2]int, pathLen)
	capacities := make([]int, pathLen)
	pickupDropoffMap := map[string][2]int{}
	departureTime := offerNode.Offer().DepartureTime()

	for i, fromPoint := range path {
		fromIdx, ok := pointIndex[fromPoint.ID()]
		if !ok {
			return nil, false, fmt.Errorf("unknown point ID: %q", fromPoint.ID())
		}

		timeWindow, capacity := calculateTimeWindow(fromPoint, departureTime, pickupDropoffMap, i)
		timeWindows[i] = timeWindow
		capacities[i] = capacity

		row := make([]int, pathLen)
		for j, toPoint := range path {
			toIdx := pointIndex[toPoint.ID()]
			duration := fullMatrix[fromIdx][toIdx]
			row[j] = scaleDownDuration(duration)
		}
		subMatrix[i] = row
	}

	var pickupAndDropoffs = make([][2]int, 0, len(pickupDropoffMap))
	for _, v := range pickupDropoffMap {
		pickupAndDropoffs = append(pickupAndDropoffs, v)
	}

	data := ortool.NewORToolData(
		subMatrix,
		timeWindows,
		capacities,
		offerNode.Offer().Capacity(),
		pickupAndDropoffs,
		scaleDownDuration(offerNode.Offer().MaxEstimatedArrivalTime().Sub(departureTime)),
	)

	data.SetMethod(p.cfg.Method)
	data.SetTimeout(p.cfg.Timeout)
	data.SetEnableGuidedLocalSearch(p.cfg.EnableGuidedLocalSearch)

	// Step 5: Call ORTool solver
	solution, err := p.orToolClient.CallPythonORToolSolver(data)

	if err != nil {
		log.Debug().
			Str("offer_id", offerNode.Offer().ID()).
			Str("request_id", requestNode.Request().ID()).
			Msg("Solver error")
		return nil, false, fmt.Errorf("ORTool solver error: %w", err)
	}

	if solution == nil || !solution.Success {
		log.Debug().
			Str("offer_id", offerNode.Offer().ID()).
			Str("request_id", requestNode.Request().ID()).
			Msg("Solver did not find a solution")
		return nil, false, nil
	}

	// Step 6: Construct result with expected arrival times
	result := make([]model.PathPoint, pathLen)
	for i, step := range solution.Route {
		point := path[step.Node]
		point.SetExpectedArrivalTime(departureTime.Add(scaleUpDuration(step.ArrivalTime)))
		result[i] = point
	}

	return result, true, nil
}

func calculateTimeWindow(point model.PathPoint, departure time.Time, pickupDropoffMap map[string][2]int, idx int) ([2]int, int) {
	window := make([]int, 2)
	capacity := 0
	if request, ok := point.Owner().AsRequest(); ok {
		window[0] = scaleDownDuration(request.EarliestDepartureTime().Sub(departure))
		window[1] = scaleDownDuration(request.LatestArrivalTime().Sub(departure))
		switch point.PointType() {
		case enums.Pickup:
			{
				capacity = request.NumberOfRiders()
				pickupDropoff := pickupDropoffMap[request.ID()]
				pickupDropoff[0] = idx
			}

		case enums.Dropoff:
			{
				capacity = -request.NumberOfRiders()
				pickupDropoff := pickupDropoffMap[request.ID()]
				pickupDropoff[1] = idx
			}
		}
	} else if offer, ok := point.Owner().AsOffer(); ok {
		window[0] = 0
		window[1] = scaleDownDuration(offer.MaxEstimatedArrivalTime().Sub(departure))
	}

	return [2]int(window), capacity
}

// scaleDownDuration converts a time.Duration to an integer with reduced precision.
// This function divides nanoseconds by seconds to create smaller integers,
// making the values more manageable when sending to external solvers.
// The precision is reduced but sufficient for routing calculations.
func scaleDownDuration(t time.Duration) int {
	return int(t.Nanoseconds() / time.Second.Nanoseconds())
}

// scaleUpDuration restores the reduced-precision integer back to time.Duration.
// This function is used to convert the integer values returned by the solver
// back to Go's native time.Duration format with second-level precision.
func scaleUpDuration(t int) time.Duration {
	return time.Duration(t) * time.Second
}

// printMatrix prints the given 2D int matrix to stdout in a aligned grid.
func printMatrix(matrix [][]int) {
	for i := range matrix {
		for j := range matrix[i] {
			// adjust %4d to a larger width if your numbers are bigger
			fmt.Printf("%6d", matrix[i][j])
		}
		fmt.Println()
	}
}
