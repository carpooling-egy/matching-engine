package service

import "matching-engine/internal/model"

type Matcher struct {
	offerNodes             []*model.OfferNode
	requestNodes           []*model.RequestNode
	potentialOfferRequests map[*model.OfferNode][]*model.RequestNode
}

func (matcher Matcher) match(offers []model.Offer, requests []model.Request) []model.MatchingResult {
	return nil

}
