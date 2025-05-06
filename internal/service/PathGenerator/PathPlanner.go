package PathGenerator

import "matching-engine/internal/model"

type PathPlanner interface {
	findFirstFeasiblePath(offerNode *model.OfferNode, requestNode *model.RequestNode) ([]*model.Point, bool, error)
}
