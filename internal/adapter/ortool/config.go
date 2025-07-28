package ortool

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

	if v, ok := os.LookupEnv("ORTOOL_HOST"); ok && v != "" {
		c.host = v
	}

	if v, ok := os.LookupEnv("ORTOOL_PORT"); ok && v != "" {
		p, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("invalid ORTOOL_PORT %q: %w", v, err)
		}
		c.port = p
	}

	return c, nil
}

func defaultConfig() *Config {
	return &Config{
		host: "localhost",
		port: 8000,
	}
}

func (c *Config) ORToolHost() string { return c.host }
func (c *Config) ORToolPort() int    { return c.port }
func (c *Config) ORToolURL() string {
	return fmt.Sprintf("http://%s:%d/solve", c.host, c.port)
}
