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
	timeMatrix, err := v.getTimeMatrix(routeParams.Waypoints(), routeParams.DepartureTime())
	if err != nil {
		return nil, err
	}
	cumulativeDurations, err := v.getCumulativeDurations(timeMatrix, len(routeParams.Waypoints()))
	if err != nil {
		return nil, err
	}

	return cumulativeDurations, nil
}

func (v *Valhalla) getCumulativeDurations(distanceTimeMatrix [][]time.Duration, pathLength int) ([]time.Duration, error) {
	cumulativeDurations := make([]time.Duration, pathLength)
	cumulativeDurations[0] = 0
	for i := 0; i < pathLength-1; i++ {
		duration := distanceTimeMatrix[i][i+1]
		if duration < 0 {
			return nil, fmt.Errorf("negative duration found between points %d and %d", i, i+1)
		}
		cumulativeDurations[i+1] = cumulativeDurations[i] + duration
	}
	return cumulativeDurations, nil
}

func (v *Valhalla) getTimeMatrix(matrixPoints []model.Coordinate, departureTime time.Time) ([][]time.Duration, error) {
	// Validate the input points
	if len(matrixPoints) < 2 {
		return nil, fmt.Errorf("not enough points to generate a distance/time matrix")
	}

	// Call the routing engine to get the distance and time matrix
	params, err := model.NewDistanceTimeMatrixParams(
		matrixPoints,
		model.ProfileAuto,
		model.WithDepartureTime(departureTime),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create distance time matrix params: %w", err)
	}

	distanceTimeMatrix, err := v.ComputeDistanceTimeMatrix(nil, params)
	if err != nil {
		return nil, fmt.Errorf("failed to compute distance time matrix: %w", err)
	}
	return distanceTimeMatrix.Times(), nil
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
