package di

import (
	"matching-engine/internal/app/di/utils"
	"matching-engine/internal/service/pickupdropoffservice"
	"matching-engine/internal/service/pickupdropoffservice/pickupdropoffcache"
	"os"
	"strconv"

	"github.com/rs/zerolog/log"
	"go.uber.org/dig"
)

// This function is exported to be called from tests until a cleaner approach is implemented.

// RegisterPickupDropoffServices registers pickup/dropoff services
func RegisterPickupDropoffServices(c *dig.Container) {
	walkingTimeEnabled := getWalkingTimeEnabled()
	if walkingTimeEnabled {
		utils.Must(c.Provide(pickupdropoffservice.NewIntersectionBasedGenerator))
	} else {
		utils.Must(c.Provide(pickupdropoffservice.NewSnappedSourceDestinationGenerator))
	}

	utils.Must(c.Provide(pickupdropoffcache.NewPickupDropoffCache))
	utils.Must(c.Provide(pickupdropoffservice.NewPickupDropoffSelector))
}

func getWalkingTimeEnabled() bool {
	walkingTimeEnabled := true // Default walking time enabled
	if v, ok := os.LookupEnv("WALKING_TIME_ENABLED"); ok {
		// Parse the string value to boolean
		if parsed, err := strconv.ParseBool(v); err == nil {
			return parsed
		} else {
			log.Warn().Msgf("Invalid WALKING_TIME_ENABLED value '%s', using default: %t", v, walkingTimeEnabled)
		}
	} else {
		log.Info().Msgf("WALKING_TIME_ENABLED environment variable not set. Using default: %t", walkingTimeEnabled)
	}

	return walkingTimeEnabled
}
