package common

import (
	"fmt"
	pb2 "matching-engine/internal/adapter/valhalla/client/pb"
	"matching-engine/internal/model"
)

func CreateLocation(lat, lng float64) *pb2.Location {
	return &pb2.Location{
		Type: DefaultLocationType,
		Ll: &pb2.LatLng{
			HasLat: &pb2.LatLng_Lat{Lat: lat},
			HasLng: &pb2.LatLng_Lng{Lng: lng},
		},
	}
}

func ToDomainDistanceUnit(pbUnit pb2.Options_Units) (model.DistanceUnit, error) {
	switch pbUnit {
	case pb2.Options_kilometers:
		return model.Kilometer, nil
	case pb2.Options_miles:
		return model.Mile, nil
	default:
		return 0, fmt.Errorf("unknown pb.Options_Units value: %v", pbUnit)
	}
}
