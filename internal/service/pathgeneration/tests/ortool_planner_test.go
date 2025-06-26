package tests

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"matching-engine/internal/adapter/ortool"
	"matching-engine/internal/enums"
	"matching-engine/internal/model"
	planner2 "matching-engine/internal/service/pathgeneration/planner"
	"matching-engine/internal/service/pickupdropoffservice/pickupdropoffcache"
	"matching-engine/internal/service/timematrix/cache"
	"testing"
	"time"
)

// MockTimeMatrixSelector implements the timematrix.DefaultSelector interface for testing
type MockTimeMatrixSelector struct {
	mock.Mock
}

func (m *MockTimeMatrixSelector) GetTimeMatrix(offer *model.OfferNode) (*cache.PathPointMappedTimeMatrix, error) {
	args := m.Called(offer)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cache.PathPointMappedTimeMatrix), args.Error(1)
}

func createSourcePathPoint(offer *model.Offer) *model.PathPoint {
	return model.NewPathPoint(
		*createDefaultCoordinate(),
		enums.Source, time.Now().Add(30*time.Minute), offer, 5*time.Minute,
	)
}

func createDestinationPathPoint(offer *model.Offer) *model.PathPoint {
	return model.NewPathPoint(
		*createDefaultCoordinate(),
		enums.Destination, time.Now().Add(30*time.Minute), offer, 5*time.Minute,
	)
}

func TestORToolPlanner_validPath(t *testing.T) {
	resetPathPointIDCounter()
	timeNow := time.Now()

	mockPickupDropoffSelector := new(MockPickupDropoffSelector)
	mockTimeMatrixSelector := new(MockTimeMatrixSelector)

	// Reset the PathPointID counter before each test
	offer := createDefaultOfferWithTime(timeNow)
	offerNode := model.NewOfferNode(offer)
	offerSource := createSourcePathPoint(offer)
	offerDestination := createDestinationPathPoint(offer)

	request := createDefaultRequestWithEarliestDepartureTime(timeNow)
	request.SetLatestArrivalTime(timeNow.Add(31 * time.Minute))
	requestNode := model.NewRequestNode(request)

	pickup, dropoff := createDefaultPickupDropoff(request)
	pickupDropoff := pickupdropoffcache.NewValue(pickup, dropoff)

	matchedRequest := createDefaultRequestWithEarliestDepartureTime(timeNow.Add(8 * time.Minute))
	matchedRequest.SetLatestArrivalTime(timeNow.Add(21 * time.Minute))

	matchedRequestPickup, matchedRequestDropoff := createDefaultPickupDropoff(matchedRequest)

	offerSource.SetExpectedArrivalTime(timeNow)
	pickup.SetExpectedArrivalTime(timeNow.Add(5 * time.Minute))
	matchedRequestPickup.SetExpectedArrivalTime(timeNow.Add(10 * time.Minute))
	matchedRequestDropoff.SetExpectedArrivalTime(timeNow.Add(20 * time.Minute))
	dropoff.SetExpectedArrivalTime(timeNow.Add(30 * time.Minute))
	offerDestination.SetExpectedArrivalTime(timeNow.Add(40 * time.Minute))

	oldPath := []model.PathPoint{*offerSource, *matchedRequestPickup, *matchedRequestDropoff, *offerDestination}
	validPath := []model.PathPoint{*offerSource, *pickup, *matchedRequestPickup, *matchedRequestDropoff, *dropoff, *offerDestination}

	offer.SetMatchedRequests([]*model.Request{matchedRequest})
	offer.SetPath(oldPath)

	mockPickupDropoffSelector.On("GetPickupDropoffPointsAndDurations", request, offer).Return(pickupDropoff, nil)

	// Create a simple time matrix
	timeMatrix := [][]time.Duration{
		{0 * time.Minute, 5 * time.Minute, 10 * time.Minute, 20 * time.Minute, 30 * time.Minute, 40 * time.Minute},
		{5 * time.Minute, 0 * time.Minute, 5 * time.Minute, 15 * time.Minute, 25 * time.Minute, 35 * time.Minute},
		{10 * time.Minute, 5 * time.Minute, 0 * time.Minute, 10 * time.Minute, 20 * time.Minute, 30 * time.Minute},
		{20 * time.Minute, 15 * time.Minute, 10 * time.Minute, 0 * time.Minute, 10 * time.Minute, 20 * time.Minute},
		{30 * time.Minute, 25 * time.Minute, 20 * time.Minute, 10 * time.Minute, 0 * time.Minute, 10 * time.Minute},
		{40 * time.Minute, 35 * time.Minute, 30 * time.Minute, 20 * time.Minute, 10 * time.Minute, 0 * time.Minute},
	}

	pointIdToIndex := map[model.PathPointID]int{
		1: 0,
		2: 5,
		3: 1,
		4: 4,
		5: 2,
		6: 3,
	}

	mappedMatrix := cache.NewPathPointMappedTimeMatrix(timeMatrix, pointIdToIndex)
	mockTimeMatrixSelector.On("GetTimeMatrix", offerNode).Return(mappedMatrix, nil)

	client, _ := ortool.NewORToolClient()
	planner := planner2.NewORToolPlanner(mockPickupDropoffSelector, mockTimeMatrixSelector, client)
	resultPath, found, err := planner.FindFirstFeasiblePath(offerNode, requestNode)

	fmt.Println("timeNow:", timeNow)
	for i, point := range resultPath {
		fmt.Printf("Point %d: %+v\n", i, point)
		fmt.Printf("Point %d ID: %v\n", i, point.ID())
		fmt.Printf("Estimated arrival time for point %d: %s\n", i, point.ExpectedArrivalTime().String())
	}

	for i, point := range validPath {
		fmt.Printf("Valid Point %d: %+v\n", i, point)
		fmt.Printf("Valid Point %d ID: %v\n", i, point.ID())
		fmt.Printf("Valid Estimated arrival time for point %d: %s\n", i, point.ExpectedArrivalTime().String())
	}

	assert.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, validPath, resultPath)

	mockPickupDropoffSelector.AssertExpectations(t)
	mockTimeMatrixSelector.AssertExpectations(t)
}

