package validator

import (
	"time"

	"matching-engine/internal/enums"
	"matching-engine/internal/model"
)

// validateCapacityAndTiming checks if the path satisfies capacity and timing constraints
func (validator *DefaultPathValidator) validateCapacityAndTiming(
	offer *model.Offer,
	path []model.PathPoint,
	cumulativeDurations []time.Duration,
	availableExtraDetour *time.Duration,
) (bool, error) {
	currentCapacity := 0
	extraAccumulatedDuration := time.Duration(0)
	for i := range path {

		// Get a reference to the actual point in the path slice.
		// This allows us to mutate the point in-place and set the expected arrival time.
		point := &path[i]
		// Apply accumulated waiting time to all future points
		cumulativeDurations[i] += extraAccumulatedDuration

		switch point.PointType() {
		case enums.Pickup:
			valid, err := validator.handlePickupPoint(
				offer,
				point, // point.expectedArrivalTime IS BEING MODIFIED BY THE HANDLER
				cumulativeDurations[i],
				&currentCapacity, // THIS VALUE IS BEING MODIFIED BY THE HANDLER
			)
			if !valid || err != nil {
				return valid, err
			}

		case enums.Dropoff:
			valid, err := validator.handleDropoffPoint(
				offer,
				point, // point.expectedArrivalTime IS BEING MODIFIED BY THE HANDLER
				cumulativeDurations[i],
				&currentCapacity, // THIS VALUE IS BEING MODIFIED BY THE HANDLER
			)
			if !valid || err != nil {
				return valid, err
			}

		case enums.Destination:
			point.SetExpectedArrivalTime(offer.DepartureTime().Add(cumulativeDurations[i]))
		}
	}

	return true, nil
}
