package helpers

import (
	"fmt"
	"matching-engine/internal/adapter/routing-engine/valhalla/client/pb"
	"matching-engine/internal/model"
)

func CreateLocation(lat, lng float64, locType pb.Location_Type) *pb.Location {
	return &pb.Location{
		Type: DefaultLocationType,
		Ll: &pb.LatLng{
			HasLat: &pb.LatLng_Lat{Lat: lat},
			HasLng: &pb.LatLng_Lng{Lng: lng},
		},
	}
}

func ToDomainDistanceUnit(pbUnit pb.Options_Units) (model.DistanceUnit, error) {
	switch pbUnit {
	case pb.Options_kilometers:
		return model.Kilometer, nil
	case pb.Options_miles:
		return model.Mile, nil
	default:
		return 0, fmt.Errorf("unknown pb.Options_Units value: %v", pbUnit)
	}
}
