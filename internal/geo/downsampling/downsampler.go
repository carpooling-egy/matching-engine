package downsampling

import (
	"matching-engine/internal/model"
)

type RouteDownSampler interface {
	DownSample(route model.LineString) (model.LineString, error)
}
