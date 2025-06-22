package tests

import (
	"fmt"
	"matching-engine/cmd/correcteness_test"
	"matching-engine/internal/adapter/routing"
	"matching-engine/internal/enums"
	"matching-engine/internal/model"
	"time"
)

func getTest4(engine routing.Engine) ([]*model.Offer, []*model.Request, map[string]*model.MatchingResult) {
	offers := make([]*model.Offer, 0)
	requests := make([]*model.Request, 0)

	offer1 := createAndAppendOffer(
		"1", "1",
		*must(model.NewCoordinate(31.26201, 29.98355)),
		*must(model.NewCoordinate(31.26201, 29.98355)),
		correcteness_test.ParseTime("10:00"),
		8*time.Minute, 4, 0,
		enums.Male, false, engine, &offers,
	)

	offer2 := createAndAppendOffer(
		"2", "2",
		*must(model.NewCoordinate(31.22973, 29.9796)),
		*must(model.NewCoordinate(31.21459, 29.94633)),
		correcteness_test.ParseTime("10:00"),
		8*time.Minute, 4, 0,
		enums.Male, false, engine, &offers,
	)

	offer3 := createAndAppendOffer(
		"3", "3",
		*must(model.NewCoordinate(31.20961, 29.90817)),
		*must(model.NewCoordinate(31.20837, 29.88277)),
		correcteness_test.ParseTime("10:00"),
		8*time.Minute, 4, 0,
		enums.Male, false, engine, &offers,
	)

	// Create request 1
	request1Source, _ := model.NewCoordinate(31.26519, 29.99721)
	request1Destination, _ := model.NewCoordinate(31.23661, 29.95622)
	request1MaxWalkingDuration := time.Duration(0) * time.Minute
	request1NumberOfRiders := 1
	request1SameGender := true
	request1Gender := enums.Male
	request1EarliestDepartureTime := correcteness_test.ParseTime("10:00")
	pickup1, _, dropoff1, _ := getRequestPointsAndDurations(engine, offer1, request1Source, request1MaxWalkingDuration, request1Destination)

	// Create request 2
	request2Source, _ := model.NewCoordinate(31.23735, 29.97093)
	request2Destination, _ := model.NewCoordinate(31.22384, 29.96977)
	request2MaxWalkingDuration := time.Duration(0) * time.Minute
	request2NumberOfRiders := 3
	request2SameGender := false
	request2Gender := enums.Male
	request2EarliestDepartureTime := correcteness_test.ParseTime("10:00") // set equal to earliest departure time of offers since there is no walking in this test case
	pickup2, _, dropoff2, _ := getRequestPointsAndDurations(engine, offer1, request2Source, request2MaxWalkingDuration, request2Destination)

	// Request 3
	request3Source := *must(model.NewCoordinate(31.22913, 29.95111))
	request3Destination := *must(model.NewCoordinate(31.24035, 29.99252))
	request3MaxWalkingDuration := 0 * time.Minute
	request3NumberOfRiders := 2
	request3SameGender := true
	request3Gender := enums.Female
	request3EarliestDepartureTime := correcteness_test.ParseTime("10:00")
	pickup3, _, dropoff3, _ := getRequestPointsAndDurations(engine, offer1, request3Source, request3MaxWalkingDuration, request3Destination)

	// Request 4
	request4Source := *must(model.NewCoordinate(31.23265, 29.94748))
	request4Destination := *must(model.NewCoordinate(31.22505, 30.00766))
	request4MaxWalkingDuration := 0 * time.Minute
	request4NumberOfRiders := 1
	request4SameGender := false
	request4Gender := enums.Male
	request4EarliestDepartureTime := correcteness_test.ParseTime("10:00")
	pickup4, _, dropoff4, _ := getRequestPointsAndDurations(engine, offer1, request4Source, request4MaxWalkingDuration, request4Destination)

	// From heereeeeeeeeeeeeeeeeeeeeeeee
	//cumulativeTimesWithoutRider := correcteness_test.GetCumulativeTimes([]model.Coordinate{*offer1.Source(), *offer1.Destination()}, offerDepartureTime1, engine)
	//cumulativeTimesWithRider := correcteness_test.GetCumulativeTimes([]model.Coordinate{*offerSource1, *pickup2, *dropoff2, *matchedRequestSource, *matchedRequestDestination, *pickup1, *dropoff1, *offerDestination1}, offerDepartureTime1, engine)

	//// overwrite offer1 detour, maxEstimated arrival time && matchedRequestLatestArrivalTime
	//fmt.Println(cumulativeTimesWithoutRider)
	//fmt.Println(cumulativeTimesWithRider)
	//offerDetourDuration1 := cumulativeTimesWithRider[7] - cumulativeTimesWithoutRider[1] + 10*time.Second // adding 1 minutes to ensure the detour is valid
	//offer1.SetDetour(offerDetourDuration1)
	//offer1.SetMaxEstimatedArrivalTime(getMaxEstimatedArrivalTime(*offerSource1, *offerDestination1, offerDepartureTime1, offerDetourDuration1, engine))
	//// Till hereeeeeeeeeeeeeee ti
	// TODO overwrite this with max estimated arrival time
	request1LatestArrivalTime := offer1.MaxEstimatedArrivalTime().Add(request1MaxWalkingDuration).Add(1 * time.Minute)
	request2LatestArrivalTime := offer1.MaxEstimatedArrivalTime().Add(request2MaxWalkingDuration).Add(1 * time.Minute)
	request3LatestArrivalTime := offer1.MaxEstimatedArrivalTime().Add(request3MaxWalkingDuration).Add(1 * time.Minute)
	request4LatestArrivalTime := offer1.MaxEstimatedArrivalTime().Add(request4MaxWalkingDuration).Add(1 * time.Minute)

	request1 := createRequest("4", "1", *request1Source, *request1Destination,
		request1EarliestDepartureTime, request1LatestArrivalTime,
		request1MaxWalkingDuration, request1NumberOfRiders,
		request1Gender, request1SameGender)

	requests = append(requests, request1)

	request2 := createRequest("5", "2", *request2Source, *request2Destination,
		request2EarliestDepartureTime, request2LatestArrivalTime,
		request2MaxWalkingDuration, request2NumberOfRiders,
		request2Gender, request2SameGender)

	// Add request 2 to the list of requests
	requests = append(requests, request2)

	request3 := createRequest("6", "3", request3Source, request3Destination,
		request3EarliestDepartureTime, offer1.MaxEstimatedArrivalTime().Add(request3MaxWalkingDuration).Add(1*time.Minute),
		request3MaxWalkingDuration, request3NumberOfRiders,
		request3Gender, request3SameGender)
	requests = append(requests, request3)

	request4 := createRequest("7", "4", request4Source, request4Destination,
		request4EarliestDepartureTime, offer1.MaxEstimatedArrivalTime().Add(request4MaxWalkingDuration).Add(1*time.Minute),
		request4MaxWalkingDuration, request4NumberOfRiders,
		request4Gender, request4SameGender)
	requests = append(requests, request4)

	// Create expected results
	expectedResults := make(map[string]*model.MatchingResult)
	pickupPoint1, dropoffPoint1 := computeRequestPickupDropoffPoints(engine, offer1, request1Source, request1MaxWalkingDuration, request1Destination, request1EarliestDepartureTime, request1, request1LatestArrivalTime)
	pickupOrder1, dropoffOrder1 := 5, 6

	pickupPoint2, dropoffPoint2 := computeRequestPickupDropoffPoints(engine, offer1, request2Source, request2MaxWalkingDuration, request2Destination, request2EarliestDepartureTime, request2, request2LatestArrivalTime)
	pickupOrder2, dropoffOrder2 := 1, 3
	// Add the requests to the offer1
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
	offerPath := addPointsToPath(engine, offer1, pointsOrder, points)

	expectedResults[offer1.ID()] = model.NewMatchingResult(
		offer1.ID(),
		offer1.UserID(),
		[]*model.Request{request1, request2},
		offerPath,
		offer1.CurrentNumberOfRequests()+2,
	)
	return offers, requests, expectedResults
}
func createAndAppendOffer(
	id string,
	driverID string,
	source, destination model.Coordinate,
	departureTime time.Time,
	detourDuration time.Duration,
	capacity int,
	currentRequests int,
	gender enums.Gender,
	sameGender bool,
	engine routing.Engine,
	offers *[]*model.Offer,
) *model.Offer {
	// Will be overwritten later in the test
	maxArrivalTime := getMaxEstimatedArrivalTime(source, destination, departureTime, detourDuration, engine)

	offer := createOffer(id, driverID, source, destination, departureTime,
		detourDuration, capacity, currentRequests, gender,
		sameGender, maxArrivalTime, []*model.Request{})

	path := []model.PathPoint{
		*model.NewPathPoint(source, enums.Source, departureTime, offer, 0),
		*model.NewPathPoint(destination, enums.Destination, maxArrivalTime, offer, 0),
	}

	offer.SetPath(path)
	*offers = append(*offers, offer)
	return offer
}
