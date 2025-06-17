package di

import (
	"go.uber.org/dig"
	"matching-engine/internal/app/di/utils"

	"matching-engine/internal/adapter/messaging/natsjetstream"
	"matching-engine/internal/adapter/valhalla"
)

// This function is exported to be called from tests until a cleaner approach is implemented.

// registerAdapters registers external adapters
func registerAdapters(c *dig.Container) {
	utils.Must(c.Provide(valhalla.NewValhalla))
	utils.Must(c.Provide(natsjetstream.NewNATSPublisher))
}
