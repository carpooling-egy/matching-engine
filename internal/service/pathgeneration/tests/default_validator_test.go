package tests

import (
	"errors"
	"matching-engine/internal/enums"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"matching-engine/internal/model"
	"matching-engine/internal/service/pathgeneration/validator"
)

// Mock dependencies
type MockTimeMatrixService struct {
	mock.Mock
}

func (m *MockTimeMatrixService) GetCumulativeTravelDurations(offerNode *model.OfferNode, path []model.PathPoint) ([]time.Duration, error) {
	args := m.Called(offerNode, path)
	if err := args.Error(1); err != nil {
		return nil, err
	}
	return args.Get(0).([]time.Duration), nil
}
func (m *MockTimeMatrixService) GetTravelDuration(offerNode *model.OfferNode, startID, endID model.PathPointID) (time.Duration, error) {
	args := m.Called(offerNode, startID, endID)
	return args.Get(0).(time.Duration), args.Error(1)
}

func (m *MockTimeMatrixService) GetCumulativeTravelTimes(offerNode *model.OfferNode, path []model.PathPoint) ([]time.Time, error) {
	panic("GetCumulativeTravelTimes should not be called in these tests")
}

// Helper functions
func createTestOffer(departureTime time.Time, detour time.Duration) *model.Offer {
	source := *must(model.NewCoordinate(42.43, 1.42))
	destination := *must(model.NewCoordinate(42.6, 1.7))
	return model.NewOffer(
		"offer1", "user1", source, destination, departureTime, detour, 1, *createTestPreferences(), time.Now().Add(1*time.Hour), 0, nil, nil,
	)
}

func createTestOfferNode(departureTime time.Time, detour time.Duration) *model.OfferNode {
	offer := createTestOffer(departureTime, detour)
	return model.NewOfferNode(offer)
}

// Helper function to create test data
func createTestPathPoints() []model.PathPoint {
	coord1 := must(model.NewCoordinate(42.43, 1.42))
	coord2 := must(model.NewCoordinate(42.6, 1.7))
	point1 := model.NewPathPoint(*coord1, enums.Pickup, time.Now(), nil, 0)
	point2 := model.NewPathPoint(*coord2, enums.Dropoff, time.Now(), nil, 0)
	return []model.PathPoint{*point1, *point2}
}

func createCumulativeTravelDurations() []time.Duration {
	return []time.Duration{
		0 * time.Minute,
		5 * time.Minute,
		10 * time.Minute,
		15 * time.Minute,
		20 * time.Minute,
		25 * time.Minute,
	}
}

func createPathPoint(owner model.Role, pointType enums.PointType, expectedArrivalTime time.Time, walkingDuration time.Duration) *model.PathPoint {
	coord := must(model.NewCoordinate(42.43, 1.42))
	return model.NewPathPoint(*coord, pointType, expectedArrivalTime, owner, walkingDuration)
}

func createTestRequest(latestDeparture time.Time, earliestArrival time.Time, numberOfRiders int) *model.Request {
	source := *must(model.NewCoordinate(42.43, 1.42))
	destination := *must(model.NewCoordinate(42.6, 1.7))
	return model.NewRequest(
		"request1", "user1", source, destination, latestDeparture, earliestArrival, 0*time.Minute, numberOfRiders, *createTestPreferences(),
	)
}

func createTestPreferences() *model.Preference {
	return model.NewPreference(enums.Male, false)
}

