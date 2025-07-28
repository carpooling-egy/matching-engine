package processor

import (
	"github.com/rs/zerolog/log"
	"matching-engine/internal/app/config"
	"matching-engine/internal/enums"
	"strconv"
)

type Config struct {
	EnablePruning      bool
	EnableDownsampling bool
	DownsamplerType    enums.DownsamplerType
}

func Load() Config {
	cfg := Config{
		EnablePruning:      config.GetEnvBool("ENABLE_PRUNING", true),
		EnableDownsampling: config.GetEnvBool("ENABLE_DOWNSAMPLING", true),
	}

	raw := config.GetEnv("DOWNSAMPLER_TYPE", string(enums.DownsamplerRDP))
	log.Debug().Msgf("DOWNSAMPLER_TYPE set to: %s", raw)

	dsType := enums.DownsamplerType(raw)
	if !dsType.IsValid() {
		log.Warn().Msgf("Invalid DOWNSAMPLER_TYPE %q, falling back to %s", raw, enums.DownsamplerRDP)
		dsType = enums.DownsamplerRDP
	}
	cfg.DownsamplerType = dsType

	return cfg
}

func (c Config) String() string {
	return "Config{" +
		"EnablePruning: " + strconv.FormatBool(c.EnablePruning) + ", " +
		"EnableDownsampling: " + strconv.FormatBool(c.EnableDownsampling) + ", " +
		"DownsamplerType: " + c.DownsamplerType.String() +
		"}"
}
