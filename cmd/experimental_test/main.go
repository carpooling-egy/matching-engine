package main

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"matching-engine/cmd/appmetrics"
	"matching-engine/internal/app"
	"matching-engine/internal/app/config"
	"matching-engine/internal/app/shutdown"
	"os"
	"path/filepath"
	"time"
)

func applicationLogic(logger zerolog.Logger) {
	if err := config.LoadEnv(); err != nil {
		log.Fatal().Err(err).Msg("Failed to load environment variables")
	}

	config.ConfigureWithLogger(logger)
	log.Info().Msg("Starting ride matcher service...")

	ctx, cancel := context.WithCancel(context.Background())
	shutdown.Setup(cancel)

	newApp := app.NewApp()
	if err := newApp.Run(ctx); err != nil {
		log.Fatal().Err(err).Msg("Matching engine failed to run")
	}

	log.Info().Msg("Matching Engine shutting down...")
}

func processDataset(datasetID string, runs int, logger zerolog.Logger) {
	var totalDuration time.Duration

	for run := 1; run <= runs; run++ {
		logger.Info().
			Int("run", run).
			Str("dataset", datasetID).
			Msg("=========================================================================")
		start := time.Now()
		logger.Info().
			Int("run", run).
			Str("dataset", datasetID).
			Msg("Started processing dataset")
		applicationLogic(logger)
		logger.Info().
			Int("run", run).
			Str("dataset", datasetID).
			Msg("========================================================================")
		duration := time.Since(start)
		logger.Info().
			Int("run", run).
			Str("dataset", datasetID).
			Str("duration", duration.String()).
			Msg("Finished processing dataset")
		totalDuration += duration
	}

	logger.Info().
		Str("dataset", datasetID).
		Int("runs", runs).
		Msg("===============================Average Metrics Calculation===========================")

	avgDuration := totalDuration / time.Duration(runs)
	logger.Info().
		Str("dataset", datasetID).
		Int("runs", runs).
		Str("average_duration", avgDuration.String()).
		Msg("Average processing time for dataset")

	// Log average counts and timings
	logger.Info().
		Str("dataset", datasetID).
		Int("runs", runs).
		Msg("================================= Average Count Metrics =================================")
	counts := appmetrics.GetAllCounts()
	for name, count := range counts {
		averageCount := count / float64(runs)
		logger.Info().
			Str("metric", name).
			Float64("average_count", averageCount).
			Msg("Average count for metric")
	}
	logger.Info().
		Str("dataset", datasetID).
		Int("runs", runs).
		Msg("================================= Average Time Metrics =================================")
	times := appmetrics.GetAllTimings()
	for name, duration := range times {
		if(name == "one_edge" || name == "detour_time_checker_duration") {
			continue // Skip these metrics as they are logged separately
		}
		averageDuration := duration / time.Duration(runs)
		logger.Info().
			Str("metric", name).
			Str("average_duration", averageDuration.String()).
			Msg("Average duration for metric")
	}
	one_edge_duration := appmetrics.GetTime("one_edge") / time.Duration(appmetrics.GetCount("one_edge"))
	logger.Info().
		Str("metric", "one_edge").
		Str("average_duration", one_edge_duration.String()).
		Msg("Average duration for metric")

	detour_time_checker_duration := appmetrics.GetTime("detour_time_checker_duration") / time.Duration(appmetrics.GetCount("detour_time_checker_count"))
	log.Info().
		Str("metric", "detour_time_checker_duration").
		Str("average_duration", detour_time_checker_duration.String()).
		Msg("Average duration for metric")
}

func main() {
	datasetIDs := []string{"nyc_rt_100_v"}
	runsPerDataset := 1

	logsDir := "logs"
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		fmt.Printf("Error creating logs directory: %v\n", err)
		return
	}

	for _, datasetID := range datasetIDs {
		appmetrics.ResetTimings()
		appmetrics.ResetCounts()
		os.Setenv("DATASET_ID", datasetID)
		logFileName := filepath.Join(logsDir, fmt.Sprintf("%s.log", datasetID))
		logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			fmt.Printf("Error opening log file: %v\n", err)
			continue
		}
		logger := zerolog.New(logFile).With().Timestamp().Logger()
		processDataset(datasetID, runsPerDataset, logger)
		logFile.Close()
	}
}
