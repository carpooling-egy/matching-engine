package generator

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
)

const (
	// DefaultK is the default number of random samples to generate
	DefaultK = 1000
)

func initEnv() {
	// IMPORTANT NOTE: replace it with your actual path to the .env file
	if err := godotenv.Load("/home/husseinkhaled/GolandProjects/matching-engine/.env"); err != nil {
		fmt.Println("Error loading .env:", err)
	}
}

func getNumberOfSamples() int {
	numberOfSamples := DefaultK // Default number of samples
	initEnv()
	if v, ok := os.LookupEnv("SAMPLES_K"); ok && v != "" {
		k, err := strconv.Atoi(v)
		if err != nil {
			log.Error().Msgf("Invalid SAMPLES_K value: %s, using default: %d", v, DefaultK)
		} else {
			numberOfSamples = k
		}
	} else {
		log.Warn().Msgf("SAMPLES_K environment variable is not set. Using default: %d", DefaultK)
	}
	return numberOfSamples
}

func getPathGeneratorType() string {
	initEnv()
	pathGeneratorType := "insertion" // Default path generator type
	if v, ok := os.LookupEnv("PATH_GENERATOR_TYPE"); ok && v != "" {
		pathGeneratorType = v
	} else {
		log.Warn().Msgf("PATH_GENERATOR_TYPE environment variable is not set. Using default: %s", pathGeneratorType)
	}
	return pathGeneratorType
}
