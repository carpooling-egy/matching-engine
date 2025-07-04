package client

import (
	"fmt"
	"matching-engine/internal/model"
	"os"
	"strings"
)

type Config struct {
	url     string
	profile model.OSRMProfile
}

var defaultProfileConfig = map[model.OSRMProfile]string{
	model.OSRMProfileCar:  "http://72.146.184.88:5000",
	model.OSRMProfileFoot: "http://72.146.184.88:5001",
}

func LoadConfig(rawProfile string) (*Config, error) {
	profile := model.OSRMProfile(rawProfile)
	def, ok := defaultProfileConfig[profile]
	if !ok {
		def = "http://localhost:5000"
	}

	cfg := &Config{
		url:     def,
		profile: profile,
	}

	profileUpper := strings.ToUpper(rawProfile)
	urlEnv := fmt.Sprintf("OSRM_%s_URL", profileUpper)

	if v, ok := os.LookupEnv(urlEnv); ok && v != "" {
		cfg.url = v
	}

	return cfg, nil
}

func (c *Config) OSRMURL() string                { return c.url }
func (c *Config) OSRMProfile() model.OSRMProfile { return c.profile }
