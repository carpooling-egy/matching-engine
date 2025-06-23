package tests

import (
	"context"
	"matching-engine/cmd/correcteness_test"
	"matching-engine/internal/adapter/routing"
	"matching-engine/internal/enums"
	"matching-engine/internal/model"
	"matching-engine/internal/service/pickupdropoffservice"
	"sort"
	"time"
)

type MatchedRequest struct {
	request      *model.Request
	pickupCoord  *model.Coordinate
	pickupOrder  int
	dropoffCoord *model.Coordinate
	dropoffOrder int
}

func AddPointsToPath(engine routing.Engine, offer *model.Offer, pointsOrder []int, points []*model.PathPoint) []model.PathPoint {
	if len(pointsOrder) != len(points) {
		panic("pointsOrder and points must have the same length")
	}

	originalPath := offer.Path()
	newLen := len(originalPath) + len(points)

	// Pair points with their insertion order
	type pointInsert struct {
		order int
		point *model.PathPoint
	}
	inserts := make([]pointInsert, len(pointsOrder))
	for i := range pointsOrder {
		order := pointsOrder[i]
		inserts[i] = pointInsert{order, points[i]}
	}

	// Sort by insertion order
	sort.Slice(inserts, func(i, j int) bool {
		return inserts[i].order < inserts[j].order
	})

	newPath := make([]model.PathPoint, 0, newLen)

	origIndex, insertIndex := 0, 0

	for i := 0; i < newLen; i++ {
		if insertIndex < len(inserts) && inserts[insertIndex].order == i {
			// Insert the new point at the correct position
			newPath = append(newPath, *inserts[insertIndex].point)
			insertIndex++
		} else {
			// Add the original path point
			if origIndex < len(originalPath) {
				newPath = append(newPath, originalPath[origIndex])
				origIndex++
			} else {
				panic("Original path has fewer points than expected")
			}
		}
	}

	// Recalculate arrival times
	newPath = CalculateExpectedArrivalTimes(newPath, offer.DepartureTime(), engine)

	return newPath
}

func GetRequestPointsAndDurations(engine routing.Engine, offer *model.Offer, source *model.Coordinate, walkingDuration time.Duration, destination *model.Coordinate) (*model.Coordinate, time.Duration, *model.Coordinate, time.Duration) {
	pickup, pickupDuration, dropoff, dropoffDuration := correcteness_test.GetPickupDropoffPointsAndDurations(
		engine, offer, source, walkingDuration, destination)
	if pickupDuration > walkingDuration {
		pickupCoord, err := engine.SnapPointToRoad(context.Background(), source)
		if err != nil {
			pickupCoord = source
		}
		pickup = pickupCoord
		pickupDuration = 0
	}
	if dropoffDuration > walkingDuration {
		dropoffCoord, err := engine.SnapPointToRoad(context.Background(), destination)
		if err != nil {
			dropoffCoord = destination
		}
		dropoff = dropoffCoord
		dropoffDuration = 0
	}
	return pickup, pickupDuration, dropoff, dropoffDuration
}

func ComputeRequestPickupDropoffPoints(engine routing.Engine, offer *model.Offer, requestSource *model.Coordinate, requestMaxWalkingDuration time.Duration, requestDestination *model.Coordinate, requestEarliestDepartureTime time.Time, request *model.Request, requestLatestArrivalTime time.Time) (*model.PathPoint, *model.PathPoint) {
	pickupCoord, pickupDuration, dropoffCoord, dropoffDuration := GetRequestPointsAndDurations(
		engine, offer, requestSource, requestMaxWalkingDuration, requestDestination)
	pickupPoint := model.NewPathPoint(
		*pickupCoord, enums.Pickup, requestEarliestDepartureTime, request, pickupDuration)
	dropoffPoint := model.NewPathPoint(
		*dropoffCoord, enums.Dropoff, requestLatestArrivalTime, request, dropoffDuration)
	return pickupPoint, dropoffPoint
}

