package di

import (
	"go.uber.org/dig"
	"matching-engine/internal/app/di/utils"

	"matching-engine/internal/service/pickupdropoffservice"
	"matching-engine/internal/service/pickupdropoffservice/pickupdropoffcache"
)

// This function is exported to be called from tests until a cleaner approach is implemented.

// RegisterPickupDropoffServices registers pickup/dropoff services
func RegisterPickupDropoffServices(c *dig.Container) {
	utils.Must(c.Provide(pickupdropoffservice.NewWalkingTimeCalculator))
	utils.Must(c.Provide(pickupdropoffservice.NewIntersectionBasedGenerator))
	utils.Must(c.Provide(pickupdropoffcache.NewPickupDropoffCache))
	utils.Must(c.Provide(pickupdropoffservice.NewPickupDropoffSelector))
}
