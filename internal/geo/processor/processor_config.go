package processor

import (
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
	enablePruning := true
	enableDownsampling := true

	if val, err := strconv.ParseBool(os.Getenv("ENABLE_PRUNING")); err == nil {
		enablePruning = val
	}
	if val, err := strconv.ParseBool(os.Getenv("ENABLE_DOWNSAMPLING")); err == nil {
		enableDownsampling = val
	}

	rawType := os.Getenv("DOWNSAMPLER_TYPE")
	downsamplerType := enums.DownsamplerType(rawType)
	if !downsamplerType.IsValid() {
		downsamplerType = enums.DownsamplerRDP
	}

	return Config{
		EnablePruning:      enablePruning,
		EnableDownsampling: enableDownsampling,
		DownsamplerType:    downsamplerType,
	}
}
