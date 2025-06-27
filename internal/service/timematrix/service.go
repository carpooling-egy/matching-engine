package timematrix

import (
	"matching-engine/internal/model"
	"time"
)

// Service provides travel-duration and cumulative-time lookups.
type Service interface {
	GetCumulativeTravelDurations(offer *model.OfferNode, requestNode *model.RequestNode, pathPoints []model.PathPoint) ([]time.Duration, error)
	GetCumulativeTravelTimes(offer *model.OfferNode, requestNode *model.RequestNode, pathPoints []model.PathPoint) ([]time.Time, error)
	GetTravelDuration(offer *model.OfferNode, requestNode *model.RequestNode, from, to model.PathPointID) (time.Duration, error)
}
