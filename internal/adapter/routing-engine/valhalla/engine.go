package valhalla

import (
	"context"
	"fmt"
	re "matching-engine/internal/adapter/routing-engine"
	"matching-engine/internal/adapter/routing-engine/valhalla/client"
	"matching-engine/internal/model"
	"time"
)

const PORT = 8002

var BaseURL = fmt.Sprintf("http://localhost:%d", PORT)

type Valhalla struct {
	client *client.ValhallaClient
	mapper Mapper
}

func NewValhalla() (*Valhalla, error) {
	c, err := client.NewValhallaClient(BaseURL, fmt.Sprint(PORT))
	if err != nil {
		return nil, fmt.Errorf("failed to create valhalla client: %w", err)
	}

	return &Valhalla{
		client: c,
		//mapper: NewValhallaMapper(),
	}, nil
}

var _ re.RoutingEngine = (*Valhalla)(nil)

func (v Valhalla) PlanDrivingRoute(ctx context.Context, routeParams model.RouteParams) (*model.Route, error) {
	//TODO implement me
	panic("implement me")
}

func (v Valhalla) ComputeWalkingDistance(ctx context.Context, walkParams model.WalkParams) (model.Distance, error) {
	//TODO implement me
	panic("implement me")
}

func (v Valhalla) ComputeDrivingDistance(ctx context.Context, routeParams model.RouteParams) (model.Distance, error) {
	//TODO implement me
	panic("implement me")
}

func (v Valhalla) ComputeDrivingTime(ctx context.Context, routeParams model.RouteParams) (time.Duration, error) {
	//TODO implement me
	panic("implement me")
}

func (v Valhalla) ComputeIsochrone(ctx context.Context, req *model.IsochroneParams) (*model.Isochrone, error) {
	//TODO implement me
	panic("implement me")
}

func (v Valhalla) ComputeDistanceTimeMatrix(ctx context.Context, req *model.DistanceTimeMatrixParams) (*model.DistanceTimeMatrix, error) {
	//TODO implement me
	panic("implement me")
}
