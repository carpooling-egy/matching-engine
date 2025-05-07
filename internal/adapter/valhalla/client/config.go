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
	c := defaultConfig()

	if v, ok := os.LookupEnv("VALHALLA_HOST"); ok && v != "" {
		c.host = v
	}

	if v, ok := os.LookupEnv("VALHALLA_PORT"); ok && v != "" {
		p, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("invalid VALHALLA_PORT %q: %w", v, err)
		}
		c.port = p
	}

	return c, nil
}

func defaultConfig() *Config {
	return &Config{
		host: "localhost",
		port: 8002,
	}
}

func (c *Config) ValhallaHost() string { return c.host }
func (c *Config) ValhallaPort() int    { return c.port }
func (c *Config) ValhallaURL() string {
	return fmt.Sprintf("http://%s:%d", c.host, c.port)
}
