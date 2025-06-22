package tests

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"go.uber.org/dig"
	"matching-engine/cmd/correcteness_test"
	"matching-engine/internal/adapter/routing"
	"matching-engine/internal/adapter/valhalla"
	"matching-engine/internal/app/config"
	"matching-engine/internal/app/di"
	"matching-engine/internal/app/di/utils"
	"matching-engine/internal/enums"
	"matching-engine/internal/model"
	matcher2 "matching-engine/internal/service/matcher"
	"matching-engine/internal/service/pickupdropoffservice"
	"sort"
	"testing"
	"time"
)

type MatchedRequest struct {
	request      *model.Request
	pickupCoord  *model.Coordinate
	pickupOrder  int
	dropoffCoord *model.Coordinate
	dropoffOrder int
}

func TestCorrecteness(t *testing.T) {
	config.ConfigureLogging()

	// Create a mock routing engine
	engine, err := valhalla.NewValhalla()
	if err != nil {
		t.Fatalf("Failed to create Valhalla engine: %v", err)
	}
	tests := []struct {
		name     string
		testFunc func(engine routing.Engine) ([]*model.Offer, []*model.Request, map[string]*model.MatchingResult)
	}{
		{
			name:     "Test1ai",
			testFunc: getTest1aiData,
		},
		{
			name:     "Test1aii",
			testFunc: getTest1aiiData,
		},
		{
			name:     "Test1b",
			testFunc: getTest1bData,
		},
		{
			name:     "Test1ci",
			testFunc: getTest1ciData,
		},
		{
			name:     "Test1cii",
			testFunc: getTest1ciiData,
		},
		{
			name:     "Test1di",
			testFunc: getTest1diData,
		},
		{
			name:     "Test1dii",
			testFunc: getTest1diiData,
		},
		{
			name:     "Test1e",
			testFunc: getTest1eData,
		},
		{
			name:     "Test1fi",
			testFunc: getTest1fiData,
		},
		{
			name:     "Test1fii",
			testFunc: getTest1fiiData,
		},
		{
			name:     "Test1fiii",
			testFunc: getTest1fiiiData,
		},
		{
			name:     "Test1fiv",
			testFunc: getTest1fivData,
		},
		{
			name:     "Test2a",
			testFunc: getTest2aData,
		},
		{
			name:     "Test3a",
			testFunc: getTest3aData,
		},
		{
			name:     "Test3b",
			testFunc: getTest3bData,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			offers, requests, expectedResults := tt.testFunc(engine)
			if len(offers) == 0 || len(requests) == 0 {
				t.Fatalf("No offers or requests generated for test %s", tt.name)
			}
			results, err := runMatcher(offers, requests)
			if err != nil {
				t.Fatalf("Matcher failed for test %s: %v", tt.name, err)
			}
			if results == nil && expectedResults == nil {
				return // Both results are nil, which is acceptable
			}
			if len(results) != len(expectedResults) {
				t.Fatalf("Expected %d results, got %d for test %s", len(expectedResults), len(results), tt.name)
			}
			if !compareResults(results, expectedResults) {
				t.Fatalf("Results do not match expected for test %s", tt.name)
			}
		})
	}
}

