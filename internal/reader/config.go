package reader

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

type Config struct {
	datasetId string
	start     time.Time
	end       time.Time
}

func DefaultConfig() Config {
	return Config{
		datasetId: "default",
		start:     time.Now(),
		end:       time.Now().Add(24 * time.Hour),
	}
}

func LoadConfig() (Config, error) {
	cfg := DefaultConfig()

	// Simple override helper
	override := func(envVar string, target *string) {
		if val := os.Getenv(envVar); val != "" {
			*target = val
		}
	}
	// Override dataset ID from environment variable if set
	override("DATASET_ID", &cfg.datasetId)

	// Parse start and end times from environment variables if set
	if err := parseTimeEnv("START", &cfg.start); err != nil {
		return cfg, err
	}
	if err := parseTimeEnv("END", &cfg.end); err != nil {
		return cfg, err
	}

	logConfig(cfg)
	return cfg, nil

}

// parseTimeEnv reads the env variable and parses it as "2006-01-02" into target if set.
// Returns error if parsing fails, does nothing if env is not set.
func parseTimeEnv(envVar string, target *time.Time) error {
	if val := os.Getenv(envVar); val != "" {
		t, err := time.Parse("2006-01-02 15:04:05", val)
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", envVar, err)
		}
		*target = t
	}
	return nil
}
func logConfig(cfg Config) {
	log.Info().
		Str("datasetId", cfg.datasetId).
		Str("from", cfg.start.Format(time.RFC3339Nano)).
		Str("to", cfg.end.Format(time.RFC3339Nano)).
		Msg("Reader configuration loaded")
}
