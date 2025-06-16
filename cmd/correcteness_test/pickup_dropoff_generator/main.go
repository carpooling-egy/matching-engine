package main

import (
	"github.com/rs/zerolog/log"
	"matching-engine/internal/adapter/valhalla"
	"matching-engine/internal/geo/downsampling"
	"matching-engine/internal/geo/processor"
	"matching-engine/internal/geo/pruning"
	"matching-engine/internal/model"
	"time"
)

func main() {
	// Need to be read from a file
	offer := &model.Offer{}
	// Need to be read from a file
	source, destination := &model.Coordinate{}, &model.Coordinate{}
	// Need to be read from a file
	walkingDuration := time.Duration(60)
	engine, err := valhalla.NewValhalla()
	if err != nil {
		log.Error().Err(err).Msg("Error creating new Valhalla")
		return
	}
	prunerFactory := pruning.CreateRTreePrunerFactory()
	downSampler := downsampling.NewRDPDownSampler()
	factory := processor.NewProcessorFactory(prunerFactory, downSampler, engine)
	proc, err := factory.CreateProcessor(offer)
	if err != nil {
		log.Error().Err(err).Msg("Error creating geospatial processor")
		return
	}
	pickup, pickupDuration, err := proc.ComputeClosestRoutePoint(source, walkingDuration)
	if err != nil {
		log.Error().Err(err).Msg("Error computing closest route point for pickup")
		return
	}
	dropoff, dropoffDuration, err := proc.ComputeClosestRoutePoint(destination, walkingDuration)
	if err != nil {
		log.Error().Err(err).Msg("Error computing closest route point for dropoff")
		return
	}
	// Need to be written to a file
	log.Info().Msgf("Pickup point: %v, Duration: %v", pickup, pickupDuration)
	log.Info().Msgf("Dropoff point: %v, Duration: %v", dropoff, dropoffDuration)
}
