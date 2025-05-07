package dto

// MatchingResultDTO is a Data Transfer Object for JSON serialization
type MatchingResultDTO struct {
	UserID                  string              `json:"userId"`
	OfferID                 string              `json:"offerId"`
	AssignedMatchedRequests []MatchedRequestDTO `json:"assignedMatchedRequests"`
	Path                    []PointDTO          `json:"path"`
}
