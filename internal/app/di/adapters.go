package di

import (
	"fmt"
	"go.uber.org/dig"
	"matching-engine/internal/adapter/messaging/natsjetstream"
	"matching-engine/internal/adapter/osrm"
	"matching-engine/internal/adapter/routing"
	"matching-engine/internal/adapter/valhalla"
	"matching-engine/internal/app/di/utils"
)

// This function is exported to be called from tests until a cleaner approach is implemented.

// registerAdapters registers external adapters
func registerAdapters(c *dig.Container) {
	utils.Must(c.Provide(provideRoutingEngine))
	utils.Must(c.Provide(natsjetstream.NewNATSPublisher))
}

// provideRoutingEngine provides the appropriate routing engine based on configuration
func provideRoutingEngine() (routing.Engine, error) {
	config, err := routing.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load routing config: %w", err)
	}

	switch {
	case config.IsValhalla():
		return valhalla.NewValhalla()
	case config.IsOSRM():
		return osrm.NewOSRM()
	default:
		return nil, fmt.Errorf("no valid routing engine configured")
	}
}
