package generator

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
)

type Config struct {
	samplesK int
}

func loadEnv() {
	// IMPORTANT NOTE: replace it with your actual path to the .env file
	if err := godotenv.Load("/home/husseinkhaled/GolandProjects/matching-engine/.env"); err != nil {
		fmt.Println("Error loading .env:", err)
	}
}

func LoadConfig() *Config {
	c := defaultConfig()
	loadEnv()
	if v, ok := os.LookupEnv("SAMPLES_K"); ok && v != "" {
		k, err := strconv.Atoi(v)
		if err != nil {
			log.Error().Msgf("PATHGEN_K environment variable must be an integer." +
				"Using the default k instead")
		} else {
			c.samplesK = k
		}
	}
	return c
}

func defaultConfig() *Config {
	return &Config{
		samplesK: DefaultK,
	}
}

func (c *Config) K() int {
	return c.samplesK
}
