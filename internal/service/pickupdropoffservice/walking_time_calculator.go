package pickupdropoffservice

import (
	"context"
	"fmt"
	"matching-engine/internal/adapter/routing"
	"matching-engine/internal/model"
	"time"
)

type WalkingTimeCalculator struct {
	engine routing.Engine
}

func NewWalkingTimeCalculator(engine routing.Engine) *WalkingTimeCalculator {
	return &WalkingTimeCalculator{engine: engine}
}

func (calculator *WalkingTimeCalculator) ComputeWalkingDurations(ctx context.Context, request *model.Request, pickup, dropoff *model.PathPoint) (pickupWalkingDuration, dropoffWalkingDuration time.Duration, err error) {
	pickupWalkingParams, err := model.NewWalkParams(request.Source(), pickup.Coordinate())
	if err != nil {
		return 0, 0, fmt.Errorf("pickup walking params: %v", err)
	}
	pickupWalkingDuration, err = calculator.engine.ComputeWalkingTime(ctx, pickupWalkingParams)
	if err != nil {
		return 0, 0, fmt.Errorf("pickup walking time: %v", err)
	}

	dropoffWalkingParams, err := model.NewWalkParams(dropoff.Coordinate(), request.Destination())
	if err != nil {
		return 0, 0, fmt.Errorf("dropoff walking params: %v", err)
	}
	dropoffWalkingDuration, err = calculator.engine.ComputeWalkingTime(ctx, dropoffWalkingParams)
	if err != nil {
		return 0, 0, fmt.Errorf("dropoff walking time: %v", err)
	}

	return pickupWalkingDuration, dropoffWalkingDuration, nil
}
