package tests

import (
	"context"
	"errors"
	"matching-engine/internal/enums"
	"matching-engine/internal/model"
	"matching-engine/internal/service/pickupdropoffservice/pickupdropoffcache"
	"matching-engine/internal/service/timematrix"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockRoutingEngine mocks the routing.Engine interface
type MockRoutingEngine struct {
	mock.Mock
}

func (m *MockRoutingEngine) PlanDrivingRoute(ctx context.Context, routeParams *model.RouteParams) (*model.Route, error) {
	return nil, errors.New("PlanDrivingRoute should not be called in this test")
}

func (m *MockRoutingEngine) ComputeDrivingTime(ctx context.Context, routeParams *model.RouteParams) ([]time.Duration, error) {
	return nil, errors.New("ComputeDrivingTime should not be called in this test")
}

func (m *MockRoutingEngine) ComputeWalkingTime(ctx context.Context, walkParams *model.WalkParams) (time.Duration, error) {
	return 0, errors.New("ComputeWalkingTime should not be called in this test")
}

func (m *MockRoutingEngine) ComputeIsochrone(ctx context.Context, req *model.IsochroneParams) (*model.Isochrone, error) {
	return nil, errors.New("ComputeIsochrone should not be called in this test")
}

func (m *MockRoutingEngine) SnapPointToRoad(ctx context.Context, point *model.Coordinate) (*model.Coordinate, error) {
	return nil, errors.New("SnapPointToRoad should not be called in this test")
}

func (m *MockRoutingEngine) ComputeDistanceTimeMatrix(ctx context.Context, req *model.DistanceTimeMatrixParams) (*model.DistanceTimeMatrix, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.DistanceTimeMatrix), args.Error(1)
}

// MockPickupDropoffSelector implements the PickupDropoffSelectorInterface interface for testing
type MockPickupDropoffSelector struct {
	mock.Mock
}

func (m *MockPickupDropoffSelector) GetPickupDropoffPointsAndDurations(request *model.Request, offer *model.Offer) (*pickupdropoffcache.Value, error) {
	args := m.Called(request, offer)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pickupdropoffcache.Value), args.Error(1)
}

func TestDefaultGenerator_Generate_Success(t *testing.T) {
	offer := createTestOffer()
	request := createTestRequest()

	sourcePoint := createPathPoint(offer.Source(), enums.Source, offer)
	destinationPoint := createPathPoint(offer.Destination(), enums.Destination, offer)
	offer.SetPath([]model.PathPoint{*sourcePoint, *destinationPoint})

	mockEngine := new(MockRoutingEngine)
	mockPickupDropoffSelector := new(MockPickupDropoffSelector)

	pickupCoord, _ := model.NewCoordinate(1.4, 1.4)
	dropoffCoord, _ := model.NewCoordinate(2.4, 2.4)
	pickup := createPathPoint(pickupCoord, enums.Pickup, request)
	dropoff := createPathPoint(dropoffCoord, enums.Dropoff, request)

	pdValue := pickupdropoffcache.NewValue(pickup, dropoff)
	mockPickupDropoffSelector.On("GetPickupDropoffPointsAndDurations", request, offer).Return(pdValue, nil)

	timeMatrix, distanceMatrix := generateRandomTimeDistanceMatrices(4)

	expectedDistanceTimeParams, _ := model.NewDistanceTimeMatrixParams(
		[]model.Coordinate{
			*offer.Source(),
			*offer.Destination(),
			*pickup.Coordinate(),
			*dropoff.Coordinate(),
		},
		model.ProfileAuto,
		model.WithDepartureTime(offer.DepartureTime()),
	)

	distanceTimeMatrix, _ := model.NewDistanceTimeMatrix(distanceMatrix, timeMatrix)

	mockEngine.On("ComputeDistanceTimeMatrix", mock.Anything, expectedDistanceTimeParams).Return(distanceTimeMatrix, nil)

	generator := timematrix.NewDefaultGenerator(
		mockEngine,
		mockPickupDropoffSelector,
	)

	result, err := generator.Generate(model.NewOfferNode(offer), []*model.RequestNode{model.NewRequestNode(request)})

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, timeMatrix, result.TimeMatrix())

	// Explicitly assert PointIdToIndex content
	pointToIndex := result.PointIdToIndex()
	assert.Equal(t, 4, len(pointToIndex))
	assert.Equal(t, 0, pointToIndex[sourcePoint.ID()])
	assert.Equal(t, 1, pointToIndex[destinationPoint.ID()])
	assert.Equal(t, 2, pointToIndex[pickup.ID()])
	assert.Equal(t, 3, pointToIndex[dropoff.ID()])

	mockPickupDropoffSelector.AssertExpectations(t)
	mockEngine.AssertExpectations(t)
}

