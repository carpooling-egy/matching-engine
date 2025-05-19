package checker

import (
	"fmt"
	"matching-engine/internal/model"
)

type OverlapChecker struct{}

// NewOverlapChecker creates a new OverlapChecker
func NewOverlapChecker() *OverlapChecker {
	return &OverlapChecker{}
}

// Check checks if the given request can be matched with the offer
func (oc *OverlapChecker) Check(offer *model.Offer, request *model.Request) (bool, error) {
	if offer == nil || request == nil {
		return false, fmt.Errorf("offer or request is nil")
	}
	// Check if the request and offer have overlapping time slots
	if request.EarliestDepartureTime().After(offer.MaxEstimatedArrivalTime()) || request.LatestArrivalTime().Before(offer.DepartureTime()) {
		return false, nil
	}
	return true, nil
}
