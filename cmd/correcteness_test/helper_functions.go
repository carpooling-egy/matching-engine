package correcteness_test

import (
	"context"
	"github.com/rs/zerolog/log"
	"matching-engine/internal/adapter/routing"
	"matching-engine/internal/geo/processor"
	"matching-engine/internal/model"
	"time"
)

func GetCumulativeTimes(coords []model.Coordinate, departureTime time.Time, engine routing.Engine) []time.Duration {
	// Build route parameters
	routeParams, err := model.NewRouteParams(coords, departureTime)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create route parameters")
	}

	// Compute driving time
	drivingTimes, err := engine.ComputeDrivingTime(context.Background(), routeParams)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to compute driving time")
	}
	return drivingTimes
}

func GetPickupDropoffPointsAndDurations(engine routing.Engine, offer *model.Offer, source *model.Coordinate, walkingDuration time.Duration, destination *model.Coordinate) (*model.Coordinate, time.Duration, *model.Coordinate, time.Duration) {
	factory := processor.NewProcessorFactory(
		engine,
	)
	proc, err := factory.CreateProcessor(offer)
	if err != nil {
		log.Fatal().Err(err).Msg("Error creating geospatial processor")
	}

	// Compute pickup and dropoff
	pickup, pickupDuration, err := proc.ComputeClosestRoutePoint(source, walkingDuration)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to compute pickup point")
	}
	dropoff, dropoffDuration, err := proc.ComputeClosestRoutePoint(destination, walkingDuration)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to compute dropoff point")
	}
	return pickup, pickupDuration, dropoff, dropoffDuration
}

// AdjustToNextDay adjusts the given time to the next day while preserving the time of day
func adjustToNextDay(t time.Time) time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day()+1, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
}

// ParseTime parses a time string in "15:04:05" or "15:04" format based on a base timestamp
func ParseTime(s string) time.Time {
	if s == "" {
		return time.Time{}
	}

	// Base timestamp: 1767225600 (Unix seconds)
	baseTimestamp := int64(1767225600)
	baseTime := time.Unix(baseTimestamp, 0).UTC()

	// Try "15:04:05" format first, then "15:04"
	var t time.Time
	var err error
	if len(s) == len("15:04:05") {
		t, err = time.Parse("15:04:05", s)
	} else if len(s) == len("15:04") {
		t, err = time.Parse("15:04", s)
	} else {
		return time.Time{}
	}
	if err != nil {
		return time.Time{}
	}

	// Add parsed hours, minutes, seconds to baseTime
	result := time.Date(
		baseTime.Year(), baseTime.Month(), baseTime.Day(),
		t.Hour(), t.Minute(), t.Second(), 0, time.UTC,
	)

	// If the result is before the current next day, adjust to the next day
	nextDay := adjustToNextDay(t)
	if result.Before(nextDay) {
		result = nextDay
	}
	return result
}
