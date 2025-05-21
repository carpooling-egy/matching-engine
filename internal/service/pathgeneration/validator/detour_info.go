package validator

import (
	"fmt"
	"time"

	"matching-engine/internal/model"
)

// calculateDetourInfo determines if the path is within detour constraints
func (validator *DefaultPathValidator) calculateDetourInfo(
	offerNode *model.OfferNode,
	path []model.PathPoint,
	cumulativeDurations []time.Duration,
) (bool, error) {
	offer := offerNode.Offer()

	totalTripDuration := cumulativeDurations[len(cumulativeDurations)-1]
	directTripDuration, err := validator.timeMatrixService.GetTravelDuration(
		offerNode, path[0].ID(), path[len(path)-1].ID(),
	)

	if err != nil {
		return false,
			fmt.Errorf("failed to calculate direct trip duration: %w", err)
	}

	tripDetour := totalTripDuration - directTripDuration
	isWithinDetourLimit := tripDetour <= offer.DetourDurationMinutes()

	return isWithinDetourLimit,
		nil
}
