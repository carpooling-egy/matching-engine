package repository

import (
	"context"
	"github.com/google/uuid"
	"matching-engine/internal/model"
	"time"
)

// Repository defines the interface for database operations
type Repository interface {
	// GetPendingRequests retrieves all pending requests
	GetPendingRequests(ctx context.Context) ([]*model.Request, error)

	// GetAvailableOffers retrieves all available offers
	GetAvailableOffers(ctx context.Context) ([]*model.Offer, error)

	// UpdateRequestStatus updates the status of a request
	UpdateRequestStatus(ctx context.Context, requestID uuid.UUID, status string) error

	// SaveMatch saves a match between a request and an offer
	SaveMatch(ctx context.Context, requestID uuid.UUID, offerID uuid.UUID) error
}

// MockRepository is a simple in-memory implementation of the Repository interface
// This is used for testing and development purposes
type MockRepository struct {
	PendingRequests []*model.Request
	AvailableOffers []*model.Offer
}

// NewMockRepository creates a new instance of MockRepository
func NewMockRepository() *MockRepository {
	// Create some test data
	now := time.Now()

	// Create coordinates for requests
	sourceCoord1, _ := model.NewCoordinate(37.7749, -122.4194)
	destCoord1, _ := model.NewCoordinate(37.7833, -122.4167)
	sourceCoord2, _ := model.NewCoordinate(37.7833, -122.4167)
	destCoord2, _ := model.NewCoordinate(37.7749, -122.4194)

	// Create coordinates for offers
	offerSourceCoord1, _ := model.NewCoordinate(37.7749, -122.4194)
	offerDestCoord1, _ := model.NewCoordinate(37.8000, -122.4300)
	offerSourceCoord2, _ := model.NewCoordinate(37.7900, -122.4100)
	offerDestCoord2, _ := model.NewCoordinate(37.7700, -122.4200)

	// Create request IDs
	requestID1 := uuid.New()
	requestID2 := uuid.New()
	userID1 := uuid.New()
	userID2 := uuid.New()

	// Create offer IDs
	offerID1 := uuid.New()
	offerID2 := uuid.New()
	driverID1 := uuid.New()
	driverID2 := uuid.New()

	// Create empty preferences for simplicity
	emptyPreference := model.Preference{}

	return &MockRepository{
		PendingRequests: []*model.Request{
			model.NewRequest(
				requestID1,
				userID1,
				*sourceCoord1,
				*destCoord1,
				now.Add(-5*time.Minute),
				now.Add(30*time.Minute),
				10*time.Minute,
				emptyPreference,
				1,
			),
			model.NewRequest(
				requestID2,
				userID2,
				*sourceCoord2,
				*destCoord2,
				now.Add(-3*time.Minute),
				now.Add(45*time.Minute),
				15*time.Minute,
				emptyPreference,
				2,
			),
		},
		AvailableOffers: []*model.Offer{
			model.NewOffer(
				offerID1,
				driverID1,
				*offerSourceCoord1,
				*offerDestCoord1,
				20*time.Minute,
				now.Add(5*time.Minute),
				[]model.MatchedRequest{},
				emptyPreference,
				[]*model.Point{},
			),
			model.NewOffer(
				offerID2,
				driverID2,
				*offerSourceCoord2,
				*offerDestCoord2,
				15*time.Minute,
				now.Add(10*time.Minute),
				[]model.MatchedRequest{},
				emptyPreference,
				[]*model.Point{},
			),
		},
	}
}

// GetPendingRequests retrieves all pending requests
func (r *MockRepository) GetPendingRequests(ctx context.Context) ([]*model.Request, error) {
	return r.PendingRequests, nil
}

// GetAvailableOffers retrieves all available offers
func (r *MockRepository) GetAvailableOffers(ctx context.Context) ([]*model.Offer, error) {
	return r.AvailableOffers, nil
}

// UpdateRequestStatus updates the status of a request
func (r *MockRepository) UpdateRequestStatus(ctx context.Context, requestID uuid.UUID, status string) error {
	// In a real implementation, this would update the status in a database
	// For the mock, we'll just log that the status was updated
	return nil
}

// SaveMatch saves a match between a request and an offer
func (r *MockRepository) SaveMatch(ctx context.Context, requestID uuid.UUID, offerID uuid.UUID) error {
	// In a real implementation, this would save to a database
	return nil
}
