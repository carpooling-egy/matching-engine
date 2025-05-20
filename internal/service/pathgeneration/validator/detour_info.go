package validator

import (
	"fmt"
	"matching-engine/internal/model"
	"time"
)

// calculateDetourInfo determines if the path is within detour constraints
func (validator *DefaultPathValidator) calculateDetourInfo(
	offerNode *model.OfferNode,
	path []model.PathPoint,
	cumulativeDurations []time.Duration,
) (bool, time.Duration, error) {
	offer := offerNode.Offer()

	totalTripDuration := cumulativeDurations[len(cumulativeDurations)-1]
	directTripDuration, err := validator.timeMatrixService.GetTravelDuration(
		offerNode, path[0].ID(), path[len(path)-1].ID(),
	)

	if err != nil {
		return false,
			0,
			fmt.Errorf("failed to calculate direct trip duration: %w", err)
	}

	tripDetour := totalTripDuration - directTripDuration
	isWithinDetourLimit := tripDetour <= offer.DetourDurationMinutes()

	return isWithinDetourLimit,
		offer.DetourDurationMinutes() - tripDetour,
		nil
}
