package pickupdropoffservice

import "matching-engine/internal/model"

type PickupDropoffGenerator interface {
	// GeneratePickupDropoffPoints generates the best pickup and dropoff points for a given request and offer
	GeneratePickupDropoffPoints(request *model.Request, offer *model.Offer) (*model.PathPoint, *model.PathPoint, error)
}
