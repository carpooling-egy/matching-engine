package main

import (
	"context"
	"matching-engine/internal/repository"
	"matching-engine/internal/service"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Configure zerolog
	// zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	// log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	// log.Info().Msg("Starting ride matcher service...")

	// // Create a context that can be cancelled
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	// // Set up signal handling for graceful shutdown
	// sigChan := make(chan os.Signal, 1)
	// signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Create a repository (using mock for now)
	repo := repository.NewMockRepository()

	// Create the matcher service
	matcherService := service.NewMatcherService(repo)

	// Define the interval for running the matching algorithm
	// Get interval from environment variable or use default
	intervalStr := os.Getenv("MATCHER_INTERVAL")
	interval := 15 * time.Minute // default interval
	if intervalStr != "" {
		customInterval, err := time.ParseDuration(intervalStr)
		if err != nil {
			log.Warn().Err(err).Str("intervalStr", intervalStr).Msg("Invalid interval format, using default")
		} else {
			interval = customInterval
			log.Info().Str("interval", intervalStr).Msg("Using custom interval from environment")
		}
	}

	// // Create a ticker that triggers at the specified interval
	// ticker := time.NewTicker(interval)
	// defer ticker.Stop()

	// log.Info().Msgf("Matcher will run every %s", interval)

	// Run the matching algorithm immediately on startup
	go func() {
		if err := matcherService.RunMatchingAlgorithm(ctx); err != nil {
			log.Error().Err(err).Msg("Error running matching algorithm")
		}
	}()

	// Main loop
	for {
		select {
		case <-ticker.C:
			// Run the matching algorithm at each tick
			go func() {
				if err := matcherService.RunMatchingAlgorithm(ctx); err != nil {
					log.Error().Err(err).Msg("Error running matching algorithm")
				}
			}()
		case sig := <-sigChan:
			log.Info().Msgf("Received signal: %v, shutting down...", sig)
			cancel()
			return
		case <-ctx.Done():
			log.Info().Msg("Context cancelled, shutting down...")
			return
		}
	}
}
