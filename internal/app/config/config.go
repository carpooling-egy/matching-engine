package config

import (
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
)

func LoadEnv() error {
	// Try to load from .env file, but don't fail if a file is missing
	if err := godotenv.Load("/home/samni/College/graduation_project/carpooling/matching-engine/.env"); err != nil {
		log.Warn().Msg("No .env file found, proceeding with environment variables")
	} else {
		log.Info().Msg("Loaded environment variables from .env file")
	}
	return nil
}

func GetEnvBool(key string, def bool) bool {
	val := os.Getenv(key)
	if val == "" {
		log.Debug().Msgf("%s not set, using default: %t", key, def)
		return def
	}
	parsed, err := strconv.ParseBool(val)
	if err != nil {
		log.Warn().Msgf("Invalid %s value %q, using default: %t", key, val, def)
		return def
	}
	log.Debug().Msgf("%s set to: %t", key, parsed)
	return parsed
}

func GetEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	log.Debug().Msgf("%s not set, using default %q", key, def)
	return def
}

func GetEnvFloat(key string, def float64) float64 {
	val := os.Getenv(key)
	if val == "" {
		log.Debug().Msgf("%s not set, using default: %f", key, def)
		return def
	}
	parsed, err := strconv.ParseFloat(val, 64)
	if err != nil {
		log.Warn().Msgf("Invalid %s value %q, using default: %f", key, val, def)
		return def
	}
	log.Debug().Msgf("%s set to: %f", key, parsed)
	return parsed
}
