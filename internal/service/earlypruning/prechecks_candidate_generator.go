package earlypruning

import (
	"fmt"
	"matching-engine/internal/model"
	"matching-engine/internal/service/checker"
)

type PreChecksCandidateGenerator struct {
	// PreChecksCandidateGenerator is a candidate generator that performs pre-checks
	checker checker.Checker
}

func NewPreChecksCandidateGenerator(checker checker.Checker) *PreChecksCandidateGenerator {
	return &PreChecksCandidateGenerator{
		checker: checker,
	}
}

func (g *PreChecksCandidateGenerator) GenerateCandidates(offers []*model.Offer, requests []*model.Request) (*CandidateIterator, error) {
	if len(offers) == 0 || len(requests) == 0 {
		return nil, fmt.Errorf("no offers or requests")
	}
	candidateIterator := NewCandidateIterator(offers, requests, g.checker)
	if candidateIterator == nil {
		return nil, fmt.Errorf("failed to create candidate iterator")
	}
	return candidateIterator, nil
}
