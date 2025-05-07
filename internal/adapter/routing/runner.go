package routing

import (
	"context"
	"fmt"
)

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
) (DomainRes, error) {
	var zero DomainRes

	request, err := mapper.ToTransport(params)
	if err != nil {
		return zero, fmt.Errorf(
			"failed to convert params to transport format: %w", err,
		)
	}

	response, err := client.Post(endpoint, request)
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
