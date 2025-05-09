package pickupdropoffservice

import "matching-engine/internal/model"

type CacheKey struct {
	OfferID   string
	RequestID string
}

type CacheValue struct {
	Pickup, Dropoff *model.PathPoint
}
