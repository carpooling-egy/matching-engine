package checker

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"matching-engine/internal/model"
)

type OverlapChecker struct{}

// NewOverlapChecker creates a new OverlapChecker
func NewOverlapChecker() Checker {
	return &OverlapChecker{}
}

// Check checks if the given request can be matched with the offer
func (oc *OverlapChecker) Check(offer *model.Offer, request *model.Request) (bool, error) {
	if offer == nil || request == nil {
		return false, fmt.Errorf("offer or request is nil")
	}
	// Check if the request and offer have overlapping time slots
	if request.EarliestDepartureTime().After(offer.MaxEstimatedArrivalTime()) || request.LatestArrivalTime().Before(offer.DepartureTime()) {
		log.Debug().
			Str("offer_id", offer.ID()).
			Str("request_id", request.ID()).
			Msg("offer and request do not overlap in time")
		return false, nil
	}
	return true, nil
}
