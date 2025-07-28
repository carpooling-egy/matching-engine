package pickupdropoffcache

import (
	"matching-engine/internal/model"
)

type Value struct {
	pickup  *model.PathPoint
	dropoff *model.PathPoint
}

func NewValue(pickup, dropoff *model.PathPoint) *Value {
	return &Value{
		pickup:  pickup,
		dropoff: dropoff,
	}
}

func (v *Value) Pickup() *model.PathPoint {
	return v.pickup
}

func (v *Value) Dropoff() *model.PathPoint {
	return v.dropoff
}

func (v *Value) SetPickup(pickup *model.PathPoint) {
	v.pickup = pickup
}

func (v *Value) SetDropoff(dropoff *model.PathPoint) {
	v.dropoff = dropoff
}
