package prechecker

import (
	"context"
	"fmt"
	"matching-engine/internal/adapter/routing"
	"matching-engine/internal/model"
	"matching-engine/internal/service/pickupdropoffservice"
)

type DetourTimeChecker struct {
	selector pickupdropoffservice.PickupDropoffSelectorInterface
	engine   routing.Engine
}

func NewDetourTimeChecker(selector pickupdropoffservice.PickupDropoffSelectorInterface, engine routing.Engine) *DetourTimeChecker {
	return &DetourTimeChecker{
		selector: selector,
		engine:   engine,
	}
}

// Check checks if the detour time is within the acceptable range and if the offer can accommodate the request
func (dtc *DetourTimeChecker) Check(offer *model.Offer, request *model.Request) (bool, error) {

	value, err := dtc.selector.GetPickupDropoffPointsAndDurations(request, offer)
	if err != nil {
		return false, fmt.Errorf("failed to get pickup and dropoff points: %w", err)
	}

	waypoints := []model.Coordinate{*offer.Source(), *value.Pickup().Coordinate(), *value.Dropoff().Coordinate(), *offer.Destination()}
	params, err := model.NewRouteParams(waypoints, offer.DepartureTime())
	if err != nil {
		return false, fmt.Errorf("failed to create route params: %w", err)
	}

	durations, err := dtc.engine.ComputeDrivingTime(context.Background(), params)
	if err != nil {
		return false, fmt.Errorf("failed to compute durations between points: %w", err)
	}

	pickupDuration := durations[1]
	dropoffDuration := durations[2]

	if offer.DepartureTime().Add(pickupDuration).Before(request.EarliestDepartureTime().Add(value.PickupWalkingDuration())) {
		return false, nil
	}

	if offer.DepartureTime().Add(dropoffDuration).After(request.LatestArrivalTime().Add(-value.DropoffWalkingDuration())) {
		return false, nil
	}
	// Check if the detour time is within the acceptable range
	totalTripDuration := durations[3]
	directWaypointParams, err := model.NewRouteParams([]model.Coordinate{*offer.Source(), *offer.Destination()}, offer.DepartureTime())
	if err != nil {
		return false, fmt.Errorf("failed to create direct route params: %w", err)
	}

	directDurations, err := dtc.engine.ComputeDrivingTime(context.Background(), directWaypointParams)
	if err != nil {
		return false, fmt.Errorf("failed to compute direct durations: %w", err)
	}
	// Check if the detour time is within the acceptable range
	if totalTripDuration > directDurations[1]+offer.DetourDurationMinutes() {
		return false, nil
	}

	return true, nil
}
