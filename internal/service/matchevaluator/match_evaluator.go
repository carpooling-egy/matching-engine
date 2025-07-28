package matchevaluator

import (
	"fmt"
	"matching-engine/internal/model"
	"matching-engine/internal/service/checker"
	"matching-engine/internal/service/pathgeneration/planner"
	"matching-engine/internal/service/timematrix"
)

type MatchEvaluator struct {
	pathPlanner                                           planner.PathPlanner
	preferenceChecker                                     checker.Checker
	timeMatrixCacheWithDriverOfferIdAndRequestIdPopulator *timematrix.CacheWithOfferIdRequestIdPopulator
}

func NewMatchEvaluator(pathPlanner planner.PathPlanner, preferenceChecker checker.Checker, timeMatrixCacheWithDriverOfferIdAndRequestIdPopulator *timematrix.CacheWithOfferIdRequestIdPopulator) Evaluator {
	return &MatchEvaluator{
		pathPlanner:       pathPlanner,
		preferenceChecker: preferenceChecker,
		timeMatrixCacheWithDriverOfferIdAndRequestIdPopulator: timeMatrixCacheWithDriverOfferIdAndRequestIdPopulator,
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

	// Populate the time matrix cache with offer ID and request ID
	err = m.timeMatrixCacheWithDriverOfferIdAndRequestIdPopulator.Populate(offerNode, []*model.RequestNode{requestNode})
	if err != nil {
		return nil, false, fmt.Errorf("failed to populate time matrix cache for offer %s and request %s: %w", offer.ID(), request.ID(), err)
	}

	// Find the first feasible path using the path planner
	path, isFeasible, err := m.pathPlanner.FindFirstFeasiblePath(offerNode, requestNode)
	if err != nil {
		return nil, false, fmt.Errorf("failed to find feasible path for offer %s and request %s: %w", offer.ID(), request.ID(), err)
	}

	if !isFeasible {
		return nil, false, nil
	}

	if len(path) < 2 || path == nil {
		return nil, false, fmt.Errorf("path is empty or has less than 2 points for offer %s and request %s", offer.ID(), request.ID())
	}

	m.timeMatrixCacheWithDriverOfferIdAndRequestIdPopulator.RemoveEntry(offerNode, []*model.RequestNode{requestNode})

	return path, true, nil
}
