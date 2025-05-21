package timematrix

import "matching-engine/internal/model"

type Populator interface {
	Populate(offer *model.OfferNode, requestNodes []*model.RequestNode) error
}
