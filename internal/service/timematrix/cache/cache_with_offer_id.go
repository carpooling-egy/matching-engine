package cache

import (
	"matching-engine/internal/collections"
)

// TimeMatrixCacheWithOfferId is a specialized cache for route data
type TimeMatrixCacheWithOfferId struct {
	cache *collections.SyncMap[string, *PathPointMappedTimeMatrix]
}

// NewTimeMatrixCacheWithOfferId creates a new TimeMatrixCacheWithOfferIdAndRequestId instance
func NewTimeMatrixCacheWithOfferId() *TimeMatrixCacheWithOfferId {
	return &TimeMatrixCacheWithOfferId{
		cache: collections.NewSyncMap[string, *PathPointMappedTimeMatrix](),
	}
}

// Get retrieves a matrix value from the cache by key
func (c *TimeMatrixCacheWithOfferId) Get(offerID string) (*PathPointMappedTimeMatrix, bool) {
	return c.cache.Get(offerID)
}

// Set stores a matrix value in the cache with the given key
func (c *TimeMatrixCacheWithOfferId) Set(offerID string, pointMappedMatrix *PathPointMappedTimeMatrix) {
	c.cache.Set(offerID, pointMappedMatrix)
}

// Delete removes a matrix value from the cache by key
func (c *TimeMatrixCacheWithOfferId) Delete(offerID string) {
	c.cache.Delete(offerID)
}

// Clear removes all items from the cache
func (c *TimeMatrixCacheWithOfferId) Clear() {
	c.cache.Clear()
}
