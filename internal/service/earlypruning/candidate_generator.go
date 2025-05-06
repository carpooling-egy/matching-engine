package earlypruning

import "matching-engine/internal/model"

type CandidateGenerator interface {
	// GenerateCandidates generates candidates for a given offer and requests
	GenerateCandidates(offers []*model.Offer, requests []*model.Request) (model.MatchCandidate, error)
}
