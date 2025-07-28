package tests

import (
	"github.com/stretchr/testify/mock"
	"matching-engine/internal/enums"
	"matching-engine/internal/model"
	"matching-engine/internal/service/pickupdropoffservice/pickupdropoffcache"
	"time"
)

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

// Default helper functions
func createDefaultCoordinate() *model.Coordinate {
	coord, _ := model.NewCoordinate(0.0, 0.0)
	return coord
}

func createDefaultOffer() *model.Offer {
	return model.NewOffer(
		"defaultOffer", "defaultDriver",
		*createDefaultCoordinate(), *createDefaultCoordinate(),
		time.Now(), time.Hour, 4,
		model.Preference{}, time.Now().Add(2*time.Hour),
		0, nil, nil,
	)
}

func createDefaultOfferWithTime(t time.Time) *model.Offer {
	return model.NewOffer(
		"defaultOffer", "defaultDriver",
		*createDefaultCoordinate(), *createDefaultCoordinate(),
		t, time.Hour, 4,
		model.Preference{}, t.Add(2*time.Hour),
		0, nil, nil,
	)
}

func createDefaultOfferWithTimeAndCapacity(t time.Time, capacity int) *model.Offer {
	return model.NewOffer(
		"defaultOffer", "defaultDriver",
		*createDefaultCoordinate(), *createDefaultCoordinate(),
		t, time.Hour, capacity,
		model.Preference{}, t.Add(2*time.Hour),
		0, nil, nil,
	)
}

func createDefaultRequest() *model.Request {
	return model.NewRequest(
		"defaultRequest", "defaultRider",
		*createDefaultCoordinate(), *createDefaultCoordinate(),
		time.Now(), time.Now().Add(time.Hour),
		10*time.Minute, 1, model.Preference{},
	)
}

func createDefaultRequestWithEarliestDepartureTime(t time.Time) *model.Request {
	return model.NewRequest(
		"defaultRequest", "defaultRider",
		*createDefaultCoordinate(), *createDefaultCoordinate(),
		t, time.Now().Add(time.Hour),
		10*time.Minute, 1, model.Preference{},
	)
}

func createDefaultPickupDropoff(request *model.Request) (*model.PathPoint, *model.PathPoint) {
	pickup := model.NewPathPoint(
		*createDefaultCoordinate(),
		enums.Pickup, time.Now(), request, 5*time.Minute,
	)
	dropoff := model.NewPathPoint(
		*createDefaultCoordinate(),
		enums.Dropoff, time.Now().Add(30*time.Minute), request, 5*time.Minute,
	)
	return pickup, dropoff
}
