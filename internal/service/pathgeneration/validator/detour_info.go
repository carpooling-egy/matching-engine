package validator

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"time"

	"matching-engine/internal/model"
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

	if !isWithinDetourLimit {
		log.Debug().
			Int("tripDetourMinutes", int(tripDetour.Minutes())).
			Int("offerDetourLimitMinutes", int(offer.DetourDurationMinutes().Minutes())).
			Msg("Trip exceeds detour limit: ")
	}

	return isWithinDetourLimit,
		offer.DetourDurationMinutes() - tripDetour,
		nil
}
