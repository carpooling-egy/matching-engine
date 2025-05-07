package converters

import (
	"matching-engine/internal/adapter/messaging/natsjetstream/dto"
	"matching-engine/internal/model"
)

// RequestConverter handles conversion between domain MatchedRequest and DTO
type RequestConverter struct {
	pointConverter *PointConverter
}

// ToDTO converts a domain MatchedRequest to a MatchedRequestDTO
func (c *RequestConverter) ToDTO(req *model.MatchedRequest) dto.MatchedRequestDTO {
	return dto.MatchedRequestDTO{
		UserID:       req.Request().UserID(),
		RequestID:    req.Request().ID(),
		PickupPoint:  c.pointConverter.ToDTO(req.Pickup()),
		DropoffPoint: c.pointConverter.ToDTO(req.Dropoff()),
	}
}

// ToMatchedRequestsDTO converts a slice of domain MatchedRequests to a slice of MatchedRequestDTOs
func (c *RequestConverter) ToMatchedRequestsDTO(requests []*model.MatchedRequest) []dto.MatchedRequestDTO {
	result := make([]dto.MatchedRequestDTO, 0, len(requests))
	for _, req := range requests {
		result = append(result, c.ToDTO(req))
	}
	return result
}

// NewRequestConverter creates a new RequestConverter
func NewRequestConverter() *RequestConverter {
	return &RequestConverter{
		pointConverter: NewPointConverter(),
	}
}
