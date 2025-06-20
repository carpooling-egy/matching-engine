package correcteness_test

import (
	"context"
	"github.com/rs/zerolog/log"
	"matching-engine/internal/adapter/routing"
	"matching-engine/internal/geo/downsampling"
	"matching-engine/internal/geo/processor"
	"matching-engine/internal/geo/pruning"
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
		pruning.NewRTreePrunerFactory(),
		downsampling.NewRDPDownSampler(),
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