func TestDefaultGenerator_Generate_NoPotentialRequests(t *testing.T) {
	// Setup
	offer := createTestOffer()
	mockEngine := new(MockRoutingEngine)
	mockPickupDropoffSelector := new(MockPickupDropoffSelector)

	// Create the generator
	generator := timematrix.NewDefaultGenerator(
		mockEngine,
		mockPickupDropoffSelector,
	)

	// Call the Generate method
	result, err := generator.Generate(model.NewOfferNode(offer), []*model.RequestNode{})

	// Assertions
	require.Error(t, err)
	assert.Nil(t, result)
}

func TestDefaultGenerator_Generate_PickupDropoffError(t *testing.T) {
	// Setup
	offer := createTestOffer()
	request := createTestRequest()

	// Create source and destination path points for the offer
	sourcePoint := createPathPoint(offer.Source(), enums.Source, offer)
	destinationPoint := createPathPoint(offer.Destination(), enums.Destination, offer)

	// Set path points on the offer
	offer.SetPath([]model.PathPoint{*sourcePoint, *destinationPoint})

	// Create mocks
	mockEngine := new(MockRoutingEngine)
	mockPickupDropoffSelector := new(MockPickupDropoffSelector)

	// Setup pickup/dropoff points error
	mockPickupDropoffSelector.On("GetPickupDropoffPointsAndDurations", request, offer).
		Return(nil, errors.New("failed to get pickup/dropoff points"))

	// Create the generator
	generator := timematrix.NewDefaultGenerator(
		mockEngine,
		mockPickupDropoffSelector,
	)

	// Call the Generate method
	result, err := generator.Generate(model.NewOfferNode(offer), []*model.RequestNode{model.NewRequestNode(request)})

	// Assertions
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get pickup/dropoff points")
	assert.Nil(t, result)

	// Verify the expected calls
	mockPickupDropoffSelector.AssertExpectations(t)
}

func TestDefaultGenerator_Generate_DistanceTimeMatrixError(t *testing.T) {
	// Setup
	offer := createTestOffer()
	request := createTestRequest()

	// Create source and destination path points for the offer
	sourcePoint := createPathPoint(offer.Source(), enums.Source, offer)
	destinationPoint := createPathPoint(offer.Destination(), enums.Destination, offer)

	// Set path points on the offer
	offer.SetPath([]model.PathPoint{*sourcePoint, *destinationPoint})

	// Create mocks
	mockEngine := new(MockRoutingEngine)
	mockPickupDropoffSelector := new(MockPickupDropoffSelector)

	// Setup pickup/dropoff points
	pickupCoord, _ := model.NewCoordinate(1.4, 1.4)
	dropoffCoord, _ := model.NewCoordinate(2.4, 2.4)

	pickup := createPathPoint(pickupCoord, enums.Pickup, request)
	dropoff := createPathPoint(dropoffCoord, enums.Dropoff, request)

	pdValue := pickupdropoffcache.NewValue(pickup, dropoff)
	mockPickupDropoffSelector.On("GetPickupDropoffPointsAndDurations", request, offer).Return(pdValue, nil)

	// Setup mock distance/time matrix error
	mockEngine.On("ComputeDistanceTimeMatrix", mock.Anything, mock.AnythingOfType("*model.DistanceTimeMatrixParams")).
		Return(nil, errors.New("routing engine error"))

	// Create the generator
	generator := timematrix.NewDefaultGenerator(
		mockEngine,
		mockPickupDropoffSelector,
	)

	// Call the Generate method
	result, err := generator.Generate(model.NewOfferNode(offer), []*model.RequestNode{model.NewRequestNode(request)})

	// Assertions
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to compute distance time matrix")
	assert.Nil(t, result)

	// Verify the expected calls
	mockPickupDropoffSelector.AssertExpectations(t)
	mockEngine.AssertExpectations(t)
}

