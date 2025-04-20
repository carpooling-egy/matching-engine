package model

import (
	"time"
)

// RiderRequest represents a request from a rider for a ride
type RiderRequest struct {
	ID        string    `json:"id"`
	RiderID   string    `json:"rider_id"`
	Pickup    Location  `json:"pickup"`
	Dropoff   Location  `json:"dropoff"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Location represents a geographical location with latitude and longitude
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// RiderRequestStatus represents the possible statuses of a rider request
const (
	StatusPending   = "pending"
	StatusMatched   = "matched"
	StatusCancelled = "cancelled"
	StatusCompleted = "completed"
)
