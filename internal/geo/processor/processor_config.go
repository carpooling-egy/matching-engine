package processor

import (
	"github.com/rs/zerolog/log"
	"matching-engine/internal/enums"
	"os"
	"strconv"
)

type Config struct {
	EnablePruning      bool
	EnableDownsampling bool
	DownsamplerType    enums.DownsamplerType
}

func Load() Config {
	cfg := Config{
		EnablePruning:      getEnvBool("ENABLE_PRUNING", true),
		EnableDownsampling: getEnvBool("ENABLE_DOWNSAMPLING", true),
	}

	raw := getEnv("DOWNSAMPLER_TYPE", string(enums.DownsamplerRDP))
	log.Debug().Msgf("DOWNSAMPLER_TYPE set to: %s", raw)

	dsType := enums.DownsamplerType(raw)
	if !dsType.IsValid() {
		log.Warn().Msgf("Invalid DOWNSAMPLER_TYPE %q, falling back to %s", raw, enums.DownsamplerRDP)
		dsType = enums.DownsamplerRDP
	}
	cfg.DownsamplerType = dsType

	return cfg
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	log.Debug().Msgf("%s not set, using default %q", key, def)
	return def
}

func getEnvBool(key string, def bool) bool {
	val := os.Getenv(key)
	if val == "" {
		log.Debug().Msgf("%s not set, using default: %t", key, def)
		return def
	}
	parsed, err := strconv.ParseBool(val)
	if err != nil {
		log.Warn().Msgf("Invalid %s value %q, using default: %t", key, val, def)
		return def
	}
	log.Debug().Msgf("%s set to: %t", key, parsed)
	return parsed
}

func (c Config) String() string {
	return "Config{" +
		"EnablePruning: " + strconv.FormatBool(c.EnablePruning) + ", " +
		"EnableDownsampling: " + strconv.FormatBool(c.EnableDownsampling) + ", " +
		"DownsamplerType: " + c.DownsamplerType.String() +
		"}"
}
