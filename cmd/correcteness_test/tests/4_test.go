package tests

import (
	"matching-engine/cmd/correcteness_test"
	"matching-engine/internal/adapter/routing"
	"matching-engine/internal/enums"
	"matching-engine/internal/model"
	"time"
)

func getTest4(engine routing.Engine) ([]*model.Offer, []*model.Request, map[string]*model.MatchingResult, map[string]*model.MatchingResult) {
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

	// Create request 2
	request2Source, _ := model.NewCoordinate(31.23735, 29.97093)
	request2Destination, _ := model.NewCoordinate(31.22384, 29.96977)
	request2MaxWalkingDuration := time.Duration(0) * time.Minute
	request2NumberOfRiders := 1
	request2SameGender := false
	request2Gender := enums.Male
	request2EarliestDepartureTime := correcteness_test.ParseTime("10:00") // set equal to earliest departure time of offers since there is no walking in this test case

	// Request 3
	request3Source := *must(model.NewCoordinate(31.22913, 29.95111))
	request3Destination := *must(model.NewCoordinate(31.24035, 29.99252))
	request3MaxWalkingDuration := 0 * time.Minute
	request3NumberOfRiders := 1
	request3SameGender := true
	request3Gender := enums.Male
	request3EarliestDepartureTime := correcteness_test.ParseTime("10:00")

	// Request 4
	request4Source := *must(model.NewCoordinate(31.23265, 29.94748))
	request4Destination := *must(model.NewCoordinate(31.22505, 30.00766))
	request4MaxWalkingDuration := 0 * time.Minute
	request4NumberOfRiders := 1
	request4SameGender := false
	request4Gender := enums.Male
	request4EarliestDepartureTime := correcteness_test.ParseTime("10:00")
	request1LatestArrivalTime := offer1.MaxEstimatedArrivalTime()
	request2LatestArrivalTime := offer1.MaxEstimatedArrivalTime()
	request3LatestArrivalTime := offer1.MaxEstimatedArrivalTime()
	request4LatestArrivalTime := offer1.MaxEstimatedArrivalTime()

	request1 := CreateRequest("4", "1", *request1Source, *request1Destination,
		request1EarliestDepartureTime, request1LatestArrivalTime,
		request1MaxWalkingDuration, request1NumberOfRiders,
		request1Gender, request1SameGender)

	requests = append(requests, request1)

	request2 := CreateRequest("5", "2", *request2Source, *request2Destination,
		request2EarliestDepartureTime, request2LatestArrivalTime,
		request2MaxWalkingDuration, request2NumberOfRiders,
		request2Gender, request2SameGender)

	// Add request 2 to the list of requests
	requests = append(requests, request2)

	request3 := CreateRequest("6", "3", request3Source, request3Destination,
		request3EarliestDepartureTime, request3LatestArrivalTime,
		request3MaxWalkingDuration, request3NumberOfRiders,
		request3Gender, request3SameGender)
	requests = append(requests, request3)

	request4 := CreateRequest("7", "4", request4Source, request4Destination,
		request4EarliestDepartureTime, request4LatestArrivalTime,
		request4MaxWalkingDuration, request4NumberOfRiders,
		request4Gender, request4SameGender)
	requests = append(requests, request4)

	expectedResults_a := make(map[string]*model.MatchingResult)
	expectedResults_b := make(map[string]*model.MatchingResult)

	pickupPoint1_a, dropoffPoint1_a := ComputeRequestPickupDropoffPoints(engine, offer1, request1Source, request1MaxWalkingDuration, request1Destination, request1EarliestDepartureTime, request1, request1LatestArrivalTime)
	pickupPoint1_b, dropoffPoint1_b := ComputeRequestPickupDropoffPoints(engine, offer1, request1Source, request1MaxWalkingDuration, request1Destination, request1EarliestDepartureTime, request1, request1LatestArrivalTime)
	pickupOrder1, dropoffOrder1 := 1, 2

	pickupPoint2_a, dropoffPoint2_a := ComputeRequestPickupDropoffPoints(engine, offer1, request2Source, request2MaxWalkingDuration, request2Destination, request2EarliestDepartureTime, request2, request2LatestArrivalTime)
	pickupPoint2_b, dropoffPoint2_b := ComputeRequestPickupDropoffPoints(engine, offer1, request2Source, request2MaxWalkingDuration, request2Destination, request2EarliestDepartureTime, request2, request2LatestArrivalTime)
	pickupOrder2, dropoffOrder2 := 1, 2

	pickupPoint3_a, dropoffPoint3_a := ComputeRequestPickupDropoffPoints(engine, offer1, &request3Source, request3MaxWalkingDuration, &request3Destination, request3EarliestDepartureTime, request3, request3LatestArrivalTime)
	pickupPoint3_b, dropoffPoint3_b := ComputeRequestPickupDropoffPoints(engine, offer2, &request3Source, request3MaxWalkingDuration, &request3Destination, request3EarliestDepartureTime, request3, request3LatestArrivalTime)
	pickupOrder3, dropoffOrder3 := 3, 4

	cumulativeTimesForDriver1With3 := correcteness_test.GetCumulativeTimes([]model.Coordinate{*offer1.Source(), *request1Source, *request1Destination, request3Source, request3Destination, *offer1.Destination()}, offer1.DepartureTime(), engine)
	cumulativeTimesForDriver1Direct := correcteness_test.GetCumulativeTimes([]model.Coordinate{*offer1.Source(), *offer1.Destination()}, offer1.DepartureTime(), engine)

	cumulativeTimesForDriver2With3 := correcteness_test.GetCumulativeTimes([]model.Coordinate{*offer2.Source(), *pickupPoint2_b.Coordinate(), *dropoffPoint2_b.Coordinate(), *pickupPoint3_b.Coordinate(), *dropoffPoint3_b.Coordinate(), *offer2.Destination()}, offer2.DepartureTime(), engine)
	cumulativeTimesForDriver2Direct := correcteness_test.GetCumulativeTimes([]model.Coordinate{*offer2.Source(), *offer2.Destination()}, offer2.DepartureTime(), engine)

	cumulativeTimesForDriver3Direct := correcteness_test.GetCumulativeTimes([]model.Coordinate{*offer3.Source(), *offer3.Destination()}, offer3.DepartureTime(), engine)
	cumulativeTimesForDriver2With4Only := correcteness_test.GetCumulativeTimes([]model.Coordinate{*offer2.Source(), request4Source, request4Destination, *offer2.Destination()}, offer2.DepartureTime(), engine)

	offerDetourDuration1 := cumulativeTimesForDriver1With3[5] - cumulativeTimesForDriver1Direct[1] + 1*time.Second // adding 1 seconds to ensure the detour is valid
	offerDetourDuration2 := cumulativeTimesForDriver2With3[5] - cumulativeTimesForDriver2Direct[1] + 1*time.Second // adding 1 seconds to ensure the detour is valid
	offerDetourDuration3 := cumulativeTimesForDriver3Direct[1] + 1*time.Second                                     // adding 1 seconds to ensure the detour is valid

	offer1.SetDetour(offerDetourDuration1 + 10*time.Second) // adding 10 seconds to ensure the detour is valid
	offer1.SetMaxEstimatedArrivalTime(GetMaxEstimatedArrivalTime(*offer1.Source(), *offer1.Destination(), offer1.DepartureTime(), offerDetourDuration1, engine))
	offer2.SetDetour(offerDetourDuration2 + 10*time.Second) // adding 10 seconds to ensure the detour is valid
	offer2.SetMaxEstimatedArrivalTime(GetMaxEstimatedArrivalTime(*offer2.Source(), *offer2.Destination(), offer2.DepartureTime(), offerDetourDuration2, engine))
	offer3.SetDetour(offerDetourDuration3)
	offer3.SetMaxEstimatedArrivalTime(GetMaxEstimatedArrivalTime(*offer3.Source(), *offer3.Destination(), offer3.DepartureTime(), offerDetourDuration3, engine))

	request1.SetLatestArrivalTime(offer1.DepartureTime().Add(request1MaxWalkingDuration).Add(cumulativeTimesForDriver1With3[2]).Add(10 * time.Second))
	request2.SetLatestArrivalTime(offer2.DepartureTime().Add(request2MaxWalkingDuration).Add(cumulativeTimesForDriver2With3[2]).Add(1 * time.Second)) // so that driver 1 is not able to pickup rider 2
	request3.SetLatestArrivalTime(offer1.MaxEstimatedArrivalTime().Add(request3MaxWalkingDuration).Add(1 * time.Hour))
	request4.SetLatestArrivalTime(offer2.DepartureTime().Add(request4MaxWalkingDuration).Add(cumulativeTimesForDriver2With4Only[2]).Add(-1 * time.Minute))

	// Create two expected results

	points1_a := []*model.PathPoint{
		pickupPoint1_a,
		dropoffPoint1_a,
		pickupPoint3_a,
		dropoffPoint3_a,
	}
	pointsOrder1_a := []int{
		pickupOrder1,
		dropoffOrder1,
		pickupOrder3,
		dropoffOrder3,
	}

	points2_a := []*model.PathPoint{
		pickupPoint2_a,
		dropoffPoint2_a,
	}
	pointsOrder2_a := []int{
		pickupOrder2,
		dropoffOrder2,
	}

	points1_b := []*model.PathPoint{
		pickupPoint1_b,
		dropoffPoint1_b,
	}
	pointsOrder1_b := []int{
		pickupOrder1,
		dropoffOrder1,
	}

	points2_b := []*model.PathPoint{
		pickupPoint2_b,
		dropoffPoint2_b,
		pickupPoint3_b,
		dropoffPoint3_b,
	}
	pointsOrder2_b := []int{
		pickupOrder2,
		dropoffOrder2,
		pickupOrder3,
		dropoffOrder3,
	}

	offerPath1_a := AddPointsToPath(engine, offer1, pointsOrder1_a, points1_a)
	offerPath2_a := AddPointsToPath(engine, offer2, pointsOrder2_a, points2_a)
	offerPath1_b := AddPointsToPath(engine, offer1, pointsOrder1_b, points1_b)
	offerPath2_b := AddPointsToPath(engine, offer2, pointsOrder2_b, points2_b)

	expectedResults_a[offer1.ID()] = model.NewMatchingResult(
		offer1.ID(),
		offer1.UserID(),
		[]*model.Request{request1, request3},
		offerPath1_a,
		2,
	)
	expectedResults_a[offer2.ID()] = model.NewMatchingResult(
		offer2.ID(),
		offer2.UserID(),
		[]*model.Request{request2},
		offerPath2_a,
		1,
	)

	expectedResults_b[offer1.ID()] = model.NewMatchingResult(
		offer1.ID(),
		offer1.UserID(),
		[]*model.Request{request1},
		offerPath1_b,
		1,
	)
	expectedResults_b[offer2.ID()] = model.NewMatchingResult(
		offer2.ID(),
		offer2.UserID(),
		[]*model.Request{request2, request3},
		offerPath2_b,
		2,
	)

	return offers, requests, expectedResults_a, expectedResults_b
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
	maxArrivalTime := GetMaxEstimatedArrivalTime(source, destination, departureTime, detourDuration, engine)

	offer := CreateOffer(id, driverID, source, destination, departureTime,
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

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