func CreateOffer(userID, id string, source, destination model.Coordinate, departureTime time.Time,
	detourDurMins time.Duration, capacity, currentNumberOfRequests int, gender enums.Gender, sameGender bool,
	maxEstimatedArrivalTime time.Time, matchedRequests []*model.Request) *model.Offer {
	preference := *model.NewPreference(gender, sameGender)
	return model.NewOffer(
		id,
		userID,
		source,
		destination,
		departureTime,
		detourDurMins,
		capacity,
		preference,
		maxEstimatedArrivalTime,
		currentNumberOfRequests,
		nil,
		matchedRequests,
	)
}

func GetMaxEstimatedArrivalTime(source model.Coordinate, destination model.Coordinate, departureTime time.Time, detour time.Duration, engine routing.Engine) time.Time {
	directCoords := []model.Coordinate{source, destination}
	directTimes := correcteness_test.GetCumulativeTimes(directCoords, departureTime, engine)
	return departureTime.Add(detour).Add(directTimes[1])
}

func CreateRequest(userID, id string, source, destination model.Coordinate, earliestDepartureTime, latestArrivalTime time.Time,
	maxWalkingDurationMinutes time.Duration, numberOfRiders int, gender enums.Gender, sameGender bool) *model.Request {
	preference := *model.NewPreference(gender, sameGender)
	return model.NewRequest(
		id,
		userID,
		source,
		destination,
		earliestDepartureTime,
		latestArrivalTime,
		maxWalkingDurationMinutes,
		numberOfRiders,
		preference,
	)
}

func CreatePath(offer *model.Offer, matchedRequests []*MatchedRequest, engine routing.Engine) []model.PathPoint {
	path := make([]model.PathPoint, len(matchedRequests)*2+2) // 2 points for the offer source and destination, 2 points for each request pickup and dropoff
	path[0] = *model.NewPathPoint(*offer.Source(), enums.Source, offer.DepartureTime(), offer, 0)
	path[len(path)-1] = *model.NewPathPoint(*offer.Destination(), enums.Destination, offer.MaxEstimatedArrivalTime(), offer, 0)
	walkingTimeCalculator := pickupdropoffservice.NewWalkingTimeCalculator(engine)
	for _, matchedReq := range matchedRequests {
		pickupPoint := model.NewPathPoint(
			*matchedReq.pickupCoord, enums.Pickup, matchedReq.request.EarliestDepartureTime(), matchedReq.request, 0)
		dropoffPoint := model.NewPathPoint(
			*matchedReq.dropoffCoord, enums.Dropoff, matchedReq.request.LatestArrivalTime(), matchedReq.request, 0)
		pickupWalkingDuration, dropoffWalkingDuration, err := walkingTimeCalculator.ComputeWalkingDurations(context.Background(), matchedReq.request, pickupPoint, dropoffPoint)
		if err != nil {
			panic("Failed to compute walking durations: " + err.Error())
		}
		pickupPoint.SetWalkingDuration(pickupWalkingDuration)
		dropoffPoint.SetWalkingDuration(dropoffWalkingDuration)
		path[matchedReq.pickupOrder] = *pickupPoint
		path[matchedReq.dropoffOrder] = *dropoffPoint
	}
	// Calculate travel times for the path
	path = CalculateExpectedArrivalTimes(path, offer.DepartureTime(), engine)
	return path
}

func CalculateExpectedArrivalTimes(path []model.PathPoint, departureTime time.Time, engine routing.Engine) []model.PathPoint {
	coords := make([]model.Coordinate, len(path))
	for i, p := range path {
		coords[i] = *p.Coordinate()
	}
	drivingTimes := correcteness_test.GetCumulativeTimes(coords, departureTime, engine)
	for i := range path {
		path[i].SetExpectedArrivalTime(departureTime.Add(drivingTimes[i]))
	}
	return path
}
