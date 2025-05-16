package earlypruning

import (
	"matching-engine/internal/model"
	"matching-engine/internal/service/iterator"
)

type CandidateGenerator interface {
	// GenerateCandidates generates candidates for a given offer and requests
	GenerateCandidates(offers []*model.Offer, requests []*model.Request) (*iterator.CandidateIterator, error)
}
