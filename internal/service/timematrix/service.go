package timematrix

import (
	"fmt"
	"matching-engine/internal/model"
	"matching-engine/internal/service/timematrix/cache"
	"time"
)

type TimeMatrixService struct {
	timeMatrixSelector Selector
}

func NewTimeMatrixService(selector Selector) *TimeMatrixService {
	return &TimeMatrixService{
		timeMatrixSelector: selector,
	}
}

func (s *TimeMatrixService) getTravelDuration(matrix *cache.PathPointMappedTimeMatrix, from, to model.PathPointID) (time.Duration, error) {
	fromIdx, fromOk := matrix.PointIdToIndex()[from]
	toIdx, toOk := matrix.PointIdToIndex()[to]

	if !fromOk || !toOk {
		return 0, fmt.Errorf("invalid path point ID: from=%v, to=%v", from, to)
	}

	if fromIdx >= len(matrix.TimeMatrix()) || toIdx >= len(matrix.TimeMatrix()[fromIdx]) {
		return 0, fmt.Errorf("index out of bounds: fromIdx=%d, toIdx=%d", fromIdx, toIdx)
	}

	return matrix.TimeMatrix()[fromIdx][toIdx], nil
}

func (s *TimeMatrixService) GetCumulativeTravelDurations(offer *model.OfferNode, pathPoints []model.PathPoint) ([]time.Duration, error) {
	if len(pathPoints) < 2 {
		return nil, fmt.Errorf("pathPointIDs must contain at least two points")
	}

	matrix, err := s.timeMatrixSelector.GetTimeMatrix(offer)
	if err != nil {
		return nil, err
	}

	cumulativeDuration := make([]time.Duration, len(pathPoints))
	cumulativeDuration[0] = 0
	for i := 0; i < len(pathPoints)-1; i++ {
		duration, err := s.getTravelDuration(matrix, pathPoints[i].ID(), pathPoints[i+1].ID())
		if err != nil {
			return nil, err
		}
		cumulativeDuration[i+1] = cumulativeDuration[i] + duration
	}

	return cumulativeDuration, nil
}

func (s *TimeMatrixService) GetCumulativeTravelTimes(offer *model.OfferNode, pathPoints []model.PathPoint) ([]time.Time, error) {
	if len(pathPoints) < 2 {
		return nil, fmt.Errorf("pathPointIDs must contain at least two points")
	}

	matrix, err := s.timeMatrixSelector.GetTimeMatrix(offer)
	if err != nil {
		return nil, err
	}

	cumulativeTimes := make([]time.Time, len(pathPoints))
	cumulativeTimes[0] = offer.Offer().DepartureTime()
	for i := 0; i < len(pathPoints)-1; i++ {
		duration, err := s.getTravelDuration(matrix, pathPoints[i].ID(), pathPoints[i+1].ID())
		if err != nil {
			return nil, err
		}
		cumulativeTimes[i+1] = cumulativeTimes[i].Add(duration)
	}

	return cumulativeTimes, nil
}

func (s *TimeMatrixService) GetTravelDuration(offer *model.OfferNode, from, to model.PathPointID) (time.Duration, error) {
	matrix, err := s.timeMatrixSelector.GetTimeMatrix(offer)

	if err != nil {
		return 0, err
	}

	return s.getTravelDuration(matrix, from, to)
}
