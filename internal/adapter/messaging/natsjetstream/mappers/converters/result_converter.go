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
		UserID:                  result.GetUserID(),
		OfferID:                 result.GetOfferID(),
		AssignedMatchedRequests: c.requestConverter.ToMatchedRequestsDTO(result.GetAssignedMatchedRequests()),
		Path:                    c.pointConverter.ToPointsDTO(result.GetNewPath()),
	}
}

// NewResultConverter creates a new ResultConverter
func NewResultConverter() *ResultConverter {
	return &ResultConverter{
		requestConverter: NewRequestConverter(),
		pointConverter:   NewPointConverter(),
	}
}
