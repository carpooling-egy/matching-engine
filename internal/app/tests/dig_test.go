package tests

import (
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"go.uber.org/dig"
	"matching-engine/internal/app/di"
	"matching-engine/internal/app/starter"
	"testing"
)

func TestDI(t *testing.T) {

	// This is just temp, ideally we should use a test database
	err := godotenv.Load("../../../.env")
	if err != nil {
		log.Warn().Msg(".env file not found or failed to load")
	}

	c := di.BuildContainer()
	// Validate container
	if err = validateContainer(c); err != nil {
		log.Fatal().Err(err).Msg("Container validation failed")
	}

}

// validateContainer checks if the container is configured correctly
func validateContainer(c *dig.Container) error {
	return c.Invoke(func(starterService *starter.StarterService) {
		log.Info().Msg("Container validation successful")
	})
}
