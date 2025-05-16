package pickupdropoffservice

import (
	"matching-engine/internal/model"
	"matching-engine/internal/service/pickupdropoffservice/pickupdropoffcache"
)

// PickupDropoffSelectorInterface defines the interface needed by DetourTimeChecker
type PickupDropoffSelectorInterface interface {
	GetPickupDropoffPointsAndDurations(request *model.Request, offer *model.Offer) (*pickupdropoffcache.Value, error)
}
