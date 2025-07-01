package client

import (
	"bytes"
	"fmt"
	"io"
	re "matching-engine/internal/adapter/routing"
	"matching-engine/internal/adapter/valhalla/client/pb"
	"net/http"
	"net/url"

	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
)

type ValhallaClient struct {
	cfg        *Config
	httpClient *http.Client
}

func NewValhallaClient() (re.Client[
	*pb.Api,
	*pb.Api,
], error) {
	cfg, err := LoadConfig()
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to load configuration for ValhallaClient")
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	httpClient := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 50,
		},
	}
	vc := &ValhallaClient{
		cfg:        cfg,
		httpClient: httpClient,
	}

	baseURL := vc.cfg.ValhallaURL()
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
		Msg("ValhallaClient initialized with the following base URL")

	return vc, nil
}

func (vc *ValhallaClient) Post(endpoint string, request *pb.Api) (*pb.Api, error) {
	data, err := vc.serializeRequest(request)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to serialize the request to protobuf format")
		return nil, fmt.Errorf("failed to serialize request: %w", err)
	}

	body, err := vc.doPost(endpoint, data)
	if err != nil {
		log.Error().
			Err(err).
			Str("endpoint", endpoint).
			Msg("Failed to send HTTP POST request to Valhalla endpoint")
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	response, err := vc.deserializeResponse(body)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to deserialize the response from protobuf format")
		return nil, fmt.Errorf("failed to deserialize response: %w", err)
	}

	// log.Debug().
	// 	Str("endpoint", endpoint).
	// 	Msg("Successfully received and processed response from Valhalla")

	return response, nil
}

func (vc *ValhallaClient) doPost(endpoint string, data []byte) ([]byte, error) {
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%v%v?format=proto", vc.cfg.ValhallaURL(), endpoint),
		bytes.NewReader(data),
	)
	if err != nil {
		log.Error().
			Err(err).
			Str("endpoint", endpoint).
			Msg("Failed to create HTTP request")
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-protobuf")

	resp, err := vc.httpClient.Do(req)
	if err != nil {
		log.Error().
			Err(err).
			Str("endpoint", endpoint).
			Msg("HTTP POST request failed")
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Warn().
				Err(err).
				Msg("Failed to close the response body")
		}
	}(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		snippet, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		log.Error().
			Int("status", resp.StatusCode).
			Str("endpoint", endpoint).
			Bytes("body_snippet", snippet).
			Msg("Received non-success HTTP response")
		return nil, fmt.Errorf(
			"unexpected HTTP status %d from %s: %q",
			resp.StatusCode, endpoint, snippet,
		)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().
			Err(err).
			Str("endpoint", endpoint).
			Msg("Failed to read response body")
		return nil, fmt.Errorf("reading response: %w", err)
	}

	// log.Debug().
	// 	Str("endpoint", endpoint).
	// 	Int("statusCode", resp.StatusCode).
	// 	Msg("Successfully received HTTP response from Valhalla")

	return body, nil
}

func (vc *ValhallaClient) serializeRequest(request *pb.Api) ([]byte, error) {
	data, err := proto.Marshal(request)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to marshal protobuf request into byte format")
		return nil, fmt.Errorf("failed to serialize request: %w", err)
	}

	// log.Debug().
	// 	Msg("Protobuf request serialized successfully")
	return data, nil
}

func (vc *ValhallaClient) deserializeResponse(body []byte) (*pb.Api, error) {
	response := &pb.Api{}
	if err := proto.Unmarshal(body, response); err != nil {
		log.Error().
			Err(err).
			Msg("Failed to unmarshal response body into protobuf format")
		return nil, fmt.Errorf("failed to deserialize response: %w", err)
	}

	// log.Debug().
	// 	Msg("Response successfully unmarshalled from protobuf format")
	return response, nil
}

func (vc *ValhallaClient) Get(endpoint string, params *pb.Api) (*pb.Api, error) {
	log.Error().
		Str("endpoint", endpoint).
		Msg("GET method is not supported by ValhallaClient")
	panic("GET method is not supported by ValhallaClient")
}
