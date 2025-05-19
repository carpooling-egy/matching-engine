package validator

import (
	"fmt"
	"matching-engine/internal/model"
	"time"
)

// detourInfo contains information about the detour calculations
type detourInfo struct {
	isWithinDetourLimit      bool
	availableExtraDuration   time.Duration
	extraAccumulatedDuration time.Duration
}

// calculateDetourInfo determines if the path is within detour constraints
func (validator *DefaultPathValidator) calculateDetourInfo(
	offerNode *model.OfferNode,
	path []model.PathPoint,
	cumulativeDurations []time.Duration,
) (*detourInfo, error) {
	offer := offerNode.Offer()

	totalTripDuration := cumulativeDurations[len(cumulativeDurations)-1]
	directTripDuration, err := validator.timeMatrixService.GetTravelDuration(
		offerNode, path[0].ID(), path[len(path)-1].ID(),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to calculate direct trip duration: %w", err)
	}

	tripDetour := totalTripDuration - directTripDuration
	isWithinDetourLimit := tripDetour <= offer.DetourDurationMinutes()

	return &detourInfo{
		isWithinDetourLimit:      isWithinDetourLimit,
		availableExtraDuration:   offer.DetourDurationMinutes() - tripDetour,
		extraAccumulatedDuration: time.Duration(0),
	}, nil
}
