package di

import (
	"github.com/rs/zerolog/log"
	"go.uber.org/dig"
)

// BuildContainer creates and configures the dependency injection container
func BuildContainer() *dig.Container {
	c := dig.New()

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

	return c
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
