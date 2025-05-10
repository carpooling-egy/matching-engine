package downsampling

import (
	"errors"
	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"
	"matching-engine/internal/geo"
	"matching-engine/internal/model"
	"time"
)

type TimeThresholdDownSampler struct {
	intervalAngle s1.Angle
}

func NewTimeThresholdDownSampler(interval time.Duration) *TimeThresholdDownSampler {
	distMeters := geo.WalkingSpeedMPS * interval.Seconds()
	angleRad := distMeters / geo.EarthRadiusInMeters
	return &TimeThresholdDownSampler{
		intervalAngle: s1.Angle(angleRad),
	}
}

func (t *TimeThresholdDownSampler) DownSample(
	route model.LineString,
) (model.LineString, error) {
	n := len(route)
	if n == 0 {
		return nil, errors.New("empty route")
	}
	if n == 1 {
		return route, nil
	}

	points := make([]s2.LatLng, n)
	for i, pt := range route {
		points[i] = s2.LatLngFromDegrees(pt.Lat(), pt.Lng())
	}

	out := make(model.LineString, 0, n/2)
	out = append(out, route[0])

	var accAngle s1.Angle
	for i := 1; i < n; i++ {
		prev := points[i-1]
		cur := points[i]
		segAngle := prev.Distance(cur)
		accAngle += segAngle
		if accAngle >= t.intervalAngle {
			out = append(out, route[i])
			accAngle -= t.intervalAngle
		}
	}

	if !out[len(out)-1].Equal(&route[n-1]) {
		out = append(out, route[n-1])
	}

	return out, nil
}
