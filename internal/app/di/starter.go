package di

import (
	"go.uber.org/dig"
	"matching-engine/internal/app/di/utils"
	"matching-engine/internal/app/starter"
)

// This function is exported to be called from tests until a cleaner approach is implemented.

// RegisterStarterService registers the starter service
func RegisterStarterService(c *dig.Container) {
	utils.Must(c.Provide(starter.NewStarterService))
}
