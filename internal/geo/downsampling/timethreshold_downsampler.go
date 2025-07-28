package downsampling

import (
	"errors"
	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"
	"github.com/rs/zerolog/log"
	"matching-engine/internal/geo"
	"matching-engine/internal/model"
	"time"
)

const DefaultInterval = 10 * time.Second

type TimeThresholdDownSampler struct {
	intervalAngle s1.Angle
}

type TimeThresholdDownSamplerOption func(*TimeThresholdDownSampler)

func WithInterval(interval time.Duration) TimeThresholdDownSamplerOption {
	return func(t *TimeThresholdDownSampler) {
		t.intervalAngle = durationToAngle(interval)
	}
}

func NewTimeThresholdDownSampler(opts ...TimeThresholdDownSamplerOption) RouteDownSampler {
	t := &TimeThresholdDownSampler{
		intervalAngle: durationToAngle(DefaultInterval),
	}
	for _, opt := range opts {
		opt(t)
	}
	return t
}

func (t *TimeThresholdDownSampler) DownSample(
	route model.LineString,
) (model.LineString, error) {
	log.Debug().Msg("TimeThresholdDownSampler.DownSample called")
	n := len(route)
	if n == 0 {
		return nil, errors.New("empty route")
	}
	if n == 1 {
		return route, nil
	}

	points := convertCoordsToS2LatLngs(route)

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

func convertCoordsToS2LatLngs(route model.LineString) []s2.LatLng {
	points := make([]s2.LatLng, len(route))
	for i, pt := range route {
		points[i] = s2.LatLngFromDegrees(pt.Lat(), pt.Lng())
	}
	return points
}

func durationToAngle(interval time.Duration) s1.Angle {
	distance := geo.WalkingSpeedMPS * interval.Seconds()
	return s1.Angle(distance / geo.EarthRadiusInMeters)
}
