package di

import (
	"go.uber.org/dig"

	"matching-engine/internal/service/pathgeneration/generator"
	"matching-engine/internal/service/pathgeneration/planner"
	"matching-engine/internal/service/pathgeneration/validator"
)

// The fn is exported to be call them from tests, until we build a cleaner approach

// RegisterPathServices registers path generation services
func RegisterPathServices(c *dig.Container) {
	must(c.Provide(generator.NewInsertionPathGenerator))
	must(c.Provide(validator.NewDefaultPathValidator))
	must(c.Provide(planner.NewDefaultPathPlanner))
}
