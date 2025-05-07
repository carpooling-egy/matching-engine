package client

import (
	"bytes"
	"fmt"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
	"io"
	re "matching-engine/internal/adapter/routing"
	"matching-engine/internal/adapter/valhalla/client/pb"
	"net/http"
	"net/url"
)

type ValhallaClient struct {
	cfg *Config
}

type Option func(*ValhallaClient)

func WithConfig(c *Config) Option {
	return func(vc *ValhallaClient) {
		if c != nil {
			vc.cfg = c
		}
	}
}

func NewValhallaClient(opts ...Option) (*ValhallaClient, error) {
	cfg, err := LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	vc := &ValhallaClient{cfg: cfg}

	for _, opt := range opts {
		opt(vc)
	}

	baseURL := vc.cfg.ValhallaURL()
	_, err = url.ParseRequestURI(baseURL)
	if err != nil {
		log.Error().
			Err(err).
			Str("baseURL", baseURL).
			Msg("invalid base URL format")
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	return vc, nil
}

var _ re.RoutingClient[
	*pb.Api,
	*pb.Api,
] = (*ValhallaClient)(nil)

func (vc *ValhallaClient) Post(endpoint string, request *pb.Api) (*pb.Api, error) {
	data, err := vc.serializeRequest(request)
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to serialize request")
		return nil, fmt.Errorf("failed to serialize request: %w", err)
	}

	body, err := vc.doPost(endpoint, data)
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to send request")
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	response, err := vc.deserializeResponse(body)
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to deserialize response")
		return nil, fmt.Errorf("failed to deserialize response: %w", err)
	}

	return response, nil
}

func (vc *ValhallaClient) doPost(endpoint string, data []byte) ([]byte, error) {
	resp, err := http.Post(
		fmt.Sprintf("%v%v?format=proto", vc.cfg.ValhallaURL(), endpoint),
		"application/x-protobuf",
		bytes.NewReader(data),
	)
	if err != nil {
		log.Error().
			Err(err).
			Str("endpoint", endpoint).
			Msg("HTTP POST failed")
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Warn().
				Err(err).
				Msg("failed to close response body")
		}
	}(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		snippet, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		log.Error().
			Int("status", resp.StatusCode).
			Str("endpoint", endpoint).
			Bytes("body_snippet", snippet).
			Msg("non-successful HTTP response")
		return nil, fmt.Errorf(
			"unexpected HTTP status %d from %s: %q",
			resp.StatusCode, endpoint, snippet,
		)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("failed to read response body")
		return nil, fmt.Errorf("reading response: %w", err)
	}
	return body, nil
}

func (vc *ValhallaClient) serializeRequest(request *pb.Api) ([]byte, error) {
	data, err := proto.Marshal(request)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal protobuf request")
		return nil, fmt.Errorf("failed to serialize request: %w", err)
	}
	return data, nil
}

func (vc *ValhallaClient) deserializeResponse(body []byte) (*pb.Api, error) {
	response := &pb.Api{}
	if err := proto.Unmarshal(body, response); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal protobuf response")
		return nil, fmt.Errorf("failed to deserialize response: %w", err)
	}
	return response, nil
}
