package tests

import (
	"github.com/rs/zerolog/log"
	"matching-engine/cmd/correcteness_test"
	"matching-engine/internal/adapter/routing"
	"matching-engine/internal/enums"
	"matching-engine/internal/model"
	"time"
)

func getTest2aData(engine routing.Engine) ([]*model.Offer, []*model.Request, map[string]*model.MatchingResult) {
	offers := make([]*model.Offer, 0)
	requests := make([]*model.Request, 0)

	// Create an offer with the specified attributes
	offerSource, _ := model.NewCoordinate(31.2460735985739, 29.9744554984058)
	offerDestination, _ := model.NewCoordinate(31.2068412085851, 29.9246876930902)
	offerDepartureTime := correcteness_test.ParseTime("10:30")
	offerDetourDuration := time.Duration(8) * time.Minute
	offerCapacity := 1
	offerCurrentNumberOfRequests := 1
	offerSameGender := false
	offerGender := enums.Male
	offerMaxEstimatedArrivalTime := GetMaxEstimatedArrivalTime(*offerSource, *offerDestination, offerDepartureTime, offerDetourDuration, engine)

	// Create a matched request for this offer
	matchedRequestSource, _ := model.NewCoordinate(31.22082087, 29.94795413)
	matchedRequestDestination, _ := model.NewCoordinate(31.208936, 29.933419)
	matchedRequestEarliestDepartureTime := correcteness_test.ParseTime("10:20:00")
	matchedRequestLatestArrivalTime := correcteness_test.ParseTime("11:20")
	matchedRequestMaxWalkingDuration := time.Duration(0) * time.Minute
	matchedRequestNumberOfRiders := 1
	matchedRequestSameGender := true
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
	requestSource, _ := model.NewCoordinate(31.208936, 29.933419)
	requestDestination, _ := model.NewCoordinate(31.20773291, 29.92687263)
	requestMaxWalkingDuration := time.Duration(0) * time.Minute

	pickup, _, dropoff, _ := GetRequestPointsAndDurations(engine, offer, requestSource, requestMaxWalkingDuration, requestDestination)
	cumulativeTimesWithoutRider := correcteness_test.GetCumulativeTimes([]model.Coordinate{*offerSource, *offerDestination}, offerDepartureTime, engine)
	cumulativeTimesWithRider := correcteness_test.GetCumulativeTimes([]model.Coordinate{*offerSource, *matchedRequestSource, *matchedRequestDestination, *pickup, *dropoff, *offerDestination}, offerDepartureTime, engine)

	// overwrite offer detour, maxEstimated arrival time && matchedRequestLatestArrivalTime
	offerDetourDuration = cumulativeTimesWithRider[5] - cumulativeTimesWithoutRider[1] + 1*time.Second // adding 5 minutes to ensure the detour is valid
	offer.SetDetour(offerDetourDuration)
	offer.SetMaxEstimatedArrivalTime(offerDepartureTime.Add(cumulativeTimesWithoutRider[1]).Add(offerDetourDuration))
	matchedRequest.SetLatestArrivalTime(offer.MaxEstimatedArrivalTime().Add(10 * time.Minute))
	log.Debug().
		Int("offerDetourDurationMinutes", int(offerDetourDuration.Minutes())).
		Str("offerMaxEstimatedArrivalTime", offer.MaxEstimatedArrivalTime().Format(time.RFC3339)).
		Str("matchedRequestLatestArrivalTime", matchedRequest.LatestArrivalTime().Format(time.RFC3339)).
		Msg("Offer and matched request details after detour adjustment")
	requestEarliestDepartureTime := offerDepartureTime.Add(-10 * time.Minute)
	requestLatestArrivalTime := offerMaxEstimatedArrivalTime.Add(10 * time.Minute)
	requestNumberOfRiders := 1
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
	pickupPoint, dropoffPoint := ComputeRequestPickupDropoffPoints(engine, offer, requestSource, requestMaxWalkingDuration, requestDestination, requestEarliestDepartureTime, request, requestLatestArrivalTime)
	pickupOrder, dropoffOrder := 3, 4
	// Compute the offer path with the new pickup and dropoff points
	points := []*model.PathPoint{pickupPoint, dropoffPoint}
	pointsOrder := []int{pickupOrder, dropoffOrder}
	offerPath := AddPointsToPath(engine, offer, pointsOrder, points)

	expectedResults[offer.ID()] = model.NewMatchingResult(
		offer.ID(),
		offer.UserID(),
		[]*model.Request{request},
		offerPath,
		offer.CurrentNumberOfRequests()+1,
	)
	return offers, requests, expectedResults
}
