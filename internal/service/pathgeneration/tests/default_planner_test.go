package tests

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"iter"
	"matching-engine/internal/model"
	"matching-engine/internal/service/pathgeneration/planner"
	"matching-engine/internal/service/pickupdropoffservice/pickupdropoffcache"
	"testing"
)

// MockPathGenerator implements the PathGenerator interface for testing
type MockPathGenerator struct {
	mock.Mock
}

func (m *MockPathGenerator) GeneratePaths(path []model.PathPoint, pickup, dropoff *model.PathPoint) (iter.Seq2[[]model.PathPoint, error], error) {
	args := m.Called(path, pickup, dropoff)

	// Handle error case
	if err := args.Error(1); err != nil {
		return nil, err
	}

	paths := args.Get(0).([][]model.PathPoint)
	iterFunc := func(yield func([]model.PathPoint, error) bool) {
		for _, p := range paths {
			if !yield(p, nil) {
				break
			}
		}
	}
	return iterFunc, nil
}

// MockPathValidator implements the PathValidator interface for testing
type MockPathValidator struct {
	mock.Mock
}

func (m *MockPathValidator) ValidatePath(offerNode *model.OfferNode, requestNode *model.RequestNode, path []model.PathPoint) (bool, error) {
	args := m.Called(offerNode, requestNode, path)
	return args.Bool(0), args.Error(1)
}

// TestFindFirstFeasiblePath_SimpleSuccess tests the happy path where a valid path is found
func TestFindFirstFeasiblePath_SimpleSuccess(t *testing.T) {
	mockGenerator := new(MockPathGenerator)
	mockValidator := new(MockPathValidator)
	mockSelector := new(MockPickupDropoffSelector)

	offer := createDefaultOffer()
	offerNode := model.NewOfferNode(offer)

	request := createDefaultRequest()
	requestNode := model.NewRequestNode(request)

	pickup, dropoff := createDefaultPickupDropoff(request)
	validPath := []model.PathPoint{*pickup, *dropoff}
	validPaths := [][]model.PathPoint{validPath}

	pickupDropoff := pickupdropoffcache.NewValue(pickup, dropoff)
	mockSelector.On("GetPickupDropoffPointsAndDurations", request, offer).Return(pickupDropoff, nil)
	mockGenerator.On("GeneratePaths", offer.Path(), pickup, dropoff).Return(validPaths, nil)
	mockValidator.On("ValidatePath", offerNode, requestNode, validPath).Return(true, nil)

	planner := planner.NewDefaultPathPlanner(mockGenerator, mockValidator, mockSelector)
	resultPath, found, err := planner.FindFirstFeasiblePath(offerNode, requestNode)

	assert.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, validPath, resultPath)

	mockGenerator.AssertExpectations(t)
	mockValidator.AssertExpectations(t)
	mockSelector.AssertExpectations(t)
}

// TestFindFirstFeasiblePath_GeneratorError tests the case where the path generator returns an error
func TestFindFirstFeasiblePath_GeneratorError(t *testing.T) {
	// Create mock dependencies
	mockGenerator := new(MockPathGenerator)
	mockValidator := new(MockPathValidator)
	mockSelector := new(MockPickupDropoffSelector)

	// Create basic test data
	offer := createDefaultOffer()
	offerNode := model.NewOfferNode(offer)

	request := createDefaultRequest()
	requestNode := model.NewRequestNode(request)

	// Create pickup/dropoff points
	pickup, dropoff := createDefaultPickupDropoff(request)

	// Setup mock selector to return our points
	pickupDropoff := pickupdropoffcache.NewValue(pickup, dropoff)
	mockSelector.On("GetPickupDropoffPointsAndDurations", request, offer).Return(pickupDropoff, nil)

	// Setup mock generator to return an error when called
	expectedErr := errors.New("generator error")
	mockGenerator.On("GeneratePaths", offer.Path(), pickup, dropoff).Return([][]model.PathPoint{}, expectedErr)

	// Create planner and run test
	planner := planner.NewDefaultPathPlanner(mockGenerator, mockValidator, mockSelector)
	resultPath, found, err := planner.FindFirstFeasiblePath(offerNode, requestNode)

	// Verify results - should have error
	assert.Error(t, err)
	assert.False(t, found)
	assert.Nil(t, resultPath)
	assert.Contains(t, err.Error(), expectedErr.Error())

	// Verify mocks were called correctly
	mockGenerator.AssertExpectations(t)
	mockSelector.AssertExpectations(t)
}

