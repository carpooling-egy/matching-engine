package downsampling

import (
	"matching-engine/internal/model"
)

type RouteDownSampler interface {
	DownSample(route model.LineString) (model.LineString, error)
}

type NoOpDownSampler struct{}

var _ RouteDownSampler = (*NoOpDownSampler)(nil)

func (NoOpDownSampler) DownSample(route model.LineString) (model.LineString, error) {
	return route, nil
}
