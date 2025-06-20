package tests

import (
	"matching-engine/cmd/correcteness_test"
	"matching-engine/internal/adapter/routing"
	"matching-engine/internal/adapter/valhalla"
	"matching-engine/internal/enums"
	"matching-engine/internal/model"
	matcher2 "matching-engine/internal/service/matcher"
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

func getTest1aiData(engine routing.Engine) ([]*model.Offer, []*model.Request, []*model.MatchingResult) {
	offers := make([]*model.Offer, 0)
	requests := make([]*model.Request, 0)

	// Create an offer with the specified attributes
	offerSource, _ := model.NewCoordinate(31.2544088039743, 29.97376045816)
	offerDestination, _ := model.NewCoordinate(31.20611644667, 29.9248733439259)
	offerDepartureTime, _ := time.Parse("15:04", "10:30")
	offerDetourDuration := time.Duration(30)
	offerCapacity := 3
	offerCurrentNumberOfRequests := 1
	offerSameGender := false
	offerGender := enums.Male
	offerMaxEstimatedArrivalTime := getMaxEstimatedArrivalTime(*offerSource, *offerDestination, offerDepartureTime, offerDetourDuration, engine)

	// Create a matched request for this offer
	matchedRequestSource, _ := model.NewCoordinate(31.2544088, 29.97376046)
	matchedRequestDestination, _ := model.NewCoordinate(31.20611645, 29.92487334)
	matchedRequestPickup, _ := model.NewCoordinate(31.2544088, 29.97376046)
	matchedRequestDropoff, _ := model.NewCoordinate(31.20611645, 29.92487334)
	matchedRequestEarliestDepartureTime, _ := time.Parse("15:04:05", "10:20:00")
	matchedRequestLatestArrivalTime, _ := time.Parse("15:04", "11:20")
	matchedRequestMaxWalkingDuration := time.Duration(0)
	matchedRequestNumberOfRiders := 1
	matchedRequestSameGender := true
	matchedRequestGender := enums.Male
	matchedRequest := createRequest("2", "1", *matchedRequestSource, *matchedRequestDestination,
		matchedRequestEarliestDepartureTime, matchedRequestLatestArrivalTime,
		matchedRequestMaxWalkingDuration, matchedRequestNumberOfRiders,
		matchedRequestGender, matchedRequestSameGender)
	offerRequests := []*model.Request{matchedRequest}

	offer := createOffer("1", "1", *offerSource, *offerDestination, offerDepartureTime,
		offerDetourDuration, offerCapacity, offerCurrentNumberOfRequests, offerGender,
		offerSameGender, offerMaxEstimatedArrivalTime, offerRequests)

	// Create a matched request with pickup and dropoff coordinates
	matchedReq := &MatchedRequest{
		request:      matchedRequest,
		pickupCoord:  matchedRequestPickup,
		pickupOrder:  1,
		dropoffCoord: matchedRequestDropoff,
		dropoffOrder: 2,
	}
	offer.SetPath(createPath(offer, []*MatchedRequest{matchedReq}, engine))
	// Add the offer to the list of offers
	offers = append(offers, offer)

	// Create another request
	requestSource, _ := model.NewCoordinate(31.2544088, 29.97376046)
	requestDestination, _ := model.NewCoordinate(31.20611645, 29.92487334)
	requestEarliestDepartureTime, _ := time.Parse("15:04:05", "09:20:00")
	requestLatestArrivalTime, _ := time.Parse("15:04", "10:20")
	requestMaxWalkingDuration := time.Duration(0)
	requestNumberOfRiders := 2
	requestSameGender := true
	requestGender := enums.Male

	request := createRequest("3", "2", *requestSource, *requestDestination,
		requestEarliestDepartureTime, requestLatestArrivalTime,
		requestMaxWalkingDuration, requestNumberOfRiders,
		requestGender, requestSameGender)

	// Add the request to the list of requests
	requests = append(requests, request)

	// TODO: Add Expected Results for the matcher
	expectedResults := make([]*model.MatchingResult, 0)
	return offers, requests, expectedResults
}

func TestCorrecteness(t *testing.T) {
	// Create a mock routing engine
	engine, err := valhalla.NewValhalla()
	if err != nil {
		t.Fatalf("Failed to create Valhalla engine: %v", err)
	}
	tests := []struct {
		name     string
		testFunc func(engine routing.Engine) ([]*model.Offer, []*model.Request, []*model.MatchingResult)
	}{
		{
			name:     "Test1ai",
			testFunc: getTest1aiData,
		},
	}

	// TODO: Create a matcher
	matcher := matcher2.NewMatcher(nil, nil, nil, nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			offers, requests, expectedResults := tt.testFunc(engine)
			if len(offers) == 0 || len(requests) == 0 {
				t.Fatalf("No offers or requests generated for test %s", tt.name)
			}
			results, err := matcher.Match(offers, requests)
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

func compareResults(results []*model.MatchingResult, expectedResults []*model.MatchingResult) bool {
	if len(results) != len(expectedResults) {
		return false
	}
	for i, result := range results {
		expectedResult := expectedResults[i]
		if result.UserID() != expectedResult.UserID() || result.OfferID() != expectedResult.OfferID() {
			return false
		}
		if len(result.AssignedMatchedRequests()) != len(expectedResult.AssignedMatchedRequests()) {
			return false
		}

	}
	return true
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
	for _, matchedReq := range matchedRequests {
		pickupPoint := model.NewPathPoint(
			*matchedReq.pickupCoord, enums.Pickup, matchedReq.request.EarliestDepartureTime(), matchedReq.request, 0)
		dropoffPoint := model.NewPathPoint(
			*matchedReq.dropoffCoord, enums.Dropoff, matchedReq.request.LatestArrivalTime(), matchedReq.request, 0)
		path[matchedReq.pickupOrder] = *pickupPoint
		path[matchedReq.dropoffOrder] = *dropoffPoint
	}
	// Calculate travel times for the path
	path = calculateTravelTimes(path, offer.DepartureTime(), engine)
	return path
}

func calculateTravelTimes(path []model.PathPoint, departureTime time.Time, engine routing.Engine) []model.PathPoint {
	coords := make([]model.Coordinate, len(path))
	for i, p := range path {
		coords[i] = *p.Coordinate()
	}
	drivingTimes := correcteness_test.GetCumulativeTimes(coords, departureTime, engine)
	for i, p := range path {
		p.SetExpectedArrivalTime(departureTime.Add(drivingTimes[i]))
	}
	return path
}
