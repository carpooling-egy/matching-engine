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
func (c *PointConverter) ToDTO(p *model.PathPoint) dto.PointDTO {
	owner := p.Owner()
	var id, ownerType string

	if offer, ok := owner.AsOffer(); ok {
		id = offer.ID()
		ownerType = string(enums.Offer)
	} else if request, ok := owner.AsRequest(); ok {
		id = request.ID()
		ownerType = string(enums.Request)
	}

	return dto.PointDTO{
		OwnerType: ownerType,
		OwnerID:   id,
		Point: dto.CoordinateDTO{
			Lat: p.Coordinate().Lat(),
			Lng: p.Coordinate().Lng(),
		},
		Time:      p.ExpectedArrivalTime().Format(time.RFC3339),
		PointType: p.PointType().String(),
	}
}

// ToPointsDTO converts a slice of domain Points to a slice of PointDTOs
func (c *PointConverter) ToPointsDTO(points []model.PathPoint) []dto.PointDTO {
	result := make([]dto.PointDTO, 0, len(points))
	for i := range points {
		result = append(result, c.ToDTO(&points[i]))
	}
	return result
}

// NewPointConverter creates a new PointConverter
func NewPointConverter() *PointConverter {
	return &PointConverter{}
}
