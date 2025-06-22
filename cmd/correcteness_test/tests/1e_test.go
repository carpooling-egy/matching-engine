package tests

import (
	"matching-engine/cmd/correcteness_test"
	"matching-engine/internal/adapter/routing"
	"matching-engine/internal/enums"
	"matching-engine/internal/model"
	"time"
)

func getTest1eData(engine routing.Engine) ([]*model.Offer, []*model.Request, map[string]*model.MatchingResult) {
	offers := make([]*model.Offer, 0)
	requests := make([]*model.Request, 0)

	// Create an offer with the specified attributes
	offerSource, _ := model.NewCoordinate(31.2544088039743, 29.97376045816)
	offerDestination, _ := model.NewCoordinate(31.20611644667, 29.9248733439259)
	offerDepartureTime := correcteness_test.ParseTime("10:30")
	offerDetourDuration := time.Duration(30) * time.Minute // will be overwritten this is a placeholder value
	offerCapacity := 3
	offerCurrentNumberOfRequests := 1
	offerSameGender := true
	offerGender := enums.Male
	offerMaxEstimatedArrivalTime := getMaxEstimatedArrivalTime(*offerSource, *offerDestination, offerDepartureTime, offerDetourDuration, engine) // will be overwritten this is a placeholder value

	// Create a matched request for this offer
	matchedRequestSource, _ := model.NewCoordinate(31.2544088, 29.97376046)
	matchedRequestDestination, _ := model.NewCoordinate(31.20611645, 29.92487334)
	matchedRequestPickup, _ := model.NewCoordinate(31.2544088, 29.97376046)
	matchedRequestDropoff, _ := model.NewCoordinate(31.20611645, 29.92487334)
	matchedRequestEarliestDepartureTime := offerDepartureTime.Add(-10 * time.Minute)
	matchedRequestLatestArrivalTime := offerMaxEstimatedArrivalTime.Add(10 * time.Minute)
	matchedRequestMaxWalkingDuration := time.Duration(0) * time.Minute
	matchedRequestNumberOfRiders := 1
	matchedRequestSameGender := false
	matchedRequestGender := enums.Male
	matchedRequest := createRequest("2", "1", *matchedRequestSource, *matchedRequestDestination,
		matchedRequestEarliestDepartureTime, matchedRequestLatestArrivalTime,
		matchedRequestMaxWalkingDuration, matchedRequestNumberOfRiders,
		matchedRequestGender, matchedRequestSameGender)
	offerRequests := []*model.Request{matchedRequest}

	// Create the offer
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
	requestSource, _ := model.NewCoordinate(31.241047, 29.962955)
	requestDestination, _ := model.NewCoordinate(31.223313, 29.933295)
	requestMaxWalkingDuration := time.Duration(0) * time.Minute

	pickup, _, dropoff, _ := getRequestPointsAndDurations(engine, offer, requestSource, requestMaxWalkingDuration, requestDestination)
	cumulativeTimesWithoutRider := correcteness_test.GetCumulativeTimes([]model.Coordinate{*offerSource, *offerDestination}, offerDepartureTime, engine)
	cumulativeTimesWithRider := correcteness_test.GetCumulativeTimes([]model.Coordinate{*offerSource, *pickup, *dropoff, *offerDestination}, offerDepartureTime, engine)

	// overwrite offer detour, maxEstimated arrival time && matchedRequestLatestArrivalTime
	offer.SetDetour((cumulativeTimesWithRider[3] - cumulativeTimesWithoutRider[1]) / 2)
	offer.SetMaxEstimatedArrivalTime(getMaxEstimatedArrivalTime(*offerSource, *offerDestination, offerDepartureTime, offer.DetourDurationMinutes(), engine))
	matchedRequest.SetLatestArrivalTime(offer.MaxEstimatedArrivalTime().Add(10 * time.Minute))

	requestEarliestDepartureTime := offerDepartureTime.Add(-requestMaxWalkingDuration).Add(-1 * time.Minute)
	requestLatestArrivalTime := offerMaxEstimatedArrivalTime.Add(requestMaxWalkingDuration).Add(1 * time.Minute)
	requestNumberOfRiders := 2
	requestSameGender := true
	requestGender := enums.Male

	request := createRequest("3", "2", *requestSource, *requestDestination,
		requestEarliestDepartureTime, requestLatestArrivalTime,
		requestMaxWalkingDuration, requestNumberOfRiders,
		requestGender, requestSameGender)

	// Add the request to the list of requests
	requests = append(requests, request)

	// Create expected results
	expectedResults := make(map[string]*model.MatchingResult)
	return offers, requests, expectedResults
}
