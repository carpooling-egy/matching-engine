package di

import (
	"go.uber.org/dig"

	"matching-engine/internal/service/timematrix"
	"matching-engine/internal/service/timematrix/cache"
)

// registerTimeMatrixServices registers time matrix services
func registerTimeMatrixServices(c *dig.Container) {
	must(c.Provide(cache.NewTimeMatrixCache))
	must(c.Provide(timematrix.NewDefaultSelector))
	must(c.Provide(timematrix.NewService))
	must(c.Provide(timematrix.NewDefaultGenerator))
	must(c.Provide(timematrix.NewDefaultPopulator))
}
