package correcteness_test

import (
	"context"
	"github.com/rs/zerolog/log"
	"matching-engine/internal/adapter/routing"
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
