package di

import (
	"go.uber.org/dig"

	"matching-engine/internal/service/pickupdropoffservice"
	"matching-engine/internal/service/pickupdropoffservice/pickupdropoffcache"
)

// The fn is exported to be call them from tests, until we build a cleaner approach

// RegisterPickupDropoffServices registers pickup/dropoff services
func RegisterPickupDropoffServices(c *dig.Container) {
	must(c.Provide(pickupdropoffservice.NewWalkingTimeCalculator))
	must(c.Provide(pickupdropoffservice.NewIntersectionBasedGenerator))
	must(c.Provide(pickupdropoffcache.NewPickupDropoffCache))
	must(c.Provide(pickupdropoffservice.NewPickupDropoffSelector))
}
