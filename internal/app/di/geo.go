package di

import (
	"go.uber.org/dig"

	"matching-engine/internal/geo/downsampling"
	"matching-engine/internal/geo/processor"
	"matching-engine/internal/geo/pruning"
)

// The fn is exported to be call them from tests, until we build a cleaner approach

// RegisterGeoServices registers geo-related services
func RegisterGeoServices(c *dig.Container) {
	must(c.Provide(pruning.NewRTreePrunerFactory))
	must(c.Provide(downsampling.NewRDPDownSampler))
	must(c.Provide(processor.NewProcessorFactory))
}
