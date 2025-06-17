package di

import (
	"context"
	"go.uber.org/dig"
	"matching-engine/internal/app/di/utils"

	"matching-engine/internal/reader"
	"matching-engine/internal/repository/postgres"
)

// registerDatabase registers the database service
func registerDatabase(c *dig.Container) {
	utils.Must(c.Provide(func() (*postgres.Database, error) {
		return postgres.NewDatabase(context.Background())
	}))
}

// This function is exported to be called from tests until a cleaner approach is implemented.

// RegisterDatabaseRepositoriesAndServices registers repositories and readers that depend on the database
func RegisterDatabaseRepositoriesAndServices(c *dig.Container) {
	utils.Must(c.Provide(postgres.NewPostgresDriverOfferRepository))
	utils.Must(c.Provide(postgres.NewPostgresRiderRequestRepo))
	utils.Must(c.Provide(reader.NewPostgresInputReader))
}
