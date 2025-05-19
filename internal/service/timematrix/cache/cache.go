package cache

import (
	"matching-engine/internal/collections"
)

// TimeMatrixCache is a specialized cache for route data
type TimeMatrixCache struct {
	cache *collections.SyncMap[string, *PathPointMappedTimeMatrix]
}

// NewTimeMatrixCache creates a new TimeMatrixCache instance
func NewTimeMatrixCache() *TimeMatrixCache {
	return &TimeMatrixCache{
		cache: collections.NewSyncMap[string, *PathPointMappedTimeMatrix](),
	}
}

// Get retrieves a matrix value from the cache by key
func (c *TimeMatrixCache) Get(offerID string) (*PathPointMappedTimeMatrix, bool) {
	return c.cache.Get(offerID)
}

// Set stores a matrix value in the cache with the given key
func (c *TimeMatrixCache) Set(offerID string, pointMappedMatrix *PathPointMappedTimeMatrix) {
	c.cache.Set(offerID, pointMappedMatrix)
}

// Delete removes a matrix value from the cache by key
func (c *TimeMatrixCache) Delete(offerID string) {
	c.cache.Delete(offerID)
}

// Clear removes all items from the cache
func (c *TimeMatrixCache) Clear() {
	c.cache.Clear()
}
