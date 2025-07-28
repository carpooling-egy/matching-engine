package cache

import (
	"matching-engine/internal/collections"
)

// cacheKey is the composite map key combining OfferID and RequestID.
type cacheKey struct {
	OfferID   string
	RequestID string
}

// TimeMatrixCacheWithOfferIdAndRequestId provides a thread-safe cache for PathPointMappedTimeMatrix
// entries keyed by both OfferID and RequestID.
type TimeMatrixCacheWithOfferIdAndRequestId struct {
	cache *collections.SyncMap[cacheKey, *PathPointMappedTimeMatrix]
}

// NewTimeMatrixCacheWithOfferIdAndRequestId initializes a new TimeMatrixCacheWithOfferIdAndRequestId.
func NewTimeMatrixCacheWithOfferIdAndRequestId() *TimeMatrixCacheWithOfferIdAndRequestId {
	return &TimeMatrixCacheWithOfferIdAndRequestId{
		cache: collections.NewSyncMap[cacheKey, *PathPointMappedTimeMatrix](),
	}
}

// Get returns the cached matrix for the given offerID and requestID.
func (c *TimeMatrixCacheWithOfferIdAndRequestId) Get(offerID, requestID string) (*PathPointMappedTimeMatrix, bool) {
	key := cacheKey{OfferID: offerID, RequestID: requestID}
	return c.cache.Get(key)
}

// Set stores the matrix in the cache under the composite key.
func (c *TimeMatrixCacheWithOfferIdAndRequestId) Set(offerID, requestID string, matrix *PathPointMappedTimeMatrix) {
	key := cacheKey{OfferID: offerID, RequestID: requestID}
	c.cache.Set(key, matrix)
}

// Delete removes the entry for the given offerID and requestID.
func (c *TimeMatrixCacheWithOfferIdAndRequestId) Delete(offerID, requestID string) {
	key := cacheKey{OfferID: offerID, RequestID: requestID}
	c.cache.Delete(key)
}

// Clear evicts all entries from the cache.
func (c *TimeMatrixCacheWithOfferIdAndRequestId) Clear() {
	c.cache.Clear()
}
