package routing_engine

import (
	"context"
	"matching-engine/internal/model"
	"time"
)

type RoutingEngine interface {
	// PlanDrivingRoute get a route between two points for a driver with a departure time
	PlanDrivingRoute(
		ctx context.Context,
		routeParams model.RouteParams,
	) (*model.Route, error)

	// ComputeWalkingDistance get the distance of a route between two points while walking
	ComputeWalkingDistance(
		ctx context.Context,
		walkParams model.WalkParams,
	) (model.Distance, error)

	// ComputeDrivingDistance get the distance of a route crossing 2/3 points
	// while driving with a departure time
	ComputeDrivingDistance(
		ctx context.Context,
		routeParams model.RouteParams,
	) (model.Distance, error)

	ComputeDrivingTime(
		ctx context.Context,
		routeParams model.RouteParams,
	) (time.Duration, error)

	// ComputeIsochrone get the walkable area (isochrone) for a given point and a distance
	ComputeIsochrone(
		ctx context.Context,
		req *model.IsochroneParams,
	) (*model.Isochrone, error)

	// ComputeDistanceTimeMatrix get the distance matrix between some points
	// while driving with a departure time
	ComputeDistanceTimeMatrix(
		ctx context.Context,
		req *model.DistanceTimeMatrixParams,
	) (*model.DistanceTimeMatrix, error)
}
