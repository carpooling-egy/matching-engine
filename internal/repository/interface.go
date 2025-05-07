package repository

import (
	"context"
	"time"

	"matching-engine/internal/model"
)

// RiderRepository defines operations for rider request persistence
type RiderRequestRepo interface {
	// GetByID fetches a rider request by its ID
	GetByID(ctx context.Context, id string) (*model.Request, error)

	// GetPendingRequests retrieves all rider requests that haven't been matched with a driver
	FindUnmatched(ctx context.Context, start, end time.Time) ([]*model.Request, error)
}

// DriverRepository defines operations for driver offer persistence
type DriverOfferRepo interface {
	// GetByID fetches a driver offer by its ID
	GetByID(ctx context.Context, id string) (*model.Offer, error)

	// GetAvailableDrivers gets drivers with capacity and matching time windows
	GetAvailable(ctx context.Context, start, end time.Time) ([]*model.Offer, error)
}
