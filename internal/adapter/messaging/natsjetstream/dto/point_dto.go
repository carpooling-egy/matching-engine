package dto

// PointDTO is a Data Transfer Object for Point
type PointDTO struct {
	OwnerType              string        `json:"ownerType"`
	OwnerID                string        `json:"ownerID"`
	Point                  CoordinateDTO `json:"point"`
	Time                   string        `json:"time"`
	PointType              string        `json:"pointType"`
	WalkingDurationMinutes int           `json:"walkingDurationMinutes"`
}
