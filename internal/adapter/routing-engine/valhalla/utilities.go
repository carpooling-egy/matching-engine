package valhalla

import "matching-engine/internal/adapter/routing-engine/valhalla/client/pb"

func CreateLocation(lat, lng float64, locType pb.Location_Type) *pb.Location {
	return &pb.Location{
		Type: DefaultLocationType,
		Ll: &pb.LatLng{
			HasLat: &pb.LatLng_Lat{Lat: lat},
			HasLng: &pb.LatLng_Lng{Lng: lng},
		},
	}
}
