package processor

import (
	"context"
	"errors"
	"matching-engine/internal/adapter/routing"
	"matching-engine/internal/geo/downsampling"
	"matching-engine/internal/geo/pruning"
	"matching-engine/internal/model"
	"time"
)

var _ GeospatialProcessor = (*processorImpl)(nil)

type processorImpl struct {
	route *model.Route
	pruning.RoutePruner
	downsampling.RouteDownSampler
	routing.Engine
}

func NewGeospatialProcessor(
	route *model.Route,
	prunerFactory pruning.RoutePrunerFactory,
	engine routing.Engine,
) (GeospatialProcessor, error) {

	if route == nil {
		return nil, errors.New("route cannot be nil")
	}
	if prunerFactory == nil {
		return nil, errors.New("pruner cannot be nil")
	}
	if engine == nil {
		return nil, errors.New("engine cannot be nil")
	}

	config := Load()

	routeCoords, err := route.Polyline().Coordinates()
	if err != nil {
		return nil, err
	}

	pruner, err := SelectPruner(routeCoords, config.EnablePruning, prunerFactory)
	if err != nil {
		return nil, err
	}

	downSampler := SelectDownsampler(config.EnableDownsampling, config.DownsamplerType)

	return &processorImpl{
		route:            route,
		RoutePruner:      pruner,
		RouteDownSampler: downSampler,
		Engine:           engine,
	}, nil
}

func (p *processorImpl) ComputeClosestRoutePoint(
	point *model.Coordinate,
	walkingTime time.Duration,
) (*model.Coordinate, time.Duration, error) {
	ctx := context.Background()

	prunedRoute, err := p.Prune(point, walkingTime)
	if err != nil {
		return nil, 0, err
	}

	downSampledRoute, err := p.DownSample(prunedRoute)
	if err != nil {
		return nil, 0, err
	}

	closestPoint, closestTime, err := p.findClosestPointOnRoute(ctx, *point, downSampledRoute)
	if err != nil {
		return nil, 0, err
	}

	return closestPoint, closestTime, nil
}

func (p *processorImpl) findClosestPointOnRoute(
	ctx context.Context,
	point model.Coordinate,
	route model.LineString,
) (*model.Coordinate, time.Duration, error) {
	timeMatrixParams, err := model.NewDistanceTimeMatrixParams(
		[]model.Coordinate{point},
		model.ProfilePedestrian,
		model.WithTargets(route),
	)
	if err != nil {
		return nil, 0, err
	}

	matrix, err := p.ComputeDistanceTimeMatrix(ctx, timeMatrixParams)
	if err != nil {
		return nil, 0, err
	}

	times := matrix.Times()[0]
	closestPointIndex, minTime := 0, times[0]
	for i, t := range times {
		if t < minTime {
			closestPointIndex, minTime = i, t
		}
	}

	return &route[closestPointIndex], minTime, nil
}
