package di

import (
	"go.uber.org/dig"

	"matching-engine/internal/adapter/messaging/natsjetstream"
	"matching-engine/internal/adapter/valhalla"
)

// The fn is exported to be call them from tests, until we build a cleaner approach

// registerAdapters registers external adapters
func registerAdapters(c *dig.Container) {
	must(c.Provide(valhalla.NewValhalla))
	must(c.Provide(natsjetstream.NewNATSPublisher))
}
