package di

import (
	"github.com/rs/zerolog/log"
	"go.uber.org/dig"
)

// BuildContainer creates and configures the dependency injection container
func BuildContainer() *dig.Container {
	c := dig.New()

	// Register modules by groups
	// Some of the functions are // This function is exported to be called from tests until a cleaner approach is implemented.
	registerAdapters(c)
	RegisterGeoServices(c)
	RegisterPickupDropoffServices(c)
	RegisterTimeMatrixServices(c)
	RegisterPathServices(c)
	RegisterCheckers(c)
	RegisterMatchingServices(c)
	registerDatabase(c)
	RegisterDatabaseRepositoriesAndServices(c)

	// Register starter service
	RegisterStarterService(c)

	return c
}

// must is a helper function to handle errors during initialization
func must(err error) {
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to configure dependency injection")
	}
}
