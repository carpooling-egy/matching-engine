package common

import (
	"fmt"
	"matching-engine/internal/adapter/valhalla/client/pb"
	"matching-engine/internal/model"
)

func CreateLocation(lat, lng float64) *pb.Location {
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

		return model.DistanceUnitKilometer, nil
	case pb.Options_miles:
		return model.DistanceUnitMile, nil
	default:
		return "", fmt.Errorf("unknown pb.Options_Units value: %v", pbUnit)
	}
}

func MapProfileToCosting(profile model.Profile) (pb.Costing_Type, *pb.Costing, error) {
	switch profile {
	case model.ProfilePedestrian:
		return pb.Costing_pedestrian, DefaultPedestrianCosting, nil
	case model.ProfileAuto:
		return pb.Costing_auto_, DefaultAutoCosting, nil
	}
	return 0, nil, fmt.Errorf("unsupported profile: %s", profile)
}
