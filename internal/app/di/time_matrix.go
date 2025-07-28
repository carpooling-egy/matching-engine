package di

import (
	"go.uber.org/dig"
	"matching-engine/internal/app/di/utils"

	"matching-engine/internal/service/timematrix"
	"matching-engine/internal/service/timematrix/cache"
)

// This function is exported to be called from tests until a cleaner approach is implemented.

// RegisterTimeMatrixServices registers time matrix services
func RegisterTimeMatrixServices(c *dig.Container) {
	utils.Must(c.Provide(cache.NewTimeMatrixCacheWithOfferId))
	utils.Must(c.Provide(cache.NewTimeMatrixCacheWithOfferIdAndRequestId))
	utils.Must(c.Provide(timematrix.NewDefaultSelector))
	utils.Must(c.Provide(timematrix.NewService))
	utils.Must(c.Provide(timematrix.NewDefaultGenerator))
	utils.Must(c.Provide(timematrix.NewCacheWithOfferIdPopulator))
	utils.Must(c.Provide(timematrix.NewCacheWithOfferIdRequestIdPopulator))
}
