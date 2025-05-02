package dto

// MatchedRequestDTO is a Data Transfer Object for MatchedRequest
type MatchedRequestDTO struct {
	RequestID     string     `json:"requestId"`
	PickupPoints  []PointDTO `json:"pickupPoints"`
	DropoffPoints []PointDTO `json:"dropoffPoints"`
}
