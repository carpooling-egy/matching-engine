package downsampling

import (
	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"
	"matching-engine/internal/collections"
	"matching-engine/internal/geo"
	"matching-engine/internal/model"
)

const DefaultEpsilonMeters = 10.0
const DefaultMinRoutePoints = 2

type RDPDownSampler struct {
	eps s1.ChordAngle
}

type RDPDownSamplerOption func(*RDPDownSampler)

func WithEpsilonMeters(epsMeters float64) RDPDownSamplerOption {
	return func(s *RDPDownSampler) {
		s.eps = metersToChord(epsMeters)
	}
}

func NewRDPDownSampler(opts ...RDPDownSamplerOption) *RDPDownSampler {
	s := &RDPDownSampler{
		eps: metersToChord(DefaultEpsilonMeters),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

var _ RouteDownSampler = (*RDPDownSampler)(nil)

func (r *RDPDownSampler) DownSample(
	route model.LineString,
) (model.LineString, error) {
	n := len(route)
	if n <= DefaultMinRoutePoints {
		out := make(model.LineString, n)
		copy(out, route)
		return out, nil
	}

	points := convertCoordsToS2Points(route)

	toKeep := make([]bool, n)
	toKeep[0], toKeep[n-1] = true, true

	stack := collections.NewStack[collections.Tuple2[int, int]]()
	stack.Push(collections.NewTuple2(0, n-1))

	for !stack.IsEmpty() {
		seg, err := stack.Pop()
		if err != nil {
			break
		}

		start, end := seg.First, seg.Second
		idx, maxChord := findMaxChordIndex(points, start, end)

		if maxChord > r.eps {
			toKeep[idx] = true
			stack.Push(collections.NewTuple2(start, idx))
			stack.Push(collections.NewTuple2(idx, end))
		}
	}

	return buildOutput(route, toKeep), nil
}

func convertCoordsToS2Points(route model.LineString) []s2.Point {
	points := make([]s2.Point, len(route))
	for i, pt := range route {
		points[i] = s2.PointFromLatLng(
			s2.LatLngFromDegrees(pt.Lat(), pt.Lng()),
		)
	}
	return points
}

func findMaxChordIndex(points []s2.Point, start, end int) (int, s1.ChordAngle) {
	var maxChord s1.ChordAngle
	idx := -1

	for i := start + 1; i < end; i++ {
		ang := s2.DistanceFromSegment(
			points[i], points[start], points[end],
		)
		chord := s1.ChordAngleFromAngle(ang)
		if chord > maxChord {
			maxChord, idx = chord, i
		}
	}

	return idx, maxChord
}

func buildOutput(route model.LineString, toKeep []bool) model.LineString {
	count := 0
	for _, keep := range toKeep {
		if keep {
			count++
		}
	}

	out := make(model.LineString, count)
	j := 0
	for i, keep := range toKeep {
		if keep {
			out[j] = route[i]
			j++
		}
	}
	return out
}

func metersToChord(epsMeters float64) s1.ChordAngle {
	epsAngle := epsMeters / geo.EarthRadiusInMeters
	return s1.ChordAngleFromAngle(s1.Angle(epsAngle))
}
