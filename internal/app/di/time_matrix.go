package di

import (
	"go.uber.org/dig"

	"matching-engine/internal/service/timematrix"
	"matching-engine/internal/service/timematrix/cache"
)

// This function is exported to be called from tests until a cleaner approach is implemented.

// RegisterTimeMatrixServices registers time matrix services
func RegisterTimeMatrixServices(c *dig.Container) {
	must(c.Provide(cache.NewTimeMatrixCache))
	must(c.Provide(timematrix.NewDefaultSelector))
	must(c.Provide(timematrix.NewService))
	must(c.Provide(timematrix.NewDefaultGenerator))
	must(c.Provide(timematrix.NewDefaultPopulator))
}
