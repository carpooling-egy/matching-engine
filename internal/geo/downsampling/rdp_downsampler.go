package downsampling

import (
	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"
	"matching-engine/internal/geo"
	"matching-engine/internal/model"
)

type RDPDownSampler struct {
	eps s1.ChordAngle
}

func NewRDPDownSampler(epsMeters float64) *RDPDownSampler {
	epsAngle := epsMeters / geo.EarthRadiusInMeters
	return &RDPDownSampler{
		eps: s1.ChordAngleFromAngle(s1.Angle(epsAngle)),
	}
}

var _ RouteDownSampler = (*RDPDownSampler)(nil)

func (r *RDPDownSampler) DownSample(
	route model.LineString,
) (model.LineString, error) {
	n := len(route)
	if n < 3 {
		out := make(model.LineString, n)
		copy(out, route)
		return out, nil
	}

	points := convertCoordsToS2Points(route)

	toKeep := make([]bool, n)
	toKeep[0], toKeep[n-1] = true, true

	type segment struct{ start, end int }
	stack := []segment{{0, n - 1}}

	for len(stack) > 0 {
		seg := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		start, end := seg.start, seg.end
		idx, maxChord := findMaxChordIndex(points, start, end)

		if maxChord > r.eps {
			toKeep[idx] = true
			stack = append(stack, segment{start, idx}, segment{idx, end})
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