func addPointsToPath(engine routing.Engine, offer *model.Offer, pointsOrder []int, points []*model.PathPoint) []model.PathPoint {
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
		// Clamp insertion index to end if it's too large
		order := pointsOrder[i]
		inserts[i] = pointInsert{order, points[i]}
	}

	// Sort by insertion order
	sort.Slice(inserts, func(i, j int) bool {
		return inserts[i].order < inserts[j].order
	})

	fmt.Println(inserts)

	newPath := make([]model.PathPoint, 0, newLen)

	origIndex, insertIndex := 0, 0

	for i := 0; i < newLen; i++ {
		//fmt.Println(i, origIndex, insertIndex, len(originalPath), len(inserts))
		if insertIndex < len(inserts) && inserts[insertIndex].order == i {
			//fmt.Println("Inserting point at index", i, ":", inserts[insertIndex].point)
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
	newPath = calculateExpectedArrivalTimes(newPath, offer.DepartureTime(), engine)

	return newPath
}

func getRequestPointsAndDurations(engine routing.Engine, offer *model.Offer, source *model.Coordinate, walkingDuration time.Duration, destination *model.Coordinate) (*model.Coordinate, time.Duration, *model.Coordinate, time.Duration) {
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

func computeRequestPickupDropoffPoints(engine routing.Engine, offer *model.Offer, requestSource *model.Coordinate, requestMaxWalkingDuration time.Duration, requestDestination *model.Coordinate, requestEarliestDepartureTime time.Time, request *model.Request, requestLatestArrivalTime time.Time) (*model.PathPoint, *model.PathPoint) {
	pickupCoord, pickupDuration, dropoffCoord, dropoffDuration := getRequestPointsAndDurations(
		engine, offer, requestSource, requestMaxWalkingDuration, requestDestination)
	pickupPoint := model.NewPathPoint(
		*pickupCoord, enums.Pickup, requestEarliestDepartureTime, request, pickupDuration)
	dropoffPoint := model.NewPathPoint(
		*dropoffCoord, enums.Dropoff, requestLatestArrivalTime, request, dropoffDuration)
	return pickupPoint, dropoffPoint
}

func compareResults(results []*model.MatchingResult, expectedResults map[string]*model.MatchingResult) bool {
	if len(results) != len(expectedResults) {
		log.Debug().Msgf("Results length mismatch: got %d, expected %d\n", len(results), len(expectedResults))
		return false
	}
	for _, result := range results {
		expectedResult := expectedResults[result.OfferID()]
		if result.UserID() != expectedResult.UserID() || result.OfferID() != expectedResult.OfferID() {
			log.Debug().Msgf("UserID or OfferID mismatch: got %s/%s, expected %s/%s\n",
				result.UserID(), result.OfferID(), expectedResult.UserID(), expectedResult.OfferID())
			return false
		}
		if len(result.AssignedMatchedRequests()) != len(expectedResult.AssignedMatchedRequests()) {
			log.Debug().Msgf("Assigned matched requests length mismatch: got %d, expected %d\n",
				len(result.AssignedMatchedRequests()), len(expectedResult.AssignedMatchedRequests()))
			return false
		}
		if len(result.NewPath()) != len(expectedResult.NewPath()) {
			log.Debug().Msgf("Path length mismatch: got %d, expected %d\n",
				len(result.NewPath()), len(expectedResult.NewPath()))
			return false
		}
		if result.CurrentNumberOfRequests() != expectedResult.CurrentNumberOfRequests() {
			log.Debug().Msgf("Number of requests mismatch: got %d, expected %d\n",
				result.CurrentNumberOfRequests(), expectedResult.CurrentNumberOfRequests())
			return false
		}
		for _, req := range result.AssignedMatchedRequests() {
			matchedRequests := false
			for _, expectedReq := range expectedResult.AssignedMatchedRequests() {
				if req.ID() == expectedReq.ID() &&
					req.Source().Equal(expectedReq.Source()) &&
					req.Destination().Equal(expectedReq.Destination()) &&
					req.EarliestDepartureTime().Equal(expectedReq.EarliestDepartureTime()) &&
					req.LatestArrivalTime().Equal(expectedReq.LatestArrivalTime()) &&
					req.MaxWalkingDurationMinutes() == expectedReq.MaxWalkingDurationMinutes() &&
					req.NumberOfRiders() == expectedReq.NumberOfRiders() &&
					req.Preferences() == expectedReq.Preferences() {
					matchedRequests = true
				}
			}
			if !matchedRequests {
				log.Debug().
					Str("requestID", req.ID()).
					Msg("Assigned matched request not found in expected results")
				return false
			}
		}
		for i, point := range result.NewPath() {
			if i >= len(expectedResult.NewPath()) {
				log.Debug().
					Int("point", i).
					Msg("Point index out of bounds in expected result")
				return false
			}
			expectedPoint := expectedResult.NewPath()[i]
			if !point.Coordinate().Equal(expectedPoint.Coordinate()) ||
				point.PointType() != expectedPoint.PointType() ||
				!checkTimeOverlap(point.ExpectedArrivalTime(), expectedPoint.ExpectedArrivalTime(), 10*time.Second) ||
				point.WalkingDuration() != expectedPoint.WalkingDuration() ||
				!checkOwnerMatch(point, expectedPoint) {
				log.Debug().
					Int("point", i).
					Msg("Point mismatch: ")
				return false
			}
		}
	}
	return true
}

func checkOwnerMatch(point model.PathPoint, expectedPoint model.PathPoint) bool {
	_, isRequest := point.Owner().AsRequest()
	_, isExpectedRequest := expectedPoint.Owner().AsRequest()
	if isRequest != isExpectedRequest {
		return false
	}
	if point.GetOwnerID() != expectedPoint.GetOwnerID() {
		return false
	}
	return true
}

func checkTimeOverlap(time1, time2 time.Time, tolerance time.Duration) bool {
	// Check if the two times are within the specified tolerance
	return time1.After(time2.Add(-tolerance)) && time1.Before(time2.Add(tolerance))
}

func createOffer(userID, id string, source, destination model.Coordinate, departureTime time.Time,
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

func getMaxEstimatedArrivalTime(source model.Coordinate, destination model.Coordinate, departureTime time.Time, detour time.Duration, engine routing.Engine) time.Time {
	directCoords := []model.Coordinate{source, destination}
	directTimes := correcteness_test.GetCumulativeTimes(directCoords, departureTime, engine)
	return departureTime.Add(detour).Add(directTimes[1])
}

func createRequest(userID, id string, source, destination model.Coordinate, earliestDepartureTime, latestArrivalTime time.Time,
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

func createPath(offer *model.Offer, matchedRequests []*MatchedRequest, engine routing.Engine) []model.PathPoint {
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
	path = calculateExpectedArrivalTimes(path, offer.DepartureTime(), engine)
	return path
}

func calculateExpectedArrivalTimes(path []model.PathPoint, departureTime time.Time, engine routing.Engine) []model.PathPoint {
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

func runMatcher(offers []*model.Offer, requests []*model.Request) ([]*model.MatchingResult, error) {

	c := dig.New()

	// register all dependencies for matching services
	di.RegisterGeoServices(c)
	di.RegisterPickupDropoffServices(c)
	di.RegisterTimeMatrixServices(c)
	di.RegisterPathServices(c)
	di.RegisterCheckers(c)
	di.RegisterMatchingServices(c)
	utils.Must(c.Provide(valhalla.NewValhalla))

	var matches []*model.MatchingResult
	var matchErr error
	err := c.Invoke(func(matcher *matcher2.Matcher) {
		matches, matchErr = matcher.Match(offers, requests)
	})
	if err != nil {
		panic("Failed to invoke matcher in the container: " + err.Error())
	}
	return matches, matchErr
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
