package di

import (
	"go.uber.org/dig"

	"matching-engine/internal/service/pathgeneration/generator"
	"matching-engine/internal/service/pathgeneration/planner"
	"matching-engine/internal/service/pathgeneration/validator"
)

// registerPathServices registers path generation services
func registerPathServices(c *dig.Container) {
	must(c.Provide(generator.NewInsertionPathGenerator))
	must(c.Provide(validator.NewDefaultPathValidator))
	must(c.Provide(planner.NewDefaultPathPlanner))
}