func TestDefaultGenerator_Generate_MultipleRequests(t *testing.T) {
	// Setup
	offer := createTestOffer()

	// Create two requests
	request1 := createTestRequest()

	source2, _ := model.NewCoordinate(1.7, 1.7)
	destination2, _ := model.NewCoordinate(2.7, 2.7)
	request2 := model.NewRequest(
		"request-789",
		"user-abc",
		*source2,
		*destination2,
		time.Now(),
		time.Now().Add(2*time.Hour),
		10*time.Minute,
		1,
		model.Preference{},
	)

	// Create source and destination path points for the offer
	sourcePoint := createPathPoint(offer.Source(), enums.Source, offer)
	destinationPoint := createPathPoint(offer.Destination(), enums.Destination, offer)

	// Set path points on the offer
	offer.SetPath([]model.PathPoint{*sourcePoint, *destinationPoint})

	// Create mocks
	mockEngine := new(MockRoutingEngine)
	mockPickupDropoffSelector := new(MockPickupDropoffSelector)

	// Setup pickup/dropoff points for request1
	pickup1Coord, _ := model.NewCoordinate(1.4, 1.4)
	dropoff1Coord, _ := model.NewCoordinate(2.4, 2.4)
	pickup1 := createPathPoint(pickup1Coord, enums.Pickup, request1)
	dropoff1 := createPathPoint(dropoff1Coord, enums.Dropoff, request1)
	pdValue1 := pickupdropoffcache.NewValue(pickup1, dropoff1)
	mockPickupDropoffSelector.On("GetPickupDropoffPointsAndDurations", request1, offer).Return(pdValue1, nil)

	// Setup pickup/dropoff points for request2
	pickup2Coord, _ := model.NewCoordinate(1.6, 1.6)
	dropoff2Coord, _ := model.NewCoordinate(2.6, 2.6)
	pickup2 := createPathPoint(pickup2Coord, enums.Pickup, request2)
	dropoff2 := createPathPoint(dropoff2Coord, enums.Dropoff, request2)
	pdValue2 := pickupdropoffcache.NewValue(pickup2, dropoff2)
	mockPickupDropoffSelector.On("GetPickupDropoffPointsAndDurations", request2, offer).Return(pdValue2, nil)

	timeMatrix, distanceMatrix := generateRandomTimeDistanceMatrices(6)

	mockEngine.On("ComputeDistanceTimeMatrix", mock.Anything, mock.AnythingOfType("*model.DistanceTimeMatrixParams")).
		Return(must(model.NewDistanceTimeMatrix(distanceMatrix, timeMatrix)), nil)

	// Create the generator
	generator := timematrix.NewDefaultGenerator(
		mockEngine,
		mockPickupDropoffSelector,
	)

	// Call the Generate method
	result, err := generator.Generate(model.NewOfferNode(offer), []*model.RequestNode{model.NewRequestNode(request1), model.NewRequestNode(request2)})

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, timeMatrix, result.TimeMatrix())
	assert.Len(t, result.PointIdToIndex(), 6) // 2 offer points + 2 pickups + 2 dropoffs

	// Verify the expected calls
	mockPickupDropoffSelector.AssertExpectations(t)
	mockEngine.AssertExpectations(t)
}

func TestDefaultGenerator_Generate_InvalidDistanceTimeMatrixParams(t *testing.T) {
	// Setup with special mock engine that fails parameter validation
	mockEngine := new(MockRoutingEngine)
	mockEngine.On("ComputeDistanceTimeMatrix", mock.Anything, mock.AnythingOfType("*model.DistanceTimeMatrixParams")).
		Return(nil, errors.New("invalid parameter: sources cannot be empty"))

	offer := createTestOffer()
	request := createTestRequest()

	// Create source and destination path points for the offer
	sourcePoint := createPathPoint(offer.Source(), enums.Source, offer)
	destinationPoint := createPathPoint(offer.Destination(), enums.Destination, offer)

	// Set path points on the offer
	offer.SetPath([]model.PathPoint{*sourcePoint, *destinationPoint})

	mockPickupDropoffSelector := new(MockPickupDropoffSelector)

	// Setup pickup/dropoff points
	pickupCoord, _ := model.NewCoordinate(1.4, 1.4)
	dropoffCoord, _ := model.NewCoordinate(2.4, 2.4)

	pickup := createPathPoint(pickupCoord, enums.Pickup, request)
	dropoff := createPathPoint(dropoffCoord, enums.Dropoff, request)

	pdValue := pickupdropoffcache.NewValue(pickup, dropoff)
	mockPickupDropoffSelector.On("GetPickupDropoffPointsAndDurations", request, offer).Return(pdValue, nil)

	// Create the generator
	generator := timematrix.NewDefaultGenerator(
		mockEngine,
		mockPickupDropoffSelector,
	)

	// Call the Generate method
	result, err := generator.Generate(model.NewOfferNode(offer), []*model.RequestNode{model.NewRequestNode(request)})

	// Assertions
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to compute distance time matrix")
	assert.Nil(t, result)

	// Verify the expected calls
	mockPickupDropoffSelector.AssertExpectations(t)
	mockEngine.AssertExpectations(t)
}
