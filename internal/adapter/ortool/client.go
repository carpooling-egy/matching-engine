package ortool

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/rs/zerolog/log"
)

type ORToolClient struct {
	cfg        *Config
	httpClient *http.Client
}

func NewORToolClient() (*ORToolClient, error) {
	cfg, err := LoadConfig()
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to load configuration for ORToolClient")
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	httpClient := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 50,
		},
	}

	client := &ORToolClient{cfg: cfg, httpClient: httpClient}

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
		log.Error().
			Err(err).
			Msg("Failed to marshal ORTool VRP data to JSON format")
		return nil, fmt.Errorf("failed to marshal VRP data: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest(
		http.MethodPost,
		orc.cfg.ORToolURL(),
		bytes.NewReader(body),
	)
	if err != nil {
		log.Error().
			Err(err).
			Str("url", orc.cfg.ORToolURL()).
			Msg("Failed to create HTTP request for ORTool solver")
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Use the configured httpClient with connection pooling
	resp, err := orc.httpClient.Do(req)
	if err != nil {
		log.Error().
			Err(err).
			Str("url", orc.cfg.ORToolURL()).
			Msg("Failed to send HTTP POST request to ORTool solver")
		return nil, fmt.Errorf("failed to call VRP solver: %w", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error().
				Err(err).
				Msg("Failed to close ORTool response body")
		}
	}(resp.Body)

	// Handle non-200 status codes
	if resp.StatusCode != http.StatusOK {
		snippet, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		log.Error().
			Int("status", resp.StatusCode).
			Str("url", orc.cfg.ORToolURL()).
			Bytes("body_snippet", snippet).
			Msg("Received non-success HTTP response from ORTool solver")

		var errMsg map[string]any
		if json.Unmarshal(snippet, &errMsg) == nil {
			return nil, fmt.Errorf("VRP solver error (%d): %v", resp.StatusCode, errMsg)
		}
		return nil, fmt.Errorf("VRP solver error (%d): %s", resp.StatusCode, string(snippet))
	}

	// Read response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().
			Err(err).
			Str("url", orc.cfg.ORToolURL()).
			Msg("Failed to read ORTool solver response body")
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Parse response
	var solution ORToolSolutionResponse
	if err := json.Unmarshal(responseBody, &solution); err != nil {
		log.Error().
			Err(err).
			Str("url", orc.cfg.ORToolURL()).
			Bytes("response_snippet", responseBody[:min(len(responseBody), 256)]).
			Msg("Failed to decode ORTool solver response")
		return nil, fmt.Errorf("failed to decode solver response: %w", err)
	}

	log.Debug().
		Str("url", orc.cfg.ORToolURL()).
		Int("status", resp.StatusCode).
		Msg("Successfully received and processed response from ORTool solver")

	return &solution, nil
}
