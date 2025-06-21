package tests

import (
	"matching-engine/internal/adapter/routing"
	"matching-engine/internal/enums"
	"matching-engine/internal/model"
	"time"
)

func getTest1aiData(engine routing.Engine) ([]*model.Offer, []*model.Request, map[string]*model.MatchingResult) {
	offers := make([]*model.Offer, 0)
	requests := make([]*model.Request, 0)

	// Create an offer with the specified attributes
	offerSource, _ := model.NewCoordinate(31.2544088039743, 29.97376045816)
	offerDestination, _ := model.NewCoordinate(31.20611644667, 29.9248733439259)
	offerDepartureTime, _ := time.Parse("15:04", "10:30")
	offerDepartureTime = adjustToNextDay(offerDepartureTime)
	offerDetourDuration := time.Duration(30)
	offerCapacity := 3
	offerCurrentNumberOfRequests := 1
	offerSameGender := false
	offerGender := enums.Male
	offerMaxEstimatedArrivalTime := getMaxEstimatedArrivalTime(*offerSource, *offerDestination, offerDepartureTime, offerDetourDuration, engine)

	// Create a matched request for this offer
	matchedRequestSource, _ := model.NewCoordinate(31.2544088, 29.97376046)
	matchedRequestDestination, _ := model.NewCoordinate(31.20611645, 29.92487334)
	matchedRequestPickup, _ := model.NewCoordinate(31.2544088, 29.97376046)
	matchedRequestDropoff, _ := model.NewCoordinate(31.20611645, 29.92487334)
	matchedRequestEarliestDepartureTime, _ := time.Parse("15:04:05", "10:20:00")
	matchedRequestLatestArrivalTime, _ := time.Parse("15:04", "11:20")
	matchedRequestEarliestDepartureTime = adjustToNextDay(matchedRequestEarliestDepartureTime)
	matchedRequestLatestArrivalTime = adjustToNextDay(matchedRequestLatestArrivalTime)
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
	requestSource, _ := model.NewCoordinate(31.2544088, 29.97376046)
	requestDestination, _ := model.NewCoordinate(31.20611645, 29.92487334)
	requestEarliestDepartureTime, _ := time.Parse("15:04:05", "09:20:00")
	requestLatestArrivalTime, _ := time.Parse("15:04", "10:20")
	requestEarliestDepartureTime = adjustToNextDay(requestEarliestDepartureTime)
	requestLatestArrivalTime = adjustToNextDay(requestLatestArrivalTime)
	requestMaxWalkingDuration := time.Duration(0)
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

// adjustToNextDay adjusts the given time to the next day while preserving the time of day
func adjustToNextDay(t time.Time) time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day()+1, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
}
