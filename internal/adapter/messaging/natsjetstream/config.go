package natsjetstream

import (
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	URL            string
	Subject        string
	ConnectTimeout time.Duration // Retry connection if the initial attempt fails
	PublishTimeout time.Duration // Set the maximum number of reconnection attempts before giving up (enabled by default)
	ReconnectWait  time.Duration // Wait duration between each reconnection attempt
	MaxReconnects  int
	ConnectionName string
	NatsUsername   string
	NatsPassword   string
}

func DefaultConfig() Config {
	return Config{
		URL:            nats.DefaultURL,
		Subject:        "matched_requests.results",
		ConnectTimeout: 10 * time.Second,
		PublishTimeout: 30 * time.Minute,
		ReconnectWait:  1 * time.Second,
		MaxReconnects:  -1,
		ConnectionName: "matching-engine-publisher",
		NatsUsername:   "publisher",
		NatsPassword:   "publisherpass",
	}
}

func LoadConfig() Config {
	cfg := DefaultConfig()

	// Simple override helper
	override := func(envVar string, target *string) {
		if val := os.Getenv(envVar); val != "" {
			*target = val
		}
	}
	override("NATS_URL", &cfg.URL)
	override("NATS_SUBJECT", &cfg.Subject)
	override("NATS_CONNECTION_NAME", &cfg.ConnectionName)
	override("NATS_USER", &cfg.NatsUsername)
	override("NATS_PASSWORD", &cfg.NatsPassword)

	cfg.ConnectTimeout = getEnvDuration("NATS_CONNECT_TIMEOUT", cfg.ConnectTimeout)
	cfg.PublishTimeout = getEnvDuration("NATS_PUBLISH_TIMEOUT", cfg.PublishTimeout)
	cfg.ReconnectWait = getEnvDuration("NATS_RECONNECT_WAIT", cfg.ReconnectWait)
	cfg.MaxReconnects = getEnvInt("NATS_MAX_RECONNECTS", cfg.MaxReconnects)

	logConfig(cfg)
	return cfg
}

func getEnvDuration(key string, defaultVal time.Duration) time.Duration {
	if valStr := os.Getenv(key); valStr != "" {
		if val, err := time.ParseDuration(valStr); err == nil {
			return val
		} else {
			log.Warn().Err(err).Str("key", key).Msg("Invalid duration format, using default")
		}
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if valStr := os.Getenv(key); valStr != "" {
		if val, err := strconv.Atoi(valStr); err == nil {
			return val
		} else {
			log.Warn().Err(err).Str("key", key).Msg("Invalid int format, using default")
		}
	}
	return defaultVal
}

func logConfig(cfg Config) {
	log.Info().
		Str("url", cfg.URL).
		Str("subject", cfg.Subject).
		Str("connectionName", cfg.ConnectionName).
		Dur("connectTimeout", cfg.ConnectTimeout).
		Dur("publishTimeout", cfg.PublishTimeout).
		Dur("reconnectWait", cfg.ReconnectWait).
		Int("maxReconnects", cfg.MaxReconnects).
		Msg("NATS configuration loaded")
}
