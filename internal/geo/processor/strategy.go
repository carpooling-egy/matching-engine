package processor

import (
	"matching-engine/internal/enums"
	"matching-engine/internal/geo/downsampling"
	"matching-engine/internal/geo/pruning"
	"matching-engine/internal/model"
)

func SelectPruner(
	routeCoords model.LineString,
	enabled bool,
	factory pruning.RoutePrunerFactory,
) (pruning.RoutePruner, error) {
	if !enabled {
		return pruning.NewNoOpPruner(routeCoords), nil
	}
	return factory.NewRoutePruner(routeCoords)
}

func SelectDownsampler(enabled bool, typ enums.DownsamplerType) downsampling.RouteDownSampler {
	if !enabled {
		return downsampling.NoOpDownSampler{}
	}
	switch typ {
	case enums.DownsamplerTimeThreshold:
		return downsampling.NewTimeThresholdDownSampler()
	case enums.DownsamplerRDP:
		fallthrough
	default:
		return downsampling.NewRDPDownSampler()
	}
}
