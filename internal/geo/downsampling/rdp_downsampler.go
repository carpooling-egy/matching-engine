package downsampling

import (
	"github.com/golang/geo/s2"
	"matching-engine/internal/geo"
	"matching-engine/internal/model"
)

type RDPDownSampler struct {
	eps float64
}

func NewRDPDownSampler(epsInMeters float64) *RDPDownSampler {
	return &RDPDownSampler{
		eps: epsInMeters,
	}
}

var _ RouteDownSampler = (*RDPDownSampler)(nil)

func (r *RDPDownSampler) DownSample(
	route model.LineString,
) (model.LineString, error) {
	if len(route) < 3 {
		out := make(model.LineString, len(route))
		copy(out, route)
		return out, nil
	}

	return r.rdp(route), nil
}

func (r *RDPDownSampler) rdp(points model.LineString) model.LineString {
	if len(points) < 3 {
		out := make(model.LineString, len(points))
		copy(out, points)
		return out
	}

	start, end := points[0], points[len(points)-1]
	maxIdx, maxD := 0, 0.0
	for i := 1; i < len(points)-1; i++ {
		if d := perpendicularDistance(points[i], start, end); d > maxD {
			maxD, maxIdx = d, i
		}
	}
	if maxD <= r.eps {
		return model.LineString{start, end}
	}
	left := r.rdp(points[:maxIdx+1])
	right := r.rdp(points[maxIdx:])
	return append(left[:len(left)-1], right...)
}

func perpendicularDistance(p, start, end model.Coordinate) float64 {
	pPoint := s2.PointFromLatLng(s2.LatLngFromDegrees(p.Lat(), p.Lng()))
	startPoint := s2.PointFromLatLng(s2.LatLngFromDegrees(start.Lat(), start.Lng()))
	endPoint := s2.PointFromLatLng(s2.LatLngFromDegrees(end.Lat(), end.Lng()))
	distanceRad := s2.DistanceFromSegment(pPoint, startPoint, endPoint)
	return distanceRad.Radians() * geo.EarthRadiusInMeters
}
