package downsampling

import (
	"github.com/rs/zerolog/log"
	"matching-engine/internal/model"
)

type RouteDownSampler interface {
	DownSample(route model.LineString) (model.LineString, error)
}

type NoOpDownSampler struct{}

var _ RouteDownSampler = (*NoOpDownSampler)(nil)

func (NoOpDownSampler) DownSample(route model.LineString) (model.LineString, error) {
	log.Debug().Msg("NoOpDownSampler.DownSample called")
	return route, nil
}
