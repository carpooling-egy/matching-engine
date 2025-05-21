package pickupdropoffcache

import (
	"matching-engine/internal/collections"
	"matching-engine/internal/model"
)

type PickupDropoffCache struct {
	store *collections.SyncMap[model.OfferRequestKey, *Value]
}

func NewPickupDropoffCache() *PickupDropoffCache {
	return &PickupDropoffCache{
		store: collections.NewSyncMap[model.OfferRequestKey, *Value](),
	}
}

func (c *PickupDropoffCache) Set(key model.OfferRequestKey, value *Value) {
	c.store.Set(key, value)
}

func (c *PickupDropoffCache) Get(key model.OfferRequestKey) (*Value, bool) {
	return c.store.Get(key)
}

func (c *PickupDropoffCache) Delete(key model.OfferRequestKey) {
	c.store.Delete(key)
}
