package di

import (
	"go.uber.org/dig"

	"matching-engine/internal/service/timematrix"
	"matching-engine/internal/service/timematrix/cache"
)

// The fn is exported to be call them from tests, until we build a cleaner approach

// RegisterTimeMatrixServices registers time matrix services
func RegisterTimeMatrixServices(c *dig.Container) {
	must(c.Provide(cache.NewTimeMatrixCache))
	must(c.Provide(timematrix.NewDefaultSelector))
	must(c.Provide(timematrix.NewService))
	must(c.Provide(timematrix.NewDefaultGenerator))
	must(c.Provide(timematrix.NewDefaultPopulator))
}
