package di

import (
	"go.uber.org/dig"

	"matching-engine/internal/service/pickupdropoffservice"
	"matching-engine/internal/service/pickupdropoffservice/pickupdropoffcache"
)

// This function is exported to be called from tests until a cleaner approach is implemented.

// RegisterPickupDropoffServices registers pickup/dropoff services
func RegisterPickupDropoffServices(c *dig.Container) {
	must(c.Provide(pickupdropoffservice.NewWalkingTimeCalculator))
	must(c.Provide(pickupdropoffservice.NewIntersectionBasedGenerator))
	must(c.Provide(pickupdropoffcache.NewPickupDropoffCache))
	must(c.Provide(pickupdropoffservice.NewPickupDropoffSelector))
}
