package main

import (
	"context"
	"github.com/rs/zerolog/log"
	"matching-engine/internal/app"
	"matching-engine/internal/app/config"
	"matching-engine/internal/app/shutdown"
	"os"
)

func main() {
	// Configure logging
	config.ConfigureLogging()
	log.Info().Msg("Starting ride matcher service...")

	// TODO: setup an actual graceful shutdown, this is a placeholder
	ctx, cancel := context.WithCancel(context.Background())
	shutdown.Setup(cancel)

	// Create and run the application
	newApp := app.NewApp()
	if err := newApp.Run(ctx); err != nil {
		log.Fatal().Err(err).Msg("Matching engine failed to run")
	}

	log.Info().Msg("Matching Engine shutting down...")
	os.Exit(0)
}
