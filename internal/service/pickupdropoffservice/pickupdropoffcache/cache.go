package pickupdropoffcache

import "matching-engine/internal/collections"

type PickupDropoffCache struct {
	store *collections.SyncMap[Key, *Value]
}

func NewPickupDropoffCache() *PickupDropoffCache {
	return &PickupDropoffCache{
		store: collections.NewSyncMap[Key, *Value](),
	}
}

func (c *PickupDropoffCache) Set(key Key, value *Value) {
	c.store.Set(key, value)
}

func (c *PickupDropoffCache) Get(key Key) (*Value, bool) {
	return c.store.Get(key)
}

func (c *PickupDropoffCache) Delete(key Key) {
	c.store.Delete(key)
}
