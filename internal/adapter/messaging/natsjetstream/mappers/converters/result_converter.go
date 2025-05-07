package converters

import (
	"matching-engine/internal/adapter/messaging/natsjetstream/dto"
	"matching-engine/internal/model"
)

// ResultConverter handles conversion between domain MatchingResult and DTO
type ResultConverter struct {
	requestConverter *RequestConverter
	pointConverter   *PointConverter
}

// ToDTO converts a domain MatchingResult to a MatchingResultDTO
func (c *ResultConverter) ToDTO(result *model.MatchingResult) dto.MatchingResultDTO {
	return dto.MatchingResultDTO{
		UserID:                  result.UserID(),
		OfferID:                 result.OfferID(),
		AssignedMatchedRequests: c.requestConverter.ToMatchedRequestsDTO(result.AssignedMatchedRequests()),
		Path:                    c.pointConverter.ToPointsDTO(result.NewPath()),
	}
}

// NewResultConverter creates a new ResultConverter
func NewResultConverter() *ResultConverter {
	return &ResultConverter{
		requestConverter: NewRequestConverter(),
		pointConverter:   NewPointConverter(),
	}
}
