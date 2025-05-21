package di

import (
	"go.uber.org/dig"

	"matching-engine/internal/service/pickupdropoffservice"
	"matching-engine/internal/service/pickupdropoffservice/pickupdropoffcache"
)

// registerPickupDropoffServices registers pickup/dropoff services
func registerPickupDropoffServices(c *dig.Container) {
	must(c.Provide(pickupdropoffservice.NewWalkingTimeCalculator))
	must(c.Provide(pickupdropoffservice.NewIntersectionBasedGenerator))
	must(c.Provide(pickupdropoffcache.NewPickupDropoffCache))
	must(c.Provide(pickupdropoffservice.NewPickupDropoffSelector))
}
