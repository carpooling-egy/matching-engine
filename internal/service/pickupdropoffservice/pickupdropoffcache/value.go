package pickupdropoffcache

import (
	"matching-engine/internal/model"
	"time"
)

type Value struct {
	pickup                 *model.PathPoint
	dropoff                *model.PathPoint
	pickupWalkingDuration  time.Duration
	dropoffWalkingDuration time.Duration
}

func NewValue(pickup, dropoff *model.PathPoint, pickupWalkingDuration, dropoffWalkingDuration time.Duration) *Value {
	return &Value{
		pickup:                 pickup,
		dropoff:                dropoff,
		pickupWalkingDuration:  pickupWalkingDuration,
		dropoffWalkingDuration: dropoffWalkingDuration,
	}
}

func (v *Value) Pickup() *model.PathPoint {
	return v.pickup
}

func (v *Value) Dropoff() *model.PathPoint {
	return v.dropoff
}

func (v *Value) PickupWalkingDuration() time.Duration {
	return v.pickupWalkingDuration
}

func (v *Value) DropoffWalkingDuration() time.Duration {
	return v.dropoffWalkingDuration
}

func (v *Value) SetPickup(pickup *model.PathPoint) {
	v.pickup = pickup
}

func (v *Value) SetDropoff(dropoff *model.PathPoint) {
	v.dropoff = dropoff
}

func (v *Value) SetPickupWalkingDuration(duration time.Duration) {
	v.pickupWalkingDuration = duration
}

func (v *Value) SetDropoffWalkingDuration(duration time.Duration) {
	v.dropoffWalkingDuration = duration
}
