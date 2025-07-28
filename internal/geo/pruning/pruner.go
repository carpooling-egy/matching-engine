package pruning

import (
	"github.com/rs/zerolog/log"
	"matching-engine/internal/model"
	"time"
)

type RoutePruner interface {
	Prune(origin *model.Coordinate, threshold time.Duration) (model.LineString, error)
}

type RoutePrunerFactory interface {
	NewRoutePruner(route model.LineString) (RoutePruner, error)
}

type NoOpPruner struct {
	fullRoute model.LineString
}

func NewNoOpPruner(route model.LineString) RoutePruner {
	return &NoOpPruner{fullRoute: route}
}

func (n *NoOpPruner) Prune(origin *model.Coordinate, threshold time.Duration) (model.LineString, error) {
	log.Debug().Msg("NoOpPruner.Prune called")
	return n.fullRoute, nil
}
