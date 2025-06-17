package di

import (
	"go.uber.org/dig"
	"matching-engine/internal/app/di/utils"

	"matching-engine/internal/service/pathgeneration/generator"
	"matching-engine/internal/service/pathgeneration/planner"
	"matching-engine/internal/service/pathgeneration/validator"
)

// This function is exported to be called from tests until a cleaner approach is implemented.

// RegisterPathServices registers path generation services
func RegisterPathServices(c *dig.Container) {
	utils.Must(c.Provide(generator.NewInsertionPathGenerator))
	utils.Must(c.Provide(validator.NewDefaultPathValidator))
	utils.Must(c.Provide(planner.NewDefaultPathPlanner))
}
