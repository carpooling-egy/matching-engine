package geo

import (
	"matching-engine/internal/model"
	"time"
)

type GeospatialProcessor interface {
	ComputeClosestRoutePoint(
		route *model.Route,
		point *model.Coordinate,
	) (*model.Coordinate, time.Duration, error)

	pruneRouteByDuration(
		route *model.Route,
		threshold time.Duration,
	) (*model.Polyline, error)

	downSampleRouteByDuration(
		route *model.Route,
		interval time.Duration,
	) (*model.Polyline, error)

	//intersectRouteWithIsochrone(
	//    route *model.Route,
	//    isochrone *model.Isochrone,
	//) ([]model.Polyline, error)
}
