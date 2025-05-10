package downsampling

import (
	"errors"
	"github.com/golang/geo/s2"
	"matching-engine/internal/geo"
	"matching-engine/internal/model"
	"time"
)

type TimeThresholdDownSampler struct {
	interval time.Duration
}

func NewTimeThresholdDownSampler(interval time.Duration) *TimeThresholdDownSampler {
	return &TimeThresholdDownSampler{
		interval: interval,
	}
}

func (t *TimeThresholdDownSampler) DownSample(
	route model.LineString,
) (model.LineString, error) {

	if len(route) == 0 {
		return nil, errors.New("empty route")
	}
	if len(route) == 1 {
		return route, nil
	}

	result := make(model.LineString, 0, len(route)/2)
	result = append(result, route[0])

	accumulatedTime := 0.0
	thresholdSec := t.interval.Seconds()

	for i := 1; i < len(route); i++ {
		prev := route[i-1]
		curr := route[i]

		distMeters := geodesicDistance(prev, curr)
		segmentTime := distMeters / geo.WalkingSpeedMPS

		accumulatedTime += segmentTime

		if accumulatedTime >= thresholdSec {
			result = append(result, curr)
			accumulatedTime -= thresholdSec
		}
	}

	if len(result) == 0 || !result[len(result)-1].Equal(&route[len(route)-1]) {
		result = append(result, route[len(route)-1])
	}

	return result, nil
}

func geodesicDistance(a, b model.Coordinate) float64 {
	p1 := s2.LatLngFromDegrees(a.Lat(), a.Lng())
	p2 := s2.LatLngFromDegrees(b.Lat(), b.Lng())
	return p1.Distance(p2).Radians() * geo.EarthRadiusInMeters
}
