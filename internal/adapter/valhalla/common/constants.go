package common

import (
	pb2 "matching-engine/internal/adapter/valhalla/client/pb"
	"matching-engine/internal/app/config"
)

const (
	DefaultUnit                  = pb2.Options_kilometers
	DefaultResponseFormat        = pb2.Options_pbf
	DefaultShapeFormat           = pb2.ShapeFormat_polyline6
	DefaultTimeFormat            = "2006-01-02T15:04"
	DefaultPedestrianMaxDistance = 100_000
	DefaultSearchRadiusInMeters  = 100
)

var DefaultAutoCosting = &pb2.Costing{
	HasOptions: &pb2.Costing_Options_{
		Options: &pb2.Costing_Options{
			HasShortest: &pb2.Costing_Options_Shortest{
				Shortest: true,
			},
			FixedSpeed: uint32(config.GetEnvFloat("FIXED_SPEED_KMH", 27.0)),
		},
	},
}

var DefaultPedestrianCosting = &pb2.Costing{
	Type: pb2.Costing_pedestrian,
	HasOptions: &pb2.Costing_Options_{
		Options: &pb2.Costing_Options{
			HasShortest: &pb2.Costing_Options_Shortest{
				Shortest: true,
			},
			HasMaxDistance: &pb2.Costing_Options_MaxDistance{
				MaxDistance: DefaultPedestrianMaxDistance,
			},
		},
	},
}
