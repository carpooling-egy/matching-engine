package config

import (
	"github.com/joho/godotenv"
	"log"
)

func LoadEnv() error {
	// Try to load from .env file, but don't fail if file is missing
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}
	return nil
}
