package utils

import "github.com/rs/zerolog/log"

// Must is a helper function to handle errors during initialization
func Must(err error) {
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to configure dependency injection")
	}
}
