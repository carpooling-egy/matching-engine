package valhalla

import "matching-engine/internal/adapter/routing-engine/valhalla/client/pb"

const (
	DefaultUnit                  = pb.Options_kilometers
	DefaultFormat                = pb.Options_pbf
	DefaultShapeFormat           = pb.ShapeFormat_polyline6
	DefaultLocationType          = pb.Location_kBreak
	DefaultPedestrianMaxDistance = 100_000
)

var DefaultAutoCosting = &pb.Costing{
	HasOptions: &pb.Costing_Options_{
		Options: &pb.Costing_Options{
			HasShortest: &pb.Costing_Options_Shortest{
				Shortest: true,
			},
		},
	},
}

var DefaultPedestrianCosting = &pb.Costing{
	Type: pb.Costing_pedestrian,
	HasOptions: &pb.Costing_Options_{
		Options: &pb.Costing_Options{
			HasShortest: &pb.Costing_Options_Shortest{
				Shortest: true,
			},
			HasMaxDistance: &pb.Costing_Options_MaxDistance{
				MaxDistance: DefaultPedestrianMaxDistance,
			},
		},
	},
}