func TestORToolPlanner_validPath_checkingMatrixOperation(t *testing.T) {
	resetPathPointIDCounter()
	timeNow := time.Now()

	mockPickupDropoffSelector := new(MockPickupDropoffSelector)
	mockTimeMatrixSelector := new(MockTimeMatrixSelector)

	// Reset the PathPointID counter before each test
	offer := createDefaultOfferWithTime(timeNow)
	offerNode := model.NewOfferNode(offer)
	offerSource := createSourcePathPoint(offer)
	offerDestination := createDestinationPathPoint(offer)

	request := createDefaultRequestWithEarliestDepartureTime(timeNow)
	request.SetLatestArrivalTime(timeNow.Add(31 * time.Minute))
	requestNode := model.NewRequestNode(request)

	pickup, dropoff := createDefaultPickupDropoff(request)
	pickupDropoff := pickupdropoffcache.NewValue(pickup, dropoff)

	matchedRequest := createDefaultRequestWithEarliestDepartureTime(timeNow.Add(8 * time.Minute))
	matchedRequest.SetLatestArrivalTime(timeNow.Add(21 * time.Minute))

	matchedRequestPickup, matchedRequestDropoff := createDefaultPickupDropoff(matchedRequest)

	offerSource.SetExpectedArrivalTime(timeNow)
	pickup.SetExpectedArrivalTime(timeNow.Add(5 * time.Minute))
	matchedRequestPickup.SetExpectedArrivalTime(timeNow.Add(10 * time.Minute))
	matchedRequestDropoff.SetExpectedArrivalTime(timeNow.Add(20 * time.Minute))
	dropoff.SetExpectedArrivalTime(timeNow.Add(30 * time.Minute))
	offerDestination.SetExpectedArrivalTime(timeNow.Add(40 * time.Minute))

	oldPath := []model.PathPoint{*offerSource, *matchedRequestPickup, *matchedRequestDropoff, *offerDestination}
	validPath := []model.PathPoint{*offerSource, *pickup, *matchedRequestPickup, *matchedRequestDropoff, *dropoff, *offerDestination}

	offer.SetMatchedRequests([]*model.Request{matchedRequest})
	offer.SetPath(oldPath)

	mockPickupDropoffSelector.On("GetPickupDropoffPointsAndDurations", request, offer).Return(pickupDropoff, nil)

	// Create a simple time matrix
	timeMatrix := [][]time.Duration{
		{0 * time.Minute, 55 * time.Minute, 5 * time.Minute, 10 * time.Minute, 20 * time.Minute, 30 * time.Minute, 40 * time.Minute},
		{55 * time.Minute, 0 * time.Minute, 55 * time.Minute, 55 * time.Minute, 15 * time.Minute, 25 * time.Minute, 35 * time.Minute},
		{5 * time.Minute, 55 * time.Minute, 0 * time.Minute, 5 * time.Minute, 15 * time.Minute, 25 * time.Minute, 35 * time.Minute},
		{10 * time.Minute, 55 * time.Minute, 5 * time.Minute, 0 * time.Minute, 10 * time.Minute, 20 * time.Minute, 30 * time.Minute},
		{20 * time.Minute, 55 * time.Minute, 15 * time.Minute, 10 * time.Minute, 0 * time.Minute, 10 * time.Minute, 20 * time.Minute},
		{30 * time.Minute, 55 * time.Minute, 25 * time.Minute, 20 * time.Minute, 10 * time.Minute, 0 * time.Minute, 10 * time.Minute},
		{40 * time.Minute, 55 * time.Minute, 35 * time.Minute, 30 * time.Minute, 20 * time.Minute, 10 * time.Minute, 0 * time.Minute},
	}

	pointIdToIndex := map[model.PathPointID]int{
		1: 0,
		2: 6,
		3: 2,
		4: 5,
		5: 3,
		6: 4,
		7: 1,
	}

	mappedMatrix := cache.NewPathPointMappedTimeMatrix(timeMatrix, pointIdToIndex)
	mockTimeMatrixSelector.On("GetTimeMatrix", offerNode).Return(mappedMatrix, nil)

	client, _ := ortool.NewORToolClient()
	planner := planner2.NewORToolPlanner(mockPickupDropoffSelector, mockTimeMatrixSelector, client)
	resultPath, found, err := planner.FindFirstFeasiblePath(offerNode, requestNode)

	fmt.Println("timeNow:", timeNow)
	for i, point := range resultPath {
		fmt.Printf("Point %d: %+v\n", i, point)
		fmt.Printf("Point %d ID: %v\n", i, point.ID())
		fmt.Printf("Estimated arrival time for point %d: %s\n", i, point.ExpectedArrivalTime().String())
	}

	for i, point := range validPath {
		fmt.Printf("Valid Point %d: %+v\n", i, point)
		fmt.Printf("Valid Point %d ID: %v\n", i, point.ID())
		fmt.Printf("Valid Estimated arrival time for point %d: %s\n", i, point.ExpectedArrivalTime().String())
	}

	assert.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, validPath, resultPath)

	mockPickupDropoffSelector.AssertExpectations(t)
	mockTimeMatrixSelector.AssertExpectations(t)
}

