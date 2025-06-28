package config

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// ConfigureLogging sets up the global logger with appropriate settings
func ConfigureLogging() {
	// Set global logging level
	level := getLogLevel()
	zerolog.SetGlobalLevel(level)

	// ConfigureLogging time format
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// Set output format
	if isDevMode() {
		// Pretty console output for development
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	} else {
		// JSON for production (easier to parse by logging systems)
		log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	}
}

func ConfigureWithLogger(logger zerolog.Logger) {
	// Set global logging level
	level := getLogLevel()
	zerolog.SetGlobalLevel(level)

	// Configure time format
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// Set output format
	log.Logger = logger
}

// getLogLevel returns the appropriate log level based on environment
func getLogLevel() zerolog.Level {
	levelStr, ok := os.LookupEnv("LOG_LEVEL")
	if !ok {
		levelStr = "info"
	}
	switch levelStr {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	default:
		// Default level based on environment
		if isDevMode() {
			return zerolog.DebugLevel
		}
		return zerolog.InfoLevel
	}
}

// isDevMode returns true if the application is running in development mode
func isDevMode() bool {
	env := os.Getenv("APP_ENV")
	return env == "development" || env == "dev" || env == ""
}
