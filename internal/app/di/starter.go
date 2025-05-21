package di

import (
	"go.uber.org/dig"
	"matching-engine/internal/app/starter"
)

// registerStarterService registers path generation services
func registerStarterService(c *dig.Container) {
	must(c.Provide(starter.NewStarterService))
}
