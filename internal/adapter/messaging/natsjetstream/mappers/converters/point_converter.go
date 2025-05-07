package converters

import (
	"matching-engine/internal/adapter/messaging/natsjetstream/dto"
	"matching-engine/internal/enums"
	"matching-engine/internal/model"
	"time"
)

// PointConverter handles conversion between domain Point and DTO
type PointConverter struct{}

// ToDTO converts a domain Point to a PointDTO
func (c *PointConverter) ToDTO(p *model.Point) dto.PointDTO {
	owner := p.GetOwner()
	var id, ownerType string

	if offer, ok := owner.AsOffer(); ok {
		id = offer.GetID()
		ownerType = string(enums.Offer)
	} else if request, ok := owner.AsRequest(); ok {
		id = request.GetID()
		ownerType = string(enums.Request)
	}

	return dto.PointDTO{
		OwnerType: ownerType,
		OwnerID:   id,
		Point: dto.CoordinateDTO{
			Lat: p.GetCoordinate().Lat(),
			Lng: p.GetCoordinate().Lng(), // Fixed bug: now using Lng() instead of Lat() twice
		},
		Time:      p.GetTime().Format(time.RFC3339),
		PointType: p.GetPointType().String(),
	}
}

// ToPointsDTO converts a slice of domain Points to a slice of PointDTOs
func (c *PointConverter) ToPointsDTO(points []*model.Point) []dto.PointDTO {
	result := make([]dto.PointDTO, 0, len(points))
	for _, p := range points {
		result = append(result, c.ToDTO(p))
	}
	return result
}

// NewPointConverter creates a new PointConverter
func NewPointConverter() *PointConverter {
	return &PointConverter{}
}
