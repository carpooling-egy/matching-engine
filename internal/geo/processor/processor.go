package processor

import (
	"matching-engine/internal/model"
	"time"
)

type GeospatialProcessor interface {
	ComputeClosestRoutePoint(
		point *model.Coordinate,
		walkingTime time.Duration,
	) (*model.Coordinate, time.Duration, error)
}
