package tests

import (
	"github.com/rs/zerolog/log"
	"matching-engine/cmd/correcteness_test"
	"matching-engine/internal/adapter/routing"
	"matching-engine/internal/enums"
	"matching-engine/internal/model"
	"time"
)

func getTest1fiiData(engine routing.Engine) ([]*model.Offer, []*model.Request, map[string]*model.MatchingResult) {
	offers := make([]*model.Offer, 0)
	requests := make([]*model.Request, 0)

	// Create an offer with the specified attributes
	offerSource, _ := model.NewCoordinate(31.2460735985739, 29.9744554984058)
	offerDestination, _ := model.NewCoordinate(31.2068412085851, 29.9246876930902)
	offerDepartureTime := correcteness_test.ParseTime("10:30")
	offerDetourDuration := time.Duration(4) * time.Minute // will be overwritten this is a placeholder value
	offerCapacity := 1
	offerCurrentNumberOfRequests := 1
	offerSameGender := false
	offerGender := enums.Male
	offerMaxEstimatedArrivalTime := offerDepartureTime.Add(15 * time.Minute) // will be overwritten this is a placeholder value

	// Create a matched request for this offer
	matchedRequestSource, _ := model.NewCoordinate(31.232139, 29.951709)
	matchedRequestDestination, _ := model.NewCoordinate(31.2208208709376, 29.9479541306202)
	matchedRequestEarliestDepartureTime := offerDepartureTime.Add(-10 * time.Minute)
	matchedRequestLatestArrivalTime := offerMaxEstimatedArrivalTime.Add(10 * time.Minute)
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
	requestDestination, _ := model.NewCoordinate(31.2077329110055, 29.9268726301741)
	requestMaxWalkingDuration := time.Duration(0) * time.Minute

	pickup, _, dropoff, _ := GetRequestPointsAndDurations(engine, offer, requestSource, requestMaxWalkingDuration, requestDestination)
	cumulativeTimesWithoutRider := correcteness_test.GetCumulativeTimes([]model.Coordinate{*offerSource, *offerDestination}, offerDepartureTime, engine)
	cumulativeTimesWithRider := correcteness_test.GetCumulativeTimes([]model.Coordinate{*offerSource, *matchedRequestSource, *matchedRequestDestination, *pickup, *dropoff, *offerDestination}, offerDepartureTime, engine)

	// overwrite offer detour, maxEstimated arrival time && matchedRequestLatestArrivalTime
	// offerDetourDuration is done as this so that it passes the early check of the detour but don't pass the detour validation in the feasiblity check
	offerDetourDuration = cumulativeTimesWithRider[5] - cumulativeTimesWithoutRider[1] - 1*time.Minute
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
	return offers, requests, expectedResults
}
