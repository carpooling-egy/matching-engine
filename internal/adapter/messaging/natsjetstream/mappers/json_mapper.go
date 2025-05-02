package mappers

import (
	"encoding/json"
	"matching-engine/internal/adapter/messaging/natsjetstream/dto"
	"matching-engine/internal/model"
	"time"
)

// JsonMapper implements the Mapper interface for JSON serialization
type JsonMapper struct {
}

func (mapper *JsonMapper) Marshal(result *model.MatchingResult) ([]byte, error) {
	resultDTO := dto.MatchingResultDTO{
		OfferID:                 result.OfferID(),
		AssignedMatchedRequests: make([]dto.MatchedRequestDTO, 0, len(result.AssignedMatchedRequests())),
		Path:                    convertPoints(result.Path()),
	}

	for _, req := range result.AssignedMatchedRequests() {
		resultDTO.AssignedMatchedRequests = append(resultDTO.AssignedMatchedRequests, convertMatchedRequest(req))
	}

	return json.Marshal(resultDTO)
}

func convertMatchedRequest(req model.MatchedRequest) dto.MatchedRequestDTO {
	return dto.MatchedRequestDTO{
		RequestID:     req.RequestID(),
		PickupPoints:  convertPoints(req.PickupPoints()),
		DropoffPoints: convertPoints(req.DropoffPoints()),
	}
}

func convertPoints(points []model.Point) []dto.PointDTO {
	result := make([]dto.PointDTO, 0, len(points))
	for _, p := range points {
		result = append(result, convertPoint(p))
	}
	return result
}

func convertPoint(p model.Point) dto.PointDTO {
	return dto.PointDTO{
		RequestID: p.RequestID(),
		Point: dto.CoordinateDTO{
			Lat: p.Coordinate().Lat(),
			Lng: p.Coordinate().Lng(),
		},
		Time:      p.Time().Format(time.RFC3339),
		PointType: int(p.PointType()),
	}
}

// NewJsonMapper creates a new JsonMapper
func NewJsonMapper() Mapper {
	return &JsonMapper{}
}
