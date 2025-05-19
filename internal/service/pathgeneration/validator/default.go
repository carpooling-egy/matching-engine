package validator

import (
	"fmt"
	"matching-engine/internal/enums"
	"matching-engine/internal/model"
	"matching-engine/internal/service/pickupdropoffservice"
	"matching-engine/internal/service/timematrix"

	"time"
)

// detourInfo contains information about the detour calculations
type detourInfo struct {
	isWithinDetourLimit      bool
	availableExtraDuration   time.Duration
	extraAccumulatedDuration time.Duration
}

type DefaultPathValidator struct {
	timeMatrixService     timematrix.Service
	pickupDropoffSelector pickupdropoffservice.PickupDropoffSelectorInterface
}

func NewDefaultPathValidator(timeMatrixService timematrix.Service, selector pickupdropoffservice.PickupDropoffSelectorInterface) *DefaultPathValidator {
	return &DefaultPathValidator{
		timeMatrixService:     timeMatrixService,
		pickupDropoffSelector: selector,
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
	detourInfo, err := validator.calculateDetourInfo(offerNode, path, cumulativeDurations)
	if err != nil {
		return false, err
	}

	if !detourInfo.isWithinDetourLimit {
		return false, nil
	}

	// Check capacity and timing constraints
	return validator.validateCapacityAndTiming(offer, path, cumulativeDurations, detourInfo)
}

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

// handlePickupPoint processes a pickup point and checks capacity and timing constraints, and updates pickup time
func (validator *DefaultPathValidator) handlePickupPoint(
	offer *model.Offer,
	point *model.PathPoint,
	cumulativeDuration time.Duration,
	currentCapacity *int,
	detourInfo *detourInfo,
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
		if waitingTime > detourInfo.availableExtraDuration {
			return false, nil
		}

		// Update detour budget and accumulated waiting time
		detourInfo.availableExtraDuration -= waitingTime
		detourInfo.extraAccumulatedDuration += waitingTime

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
