package osrm

import (
	"context"
	"fmt"
	"matching-engine/internal/adapter/osrm/client"
	"matching-engine/internal/adapter/osrm/common"
	re "matching-engine/internal/adapter/routing"
	"matching-engine/internal/model"
	"time"
)

type OSRM struct {
	carClient  re.Client[model.OSRMTransport, map[string]any]
	footClient re.Client[model.OSRMTransport, map[string]any]
	mapper     *Mapper
}

func NewOSRM() (re.Engine, error) {
	carClient, err := client.NewOSRMClient("car")
	if err != nil {
		return nil, fmt.Errorf("failed to create OSRM car client: %w", err)
	}

	footClient, err := client.NewOSRMClient("foot")
	if err != nil {
		return nil, fmt.Errorf("failed to create OSRM foot client: %w", err)
	}

	return &OSRM{
		carClient:  carClient,
		footClient: footClient,
		mapper:     NewMapper(),
	}, nil
}

func (o *OSRM) PlanDrivingRoute(ctx context.Context, routeParams *model.RouteParams) (*model.Route, error) {
	if routeParams == nil {
		return nil, fmt.Errorf("routeParams cannot be nil")
	}

	osrmProfile, err := model.ToOSRMProfile(model.ProfileCar)
	if err != nil {
		return nil, err
	}

	osrmClient, err := o.selectClient(osrmProfile)
	if err != nil {
		return nil, err
	}

	return runOperationWithGet(
		ctx,
		osrmClient,
		"/route/v1/car",
		routeParams,
		o.mapper.RouteMapper,
	)
}

func (o *OSRM) ComputeDistanceTimeMatrix(ctx context.Context, req *model.DistanceTimeMatrixParams) (*model.DistanceTimeMatrix, error) {
	if req == nil {
		return nil, fmt.Errorf("DistanceTimeMatrixParams cannot be nil")
	}

	osrmProfile, err := model.ToOSRMProfile(req.Profile())
	if err != nil {
		return nil, err
	}

	osrmClient, err := o.selectClient(osrmProfile)
	if err != nil {
		return nil, err
	}

	return runOperationWithGet(
		ctx,
		osrmClient, "/table/v1/"+osrmProfile.String(),
		req,
		o.mapper.MatrixMapper,
	)
}

func (o *OSRM) SnapPointToRoad(ctx context.Context, point *model.Coordinate) (*model.Coordinate, error) {
	if point == nil {
		return nil, fmt.Errorf("point cannot be nil")
	}

	osrmProfile, err := model.ToOSRMProfile(model.ProfileCar)
	if err != nil {
		return nil, err
	}
	osrmClient, err := o.selectClient(osrmProfile)
	if err != nil {
		return nil, err
	}
	return runOperationWithGet(
		ctx,
		osrmClient,
		"/nearest/v1/car", // should be car here to snap to a road not a foot
		point,
		o.mapper.SnapToRoadMapper,
	)
}

func (o *OSRM) ComputeDrivingTime(ctx context.Context, routeParams *model.RouteParams) ([]time.Duration, error) {
	if routeParams == nil {
		return nil, fmt.Errorf("routeParams cannot be nil")
	}

	params, err := model.NewDistanceTimeMatrixParams(
		routeParams.Waypoints(),
		model.ProfileCar,
		model.WithDepartureTime(routeParams.DepartureTime()),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create distance time matrix params: %w", err)
	}

	matrix, err := o.ComputeDistanceTimeMatrix(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to compute distance time matrix: %w", err)
	}

	cumulativeDurations, err := common.GetCumulativeDurations(matrix.Times(), len(routeParams.Waypoints()))
	if err != nil {
		return nil, err
	}
	return cumulativeDurations, nil
}

func (o *OSRM) ComputeWalkingTime(ctx context.Context, walkParams *model.WalkParams) (time.Duration, error) {
	if walkParams == nil {
		return 0, fmt.Errorf("walkParams cannot be nil")
	}

	params, err := model.NewRouteParams(
		[]model.Coordinate{*walkParams.Origin(), *walkParams.Destination()},
		time.Now().Add(time.Minute),
	)
	if err != nil {
		return 0, fmt.Errorf("failed to create route params: %w", err)
	}

	route, err := runOperationWithGet(ctx, o.footClient, "/route/v1/foot", params, o.mapper.RouteMapper)
	if err != nil {
		return 0, err
	}

	return route.Time(), nil
}

func (o *OSRM) ComputeIsochrone(ctx context.Context, req *model.IsochroneParams) (*model.Isochrone, error) {
	panic("ComputeIsochrone is unsupported in OSRM engine")
}

func (o *OSRM) selectClient(profile model.OSRMProfile) (re.Client[model.OSRMTransport, map[string]any], error) {
	switch profile {
	case model.OSRMProfileCar:
		return o.carClient, nil
	case model.OSRMProfileFoot:
		return o.footClient, nil
	default:
		return nil, fmt.Errorf("unsupported OSRM profile: %s", profile)
	}
}

func runOperationWithGet[
	DomainReq any,
	DomainRes any,
	TransReq any,
	TransRes any,
](
	ctx context.Context,
	client re.Client[TransReq, TransRes],
	endpoint string,
	params DomainReq,
	mapper re.OperationMapper[DomainReq, DomainRes, TransReq, TransRes],
) (DomainRes, error) {
	return re.RunOperation(ctx, client, endpoint, params, mapper, re.MethodGet)
}
