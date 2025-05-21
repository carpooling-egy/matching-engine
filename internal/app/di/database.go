package di

import (
	"context"
	"go.uber.org/dig"

	"matching-engine/internal/reader"
	"matching-engine/internal/repository/postgres"
)

// registerDatabaseServices registers database-related services
func registerDatabaseServices(c *dig.Container) {
	must(c.Provide(func() (*postgres.Database, error) {
		return postgres.NewDatabase(context.Background())
	}))
	must(c.Provide(postgres.NewPostgresDriverOfferRepository))
	must(c.Provide(postgres.NewPostgresRiderRequestRepo))
	must(c.Provide(reader.NewPostgresInputReader))
}
