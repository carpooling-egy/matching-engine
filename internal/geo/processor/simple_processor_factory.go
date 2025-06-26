package processor

import (
	"context"
	"fmt"
	"matching-engine/internal/adapter/routing"
	"matching-engine/internal/geo/pruning"
	"matching-engine/internal/model"
)

// Factory implements the ProcessorFactory interface.
type Factory struct {
	prunerFactory pruning.RoutePrunerFactory
	engine        routing.Engine
}

// NewProcessorFactory creates a new ProcessorFactory instance.
func NewProcessorFactory(
	prunerFactory pruning.RoutePrunerFactory,
	engine routing.Engine,
) ProcessorFactory {
	return &Factory{
		prunerFactory: prunerFactory,
		engine:        engine,
	}
}

// CreateProcessor creates a GeospatialProcessor for the given offer.
func (f *Factory) CreateProcessor(offer *model.Offer) (GeospatialProcessor, error) {
	coords := make([]model.Coordinate, len(offer.PathPoints()))
	for i, point := range offer.PathPoints() {
		coords[i] = *point.Coordinate()
	}
	routeParams, err := model.NewRouteParams(coords, offer.DepartureTime())
	if err != nil {
		return nil, fmt.Errorf("failed to create route params: %w", err)
	}
	route, err := f.engine.PlanDrivingRoute(context.Background(), routeParams)
	if err != nil {
		return nil, fmt.Errorf("failed to plan route: %w", err)
	}
	geospatialProcessor, err := NewGeospatialProcessor(route, f.prunerFactory, f.engine)
	if err != nil {
		return nil, fmt.Errorf("failed to create geospatial processor: %w", err)
	}
	return geospatialProcessor, nil
}
