package dto

// MatchedRequestDTO is a Data Transfer Object for MatchedRequest
type MatchedRequestDTO struct {
	UserID       string   `json:"userId"`
	RequestID    string   `json:"requestId"`
	PickupPoint  PointDTO `json:"pickupPoint"`
	DropoffPoint PointDTO `json:"dropoffPoint"`
}
