package valhalla

import (
	"context"
	"fmt"
	re "matching-engine/internal/adapter/routing-engine"
	"matching-engine/internal/adapter/valhalla/client"
	"matching-engine/internal/model"
	"time"
)

const PORT = 8002

var BaseURL = fmt.Sprintf("http://localhost:%d", PORT)

type Valhalla struct {
	client *client.ValhallaClient
	mapper *Mapper
}

func NewValhalla() (*Valhalla, error) {
	c, err := client.NewValhallaClient(BaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create valhalla client: %w", err)
	}

	return &Valhalla{
		client: c,
		mapper: NewMapper(),
	}, nil
}

var _ re.RoutingEngine = (*Valhalla)(nil)

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

func (v *Valhalla) ComputeDrivingDistance(
	ctx context.Context,
	routeParams *model.RouteParams,
) (*model.Distance, error) {
	distance, err := re.RunOperation(
		ctx,
		v.client,
		"/route",
		routeParams,
		v.mapper.DrivingDistanceMapper,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to compute distance: %w", err)
	}

	return distance, nil
}

func (v *Valhalla) ComputeDrivingTime(
	ctx context.Context,
	routeParams *model.RouteParams,
) (time.Duration, error) {
	duration, err := re.RunOperation(
		ctx,
		v.client,
		"/route",
		routeParams,
		v.mapper.DrivingTimeMapper,
	)

	if err != nil {
		return 0, fmt.Errorf("failed to compute time: %w", err)
	}

	return duration, nil
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