// Tests
func TestDefaultPathValidator_ValidatePath(t *testing.T) {
	mockTimeMatrix := new(MockTimeMatrixService)
	pathValidator := validator.NewDefaultPathValidator(mockTimeMatrix)

	timeNow := time.Now()

	t.Run("Valid - no walking durations set", func(t *testing.T) {

		offerNode := createTestOfferNode(timeNow, 15*time.Minute)
		sourcePoint := createPathPoint(offerNode.Offer(), enums.Source, timeNow, 0)
		destinationPoint := createPathPoint(offerNode.Offer(), enums.Destination, timeNow.Add(1*time.Hour), 0)

		r1 := createTestRequest(timeNow, timeNow.Add(25*time.Minute), 1)
		r2 := createTestRequest(timeNow, timeNow.Add(25*time.Minute), 1)
		p1 := createPathPoint(r1, enums.Pickup, timeNow, 0) //associated expected arrival time should change
		d1 := createPathPoint(r1, enums.Dropoff, timeNow.Add(25*time.Minute), 0)
		p2 := createPathPoint(r2, enums.Pickup, timeNow, 0) //associated expected arrival time should change
		d2 := createPathPoint(r2, enums.Dropoff, timeNow.Add(25*time.Minute), 0)
		path := []model.PathPoint{*sourcePoint, *p1, *d1, *p2, *d2, *destinationPoint}

		mockTimeMatrix.On("GetCumulativeTravelDurations", offerNode, path).Return(createCumulativeTravelDurations(), nil)
		mockTimeMatrix.On("GetTravelDuration", offerNode, path[0].ID(), path[5].ID()).Return(15*time.Minute, nil)

		valid, err := pathValidator.ValidatePath(offerNode, path)

		assert.NoError(t, err)
		assert.True(t, valid)
		mockTimeMatrix.AssertExpectations(t)
		assert.Equal(t, path[1].ExpectedArrivalTime(), timeNow.Add(5*time.Minute))
		assert.Equal(t, path[2].ExpectedArrivalTime(), timeNow.Add(10*time.Minute))
		assert.Equal(t, path[3].ExpectedArrivalTime(), timeNow.Add(15*time.Minute))
		assert.Equal(t, path[4].ExpectedArrivalTime(), timeNow.Add(20*time.Minute))
		assert.Equal(t, path[5].ExpectedArrivalTime(), timeNow.Add(25*time.Minute))

	})

	t.Run("Valid - walking times set", func(t *testing.T) {

		offerNode := createTestOfferNode(timeNow, 15*time.Minute)
		sourcePoint := createPathPoint(offerNode.Offer(), enums.Source, timeNow, 0)
		destinationPoint := createPathPoint(offerNode.Offer(), enums.Destination, timeNow.Add(1*time.Hour), 0)

		r1 := createTestRequest(timeNow, timeNow.Add(25*time.Minute), 1)
		r2 := createTestRequest(timeNow, timeNow.Add(25*time.Minute), 1)
		p1 := createPathPoint(r1, enums.Pickup, timeNow.Add(5*time.Minute), 5*time.Minute) //associated expected arrival time should change
		d1 := createPathPoint(r1, enums.Dropoff, timeNow.Add(25*time.Minute).Add(-15*time.Minute), 15*time.Minute)
		p2 := createPathPoint(r2, enums.Pickup, timeNow.Add(10*time.Minute), 10*time.Minute) //associated expected arrival time should change
		d2 := createPathPoint(r2, enums.Dropoff, timeNow.Add(25*time.Minute).Add(-5*time.Minute), 5*time.Minute)

		path := []model.PathPoint{*sourcePoint, *p1, *d1, *p2, *d2, *destinationPoint}

		mockTimeMatrix.On("GetCumulativeTravelDurations", offerNode, path).Return(createCumulativeTravelDurations(), nil)
		mockTimeMatrix.On("GetTravelDuration", offerNode, path[0].ID(), path[5].ID()).Return(15*time.Minute, nil)

		valid, err := pathValidator.ValidatePath(offerNode, path)

		assert.NoError(t, err)
		assert.True(t, valid)
		mockTimeMatrix.AssertExpectations(t)
	})

	t.Run("Valid - detour enough for driver to wait", func(t *testing.T) {

		offerNode := createTestOfferNode(timeNow, 15*time.Minute)
		sourcePoint := createPathPoint(offerNode.Offer(), enums.Source, timeNow, 0)
		destinationPoint := createPathPoint(offerNode.Offer(), enums.Destination, timeNow.Add(1*time.Hour), 0)

		r1 := createTestRequest(timeNow, timeNow.Add(25*time.Minute), 1)
		r2 := createTestRequest(timeNow, timeNow.Add(25*time.Minute), 1)
		p1 := createPathPoint(r1, enums.Pickup, timeNow.Add(10*time.Minute), 10*time.Minute) //associated expected arrival time should change
		d1 := createPathPoint(r1, enums.Dropoff, timeNow.Add(25*time.Minute), 0)
		p2 := createPathPoint(r2, enums.Pickup, timeNow, 0) //associated expected arrival time should change
		d2 := createPathPoint(r2, enums.Dropoff, timeNow.Add(25*time.Minute), 0)

		path := []model.PathPoint{*sourcePoint, *p1, *d1, *p2, *d2, *destinationPoint}

		mockTimeMatrix.On("GetCumulativeTravelDurations", offerNode, path).Return(createCumulativeTravelDurations(), nil)
		mockTimeMatrix.On("GetTravelDuration", offerNode, path[0].ID(), path[5].ID()).Return(15*time.Minute, nil)

		valid, err := pathValidator.ValidatePath(offerNode, path)

		assert.NoError(t, err)
		assert.True(t, valid)
		mockTimeMatrix.AssertExpectations(t)
	})

	t.Run("Invalid - not enough detour", func(t *testing.T) {

		offerNode := createTestOfferNode(timeNow, 5*time.Minute)
		sourcePoint := createPathPoint(offerNode.Offer(), enums.Source, timeNow, 0)
		destinationPoint := createPathPoint(offerNode.Offer(), enums.Destination, timeNow.Add(1*time.Hour), 0)

		r1 := createTestRequest(timeNow, timeNow.Add(25*time.Minute), 1)
		r2 := createTestRequest(timeNow, timeNow.Add(25*time.Minute), 1)
		p1 := createPathPoint(r1, enums.Pickup, timeNow, 0) //associated expected arrival time should change
		d1 := createPathPoint(r1, enums.Dropoff, timeNow.Add(25*time.Minute), 0)
		p2 := createPathPoint(r2, enums.Pickup, timeNow, 0) //associated expected arrival time should change
		d2 := createPathPoint(r2, enums.Dropoff, timeNow.Add(25*time.Minute), 0)
		path := []model.PathPoint{*sourcePoint, *p1, *d1, *p2, *d2, *destinationPoint}

		mockTimeMatrix.On("GetCumulativeTravelDurations", offerNode, path).Return(createCumulativeTravelDurations(), nil)
		mockTimeMatrix.On("GetTravelDuration", offerNode, path[0].ID(), path[5].ID()).Return(15*time.Minute, nil)

		valid, err := pathValidator.ValidatePath(offerNode, path)

		assert.False(t, valid)
		assert.Nil(t, err)
		mockTimeMatrix.AssertExpectations(t)
	})

	t.Run("Invalid - invalid times, driver detour does not accommodate waiting", func(t *testing.T) {

		offerNode := createTestOfferNode(timeNow, 10*time.Minute)
		sourcePoint := createPathPoint(offerNode.Offer(), enums.Source, timeNow, 0)
		destinationPoint := createPathPoint(offerNode.Offer(), enums.Destination, timeNow.Add(1*time.Hour), 0)

		r1 := createTestRequest(timeNow, timeNow.Add(25*time.Minute), 1)
		r2 := createTestRequest(timeNow, timeNow.Add(25*time.Minute), 1)
		p1 := createPathPoint(r1, enums.Pickup, timeNow.Add(10*time.Minute), 10*time.Minute) //associated expected arrival time should change
		d1 := createPathPoint(r1, enums.Dropoff, timeNow.Add(25*time.Minute), 0)
		p2 := createPathPoint(r2, enums.Pickup, timeNow, 0) //associated expected arrival time should change
		d2 := createPathPoint(r2, enums.Dropoff, timeNow.Add(25*time.Minute), 0)

		path := []model.PathPoint{*sourcePoint, *p1, *d1, *p2, *d2, *destinationPoint}

		mockTimeMatrix.On("GetCumulativeTravelDurations", offerNode, path).Return(createCumulativeTravelDurations(), nil)
		mockTimeMatrix.On("GetTravelDuration", offerNode, path[0].ID(), path[5].ID()).Return(15*time.Minute, nil)

		valid, err := pathValidator.ValidatePath(offerNode, path)

		assert.NoError(t, err)
		assert.False(t, valid)
		mockTimeMatrix.AssertExpectations(t)
	})

	t.Run("Invalid - invalid capacity ", func(t *testing.T) {

		offerNode := createTestOfferNode(timeNow, 15*time.Minute)
		sourcePoint := createPathPoint(offerNode.Offer(), enums.Source, timeNow, 0)
		destinationPoint := createPathPoint(offerNode.Offer(), enums.Destination, timeNow.Add(1*time.Hour), 0)

		r1 := createTestRequest(timeNow, timeNow.Add(25*time.Minute), 1)
		r2 := createTestRequest(timeNow, timeNow.Add(25*time.Minute), 1)
		p1 := createPathPoint(r1, enums.Pickup, timeNow, 0) //associated expected arrival time should change
		p2 := createPathPoint(r2, enums.Pickup, timeNow, 0)
		d1 := createPathPoint(r1, enums.Dropoff, timeNow.Add(25*time.Minute), 0) //associated expected arrival time should change
		d2 := createPathPoint(r2, enums.Dropoff, timeNow.Add(25*time.Minute), 0)

		path := []model.PathPoint{*sourcePoint, *p1, *p2, *d1, *d2, *destinationPoint}

		mockTimeMatrix.On("GetCumulativeTravelDurations", offerNode, path).Return(createCumulativeTravelDurations(), nil)
		mockTimeMatrix.On("GetTravelDuration", offerNode, path[0].ID(), path[5].ID()).Return(15*time.Minute, nil)

		valid, err := pathValidator.ValidatePath(offerNode, path)

		assert.NoError(t, err)
		assert.False(t, valid)
		mockTimeMatrix.AssertExpectations(t)
	})

	t.Run("Error - System error from time matrix service, GetCumulativeTravelDurations", func(t *testing.T) {

		offerNode := createTestOfferNode(timeNow, 5*time.Minute)
		sourcePoint := createPathPoint(offerNode.Offer(), enums.Source, timeNow, 0)
		destinationPoint := createPathPoint(offerNode.Offer(), enums.Destination, timeNow.Add(1*time.Hour), 0)

		r1 := createTestRequest(timeNow, timeNow.Add(25*time.Minute), 1)
		r2 := createTestRequest(timeNow, timeNow.Add(25*time.Minute), 1)
		p1 := createPathPoint(r1, enums.Pickup, timeNow, 0) //associated expected arrival time should change
		d1 := createPathPoint(r1, enums.Dropoff, timeNow.Add(25*time.Minute), 0)
		p2 := createPathPoint(r2, enums.Pickup, timeNow, 0) //associated expected arrival time should change
		d2 := createPathPoint(r2, enums.Dropoff, timeNow.Add(25*time.Minute), 0)
		path := []model.PathPoint{*sourcePoint, *p1, *d1, *p2, *d2, *destinationPoint}

		mockTimeMatrix.On("GetCumulativeTravelDurations", offerNode, path).Return(nil, errors.New("service error"))

		valid, err := pathValidator.ValidatePath(offerNode, path)

		assert.Error(t, err)
		assert.False(t, valid)
		mockTimeMatrix.AssertExpectations(t)
	})
	t.Run("Error - System error from time matrix service, GetTravelDuration", func(t *testing.T) {

		offerNode := createTestOfferNode(timeNow, 5*time.Minute)
		sourcePoint := createPathPoint(offerNode.Offer(), enums.Source, timeNow, 0)
		destinationPoint := createPathPoint(offerNode.Offer(), enums.Destination, timeNow.Add(1*time.Hour), 0)

		r1 := createTestRequest(timeNow, timeNow.Add(25*time.Minute), 1)
		r2 := createTestRequest(timeNow, timeNow.Add(25*time.Minute), 1)
		p1 := createPathPoint(r1, enums.Pickup, timeNow, 0) //associated expected arrival time should change
		d1 := createPathPoint(r1, enums.Dropoff, timeNow, 0)
		p2 := createPathPoint(r2, enums.Pickup, timeNow, 0) //associated expected arrival time should change
		d2 := createPathPoint(r2, enums.Dropoff, timeNow, 0)
		path := []model.PathPoint{*sourcePoint, *p1, *d1, *p2, *d2, *destinationPoint}

		mockTimeMatrix.On("GetCumulativeTravelDurations", offerNode, path).Return(createCumulativeTravelDurations(), nil)
		mockTimeMatrix.On("GetTravelDuration", offerNode, path[0].ID(), path[5].ID()).Return(time.Duration(0), errors.New("service error"))

		valid, err := pathValidator.ValidatePath(offerNode, path)

		assert.Error(t, err)
		assert.False(t, valid)
		mockTimeMatrix.AssertExpectations(t)
	})

	t.Run("Error - offer passed in path points instead of request", func(t *testing.T) {

		offerNode := createTestOfferNode(timeNow, 15*time.Minute)
		sourcePoint := createPathPoint(offerNode.Offer(), enums.Source, timeNow, 0)
		destinationPoint := createPathPoint(offerNode.Offer(), enums.Destination, timeNow.Add(1*time.Hour), 0)

		r1 := createTestRequest(timeNow, timeNow.Add(25*time.Minute), 1)
		r2 := createTestRequest(timeNow, timeNow.Add(25*time.Minute), 1)
		p1 := createPathPoint(offerNode.Offer(), enums.Pickup, timeNow, 0) //associated expected arrival time should change
		d1 := createPathPoint(r1, enums.Dropoff, timeNow, 0)
		p2 := createPathPoint(r2, enums.Pickup, timeNow, 0) //associated expected arrival time should change
		d2 := createPathPoint(r2, enums.Dropoff, timeNow, 0)
		path := []model.PathPoint{*sourcePoint, *p1, *d1, *p2, *d2, *destinationPoint}

		mockTimeMatrix.On("GetCumulativeTravelDurations", offerNode, path).Return(createCumulativeTravelDurations(), nil)
		mockTimeMatrix.On("GetTravelDuration", offerNode, path[0].ID(), path[5].ID()).Return(15*time.Minute, nil)

		valid, err := pathValidator.ValidatePath(offerNode, path)

		assert.Error(t, err)
		assert.False(t, valid)
		mockTimeMatrix.AssertExpectations(t)
	})

}
