package di

import (
	"go.uber.org/dig"

	"matching-engine/internal/adapter/messaging/natsjetstream"
	"matching-engine/internal/adapter/valhalla"
)

// registerAdapters registers external adapters
func registerAdapters(c *dig.Container) {
	must(c.Provide(valhalla.NewValhalla))
	must(c.Provide(natsjetstream.NewNATSPublisher))
}
