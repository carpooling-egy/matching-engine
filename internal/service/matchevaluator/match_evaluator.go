package matchevaluator

import (
	"fmt"
	"matching-engine/internal/model"
	"matching-engine/internal/service/checker"
	"matching-engine/internal/service/pathgeneration/planner"
)

type MatchEvaluator struct {
	pathPlanner       planner.PathPlanner
	preferenceChecker checker.Checker
}

func NewMatchEvaluator(pathPlanner planner.PathPlanner, preferenceChecker checker.Checker) *MatchEvaluator {
	return &MatchEvaluator{
		pathPlanner:       pathPlanner,
		preferenceChecker: preferenceChecker,
	}
}

func (m *MatchEvaluator) Evaluate(offerNode *model.OfferNode, requestNode *model.RequestNode) ([]model.PathPoint, error) {

	offer := offerNode.Offer()
	request := requestNode.Request()

	preferenceChecker := checker.NewPreferenceChecker()
	valid, err := preferenceChecker.Check(offerNode.Offer(), requestNode.Request())
	if err != nil || !valid {
		return nil, fmt.Errorf("preference check failed for offer %s and request %s: %w", offer.ID(), request.ID(), err)
	}

	path, isFeasible, err := m.pathPlanner.FindFirstFeasiblePath(offerNode, requestNode)
	if err != nil || !isFeasible || len(path) < 2 || path == nil {
		return nil, fmt.Errorf("failed to find feasible path for offer %s and request %s: %w", offer.ID(), request.ID(), err)
	}

	return path, nil
}
