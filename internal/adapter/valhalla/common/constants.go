package common

import (
	pb2 "matching-engine/internal/adapter/valhalla/client/pb"
)

const (
	DefaultUnit                  = pb2.Options_kilometers
	DefaultResponseFormat        = pb2.Options_pbf
	DefaultShapeFormat           = pb2.ShapeFormat_polyline6
	DefaultTimeFormat            = "2006-01-02T15:04"
	DefaultLocationType          = pb2.Location_kBreak
	DefaultPedestrianMaxDistance = 100_000
	DefaultSearchRadiusInMeters  = 100
)

var DefaultAutoCosting = &pb2.Costing{
	HasOptions: &pb2.Costing_Options_{
		Options: &pb2.Costing_Options{
			HasShortest: &pb2.Costing_Options_Shortest{
				Shortest: true,
			},
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
