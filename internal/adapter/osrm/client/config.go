package client

import (
	"fmt"
	"matching-engine/internal/model"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	host    string
	port    int
	profile model.OSRMProfile
}

var defaultProfileConfig = map[model.OSRMProfile]struct {
	host string
	port int
}{
	model.OSRMProfileCar:  {"20.46.49.38", 5000},
	model.OSRMProfileFoot: {"20.46.49.38", 5001},
}

func LoadConfig(rawProfile string) (*Config, error) {
	profile := model.OSRMProfile(rawProfile)
	def, ok := defaultProfileConfig[profile]
	if !ok {
		def = struct {
			host string
			port int
		}{"localhost", 5000}
	}

	cfg := &Config{
		host:    def.host,
		port:    def.port,
		profile: profile,
	}

	profileUpper := strings.ToUpper(rawProfile)

	hostEnv := fmt.Sprintf("OSRM_%s_HOST", profileUpper)
	portEnv := fmt.Sprintf("OSRM_%s_PORT", profileUpper)

	if v, ok := os.LookupEnv(hostEnv); ok && v != "" {
		cfg.host = v
	}

	if v, ok := os.LookupEnv(portEnv); ok && v != "" {
		port, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("invalid %s %q: %w", portEnv, v, err)
		}
		cfg.port = port
	}

	return cfg, nil
}

func (c *Config) OSRMHost() string               { return c.host }
func (c *Config) OSRMPort() int                  { return c.port }
func (c *Config) OSRMProfile() model.OSRMProfile { return c.profile }
func (c *Config) OSRMURL() string {
	return fmt.Sprintf("http://%s:%d", c.host, c.port)
}
