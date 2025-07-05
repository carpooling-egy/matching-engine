package checker

import (
	"context"
	"fmt"
	"matching-engine/cmd/appmetrics"
	"matching-engine/internal/adapter/routing"
	"matching-engine/internal/model"
	"matching-engine/internal/service/pickupdropoffservice"
	"time"

	"github.com/rs/zerolog/log"
)

type DetourTimeChecker struct {
	selector pickupdropoffservice.PickupDropoffSelectorInterface
	engine   routing.Engine
}

func NewDetourTimeChecker(selector pickupdropoffservice.PickupDropoffSelectorInterface, engine routing.Engine) Checker {
	return &DetourTimeChecker{
		selector: selector,
		engine:   engine,
	}
}

// Check checks if the detour time is within the acceptable range and if the offer can accommodate the request
func (dtc *DetourTimeChecker) Check(offer *model.Offer, request *model.Request) (bool, error) {
	start := time.Now()
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

	dropoffDuration := durations[2]

	if offer.DepartureTime().Add(dropoffDuration).After(request.LatestArrivalTime().Add(-value.Dropoff().WalkingDuration())) {
		log.Debug().
			Str("offer_id", offer.ID()).
			Str("request_id", request.ID()).
			Msg("offer arrival time at dropoff after request latest arrival time with dropoff walking duration")
		return false, nil
	}
	// Check if the detour time is within the acceptable range
	totalTripDuration := durations[3]
	if totalTripDuration > offer.MaxEstimatedArrivalTime().Sub(offer.DepartureTime()) {
		log.Debug().
			Str("offer_id", offer.ID()).
			Str("request_id", request.ID()).
			Msg("total trip duration exceeds the maximum estimated arrival time")
		return false, nil
	}
	timeTaken := time.Since(start)
	appmetrics.TrackTime("detour_time_checker_duration", timeTaken)
	appmetrics.IncrementCount("detour_time_checker_count", 1)
	return true, nil
}
