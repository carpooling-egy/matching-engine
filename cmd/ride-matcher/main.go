package main

import (
	"context"
	"github.com/rs/zerolog/log"
	"matching-engine/internal/app"
	"matching-engine/internal/app/config"
	"matching-engine/internal/app/shutdown"
	"os"
	"runtime/trace"
)

func main() {
	// Configure logging
	config.ConfigureLogging()
	log.Info().Msg("Starting ride matcher service...")

	// TODO: setup an actual graceful shutdown, this is a placeholder
	ctx, cancel := context.WithCancel(context.Background())
	shutdown.Setup(cancel)

	// Load environment variables
	if err := config.LoadEnv(); err != nil {
		log.Fatal().Err(err).Msg("Failed to load environment variables")
	}

	// Initialize environment variables
	if err := config.LoadEnv(); err != nil {
		log.Error().Err(err).Msg("Failed to load environment variables")
		os.Exit(1)
	}

	f, _ := os.Create("trace.out")
	defer f.Close()
	trace.Start(f)
	defer trace.Stop()

	// Create and run the application
	newApp := app.NewApp()
	if err := newApp.Run(ctx); err != nil {
		log.Fatal().Err(err).Msg("Matching engine failed to run")
	}

	log.Info().Msg("Matching Engine shutting down...")
	os.Exit(0)
}
