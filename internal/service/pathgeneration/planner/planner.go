package planner

import "matching-engine/internal/model"

type PathPlanner interface {
	FindFirstFeasiblePath(offerNode *model.OfferNode, requestNode *model.RequestNode) ([]model.PathPoint, bool, error)
}
