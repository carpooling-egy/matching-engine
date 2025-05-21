package di

import (
	"context"
	"go.uber.org/dig"

	"matching-engine/internal/reader"
	"matching-engine/internal/repository/postgres"
)

// registerDatabase registers the database service
func registerDatabase(c *dig.Container) {
	must(c.Provide(func() (*postgres.Database, error) {
		return postgres.NewDatabase(context.Background())
	}))
}

// This function is exported to be called from tests until a cleaner approach is implemented.

// RegisterDatabaseRepositoriesAndServices registers repositories and readers that depend on the database
func RegisterDatabaseRepositoriesAndServices(c *dig.Container) {
	must(c.Provide(postgres.NewPostgresDriverOfferRepository))
	must(c.Provide(postgres.NewPostgresRiderRequestRepo))
	must(c.Provide(reader.NewPostgresInputReader))
}
