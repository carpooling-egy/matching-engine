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

func NewMatchEvaluator(pathPlanner planner.PathPlanner, preferenceChecker checker.Checker) Evaluator {
	return &MatchEvaluator{
		pathPlanner:       pathPlanner,
		preferenceChecker: preferenceChecker,
	}
}

func (m *MatchEvaluator) Evaluate(offerNode *model.OfferNode, requestNode *model.RequestNode) ([]model.PathPoint, bool, error) {

	offer := offerNode.Offer()
	request := requestNode.Request()

	valid, err := m.preferenceChecker.Check(offerNode.Offer(), requestNode.Request())
	if err != nil {
		return nil, false, fmt.Errorf("preference check failed for offer %s and request %s: %w", offer.ID(), request.ID(), err)
	}
	if !valid {
		return nil, false, nil
	}

	path, isFeasible, err := m.pathPlanner.FindFirstFeasiblePath(offerNode, requestNode)
	if err != nil {
		return nil, false, fmt.Errorf("failed to find feasible path for offer %s and request %s: %w", offer.ID(), request.ID(), err)
	}

	if len(path) < 2 || path == nil {
		return nil, false, fmt.Errorf("path is empty or has less than 2 points for offer %s and request %s", offer.ID(), request.ID())
	}

	if !isFeasible {
		return nil, false, nil
	}

	return path, true, nil
}
