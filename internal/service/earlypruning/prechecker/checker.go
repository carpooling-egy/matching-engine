package prechecker

import "matching-engine/internal/model"

type Checker interface {
	// Check checks if the given request can be matched with the offer
	Check(offer *model.Offer, request *model.Request) (bool, error)
}
