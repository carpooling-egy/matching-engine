package mappers

import (
	"encoding/json"
	"matching-engine/internal/adapter/messaging/natsjetstream/dto"
	"matching-engine/internal/enums"
	"matching-engine/internal/model"
	"time"
)

// JsonMapper implements the Mapper interface for JSON serialization
type JsonMapper struct {
}

func (mapper *JsonMapper) Marshal(result *model.MatchingResult) ([]byte, error) {
	resultDTO := dto.MatchingResultDTO{
		UserID:                  result.GetUserID(),
		OfferID:                 result.GetOfferID(),
		AssignedMatchedRequests: make([]dto.MatchedRequestDTO, 0, len(result.GetAssignedMatchedRequests())),
		Path:                    convertPoints(result.GetNewPath()),
	}

	for _, req := range result.GetAssignedMatchedRequests() {
		resultDTO.AssignedMatchedRequests = append(resultDTO.AssignedMatchedRequests, convertMatchedRequest(req))
	}

	return json.Marshal(resultDTO)
}

func convertMatchedRequest(req *model.MatchedRequest) dto.MatchedRequestDTO {
	return dto.MatchedRequestDTO{
		UserID:       req.GetRequest().GetUserID(),
		RequestID:    req.GetRequest().GetID(),
		PickupPoint:  convertPoint(req.GetPickup()),
		DropoffPoint: convertPoint(req.GetDropoff()),
	}
}

func convertPoints(points []*model.Point) []dto.PointDTO {
	result := make([]dto.PointDTO, 0, len(points))
	for _, p := range points {
		result = append(result, convertPoint(p))
	}
	return result
}

func convertPoint(p *model.Point) dto.PointDTO {
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
			Lng: p.GetCoordinate().Lat(),
		},
		Time:      p.GetTime().Format(time.RFC3339),
		PointType: p.GetPointType().String(),
	}
}

// NewJsonMapper creates a new JsonMapper
func NewJsonMapper() Mapper {
	return &JsonMapper{}
}
