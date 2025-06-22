package tests

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"matching-engine/cmd/correcteness_test"
	"matching-engine/internal/adapter/routing"
	"matching-engine/internal/enums"
	"matching-engine/internal/model"
	"time"
)

func getTest1fivData(engine routing.Engine) ([]*model.Offer, []*model.Request, map[string]*model.MatchingResult) {
	offers := make([]*model.Offer, 0)
	requests := make([]*model.Request, 0)

	// Create an offer with the specified attributes
	offerSource, _ := model.NewCoordinate(31.2460735985739, 29.9744554984058)
	offerDestination, _ := model.NewCoordinate(31.2068412085851, 29.9246876930902)
	offerDepartureTime := correcteness_test.ParseTime("10:30")
	offerDetourDuration := time.Duration(8) * time.Minute // will be overwritten this is a placeholder value
	offerCapacity := 1
	offerCurrentNumberOfRequests := 1
	offerSameGender := false
	offerGender := enums.Male
	offerMaxEstimatedArrivalTime := offerDepartureTime.Add(15 * time.Minute) // will be overwritten this is a placeholder value

	// Create a matched request for this offer
	matchedRequestSource, _ := model.NewCoordinate(31.208936, 29.933419)
	matchedRequestDestination, _ := model.NewCoordinate(31.2077329110055, 29.9268726301741)
	matchedRequestPickup, _ := model.NewCoordinate(31.208936, 29.933419)
	matchedRequestDropoff, _ := model.NewCoordinate(31.2077329110055, 29.9268726301741)
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

	// Create another request
	requestSource, _ := model.NewCoordinate(31.232139, 29.951709)
	requestDestination, _ := model.NewCoordinate(31.2208208709376, 29.9479541306202)
	requestMaxWalkingDuration := time.Duration(0) * time.Minute

	pickup, _, dropoff, _ := getRequestPointsAndDurations(engine, offer, requestSource, requestMaxWalkingDuration, requestDestination)
	cumulativeTimesWithoutRider := correcteness_test.GetCumulativeTimes([]model.Coordinate{*offerSource, *offerDestination}, offerDepartureTime, engine)
	cumulativeTimesWithMatchedRequest := correcteness_test.GetCumulativeTimes([]model.Coordinate{*offerSource, *matchedRequestSource, *matchedRequestDestination, *offerDestination}, offerDepartureTime, engine)
	cumulativeTimesWithRider := correcteness_test.GetCumulativeTimes([]model.Coordinate{*offerSource, *pickup, *dropoff, *matchedRequestSource, *matchedRequestDestination, *offerDestination}, offerDepartureTime, engine)

	fmt.Println(cumulativeTimesWithRider)
	fmt.Println(cumulativeTimesWithMatchedRequest)
	fmt.Println(cumulativeTimesWithoutRider)

	// overwrite offer detour, maxEstimated arrival time && matchedRequestLatestArrivalTime
	offerDetourDuration = cumulativeTimesWithRider[5] - cumulativeTimesWithoutRider[1] + 5*time.Minute // adding 5 minutes to ensure the detour is valid
	offer.SetDetour(offerDetourDuration)
	offer.SetMaxEstimatedArrivalTime(offerDepartureTime.Add(cumulativeTimesWithoutRider[1]).Add(offerDetourDuration))
	matchedRequest.SetLatestArrivalTime(offerDepartureTime.Add(cumulativeTimesWithMatchedRequest[2]).Add(1 * time.Minute)) // setting matchedRequestLatestArrivalTime to the arrival time before adding a new request to ensure it will not arrive before it's latest arrival time
	fmt.Println(matchedRequest.LatestArrivalTime())

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
