package di

import (
	"go.uber.org/dig"

	"matching-engine/internal/service/pathgeneration/generator"
	"matching-engine/internal/service/pathgeneration/planner"
	"matching-engine/internal/service/pathgeneration/validator"
)

// This function is exported to be called from tests until a cleaner approach is implemented.

// RegisterPathServices registers path generation services
func RegisterPathServices(c *dig.Container) {
	must(c.Provide(generator.NewInsertionPathGenerator))
	must(c.Provide(validator.NewDefaultPathValidator))
	must(c.Provide(planner.NewDefaultPathPlanner))
}
