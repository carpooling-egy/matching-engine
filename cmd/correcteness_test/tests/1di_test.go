package tests

import (
	"matching-engine/cmd/correcteness_test"
	"matching-engine/internal/adapter/routing"
	"matching-engine/internal/enums"
	"matching-engine/internal/model"
	"time"
)

func getTest1diData(engine routing.Engine) ([]*model.Offer, []*model.Request, map[string]*model.MatchingResult) {
	offers := make([]*model.Offer, 0)
	requests := make([]*model.Request, 0)

	// Create an offer with the specified attributes
	offerSource, _ := model.NewCoordinate(31.2544088039743, 29.97376045816)
	offerDestination, _ := model.NewCoordinate(31.20611644667, 29.9248733439259)
	offerDepartureTime := correcteness_test.ParseTime("10:30")
	offerDetourDuration := time.Duration(30) * time.Minute
	offerCapacity := 3
	offerCurrentNumberOfRequests := 1
	offerSameGender := true
	offerGender := enums.Male
	offerMaxEstimatedArrivalTime := GetMaxEstimatedArrivalTime(*offerSource, *offerDestination, offerDepartureTime, offerDetourDuration, engine)

	// Create a matched request for this offer
	matchedRequestSource, _ := model.NewCoordinate(31.2544088, 29.97376046)
	matchedRequestDestination, _ := model.NewCoordinate(31.20611645, 29.92487334)
	matchedRequestEarliestDepartureTime := offerDepartureTime.Add(-10 * time.Minute)
	matchedRequestLatestArrivalTime := offerMaxEstimatedArrivalTime.Add(10 * time.Minute)
	matchedRequestMaxWalkingDuration := time.Duration(0) * time.Minute
	matchedRequestNumberOfRiders := 1
	matchedRequestSameGender := false
	matchedRequestGender := enums.Male
	matchedRequest := CreateRequest("2", "1", *matchedRequestSource, *matchedRequestDestination,
		matchedRequestEarliestDepartureTime, matchedRequestLatestArrivalTime,
		matchedRequestMaxWalkingDuration, matchedRequestNumberOfRiders,
		matchedRequestGender, matchedRequestSameGender)
	offerRequests := []*model.Request{matchedRequest}

	offer := CreateOffer("1", "1", *offerSource, *offerDestination, offerDepartureTime,
		offerDetourDuration, offerCapacity, offerCurrentNumberOfRequests, offerGender,
		offerSameGender, offerMaxEstimatedArrivalTime, offerRequests)

	path := []model.PathPoint{
		*model.NewPathPoint(*offerSource, enums.Source, offerDepartureTime, offer, 0),
		*model.NewPathPoint(*offerDestination, enums.Destination, offerMaxEstimatedArrivalTime, offer, 0)}

	offer.SetPath(path)

	matchedRequestPickup, pickupDuration, matchedRequestDropoff, dropoffDuartion := GetRequestPointsAndDurations(engine, offer, matchedRequestSource, matchedRequestMaxWalkingDuration, matchedRequestDestination)

	// Create a matched request with pickup and dropoff coordinates
	matchedReq := &MatchedRequest{
		request:         matchedRequest,
		pickupCoord:     matchedRequestPickup,
		pickupDuration:  pickupDuration,
		pickupOrder:     1,
		dropoffCoord:    matchedRequestDropoff,
		dropoffDuration: dropoffDuartion,
		dropoffOrder:    2,
	}
	offer.SetPath(CreatePath(offer, []*MatchedRequest{matchedReq}, engine))
	// Add the offer to the list of offers
	offers = append(offers, offer)

	// Create another request
	requestSource, _ := model.NewCoordinate(31.241047, 29.962955)
	requestDestination, _ := model.NewCoordinate(31.223313, 29.933295)
	requestMaxWalkingDuration := time.Duration(0) * time.Minute

	pickup, _, _, _ := GetRequestPointsAndDurations(engine, offer, requestSource, requestMaxWalkingDuration, requestDestination)
	cumulativeTimes := correcteness_test.GetCumulativeTimes([]model.Coordinate{*offerSource, *pickup}, offerDepartureTime, engine)
	driverToPickupDuration := cumulativeTimes[1]
	requestEarliestDepartureTime := offerDepartureTime.Add(driverToPickupDuration).Add(5 * time.Minute) // 5 min is just an offset
	requestLatestArrivalTime := offerMaxEstimatedArrivalTime.Add(requestMaxWalkingDuration).Add(5 * time.Minute)
	requestNumberOfRiders := 2
	requestSameGender := true
	requestGender := enums.Male

	request := CreateRequest("3", "2", *requestSource, *requestDestination,
		requestEarliestDepartureTime, requestLatestArrivalTime,
		requestMaxWalkingDuration, requestNumberOfRiders,
		requestGender, requestSameGender)

	// Add the request to the list of requests
	requests = append(requests, request)

	// Create expected results
	expectedResults := make(map[string]*model.MatchingResult)
	return offers, requests, expectedResults
}
