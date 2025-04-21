package client

import (
	"bytes"
	"fmt"
	"google.golang.org/protobuf/proto"
	"io"
	"matching-engine/internal/adapter/routing-engine/valhalla/client/pb"
	"net/http"
)

type ValhallaClient struct {
	baseUrl    string
	portNumber string
}

func NewValhallaClient(baseUrl string, portNumber string) (*ValhallaClient, error) {
	if baseUrl == "" {
		return nil, fmt.Errorf("base URL cannot be empty")
	}
	if portNumber == "" {
		return nil, fmt.Errorf("port number cannot be empty")
	}

	return &ValhallaClient{
		baseUrl:    baseUrl,
		portNumber: portNumber,
	}, nil
}

func (vc ValhallaClient) SendPostRequest(endpoint string, request *pb.Api) (*pb.Api, error) {
	data, err := vc.serializeRequest(request)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize request: %w", err)
	}

	body, err := vc.doPost(endpoint, data)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	response, err := vc.deserializeResponse(body)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize response: %w", err)
	}

	return response, nil
}

func (vc ValhallaClient) doPost(endpoint string, data []byte) ([]byte, error) {
	resp, err := http.Post(
		fmt.Sprintf(
			vc.baseUrl,
			vc.portNumber,
			endpoint,
			"?format=proto",
		),
		"application/x-protobuf",
		bytes.NewReader(data),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("warning: failed to close response body: %v\n", err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}
	return body, nil
}

func (vc ValhallaClient) serializeRequest(request *pb.Api) ([]byte, error) {
	data, err := proto.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize request: %w", err)
	}
	return data, nil
}

func (vc ValhallaClient) deserializeResponse(body []byte) (*pb.Api, error) {
	response := &pb.Api{}
	if err := proto.Unmarshal(body, response); err != nil {
		return nil, fmt.Errorf("failed to deserialize response: %w", err)
	}
	return response, nil
}
