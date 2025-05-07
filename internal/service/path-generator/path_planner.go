package path_generator

import "matching-engine/internal/model"

type PathPlanner interface {
	findFirstFeasiblePath(offerNode *model.OfferNode, requestNode *model.RequestNode) ([]*model.PathPoint, bool, error)
}
