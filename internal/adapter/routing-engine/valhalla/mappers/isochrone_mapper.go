package mappers

import (
	"fmt"
	re "matching-engine/internal/adapter/routing-engine"
	"matching-engine/internal/adapter/routing-engine/valhalla/client/pb"
	"matching-engine/internal/adapter/routing-engine/valhalla/helpers"
	"matching-engine/internal/model"
)

type IsochroneMapper struct{}

var _ re.OperationMapper[
	*model.IsochroneParams,
	*model.Isochrone,
	*pb.Api,
	*pb.Api,
] = IsochroneMapper{}

func (IsochroneMapper) ToTransport(params *model.IsochroneParams) (*pb.Api, error) {
	origin := helpers.CreateLocation(
		params.Origin().Lat(),
		params.Origin().Lng(),
		helpers.DefaultLocationType,
	)

	var (
		costingType pb.Costing_Type
		costing     *pb.Costing
	)

	if params.Profile() == model.Pedestrian {
		costingType = pb.Costing_pedestrian
		costing = helpers.DefaultPedestrianCosting
	} else { // assume auto
		costingType = pb.Costing_auto_
		costing = helpers.DefaultAutoCosting
	}

	return &pb.Api{
		Options: &pb.Options{
			Action:      pb.Options_isochrone,
			Format:      pb.Options_pbf,
			CostingType: costingType,
			Locations:   []*pb.Location{origin},
			Costings: map[int32]*pb.Costing{
				int32(pb.Costing_pedestrian): costing,
			},
			Contours: []*pb.Contour{
				{
					HasDistance: &pb.Contour_Distance{
						Distance: params.Distance().Value(),
					},
				},
			},
			HasPolygons: &pb.Options_Polygons{
				Polygons: false,
			},
		},
	}, nil
}

func (IsochroneMapper) FromTransport(response *pb.Api) (*model.Isochrone, error) {
	isochrone := response.GetIsochrone()
	if isochrone == nil {
		return nil, fmt.Errorf("no isochrone data found")
	}

	intervals := isochrone.GetIntervals()
	if len(intervals) == 0 {
		return nil, fmt.Errorf("no isochrone intervals found")
	}

	// we request a single contour, so we can safely access and use the first one only
	interval := intervals[0]

	value := interval.GetMetricValue()
	unit := interval.GetMetric().String()
	contour, err := model.NewContour(value, unit)
	if err != nil {
		return nil, err
	}

	if len(interval.GetContours()) == 0 || len(interval.GetContours()[0].GetGeometries()) == 0 {
		return nil, fmt.Errorf("no isochrone contours found")
	}

	rawContour := interval.GetContours()[0]
	if len(rawContour.GetGeometries()) == 0 {
		return nil, fmt.Errorf("no isochrone geometries found")
	}

	geometry := rawContour.GetGeometries()[0]
	rawCoords := geometry.GetCoords()

	if len(rawCoords) == 0 || len(rawCoords)%2 != 0 {
		return nil, fmt.Errorf("invalid isochrone coordinates")
	}

	// Valhalla packs each coordinate as an integer = degrees * coordScale.
	// coordScale = 1e5 gives sub‐meter precision (0.00001° ≈ 1 m).
	const coordScale = 1e5

	var ring model.LineString
	ring = make(model.LineString, 0, len(rawCoords)/2)

	// Decode coords: [lat1, lon1, lat2, lon2, …] as integer 1e5‐degree units
	for k := 0; k < len(rawCoords); k += 2 {
		// coords[k] is lat * coordScale, coords[k+1] is lon * coordScale
		lat := float64(rawCoords[k]) / coordScale
		lng := float64(rawCoords[k+1]) / coordScale

		coord, err := model.NewCoordinate(lat, lng)
		if err != nil {
			return nil, fmt.Errorf("invalid coord [%f,%f]: %w", lat, lng, err)
		}

		ring = append(ring, *coord)
	}

	return model.NewIsochrone(contour, &ring)
}
