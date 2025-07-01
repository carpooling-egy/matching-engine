package routing

import (
	"context"
	"fmt"
)

// selectOperation calls the correct client method based on HTTPMethod.
func selectOperation[TransReq any, TransRes any](
	client Client[TransReq, TransRes],
	endpoint string,
	request TransReq,
	method HTTPMethod,
) (TransRes, error) {
	switch method {
	case MethodGet:
		return client.Get(endpoint, request)
	case MethodPost:
		return client.Post(endpoint, request)
	default:
		var zero TransRes
		return zero, fmt.Errorf("unsupported HTTP method: %v", method)
	}
}

func RunOperation[
	DomainReq any,
	DomainRes any,
	TransReq any,
	TransRes any,
](
	ctx context.Context,
	client Client[TransReq, TransRes],
	endpoint string,
	params DomainReq,
	mapper OperationMapper[DomainReq, DomainRes, TransReq, TransRes],
	method HTTPMethod,
) (DomainRes, error) {
	var zero DomainRes

	request, err := mapper.ToTransport(params)
	if err != nil {
		return zero, fmt.Errorf(
			"failed to convert params to transport format: %w", err,
		)
	}

	response, err := selectOperation(client, endpoint, request, method)
	if err != nil {
		return zero, fmt.Errorf(
			"failed to send request to routing backend: %w", err,
		)
	}

	result, err := mapper.FromTransport(response)
	if err != nil {
		return zero, fmt.Errorf(
			"failed to convert response from transport format: %w", err,
		)
	}

	return result, nil
}
