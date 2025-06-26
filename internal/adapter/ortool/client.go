package ortool

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"net/url"
)

type ORToolClient struct {
	cfg *Config
}

func NewORToolClient() (*ORToolClient, error) {
	cfg, err := LoadConfig()
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to load configuration for ORToolClient")
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	client := &ORToolClient{cfg: cfg}

	baseURL := client.cfg.ORToolURL()
	if _, err := url.ParseRequestURI(baseURL); err != nil {
		log.Error().
			Err(err).
			Str("baseURL", baseURL).
			Msg("Invalid base URL format provided in configuration")
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	log.Info().
		Str("baseURL", baseURL).
		Msg("ORToolClient initialized with the following base URL")

	return client, nil
}

func (orc *ORToolClient) CallPythonORToolSolver(payload *ORToolData) (*ORToolSolutionResponse, error) {
	// Marshal data
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal VRP data: %w", err)
	}

	// Send POST request
	resp, err := http.Post(orc.cfg.ORToolURL(), "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to call VRP solver: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error().Err(err).Msg("failed to close response body")
		}
	}(resp.Body)

	// Handle non-200
	if resp.StatusCode != http.StatusOK {
		var errMsg map[string]any
		err := json.NewDecoder(resp.Body).Decode(&errMsg)
		if err != nil {
			return nil, fmt.Errorf("failed to decode error response: %w", err)
		}
		return nil, fmt.Errorf("VRP solver error (%d): %v", resp.StatusCode, errMsg)
	}

	// Parse response
	var solution ORToolSolutionResponse
	if err := json.NewDecoder(resp.Body).Decode(&solution); err != nil {
		return nil, fmt.Errorf("failed to decode solver response: %w", err)
	}

	return &solution, nil
}
