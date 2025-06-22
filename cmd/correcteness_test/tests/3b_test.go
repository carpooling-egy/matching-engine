package tests

import (
	"fmt"
	"matching-engine/cmd/correcteness_test"
	"matching-engine/internal/adapter/routing"
	"matching-engine/internal/enums"
	"matching-engine/internal/model"
	"time"
)

func getTest3b(engine routing.Engine) ([]*model.Offer, []*model.Request, map[string]*model.MatchingResult) {
	offers := make([]*model.Offer, 0)
	requests := make([]*model.Request, 0)

	// Create an offer with the specified attributes
	offerSource, _ := model.NewCoordinate(31.2460735985739, 29.9744554984058)
	offerDestination, _ := model.NewCoordinate(31.2068412085851, 29.9246876930902)
	offerDepartureTime := correcteness_test.ParseTime("10:30")
	offerDetourDuration := time.Duration(8) * time.Minute
	offerCapacity := 3
	offerCurrentNumberOfRequests := 1
	offerSameGender := false
	offerGender := enums.Male
	offerMaxEstimatedArrivalTime := getMaxEstimatedArrivalTime(*offerSource, *offerDestination, offerDepartureTime, offerDetourDuration, engine)

	// Create a matched request for this offer
	matchedRequestSource, _ := model.NewCoordinate(31.22082087, 29.94795413)
	matchedRequestDestination, _ := model.NewCoordinate(31.208936, 29.933419)
	matchedRequestPickup, _ := model.NewCoordinate(31.22082087, 29.94795413)
	matchedRequestDropoff, _ := model.NewCoordinate(31.208936, 29.933419)
	matchedRequestEarliestDepartureTime := offerDepartureTime.Add(-10 * time.Minute)
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
	request1Source, _ := model.NewCoordinate(31.208789314386852, 29.932685590743088)
	request1Destination, _ := model.NewCoordinate(31.2077329110055, 29.9268726301741)
	request1MaxWalkingDuration := time.Duration(0) * time.Minute
	request1NumberOfRiders := 1
	request1SameGender := true
	request1Gender := enums.Male
	pickup1, _, dropoff1, _ := getRequestPointsAndDurations(engine, offer, request1Source, request1MaxWalkingDuration, request1Destination)
	// Create request 2
	request2Source, _ := model.NewCoordinate(31.2398552645433, 29.9646946591899)
	request2Destination, _ := model.NewCoordinate(31.2251143325977, 29.9474103945877)
	request2MaxWalkingDuration := time.Duration(0) * time.Minute
	request2NumberOfRiders := 3
	request2SameGender := false
	request2Gender := enums.Male
	pickup2, _, dropoff2, _ := getRequestPointsAndDurations(engine, offer, request2Source, request2MaxWalkingDuration, request2Destination)

	cumulativeTimesWithoutRider := correcteness_test.GetCumulativeTimes([]model.Coordinate{*offerSource, *offerDestination}, offerDepartureTime, engine)
	cumulativeTimesWithRider := correcteness_test.GetCumulativeTimes([]model.Coordinate{*offerSource, *pickup2, *dropoff2, *matchedRequestSource, *matchedRequestDestination, *pickup1, *dropoff1, *offerDestination}, offerDepartureTime, engine)

	// overwrite offer detour, maxEstimated arrival time && matchedRequestLatestArrivalTime
	fmt.Println(cumulativeTimesWithoutRider)
	fmt.Println(cumulativeTimesWithRider)
	offerDetourDuration = cumulativeTimesWithRider[7] - cumulativeTimesWithoutRider[1] + 10*time.Second // adding 1 minutes to ensure the detour is valid
	offer.SetDetour(offerDetourDuration)
	offer.SetMaxEstimatedArrivalTime(getMaxEstimatedArrivalTime(*offerSource, *offerDestination, offerDepartureTime, offerDetourDuration, engine))
	matchedRequest.SetLatestArrivalTime(offer.MaxEstimatedArrivalTime().Add(10 * time.Minute))
	request1EarliestDepartureTime := offerDepartureTime.Add(-request1MaxWalkingDuration).Add(-1 * time.Minute)
	request1LatestArrivalTime := offerMaxEstimatedArrivalTime.Add(request1MaxWalkingDuration).Add(1 * time.Minute)

	request2EarliestDepartureTime := offerDepartureTime.Add(-request2MaxWalkingDuration).Add(-1 * time.Minute)
	request2LatestArrivalTime := offerMaxEstimatedArrivalTime.Add(request2MaxWalkingDuration).Add(1 * time.Minute)

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
	pickupOrder1, dropoffOrder1 := 5, 6

	pickupPoint2, dropoffPoint2 := computeRequestPickupDropoffPoints(engine, offer, request2Source, request2MaxWalkingDuration, request2Destination, request2EarliestDepartureTime, request2, request2LatestArrivalTime)
	pickupOrder2, dropoffOrder2 := 1, 2
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
