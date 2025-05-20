package validator

import (
	"fmt"
	"time"

	"matching-engine/internal/model"
)

// handlePickupPoint processes a pickup point and checks capacity and timing constraints, and updates pickup time
func (validator *DefaultPathValidator) handlePickupPoint(
	offer *model.Offer,
	point *model.PathPoint,
	cumulativeDuration time.Duration,
	currentCapacity *int,
	availableExtraDuration *time.Duration,
	extraAccumulatedDuration *time.Duration,
) (bool, error) {
	request, ok := point.Owner().AsRequest()
	if !ok {
		return false, fmt.Errorf("PathPoint is a pickup and Owner isn't a rider")
	}

	// Check capacity constraint
	*currentCapacity += request.NumberOfRiders()
	if *currentCapacity > offer.Capacity() {
		return false, nil
	}

	// Check timing constraints
	driverArrivalTime := offer.DepartureTime().Add(cumulativeDuration)
	riderEarliestPickupTime := request.EarliestDepartureTime().Add(point.WalkingDuration()) // also equivalent to point.ExpectedArrivalTime()

	if driverArrivalTime.Before(riderEarliestPickupTime) {
		// Driver needs to wait for the rider
		waitingTime := riderEarliestPickupTime.Sub(driverArrivalTime)

		// Check if waiting is possible within detour constraints
		if waitingTime > *availableExtraDuration {
			return false, nil
		}

		// Update detour budget and accumulated waiting time
		*availableExtraDuration -= waitingTime
		*extraAccumulatedDuration += waitingTime

		// Update expected arrival time in the path point
		point.SetExpectedArrivalTime(riderEarliestPickupTime)
	} else {
		// Driver arrives after rider's earliest pickup time, no waiting needed
		point.SetExpectedArrivalTime(driverArrivalTime)
	}

	return true, nil
}

// handleDropoffPoint processes a dropoff point and checks timing constraints
func (validator *DefaultPathValidator) handleDropoffPoint(
	offer *model.Offer,
	point *model.PathPoint,
	cumulativeDuration time.Duration,
	currentCapacity *int,
) (bool, error) {
	request, ok := point.Owner().AsRequest()
	if !ok {
		return false, fmt.Errorf("PathPoint is a dropoff and Owner isn't a rider")
	}

	// Check timing constraints
	driverArrivalTime := offer.DepartureTime().Add(cumulativeDuration)
	riderLatestDropoffTime := request.LatestArrivalTime().Add(-point.WalkingDuration())

	if driverArrivalTime.After(riderLatestDropoffTime) {
		// Driver would arrive too late
		return false, nil
	}

	// Update capacity
	*currentCapacity -= request.NumberOfRiders()

	// Set expected arrival time for dropoff
	point.SetExpectedArrivalTime(driverArrivalTime)

	return true, nil
}
