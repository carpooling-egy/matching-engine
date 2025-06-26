package planner

import (
	"os"
	"strconv"
)

type ORToolConfig struct {
	Timeout                 int
	Method                  string
	EnableGuidedLocalSearch bool
}

func NewORToolConfig() *ORToolConfig {
	// Default values
	defaultTimeout := 100 // in milliseconds
	defaultMethod := "parallel_cheapest_insertion"
	defaultEnableGuidedLocalSearch := false

	// Read from environment variables
	timeout, err := strconv.Atoi(os.Getenv("ORTOOL_TIMEOUT"))
	if err != nil {
		timeout = defaultTimeout
	}

	method := os.Getenv("ORTOOL_METHOD")
	if method == "" {
		method = defaultMethod
	}

	enableGuidedLocalSearch, err := strconv.ParseBool(os.Getenv("ORTOOL_ENABLE_GUIDED_LOCAL_SEARCH"))
	if err != nil {
		enableGuidedLocalSearch = defaultEnableGuidedLocalSearch
	}

	return &ORToolConfig{
		Timeout:                 timeout,
		Method:                  method,
		EnableGuidedLocalSearch: enableGuidedLocalSearch,
	}
}
