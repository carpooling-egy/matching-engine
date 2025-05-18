package tests

import (
	"errors"
	"matching-engine/internal/enums"
	"matching-engine/internal/model"
	"matching-engine/internal/service/timematrix"
	"matching-engine/internal/service/timematrix/cache"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
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

func TestTimeMatrixService_GetTravelDuration_Success(t *testing.T) {
	// Reset the PathPointID counter before each test
	resetPathPointIDCounter()

	// Setup
	offer := createTestOffer()
	offerNode := model.NewOfferNode(offer)

	// Create mock selector
	mockSelector := new(MockTimeMatrixSelector)

	// Create a simple time matrix
	timeMatrix := [][]time.Duration{
		{0 * time.Minute, 10 * time.Minute, 15 * time.Minute},
		{10 * time.Minute, 0 * time.Minute, 5 * time.Minute},
		{15 * time.Minute, 5 * time.Minute, 0 * time.Minute},
	}

	pointIdToIndex := map[model.PathPointID]int{
		1: 0,
		2: 1,
		3: 2,
	}

	mappedMatrix := cache.NewPathPointMappedTimeMatrix(timeMatrix, pointIdToIndex)
	mockSelector.On("GetTimeMatrix", offerNode).Return(mappedMatrix, nil)

	// Create service with mock selector
	service := timematrix.NewTimeMatrixService(mockSelector)

	// Call the method to test
	duration, err := service.GetTravelDuration(offerNode, 1, 3)

	// Assertions
	require.NoError(t, err)
	assert.Equal(t, 15*time.Minute, duration)

	// Verify expectations
	mockSelector.AssertExpectations(t)
}

func TestTimeMatrixService_GetTravelDuration_SelectorError(t *testing.T) {
	// Reset the PathPointID counter before each test
	resetPathPointIDCounter()

	// Setup
	offer := createTestOffer()
	offerNode := model.NewOfferNode(offer)

	// Create mock selector with error
	mockSelector := new(MockTimeMatrixSelector)
	mockSelector.On("GetTimeMatrix", offerNode).Return(nil, errors.New("failed to get time matrix"))

	// Create service with mock selector
	service := timematrix.NewTimeMatrixService(mockSelector)

	// Call the method to test
	duration, err := service.GetTravelDuration(offerNode, 1, 3)

	// Assertions
	require.Error(t, err)
	assert.Equal(t, time.Duration(0), duration)
	assert.Contains(t, err.Error(), "failed to get time matrix")

	// Verify expectations
	mockSelector.AssertExpectations(t)
}

func TestTimeMatrixService_GetTravelDuration_InvalidPointID(t *testing.T) {
	// Reset the PathPointID counter before each test
	resetPathPointIDCounter()

	// Setup
	offer := createTestOffer()
	offerNode := model.NewOfferNode(offer)

	// Create mock selector
	mockSelector := new(MockTimeMatrixSelector)

	// Create a simple time matrix
	timeMatrix := [][]time.Duration{
		{0 * time.Minute, 10 * time.Minute, 15 * time.Minute},
		{10 * time.Minute, 0 * time.Minute, 5 * time.Minute},
		{15 * time.Minute, 5 * time.Minute, 0 * time.Minute},
	}

	pointIdToIndex := map[model.PathPointID]int{
		1: 0,
		2: 1,
		3: 2,
	}

	mappedMatrix := cache.NewPathPointMappedTimeMatrix(timeMatrix, pointIdToIndex)
	mockSelector.On("GetTimeMatrix", offerNode).Return(mappedMatrix, nil)

	// Create service with mock selector
	service := timematrix.NewTimeMatrixService(mockSelector)

	// Call the method to test with invalid point ID
	duration, err := service.GetTravelDuration(offerNode, 1, 4)

	// Assertions
	require.Error(t, err)
	assert.Equal(t, time.Duration(0), duration)
	assert.Contains(t, err.Error(), "invalid path point ID")

	// Verify expectations
	mockSelector.AssertExpectations(t)
}

func TestTimeMatrixService_GetTravelDuration_IndexOutOfBounds(t *testing.T) {
	// Reset the PathPointID counter before each test
	resetPathPointIDCounter()

	// Setup
	offer := createTestOffer()
	offerNode := model.NewOfferNode(offer)

	// Create mock selector
	mockSelector := new(MockTimeMatrixSelector)

	// Create a malformed time matrix with incorrect dimensions
	timeMatrix := [][]time.Duration{
		{0 * time.Minute},
		{10 * time.Minute, 0 * time.Minute},
	}

	pointIdToIndex := map[model.PathPointID]int{
		1: 0,
		2: 1, // This index is out of bounds for the first row
	}

	mappedMatrix := cache.NewPathPointMappedTimeMatrix(timeMatrix, pointIdToIndex)
	mockSelector.On("GetTimeMatrix", offerNode).Return(mappedMatrix, nil)

	// Create service with mock selector
	service := timematrix.NewTimeMatrixService(mockSelector)

	// Call the method to test
	duration, err := service.GetTravelDuration(offerNode, 1, 2)

	// Assertions
	require.Error(t, err)
	assert.Equal(t, time.Duration(0), duration)
	assert.Contains(t, err.Error(), "index out of bounds")

	// Verify expectations
	mockSelector.AssertExpectations(t)
}

func TestTimeMatrixService_GetCumulativeTravelDurations_Success(t *testing.T) {
	// Reset the PathPointID counter before each test
	resetPathPointIDCounter()

	// Setup
	offer := createTestOffer()
	offerNode := model.NewOfferNode(offer)

	// Create mock selector
	mockSelector := new(MockTimeMatrixSelector)

	// Create a simple time matrix
	timeMatrix := [][]time.Duration{
		{0 * time.Minute, 10 * time.Minute, 15 * time.Minute, 20 * time.Minute},
		{10 * time.Minute, 0 * time.Minute, 5 * time.Minute, 12 * time.Minute},
		{15 * time.Minute, 5 * time.Minute, 0 * time.Minute, 8 * time.Minute},
		{20 * time.Minute, 12 * time.Minute, 8 * time.Minute, 0 * time.Minute},
	}

	pointIdToIndex := map[model.PathPointID]int{
		1: 0,
		2: 1,
		3: 2,
		4: 3,
	}

	mappedMatrix := cache.NewPathPointMappedTimeMatrix(timeMatrix, pointIdToIndex)
	mockSelector.On("GetTimeMatrix", offerNode).Return(mappedMatrix, nil)

	// Create path points for the test
	source, _ := model.NewCoordinate(1.0, 1.0)
	p1 := createPathPoint(source, enums.Source, offer)

	mid1, _ := model.NewCoordinate(1.5, 1.5)
	p2 := createPathPoint(mid1, enums.Pickup, offer)

	mid2, _ := model.NewCoordinate(1.8, 1.8)
	p3 := createPathPoint(mid2, enums.Dropoff, offer)

	dest, _ := model.NewCoordinate(2.0, 2.0)
	p4 := createPathPoint(dest, enums.Destination, offer)

	pathPoints := []model.PathPoint{*p1, *p2, *p3, *p4}

	// Create service with mock selector
	service := timematrix.NewTimeMatrixService(mockSelector)

	// Call the method to test
	cumulativeDurations, err := service.GetCumulativeTravelDurations(offerNode, pathPoints)

	// Assertions
	require.NoError(t, err)
	assert.Equal(t, 4, len(cumulativeDurations))
	assert.Equal(t, time.Duration(0), cumulativeDurations[0])
	assert.Equal(t, 10*time.Minute, cumulativeDurations[1])
	assert.Equal(t, 15*time.Minute, cumulativeDurations[2]) // 10 + 5
	assert.Equal(t, 23*time.Minute, cumulativeDurations[3]) // 10 + 5 + 8

	// Verify expectations
	mockSelector.AssertExpectations(t)
}

func TestTimeMatrixService_GetCumulativeTravelDurations_SelectorError(t *testing.T) {
	// Reset the PathPointID counter before each test
	resetPathPointIDCounter()

	// Setup
	offer := createTestOffer()
	offerNode := model.NewOfferNode(offer)

	// Create mock selector with error
	mockSelector := new(MockTimeMatrixSelector)
	mockSelector.On("GetTimeMatrix", offerNode).Return(nil, errors.New("failed to get time matrix"))

	// Create path points for the test
	source, _ := model.NewCoordinate(1.0, 1.0)
	p1 := createPathPoint(source, enums.Source, offer)

	dest, _ := model.NewCoordinate(2.0, 2.0)
	p2 := createPathPoint(dest, enums.Destination, offer)

	pathPoints := []model.PathPoint{*p1, *p2}

	// Create service with mock selector
	service := timematrix.NewTimeMatrixService(mockSelector)

	// Call the method to test
	cumulativeDurations, err := service.GetCumulativeTravelDurations(offerNode, pathPoints)

	// Assertions
	require.Error(t, err)
	assert.Nil(t, cumulativeDurations)
	assert.Contains(t, err.Error(), "failed to get time matrix")

	// Verify expectations
	mockSelector.AssertExpectations(t)
}

func TestTimeMatrixService_GetCumulativeTravelDurations_TooFewPoints(t *testing.T) {
	// Reset the PathPointID counter before each test
	resetPathPointIDCounter()

	// Setup
	offer := createTestOffer()
	offerNode := model.NewOfferNode(offer)

	// Create mock selector
	mockSelector := new(MockTimeMatrixSelector)

	// Create path points for the test - only one point
	source, _ := model.NewCoordinate(1.0, 1.0)
	p1 := createPathPoint(source, enums.Source, offer)

	pathPoints := []model.PathPoint{*p1}

	// Create service with mock selector
	service := timematrix.NewTimeMatrixService(mockSelector)

	// Call the method to test
	cumulativeDurations, err := service.GetCumulativeTravelDurations(offerNode, pathPoints)

	// Assertions
	require.Error(t, err)
	assert.Nil(t, cumulativeDurations)
	assert.Contains(t, err.Error(), "must contain at least two points")

	// Mock is not called because we check point count first
	mockSelector.AssertNotCalled(t, "GetTimeMatrix")
}

func TestTimeMatrixService_GetCumulativeTravelDurations_InvalidPointID(t *testing.T) {
	// Setup
	offer := createTestOffer()
	offerNode := model.NewOfferNode(offer)

	// Create mock selector
	mockSelector := new(MockTimeMatrixSelector)

	// Create a simple time matrix with only two points
	timeMatrix := [][]time.Duration{
		{0 * time.Minute, 10 * time.Minute},
		{10 * time.Minute, 0 * time.Minute},
	}

	pointIdToIndex := map[model.PathPointID]int{
		1: 0,
		2: 1,
	}

	mappedMatrix := cache.NewPathPointMappedTimeMatrix(timeMatrix, pointIdToIndex)
	mockSelector.On("GetTimeMatrix", offerNode).Return(mappedMatrix, nil)

	// Create path points for the test - including an unknown point
	source, _ := model.NewCoordinate(1.0, 1.0)
	p1 := createPathPoint(source, enums.Source, offer)

	mid, _ := model.NewCoordinate(1.5, 1.5)
	p2 := createPathPoint(mid, enums.Pickup, offer) // This ID is not in the matrix

	dest, _ := model.NewCoordinate(2.0, 2.0)
	p3 := createPathPoint(dest, enums.Destination, offer)

	pathPoints := []model.PathPoint{*p1, *p2, *p3}

	// Create service with mock selector
	service := timematrix.NewTimeMatrixService(mockSelector)

	// Call the method to test
	cumulativeDurations, err := service.GetCumulativeTravelDurations(offerNode, pathPoints)

	// Assertions
	require.Error(t, err)
	assert.Nil(t, cumulativeDurations)
	assert.Contains(t, err.Error(), "invalid path point ID")

	// Verify expectations
	mockSelector.AssertExpectations(t)
}

func TestTimeMatrixService_GetCumulativeTravelTimes_Success(t *testing.T) {
	// Reset the PathPointID counter before each test
	resetPathPointIDCounter()

	// Setup
	departureTime := time.Date(2023, 1, 1, 10, 0, 0, 0, time.UTC)
	source, _ := model.NewCoordinate(1.0, 1.0)
	destination, _ := model.NewCoordinate(2.0, 2.0)

	offer := model.NewOffer(
		"offer-123",
		"user-abc",
		*source,
		*destination,
		departureTime,
		15*time.Minute,
		4,
		model.Preference{},
		departureTime.Add(1*time.Hour),
		0,
		nil,
		nil,
	)
	offerNode := model.NewOfferNode(offer)

	// Create mock selector
	mockSelector := new(MockTimeMatrixSelector)

	// Create a simple time matrix
	timeMatrix := [][]time.Duration{
		{0 * time.Minute, 10 * time.Minute, 15 * time.Minute, 20 * time.Minute},
		{10 * time.Minute, 0 * time.Minute, 5 * time.Minute, 12 * time.Minute},
		{15 * time.Minute, 5 * time.Minute, 0 * time.Minute, 8 * time.Minute},
		{20 * time.Minute, 12 * time.Minute, 8 * time.Minute, 0 * time.Minute},
	}

	pointIdToIndex := map[model.PathPointID]int{
		1: 0,
		2: 1,
		3: 2,
		4: 3,
	}

	mappedMatrix := cache.NewPathPointMappedTimeMatrix(timeMatrix, pointIdToIndex)
	mockSelector.On("GetTimeMatrix", offerNode).Return(mappedMatrix, nil)

	// Create path points for the test
	p1 := createPathPoint(source, enums.Source, offer)

	mid1, _ := model.NewCoordinate(1.5, 1.5)
	p2 := createPathPoint(mid1, enums.Pickup, offer)

	mid2, _ := model.NewCoordinate(1.8, 1.8)
	p3 := createPathPoint(mid2, enums.Dropoff, offer)

	p4 := createPathPoint(destination, enums.Destination, offer)

	pathPoints := []model.PathPoint{*p1, *p2, *p3, *p4}

	// Create service with mock selector
	service := timematrix.NewTimeMatrixService(mockSelector)

	// Call the method to test
	cumulativeTimes, err := service.GetCumulativeTravelTimes(offerNode, pathPoints)

	// Assertions
	require.NoError(t, err)
	assert.Equal(t, 4, len(cumulativeTimes))
	assert.Equal(t, departureTime, cumulativeTimes[0])
	assert.Equal(t, departureTime.Add(10*time.Minute), cumulativeTimes[1])
	assert.Equal(t, departureTime.Add(15*time.Minute), cumulativeTimes[2]) // departure + 10 + 5
	assert.Equal(t, departureTime.Add(23*time.Minute), cumulativeTimes[3]) // departure + 10 + 5 + 8

	// Verify expectations
	mockSelector.AssertExpectations(t)
}

func TestTimeMatrixService_GetCumulativeTravelTimes_SelectorError(t *testing.T) {
	// Reset the PathPointID counter before each test
	resetPathPointIDCounter()

	// Setup
	offer := createTestOffer()
	offerNode := model.NewOfferNode(offer)

	// Create mock selector with error
	mockSelector := new(MockTimeMatrixSelector)
	mockSelector.On("GetTimeMatrix", offerNode).Return(nil, errors.New("failed to get time matrix"))

	// Create path points for the test
	source, _ := model.NewCoordinate(1.0, 1.0)
	p1 := createPathPoint(source, enums.Source, offer)

	dest, _ := model.NewCoordinate(2.0, 2.0)
	p2 := createPathPoint(dest, enums.Destination, offer)

	pathPoints := []model.PathPoint{*p1, *p2}

	// Create service with mock selector
	service := timematrix.NewTimeMatrixService(mockSelector)

	// Call the method to test
	cumulativeTimes, err := service.GetCumulativeTravelTimes(offerNode, pathPoints)

	// Assertions
	require.Error(t, err)
	assert.Nil(t, cumulativeTimes)
	assert.Contains(t, err.Error(), "failed to get time matrix")

	// Verify expectations
	mockSelector.AssertExpectations(t)
}

func TestTimeMatrixService_GetCumulativeTravelTimes_TooFewPoints(t *testing.T) {
	// Reset the PathPointID counter before each test
	resetPathPointIDCounter()

	// Setup
	offer := createTestOffer()
	offerNode := model.NewOfferNode(offer)

	// Create mock selector
	mockSelector := new(MockTimeMatrixSelector)

	// Create path points for the test - only one point
	source, _ := model.NewCoordinate(1.0, 1.0)
	p1 := createPathPoint(source, enums.Source, offer)

	pathPoints := []model.PathPoint{*p1}

	// Create service with mock selector
	service := timematrix.NewTimeMatrixService(mockSelector)

	// Call the method to test
	cumulativeTimes, err := service.GetCumulativeTravelTimes(offerNode, pathPoints)

	// Assertions
	require.Error(t, err)
	assert.Nil(t, cumulativeTimes)
	assert.Contains(t, err.Error(), "must contain at least two points")

	// Mock is not called because we check point count first
	mockSelector.AssertNotCalled(t, "GetTimeMatrix")
}

func TestTimeMatrixService_GetCumulativeTravelTimes_InvalidPointID(t *testing.T) {
	// Reset the PathPointID counter before each test
	resetPathPointIDCounter()

	// Setup
	offer := createTestOffer()
	offerNode := model.NewOfferNode(offer)

	// Create mock selector
	mockSelector := new(MockTimeMatrixSelector)

	// Create a simple time matrix with only two points
	timeMatrix := [][]time.Duration{
		{0 * time.Minute, 10 * time.Minute},
		{10 * time.Minute, 0 * time.Minute},
	}

	pointIdToIndex := map[model.PathPointID]int{
		1: 0,
		2: 1,
	}

	mappedMatrix := cache.NewPathPointMappedTimeMatrix(timeMatrix, pointIdToIndex)
	mockSelector.On("GetTimeMatrix", offerNode).Return(mappedMatrix, nil)

	// Create path points for the test - including an unknown point
	source, _ := model.NewCoordinate(1.0, 1.0)
	p1 := createPathPoint(source, enums.Source, offer)

	mid, _ := model.NewCoordinate(1.5, 1.5)
	p2 := createPathPoint(mid, enums.Pickup, offer)

	dest, _ := model.NewCoordinate(2.0, 2.0)
	p3 := createPathPoint(dest, enums.Destination, offer) // This ID is not in the matrix

	pathPoints := []model.PathPoint{*p1, *p2, *p3}

	// Create service with mock selector
	service := timematrix.NewTimeMatrixService(mockSelector)

	// Call the method to test
	cumulativeTimes, err := service.GetCumulativeTravelTimes(offerNode, pathPoints)

	// Assertions
	require.Error(t, err)
	assert.Nil(t, cumulativeTimes)
	assert.Contains(t, err.Error(), "invalid path point ID")

	// Verify expectations
	mockSelector.AssertExpectations(t)
}
