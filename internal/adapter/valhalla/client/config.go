package client

import (
	"os"
)

type Config struct {
	url string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		url: "http://68.221.112.34:8003",
	}

	if v, ok := os.LookupEnv("VALHALLA_URL"); ok && v != "" {
		cfg.url = v
	}

	return cfg, nil
}

func (c *Config) ValhallaURL() string { return c.url }
