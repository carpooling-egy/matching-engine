package dto

// PointDTO is a Data Transfer Object for Point
type PointDTO struct {
	RequestID string        `json:"requestId"`
	Point     CoordinateDTO `json:"point"`
	Time      string        `json:"time"`
	PointType int           `json:"pointType"`
}
