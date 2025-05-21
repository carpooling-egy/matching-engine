package di

import (
	"go.uber.org/dig"
	"matching-engine/internal/app/starter"
)

// The fn is exported to be call them from tests, until we build a cleaner approach

// RegisterStarterService registers path generation services
func RegisterStarterService(c *dig.Container) {
	must(c.Provide(starter.NewStarterService))
}
