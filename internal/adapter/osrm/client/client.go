package client

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"matching-engine/internal/adapter/routing"
	"matching-engine/internal/model"
	"net/http"
	"net/url"
	"strings"
)

// URL Schema: {baseURL}/{service}/{version}/{profile}/{coordinates}?option=value&option=value
// endpoint = {baseURL}/{service}/{version}/{profile}/

type OSRMClient struct {
	cfg        *Config
	httpClient *http.Client
}

func NewOSRMClient(profile string) (routing.Client[model.OSRMTransport, map[string]any], error) {
	cfg, err := LoadConfig(profile)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to load configuration for OSRMClient")
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 50,
		},
	}

	client := &OSRMClient{
		cfg:        cfg,
		httpClient: httpClient,
	}

	baseURL := client.cfg.OSRMURL()
	_, err = url.ParseRequestURI(baseURL)
	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", baseURL).
			Msg("Invalid base URL format provided in configuration")
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	log.Info().
		Str("baseURL", baseURL).
		Msg("OSRMClient initialized with the following base URL")

	return client, nil
}

func (c *OSRMClient) Get(endpoint string, req model.OSRMTransport) (map[string]any, error) {
	fullURL, err := c.buildURL(endpoint, req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to build URL in Get")
		return nil, err
	}

	resp, err := c.doGetRequest(fullURL)
	if err != nil {
		log.Error().Err(err).Msg("Failed to perform GET request")
		return nil, err
	}

	result, err := c.decodeResponse(resp)
	if err != nil {
		log.Error().Err(err).Msg("Failed to decode response")
		return nil, err
	}
	return result, nil
}

func (c *OSRMClient) buildURL(endpoint string, req model.OSRMTransport) (string, error) {
	baseURL := c.cfg.OSRMURL()
	path := endpoint

	appendPathVars(req.PathVars, &path)

	fullURL := baseURL + path
	if len(req.QueryParams) > 0 {
		u, err := url.Parse(fullURL)
		if err != nil {
			log.Error().Err(err).Str("fullURL", fullURL).Msg("Invalid URL in buildURL")
			return "", fmt.Errorf("invalid URL: %w", err)
		}

		q := u.Query()
		for k, v := range req.QueryParams {
			q.Set(k, fmt.Sprintf("%v", v))
		}
		u.RawQuery = q.Encode()
		fullURL = u.String()
	}
	return fullURL, nil
}

func appendPathVars(vars map[string]string, path *string) {
	order := []string{"coordinates"}

	for _, key := range order {
		if val, ok := vars[key]; ok && val != "" {
			appendVar(path, val)
		}
	}
}

func appendVar(path *string, segment string) {
	trimmed := strings.TrimRight(*path, "/")
	esc := url.PathEscape(segment)
	*path = trimmed + "/" + esc
}

func (c *OSRMClient) doGetRequest(fullURL string) (*http.Response, error) {
	request, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		log.Error().Err(err).Str("fullURL", fullURL).Msg("Failed to create GET request")
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	resp, err := c.httpClient.Do(request)
	if err != nil {
		log.Error().Err(err).Str("fullURL", fullURL).Msg("Failed to send GET request")
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		errClose := resp.Body.Close()
		if errClose != nil {
			log.Error().Err(errClose).Msg("Failed to close response body after non-2xx status")
		}
		log.Error().Int("status", resp.StatusCode).Str("body", string(b)).Msg("Non-2xx HTTP status in doGetRequest")
		return nil, fmt.Errorf("unexpected HTTP status %d: %q", resp.StatusCode, b)
	}
	return resp, nil
}

func (c *OSRMClient) decodeResponse(resp *http.Response) (map[string]any, error) {
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Error().Err(err).Msg("Failed to close response body after decoding")
		}
	}()
	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Error().Err(err).Msg("Failed to decode JSON response body")
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return result, nil
}

func (c *OSRMClient) Post(endpoint string, req model.OSRMTransport) (map[string]any, error) {
	panic("OSRMClient does not support POST requests; use Get instead.")
}
