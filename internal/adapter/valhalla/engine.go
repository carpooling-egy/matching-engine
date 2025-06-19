package valhalla

import (
	"context"
	"fmt"
	re "matching-engine/internal/adapter/routing"
	"matching-engine/internal/adapter/valhalla/client"
	"matching-engine/internal/model"
	"time"
)

type Valhalla struct {
	client *client.ValhallaClient
	mapper *Mapper
}

func NewValhalla(clientOpts ...client.Option) (re.Engine, error) {
	c, err := client.NewValhallaClient(clientOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create valhalla client: %w", err)
	}

	return &Valhalla{
		client: c,
		mapper: NewMapper(),
	}, nil
}

func (v *Valhalla) PlanDrivingRoute(
	ctx context.Context,
	routeParams *model.RouteParams,
) (*model.Route, error) {
	route, err := re.RunOperation(
		ctx,
		v.client,
		"/route",
		routeParams,
		v.mapper.RouteMapper,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to compute route: %w", err)
	}

	return route, nil
}

func (v *Valhalla) ComputeDrivingTime(
	ctx context.Context,
	routeParams *model.RouteParams,
) ([]time.Duration, error) {
	durations, err := re.RunOperation(
		ctx,
		v.client,
		"/route",
		routeParams,
		v.mapper.DrivingTimeMapper,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to compute time: %w", err)
	}

	return durations, nil
}

func (v *Valhalla) ComputeWalkingTime(
	ctx context.Context,
	walkParams *model.WalkParams,
) (time.Duration, error) {
	duration, err := re.RunOperation(
		ctx,
		v.client,
		"/route",
		walkParams,
		v.mapper.WalkingTimeMapper,
	)

	if err != nil {
		return 0, fmt.Errorf("failed to compute time: %w", err)
	}

	return duration, nil
}

func (v *Valhalla) ComputeIsochrone(
	ctx context.Context,
	req *model.IsochroneParams,
) (*model.Isochrone, error) {
	isochrone, err := re.RunOperation(
		ctx,
		v.client,
		"/isochrone",
		req,
		v.mapper.IsochroneMapper,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to compute isochrone: %w", err)
	}

	return isochrone, nil
}

func (v *Valhalla) ComputeDistanceTimeMatrix(
	ctx context.Context,
	req *model.DistanceTimeMatrixParams,
) (*model.DistanceTimeMatrix, error) {
	matrix, err := re.RunOperation(
		ctx,
		v.client,
		"/matrix",
		req,
		v.mapper.MatrixMapper,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to compute distance time matrix: %w", err)
	}

	return matrix, nil
}

func (v *Valhalla) SnapPointToRoad(
	ctx context.Context,
	point *model.Coordinate,
) (*model.Coordinate, error) {
	snappedPoint, err := re.RunOperation(
		ctx,
		v.client,
		"/trace_attributes",
		point,
		v.mapper.SnapToRoadMapper,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to snap point to road: %w", err)
	}

	return snappedPoint, nil
}
