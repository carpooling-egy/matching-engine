package tests

import (
	"matching-engine/internal/adapter/routing"
	"matching-engine/internal/enums"
	"matching-engine/internal/model"
	"time"
)

func getTest2a(engine routing.Engine) ([]*model.Offer, []*model.Request, map[string]*model.MatchingResult) {
	offers := make([]*model.Offer, 0)
	requests := make([]*model.Request, 0)

	// Create an offer with the specified attributes
	offerSource, _ := model.NewCoordinate(31.2460735985739, 29.9744554984058)
	offerDestination, _ := model.NewCoordinate(31.2068412085851, 29.9246876930902)
	offerDepartureTime, _ := time.Parse("15:04", "10:30")
	offerDetourDuration := time.Duration(8)
	offerCapacity := 1
	offerCurrentNumberOfRequests := 1
	offerSameGender := false
	offerGender := enums.Male
	offerMaxEstimatedArrivalTime := getMaxEstimatedArrivalTime(*offerSource, *offerDestination, offerDepartureTime, offerDetourDuration, engine)

	// Create a matched request for this offer
	matchedRequestSource, _ := model.NewCoordinate(31.22082087, 29.94795413)
	matchedRequestDestination, _ := model.NewCoordinate(31.208936, 29.933419)
	matchedRequestPickup, _ := model.NewCoordinate(31.22082087, 29.94795413)
	matchedRequestDropoff, _ := model.NewCoordinate(31.208936, 29.933419)
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
	requestSource, _ := model.NewCoordinate(31.208936, 29.933419)
	requestDestination, _ := model.NewCoordinate(31.20773291, 29.92687263)
	requestEarliestDepartureTime, _ := time.Parse("15:04:05", "10:20:00")
	requestLatestArrivalTime, _ := time.Parse("15:04", "11:20")
	requestMaxWalkingDuration := time.Duration(0)
	requestNumberOfRiders := 1
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
	pickupPoint, dropoffPoint := computeRequestPickupDropoffPoints(engine, offer, requestSource, requestMaxWalkingDuration, requestDestination, requestEarliestDepartureTime, request, requestLatestArrivalTime)
	pickupOrder, dropoffOrder := 4, 5
	offerPath := addPointsToPath(engine, offer, pickupOrder, dropoffOrder, pickupPoint, dropoffPoint)

	expectedResults[offer.ID()] = model.NewMatchingResult(
		offer.ID(),
		offer.UserID(),
		[]*model.Request{request},
		offerPath,
		offer.CurrentNumberOfRequests()+1,
	)
	return offers, requests, expectedResults
}
