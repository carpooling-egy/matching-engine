package validator

import (
	"fmt"
	"matching-engine/internal/model"
	"matching-engine/internal/service/timematrix"
)

type DefaultPathValidator struct {
	timeMatrixService timematrix.Service
}

func NewDefaultPathValidator(timeMatrixService timematrix.Service) *DefaultPathValidator {
	return &DefaultPathValidator{
		timeMatrixService: timeMatrixService,
	}
}

// ValidatePath checks if the given path satisfies all constraints.
// It returns true if the path is valid, false otherwise.
// An error is returned only for system errors, not for validation failures.
//
// Note: This method modifies the provided path by setting expected arrival times.
func (validator *DefaultPathValidator) ValidatePath(
	offerNode *model.OfferNode,
	path []model.PathPoint,
) (bool, error) {
	if len(path) < 2 {
		return false, fmt.Errorf("path must contain at least two points")
	}

	offer := offerNode.Offer()

	// Get travel duration information
	cumulativeDurations, err := validator.timeMatrixService.GetCumulativeTravelDurations(offerNode, path)
	if err != nil {
		return false, fmt.Errorf("failed to calculate travel durations: %w", err)
	}

	// Check if path satisfies detour constraints
	isWithinDetourLimit, availableExtraDetour, err := validator.calculateDetourInfo(offerNode, path, cumulativeDurations)
	if err != nil {
		return false, err
	}

	if !isWithinDetourLimit {
		return false, nil
	}

	// Check capacity and timing constraints
	return validator.validateCapacityAndTiming(offer, path, cumulativeDurations, &availableExtraDetour)
}
