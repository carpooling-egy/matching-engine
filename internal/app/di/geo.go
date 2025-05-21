package di

import (
	"go.uber.org/dig"

	"matching-engine/internal/geo/downsampling"
	"matching-engine/internal/geo/processor"
	"matching-engine/internal/geo/pruning"
)

// registerGeoServices registers geo-related services
func registerGeoServices(c *dig.Container) {
	must(c.Provide(pruning.NewRTreePrunerFactory))
	must(c.Provide(downsampling.NewRDPDownSampler))
	must(c.Provide(processor.NewProcessorFactory))
}
