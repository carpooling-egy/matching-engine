package tests

import (
	"matching-engine/cmd/correcteness_test"
	"matching-engine/internal/adapter/routing"
	"matching-engine/internal/enums"
	"matching-engine/internal/model"
	"time"
)

func getTest3bData(engine routing.Engine) ([]*model.Offer, []*model.Request, map[string]*model.MatchingResult) {
	offers := make([]*model.Offer, 0)
	requests := make([]*model.Request, 0)

	// Create an offer with the specified attributes
	offerSource, _ := model.NewCoordinate(31.24587, 29.97458)
	offerDestination, _ := model.NewCoordinate(31.196, 29.90127)
	offerDepartureTime := correcteness_test.ParseTime("10:30")
	offerDetourDuration := time.Duration(8) * time.Minute // will be overwritten later
	offerCapacity := 1
	offerCurrentNumberOfRequests := 1
	offerSameGender := false
	offerGender := enums.Male
	// offerMaxEstimatedArrivalTime will be overwritten later
	offerMaxEstimatedArrivalTime := getMaxEstimatedArrivalTime(*offerSource, *offerDestination, offerDepartureTime, offerDetourDuration, engine)

	// Create a matched request for this offer
	// TODO correct pick up and dropoff coordinates
	matchedRequestSource, _ := model.NewCoordinate(31.23985, 29.96469)
	matchedRequestDestination, _ := model.NewCoordinate(31.23213, 29.9517)

	matchedRequestEarliestDepartureTime := offerDepartureTime.Add(-10 * time.Minute)
	// matchedRequestLatestArrivalTime will be overwritten later
	matchedRequestLatestArrivalTime := offerMaxEstimatedArrivalTime.Add(10 * time.Minute)
	matchedRequestMaxWalkingDuration := time.Duration(0) * time.Minute
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

	path := []model.PathPoint{
		*model.NewPathPoint(*offerSource, enums.Source, offerDepartureTime, offer, 0),
		*model.NewPathPoint(*offerDestination, enums.Destination, offerMaxEstimatedArrivalTime, offer, 0)}

	offer.SetPath(path)

	matchedRequestPickup, _, matchedRequestDropoff, _ := getRequestPointsAndDurations(engine, offer, matchedRequestSource, matchedRequestMaxWalkingDuration, matchedRequestDestination)

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

	// Create request 1
	request1Source, _ := model.NewCoordinate(31.22082, 29.94795)
	request1Destination, _ := model.NewCoordinate(31.21261, 29.9416)
	request1MaxWalkingDuration := time.Duration(0) * time.Minute
	request1NumberOfRiders := 1
	request1SameGender := true
	request1Gender := enums.Male
	pickup1, _, dropoff1, _ := getRequestPointsAndDurations(engine, offer, request1Source, request1MaxWalkingDuration, request1Destination)
	// Create request 2
	request2Source, _ := model.NewCoordinate(31.19699, 29.90388)
	request2Destination, _ := model.NewCoordinate(31.19661, 29.90294)
	request2MaxWalkingDuration := time.Duration(0) * time.Minute
	request2NumberOfRiders := 1
	request2SameGender := false
	request2Gender := enums.Male
	pickup2, _, dropoff2, _ := getRequestPointsAndDurations(engine, offer, request2Source, request2MaxWalkingDuration, request2Destination)

	cumulativeTimesWithoutRider := correcteness_test.GetCumulativeTimes([]model.Coordinate{*offerSource, *offerDestination}, offerDepartureTime, engine)
	cumulativeTimesWithRider := correcteness_test.GetCumulativeTimes([]model.Coordinate{*offerSource, *matchedRequestPickup, *matchedRequestDropoff, *pickup1, *dropoff1, *pickup2, *dropoff2, *offerDestination}, offerDepartureTime, engine)

	offerDetourDuration = cumulativeTimesWithRider[7] - cumulativeTimesWithoutRider[1] + 1*time.Second // adding 1 second to ensure the detour is valid

	offer.SetDetour(offerDetourDuration)
	offer.SetMaxEstimatedArrivalTime(offerDepartureTime.Add(cumulativeTimesWithoutRider[1]).Add(offerDetourDuration))
	matchedRequest.SetLatestArrivalTime(offer.MaxEstimatedArrivalTime().Add(10 * time.Minute))
	request1EarliestDepartureTime := offerDepartureTime.Add(-request1MaxWalkingDuration).Add(-1 * time.Minute)
	request1LatestArrivalTime := offerMaxEstimatedArrivalTime.Add(request1MaxWalkingDuration).Add(100 * time.Minute)

	request2EarliestDepartureTime := offerDepartureTime.Add(-request2MaxWalkingDuration).Add(-1 * time.Minute)
	request2LatestArrivalTime := offerMaxEstimatedArrivalTime.Add(request2MaxWalkingDuration).Add(100 * time.Minute)

	request1 := createRequest("3", "2", *request1Source, *request1Destination,
		request1EarliestDepartureTime, request1LatestArrivalTime,
		request1MaxWalkingDuration, request1NumberOfRiders,
		request1Gender, request1SameGender)

	// Add request 1 to the list of requests
	requests = append(requests, request1)

	request2 := createRequest("4", "3", *request2Source, *request2Destination,
		request2EarliestDepartureTime, request2LatestArrivalTime,
		request2MaxWalkingDuration, request2NumberOfRiders,
		request2Gender, request2SameGender)

	// Add request 2 to the list of requests
	requests = append(requests, request2)

	// Create expected results
	expectedResults := make(map[string]*model.MatchingResult)
	pickupPoint1, dropoffPoint1 := computeRequestPickupDropoffPoints(engine, offer, request1Source, request1MaxWalkingDuration, request1Destination, request1EarliestDepartureTime, request1, request1LatestArrivalTime)
	pickupOrder1, dropoffOrder1 := 3, 4

	pickupPoint2, dropoffPoint2 := computeRequestPickupDropoffPoints(engine, offer, request2Source, request2MaxWalkingDuration, request2Destination, request2EarliestDepartureTime, request2, request2LatestArrivalTime)
	pickupOrder2, dropoffOrder2 := 5, 6
	// Add the requests to the offer
	points := []*model.PathPoint{
		pickupPoint1,
		dropoffPoint1,
		pickupPoint2,
		dropoffPoint2,
	}
	pointsOrder := []int{
		pickupOrder1,
		dropoffOrder1,
		pickupOrder2,
		dropoffOrder2,
	}
	offerPath := addPointsToPath(engine, offer, pointsOrder, points)

	expectedResults[offer.ID()] = model.NewMatchingResult(
		offer.ID(),
		offer.UserID(),
		[]*model.Request{request1, request2},
		offerPath,
		offer.CurrentNumberOfRequests()+2,
	)
	return offers, requests, expectedResults
}
