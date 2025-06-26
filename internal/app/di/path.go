package di

import (
	"github.com/rs/zerolog/log"
	"go.uber.org/dig"
	"matching-engine/internal/adapter/ortool"
	"matching-engine/internal/app/di/utils"
	"matching-engine/internal/service/pathgeneration/generator"
	"matching-engine/internal/service/pathgeneration/planner"
	"matching-engine/internal/service/pathgeneration/validator"
	"os"
)

// This function is exported to be called from tests until a cleaner approach is implemented.

// RegisterPathServices registers path generation services
func RegisterPathServices(c *dig.Container) {

	plannerType := getPathPlannerType()
	switch plannerType {
	case "ortool":
		utils.Must(c.Provide(planner.NewORToolPlanner))
		utils.Must(c.Provide(ortool.NewORToolClient))
	default:
		utils.Must(c.Provide(generator.NewPathGenerator))
		utils.Must(c.Provide(validator.NewDefaultPathValidator))
		utils.Must(c.Provide(planner.NewDefaultPathPlanner))
	}

}

func getPathPlannerType() string {
	pathPlannerType := "default" // Default path planner type
	if v, ok := os.LookupEnv("PATH_PLANNER_TYPE"); ok && v != "" {
		pathPlannerType = v
	} else {
		log.Warn().Msgf("PATH_PLANNER_TYPE environment variable is not set. Using default: %s", pathPlannerType)
	}
	return pathPlannerType
}
