package client

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	host string
	port int
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		host: "68.221.112.34",
		port: 8004,
	}

	if v, ok := os.LookupEnv("VALHALLA_HOST"); ok && v != "" {
		cfg.host = v
	}

	if v, ok := os.LookupEnv("VALHALLA_PORT"); ok && v != "" {
		port, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("invalid VALHALLA_PORT %q: %w", v, err)
		}
		cfg.port = port
	}

	return cfg, nil
}

func (c *Config) ValhallaHost() string { return c.host }
func (c *Config) ValhallaPort() int    { return c.port }
func (c *Config) ValhallaURL() string {
	return fmt.Sprintf("http://%s:%d", c.host, c.port)
}
