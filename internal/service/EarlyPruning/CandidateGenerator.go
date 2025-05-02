package EarlyPruning

import "matching-engine/internal/model"

type CandidateGenerator interface {
	// GenerateCandidates generates candidates for a given offer and requests
	GenerateCandidates(offerID, requestID string) (model.MatchCandidate, error)
}
