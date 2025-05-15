package earlypruning

import (
	"fmt"
	"matching-engine/internal/model"
	"matching-engine/internal/service/earlypruning/prechecker"
	"matching-engine/internal/service/iterator"
)

type PreChecksCandidateGenerator struct {
	// PreChecksCandidateGenerator is a candidate generator that performs pre-checks
	checker prechecker.Checker
}

func NewPreChecksCandidateGenerator(checker prechecker.Checker) *PreChecksCandidateGenerator {
	return &PreChecksCandidateGenerator{
		checker: checker,
	}
}

func (g *PreChecksCandidateGenerator) GenerateCandidates(offers []*model.Offer, requests []*model.Request) (*iterator.CandidateIterator, error) {
	if len(offers) == 0 || len(requests) == 0 {
		return nil, fmt.Errorf("no offers or requests")
	}
	candidateIterator := iterator.NewCandidateIterator(offers, requests, g.checker)
	if candidateIterator == nil {
		return nil, fmt.Errorf("failed to create candidate iterator")
	}
	return candidateIterator, nil
}