func TestORToolPlanner_inValidPathDueToInvalidCapacity(t *testing.T) {
	resetPathPointIDCounter()
	timeNow := time.Now()

	mockPickupDropoffSelector := new(MockPickupDropoffSelector)
	mockTimeMatrixSelector := new(MockTimeMatrixSelector)

	// Reset the PathPointID counter before each test
	offer := createDefaultOfferWithTimeAndCapacity(timeNow, 1)
	offerNode := model.NewOfferNode(offer)
	offerSource := createSourcePathPoint(offer)
	offerDestination := createDestinationPathPoint(offer)

	request := createDefaultRequestWithEarliestDepartureTime(timeNow)
	request.SetLatestArrivalTime(timeNow.Add(31 * time.Minute))
	requestNode := model.NewRequestNode(request)

	pickup, dropoff := createDefaultPickupDropoff(request)
	pickupDropoff := pickupdropoffcache.NewValue(pickup, dropoff)

	matchedRequest := createDefaultRequestWithEarliestDepartureTime(timeNow.Add(8 * time.Minute))
	matchedRequest.SetLatestArrivalTime(timeNow.Add(21 * time.Minute))

	matchedRequestPickup, matchedRequestDropoff := createDefaultPickupDropoff(matchedRequest)

	offerSource.SetExpectedArrivalTime(timeNow)
	pickup.SetExpectedArrivalTime(timeNow.Add(5 * time.Minute))
	matchedRequestPickup.SetExpectedArrivalTime(timeNow.Add(10 * time.Minute))
	matchedRequestDropoff.SetExpectedArrivalTime(timeNow.Add(20 * time.Minute))
	dropoff.SetExpectedArrivalTime(timeNow.Add(30 * time.Minute))
	offerDestination.SetExpectedArrivalTime(timeNow.Add(40 * time.Minute))

	oldPath := []model.PathPoint{*offerSource, *matchedRequestPickup, *matchedRequestDropoff, *offerDestination}

	offer.SetMatchedRequests([]*model.Request{matchedRequest})
	offer.SetPath(oldPath)

	mockPickupDropoffSelector.On("GetPickupDropoffPointsAndDurations", request, offer).Return(pickupDropoff, nil)

	// Create a simple time matrix
	timeMatrix := [][]time.Duration{
		{0 * time.Minute, 5 * time.Minute, 10 * time.Minute, 20 * time.Minute, 30 * time.Minute, 40 * time.Minute},
		{5 * time.Minute, 0 * time.Minute, 5 * time.Minute, 15 * time.Minute, 25 * time.Minute, 35 * time.Minute},
		{10 * time.Minute, 5 * time.Minute, 0 * time.Minute, 10 * time.Minute, 20 * time.Minute, 30 * time.Minute},
		{20 * time.Minute, 15 * time.Minute, 10 * time.Minute, 0 * time.Minute, 10 * time.Minute, 20 * time.Minute},
		{30 * time.Minute, 25 * time.Minute, 20 * time.Minute, 10 * time.Minute, 0 * time.Minute, 10 * time.Minute},
		{40 * time.Minute, 35 * time.Minute, 30 * time.Minute, 20 * time.Minute, 10 * time.Minute, 0 * time.Minute},
	}

	pointIdToIndex := map[model.PathPointID]int{
		1: 0,
		2: 5,
		3: 1,
		4: 4,
		5: 2,
		6: 3,
	}

	mappedMatrix := cache.NewPathPointMappedTimeMatrix(timeMatrix, pointIdToIndex)
	mockTimeMatrixSelector.On("GetTimeMatrix", offerNode).Return(mappedMatrix, nil)

	client, _ := ortool.NewORToolClient()
	planner := planner2.NewORToolPlanner(mockPickupDropoffSelector, mockTimeMatrixSelector, client)
	_, found, err := planner.FindFirstFeasiblePath(offerNode, requestNode)

	assert.NoError(t, err)
	assert.False(t, found)

	mockPickupDropoffSelector.AssertExpectations(t)
	mockTimeMatrixSelector.AssertExpectations(t)
}

