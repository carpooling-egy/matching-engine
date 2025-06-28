package timematrix

import (
	"os"
	"strconv"
)

// Default caching bound value if environment variable is not set
const DefaultCachingBound = 40

// GetCachingBound reads the CACHING_BOUND environment variable and returns its value as an integer.
// If the environment variable is not set or cannot be parsed, it returns the default value.
func GetCachingBound() int {
	cachingBoundStr := os.Getenv("CACHING_BOUND")
	if cachingBoundStr == "" {
		return DefaultCachingBound
	}

	cachingBound, err := strconv.Atoi(cachingBoundStr)
	if err != nil {
		return DefaultCachingBound
	}

	return cachingBound
}
