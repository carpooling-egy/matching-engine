package di

import (
	"go.uber.org/dig"

	"matching-engine/internal/geo/downsampling"
	"matching-engine/internal/geo/processor"
	"matching-engine/internal/geo/pruning"
)

// This function is exported to be called from tests until a cleaner approach is implemented.

// RegisterGeoServices registers geo-related services
func RegisterGeoServices(c *dig.Container) {
	must(c.Provide(pruning.NewRTreePrunerFactory))
	must(c.Provide(downsampling.NewRDPDownSampler))
	must(c.Provide(processor.NewProcessorFactory))
}
