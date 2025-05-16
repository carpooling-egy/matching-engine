package iterator

import (
	"fmt"
	"iter"
	"matching-engine/internal/model"
	"matching-engine/internal/service/earlypruning/prechecker"
)

type CandidateIterator struct {
	offers   []*model.Offer
	requests []*model.Request
	checker  prechecker.Checker
}

func NewCandidateIterator(offers []*model.Offer, requests []*model.Request, checker prechecker.Checker) *CandidateIterator {
	return &CandidateIterator{
		offers:   offers,
		requests: requests,
		checker:  checker,
	}
}

func (ci *CandidateIterator) Candidates() iter.Seq2[*model.MatchCandidate, error] {
	return func(yield func(*model.MatchCandidate, error) bool) {
		for _, offer := range ci.offers {
			for _, request := range ci.requests {
				// Check if the offer and request can be matched
				isPotential, err := ci.checker.Check(offer, request)
				if err != nil {
					if !yield(nil, fmt.Errorf("checker failed: %w", err)) {
						// If the yield function returns false, stop iterating
						return
					}
					continue
				}
				if isPotential {
					if !yield(model.NewMatchCandidate(request, offer), nil) {
						// If the yield function returns false, stop iterating
						return
					}
				}
			}
		}
	}
}
