package di

import (
	"github.com/rs/zerolog/log"
	"go.uber.org/dig"

	"matching-engine/internal/service/matcher"
)

// BuildContainer creates and configures the dependency injection container
func BuildContainer() *dig.Container {
	c := dig.New()

	// Register application services
	must(c.Provide(NewStarterService))

	// Register modules by groups
	registerAdapters(c)
	registerGeoServices(c)
	registerPickupDropoffServices(c)
	registerTimeMatrixServices(c)
	registerPathServices(c)
	registerCheckers(c)
	registerMatchingServices(c)
	registerDatabaseServices(c)

	// Register starter service
	registerStarterService(c)

	// Validate container
	if err := validateContainer(c); err != nil {
		log.Fatal().Err(err).Msg("Container validation failed")
	}

	return c
}

// validateContainer checks if the container is configured correctly
func validateContainer(c *dig.Container) error {
	return c.Invoke(func(matcher *matcher.Matcher) {
		log.Info().Msg("Container validation successful")
	})
}

// must is a helper function to handle errors during initialization
func must(err error) {
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to configure dependency injection")
	}
}

// NewStarterService is a placeholder to make the compiler happy until we build the real StarterService
// in the app package. This function will be overridden by the actual implementation.
func NewStarterService() interface{} {
	// This is just a placeholder
	return nil
}
