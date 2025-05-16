package prechecker

import (
	"fmt"
	"matching-engine/internal/model"
)

type CapacityChecker struct {
}

// NewCapacityChecker creates a new CapacityChecker
func NewCapacityChecker() *CapacityChecker {
	return &CapacityChecker{}
}

// Check checks if the given request can be matched with the offer
func (cc *CapacityChecker) Check(offer *model.Offer, request *model.Request) (bool, error) {
	if offer == nil || request == nil {
		return false, fmt.Errorf("offer or request is nil")
	}
	// Check if the offer has enough capacity to accommodate the request
	if offer.Capacity() < request.NumberOfRiders() {
		return false, nil
	}
	return true, nil
}