// TestFindFirstFeasiblePath_PickupDropoffSelectorError tests the case where the pickup/dropoff selector returns an error
func TestFindFirstFeasiblePath_PickupDropoffSelectorError(t *testing.T) {
	// Create mock dependencies
	mockGenerator := new(MockPathGenerator)
	mockValidator := new(MockPathValidator)
	mockSelector := new(MockPickupDropoffSelector)

	// Create basic test data
	offer := createDefaultOffer()
	offerNode := model.NewOfferNode(offer)

	request := createDefaultRequest()
	requestNode := model.NewRequestNode(request)

	// Setup mock selector to return an error when called
	expectedErr := errors.New("pickup/dropoff selector error")
	mockSelector.On("GetPickupDropoffPointsAndDurations", request, offer).Return(nil, expectedErr)

	// Create planner and run test
	planner := planner.NewDefaultPathPlanner(mockGenerator, mockValidator, mockSelector)
	resultPath, found, err := planner.FindFirstFeasiblePath(offerNode, requestNode)

	// Verify results - should have error
	assert.Error(t, err)
	assert.False(t, found)
	assert.Nil(t, resultPath)
	assert.Contains(t, err.Error(), expectedErr.Error())

	// Verify mocks were called correctly
	mockSelector.AssertExpectations(t)
	// Generator and validator should not be called since selector failed
	mockGenerator.AssertNotCalled(t, "GeneratePaths")
	mockValidator.AssertNotCalled(t, "ValidatePath")
}

// TestFindFirstFeasiblePath_NoValidPaths tests the case where no valid paths are found
func TestFindFirstFeasiblePath_NoValidPaths(t *testing.T) {
	// Create mock dependencies
	mockGenerator := new(MockPathGenerator)
	mockValidator := new(MockPathValidator)
	mockSelector := new(MockPickupDropoffSelector)

	// Create test data
	offer := createDefaultOffer()
	offerNode := model.NewOfferNode(offer)

	request := createDefaultRequest()
	requestNode := model.NewRequestNode(request)

	// Create pickup/dropoff points
	pickup, dropoff := createDefaultPickupDropoff(request)

	// Create paths that will all be invalid
	path1 := []model.PathPoint{*pickup, *dropoff}
	path2 := []model.PathPoint{*pickup, *dropoff, *pickup} // Invalid path for testing

	// We'll have our iterator yield these paths
	candidatePaths := [][]model.PathPoint{path1, path2}

	// Setup mock selector
	pickupDropoff := pickupdropoffcache.NewValue(pickup, dropoff)
	mockSelector.On("GetPickupDropoffPointsAndDurations", request, offer).Return(pickupDropoff, nil)

	// Setup mock generator
	mockGenerator.On("GeneratePaths", offer.Path(), pickup, dropoff).Return(candidatePaths, nil)

	// Setup mock validator to reject all paths
	mockValidator.On("ValidatePath", offerNode, requestNode, path1).Return(false, nil)
	mockValidator.On("ValidatePath", offerNode, requestNode, path2).Return(false, nil)

	// Create planner and run test
	planner := planner.NewDefaultPathPlanner(mockGenerator, mockValidator, mockSelector)
	resultPath, found, err := planner.FindFirstFeasiblePath(offerNode, requestNode)

	// Verify results - no error, but no path found
	assert.NoError(t, err)
	assert.False(t, found)
	assert.Nil(t, resultPath)

	// Verify mocks were called correctly
	mockGenerator.AssertExpectations(t)
	mockValidator.AssertExpectations(t)
	mockSelector.AssertExpectations(t)
}

