package di

import (
	"go.uber.org/dig"
	"matching-engine/internal/app/di/utils"

	"matching-engine/internal/geo/downsampling"
	"matching-engine/internal/geo/processor"
	"matching-engine/internal/geo/pruning"
)

// This function is exported to be called from tests until a cleaner approach is implemented.

// RegisterGeoServices registers geo-related services
func RegisterGeoServices(c *dig.Container) {
	utils.Must(c.Provide(pruning.NewRTreePrunerFactory))
	utils.Must(c.Provide(downsampling.NewRDPDownSampler))
	utils.Must(c.Provide(processor.NewProcessorFactory))
}
