package planner

import (
	"github.com/rs/zerolog/log"
	"os"
)

func getPathPlannerType() string {
	pathPlannerType := "default" // Default path planner type
	if v, ok := os.LookupEnv("PATH_PLANNER_TYPE"); ok && v != "" {
		pathPlannerType = v
	} else {
		log.Warn().Msgf("PATH_PLANNER_TYPE environment variable is not set. Using default: %s", pathPlannerType)
	}
	return pathPlannerType
}
