package routing

import (
	"fmt"
	"os"
	"strings"
)

// EngineType represents the available routing engines
type EngineType string

const (
	EngineTypeValhalla EngineType = "valhalla"
	EngineTypeOSRM     EngineType = "osrm"
)

type Config struct {
	Engine EngineType
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		Engine: EngineTypeValhalla,
	}

	if v, ok := os.LookupEnv("ROUTING_ENGINE"); ok && v != "" {
		engine := EngineType(strings.ToLower(v))
		if !engine.IsValid() {
			return nil, fmt.Errorf("invalid ROUTING_ENGINE %q: must be 'valhalla' or 'osrm'", v)
		}
		cfg.Engine = engine
	}

	return cfg, nil
}

func (e EngineType) IsValid() bool {
	return e == EngineTypeValhalla || e == EngineTypeOSRM
}

func (e EngineType) String() string {
	return string(e)
}

func (c *Config) IsValhalla() bool {
	return c.Engine == EngineTypeValhalla
}

func (c *Config) IsOSRM() bool {
	return c.Engine == EngineTypeOSRM
}