func TestORToolPlanner_inValidPathDueToInvalidDetour(t *testing.T) {
	resetPathPointIDCounter()
	timeNow := time.Now()

	mockPickupDropoffSelector := new(MockPickupDropoffSelector)
	mockTimeMatrixSelector := new(MockTimeMatrixSelector)

	// Reset the PathPointID counter before each test
	offer := createDefaultOfferWithTime(timeNow)
	offerNode := model.NewOfferNode(offer)
	offerSource := createSourcePathPoint(offer)
	offerDestination := createDestinationPathPoint(offer)

	request := createDefaultRequestWithEarliestDepartureTime(timeNow)
	request.SetLatestArrivalTime(timeNow.Add(21 * time.Minute))
	requestNode := model.NewRequestNode(request)

	pickup, dropoff := createDefaultPickupDropoff(request)
	pickupDropoff := pickupdropoffcache.NewValue(pickup, dropoff)

	matchedRequest := createDefaultRequestWithEarliestDepartureTime(timeNow.Add(8 * time.Minute))
	matchedRequest.SetLatestArrivalTime(timeNow.Add(21 * time.Minute))

	matchedRequestPickup, matchedRequestDropoff := createDefaultPickupDropoff(matchedRequest)

	offerSource.SetExpectedArrivalTime(timeNow)
	pickup.SetExpectedArrivalTime(timeNow.Add(5 * time.Minute))
	matchedRequestPickup.SetExpectedArrivalTime(timeNow.Add(10 * time.Minute))
	matchedRequestDropoff.SetExpectedArrivalTime(timeNow.Add(20 * time.Minute))
	dropoff.SetExpectedArrivalTime(timeNow.Add(30 * time.Minute))
	offerDestination.SetExpectedArrivalTime(timeNow.Add(40 * time.Minute))

	oldPath := []model.PathPoint{*offerSource, *matchedRequestPickup, *matchedRequestDropoff, *offerDestination}
	
	offer.SetMatchedRequests([]*model.Request{matchedRequest})
	offer.SetPath(oldPath)

	mockPickupDropoffSelector.On("GetPickupDropoffPointsAndDurations", request, offer).Return(pickupDropoff, nil)

	// Create a simple time matrix
	timeMatrix := [][]time.Duration{
		{0 * time.Minute, 5 * time.Minute, 10 * time.Minute, 20 * time.Minute, 30 * time.Minute, 40 * time.Minute},
		{5 * time.Minute, 0 * time.Minute, 5 * time.Minute, 15 * time.Minute, 25 * time.Minute, 35 * time.Minute},
		{10 * time.Minute, 5 * time.Minute, 0 * time.Minute, 10 * time.Minute, 20 * time.Minute, 30 * time.Minute},
		{20 * time.Minute, 15 * time.Minute, 10 * time.Minute, 0 * time.Minute, 10 * time.Minute, 20 * time.Minute},
		{30 * time.Minute, 25 * time.Minute, 20 * time.Minute, 10 * time.Minute, 0 * time.Minute, 10 * time.Minute},
		{40 * time.Minute, 35 * time.Minute, 30 * time.Minute, 20 * time.Minute, 10 * time.Minute, 0 * time.Minute},
	}

	pointIdToIndex := map[model.PathPointID]int{
		1: 0,
		2: 5,
		3: 1,
		4: 4,
		5: 2,
		6: 3,
	}

	mappedMatrix := cache.NewPathPointMappedTimeMatrix(timeMatrix, pointIdToIndex)
	mockTimeMatrixSelector.On("GetTimeMatrix", offerNode).Return(mappedMatrix, nil)

	client, _ := ortool.NewORToolClient()
	planner := planner2.NewORToolPlanner(mockPickupDropoffSelector, mockTimeMatrixSelector, client)
	_, found, err := planner.FindFirstFeasiblePath(offerNode, requestNode)

	assert.NoError(t, err)
	assert.False(t, found)

	mockPickupDropoffSelector.AssertExpectations(t)
	mockTimeMatrixSelector.AssertExpectations(t)
}

func resetPathPointIDCounter() {
	model.NextPointID = 1
}
