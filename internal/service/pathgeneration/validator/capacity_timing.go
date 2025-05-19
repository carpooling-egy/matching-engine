package validator

import (
	"fmt"
	"matching-engine/internal/enums"
	"matching-engine/internal/model"
	"time"
)

// validateCapacityAndTiming checks if the path satisfies capacity and timing constraints
func (validator *DefaultPathValidator) validateCapacityAndTiming(
	offer *model.Offer,
	path []model.PathPoint,
	cumulativeDurations []time.Duration,
	detourInfo *detourInfo,
) (bool, error) {
	currentCapacity := 0

	for i := range path {

		point := &path[i]
		// Apply accumulated waiting time to all future points
		cumulativeDurations[i] += detourInfo.extraAccumulatedDuration

		switch point.PointType() {
		case enums.Pickup:
			valid, err := validator.handlePickupPoint(
				offer,
				point,
				cumulativeDurations[i],
				&currentCapacity,
				detourInfo,
			)
			if !valid || err != nil {
				return valid, err
			}

		case enums.Dropoff:
			valid, err := validator.handleDropoffPoint(
				offer,
				point,
				cumulativeDurations[i],
				&currentCapacity,
			)
			if !valid || err != nil {
				return valid, err
			}
		}
	}

	return true, nil
}
