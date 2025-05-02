package service

import (
	"fmt"
	"matching-engine/internal/model"
	"matching-engine/internal/service/EarlyPruning"
	"matching-engine/internal/service/PathGenerator"
)

type Matcher struct {
	offerNodes             []*model.OfferNode
	requestNodes           []*model.RequestNode
	potentialOfferRequests map[*model.OfferNode][]*model.RequestNode
	pathPlanner            PathGenerator.PathPlanner
	candidateGenerator     EarlyPruning.CandidateGenerator
}

// NewMatcher creates a new Matcher instance
func NewMatcher(planner PathGenerator.PathPlanner, generator EarlyPruning.CandidateGenerator) *Matcher {
	return &Matcher{
		offerNodes:             make([]*model.OfferNode, 0),
		requestNodes:           make([]*model.RequestNode, 0),
		potentialOfferRequests: make(map[*model.OfferNode][]*model.RequestNode),
		pathPlanner:            planner,
		candidateGenerator:     generator,
	}
}

func (matcher *Matcher) Match(offers []*model.Offer, requests []*model.Request) ([]model.MatchingResult, error) {
	// TODO: Implement the matching logic
	return nil, fmt.Errorf("not implemented")
}