// TestFindFirstFeasiblePath_ValidatorError tests the case where the path validator returns an error
func TestFindFirstFeasiblePath_ValidatorError(t *testing.T) {
	// Create mock dependencies
	mockGenerator := new(MockPathGenerator)
	mockValidator := new(MockPathValidator)
	mockSelector := new(MockPickupDropoffSelector)

	// Create test data
	offer := createDefaultOffer()
	offerNode := model.NewOfferNode(offer)

	request := createDefaultRequest()
	requestNode := model.NewRequestNode(request)

	// Create pickup/dropoff points
	pickup, dropoff := createDefaultPickupDropoff(request)

	// Create a test path
	candidatePath := []model.PathPoint{*pickup, *dropoff}

	// Setup mock selector
	pickupDropoff := pickupdropoffcache.NewValue(pickup, dropoff)
	mockSelector.On("GetPickupDropoffPointsAndDurations", request, offer).Return(pickupDropoff, nil)

	// Setup mock generator
	mockGenerator.On("GeneratePaths", offer.Path(), pickup, dropoff).Return([][]model.PathPoint{candidatePath}, nil)

	// Setup mock validator to return an error
	expectedErr := errors.New("validation error")
	mockValidator.On("ValidatePath", offerNode, requestNode, candidatePath).Return(false, expectedErr)

	// Create planner and run test
	planner := planner.NewDefaultPathPlanner(mockGenerator, mockValidator, mockSelector)
	resultPath, found, err := planner.FindFirstFeasiblePath(offerNode, requestNode)

	// Verify results - should have error
	assert.Error(t, err)
	assert.False(t, found)
	assert.Nil(t, resultPath)
	assert.Contains(t, err.Error(), expectedErr.Error())

	// Verify mocks were called correctly
	mockGenerator.AssertExpectations(t)
	mockValidator.AssertExpectations(t)
	mockSelector.AssertExpectations(t)
}

// TestFindFirstFeasiblePath_MultiplePathsFirstInvalid tests the case where multiple paths are generated
// but the first one is invalid and the second one is valid
func TestFindFirstFeasiblePath_MultiplePathsFirstInvalid(t *testing.T) {
	// Create mock dependencies
	mockGenerator := new(MockPathGenerator)
	mockValidator := new(MockPathValidator)
	mockSelector := new(MockPickupDropoffSelector)

	// Create test data
	offer := createDefaultOffer()
	offerNode := model.NewOfferNode(offer)

	request := createDefaultRequest()
	requestNode := model.NewRequestNode(request)

	// Create pickup/dropoff points
	pickup, dropoff := createDefaultPickupDropoff(request)

	// Create paths - first invalid, second valid
	invalidPath := []model.PathPoint{*pickup, *dropoff, *dropoff}
	validPath := []model.PathPoint{*pickup, *dropoff, *pickup} // Just for testing differentiation

	// We'll have our iterator yield these paths
	candidatePaths := [][]model.PathPoint{invalidPath, validPath}

	// Setup mock selector
	pickupDropoff := pickupdropoffcache.NewValue(pickup, dropoff)
	mockSelector.On("GetPickupDropoffPointsAndDurations", request, offer).Return(pickupDropoff, nil)

	// Setup mock generator
	mockGenerator.On("GeneratePaths", offer.Path(), pickup, dropoff).Return(candidatePaths, nil)

	// Setup mock validator - first path invalid, second path valid
	mockValidator.On("ValidatePath", offerNode, requestNode, invalidPath).Return(false, nil)
	mockValidator.On("ValidatePath", offerNode, requestNode, validPath).Return(true, nil)

	// Create planner and run test
	planner := planner.NewDefaultPathPlanner(mockGenerator, mockValidator, mockSelector)
	resultPath, found, err := planner.FindFirstFeasiblePath(offerNode, requestNode)

	// Verify results - valid path found
	assert.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, validPath, resultPath)

	// Verify mocks were called correctly
	mockGenerator.AssertExpectations(t)
	mockValidator.AssertExpectations(t)
	mockSelector.AssertExpectations(t)
}
