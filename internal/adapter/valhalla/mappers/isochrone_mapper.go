package mappers

import (
	"fmt"
	re "matching-engine/internal/adapter/routing-engine"
	"matching-engine/internal/adapter/valhalla/client/pb"
	"matching-engine/internal/adapter/valhalla/common"
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
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}

	origin := common.CreateLocation(params.Origin().Lat(), params.Origin().Lng())

	var (
		costingType pb.Costing_Type
		costing     *pb.Costing
	)

	if params.Profile() == model.Pedestrian {
		costingType = pb.Costing_pedestrian
		costing = common.DefaultPedestrianCosting
	} else { // assume auto
		costingType = pb.Costing_auto_
		costing = common.DefaultAutoCosting
	}

	return &pb.Api{
		Options: &pb.Options{
			Action:      pb.Options_isochrone,
			Format:      pb.Options_pbf,
			CostingType: costingType,
			Costings: map[int32]*pb.Costing{
				int32(pb.Costing_pedestrian): costing,
			},
			Locations: []*pb.Location{origin},
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
	if response == nil {
		return nil, fmt.Errorf("response cannot be nil")
	}

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

	fmt.Println(rawCoords)

	if len(rawCoords) == 0 || len(rawCoords)%2 != 0 {
		return nil, fmt.Errorf("invalid isochrone coordinates")
	}

	// Valhalla packs each coordinate as an integer = degrees * coordScale.
	// coordScale = 1e6 gives sub‐meter precision (0.000001° ≈ 1 m).
	const coordScale = 1e6

	var ring model.LineString
	ring = make(model.LineString, 0, len(rawCoords)/2)

	// Decode coords: [lon1, lat1, lon2, lat2, …] as integer 1e5‐degree units
	for k := 0; k < len(rawCoords); k += 2 {
		// coords[k] is lng * coordScale, coords[k+1] is lat * coordScale
		lng := float64(rawCoords[k]) / coordScale
		lat := float64(rawCoords[k+1]) / coordScale

		coord, err := model.NewCoordinate(lat, lng)
		if err != nil {
			return nil, fmt.Errorf("invalid coord [%f,%f]: %w", lat, lng, err)
		}

		ring = append(ring, *coord)
	}

	return model.NewIsochrone(contour, &ring)
}
