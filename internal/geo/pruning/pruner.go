package pruning

import (
	"matching-engine/internal/model"
	"time"
)

type RoutePruner interface {
	Prune(threshold time.Duration) (model.LineString, error)
}

type RoutePrunerFactory interface {
	NewRoutePruner(route model.LineString) (RoutePruner, error)
}
