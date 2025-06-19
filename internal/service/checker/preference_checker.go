package checker

import (
	"fmt"
	"matching-engine/internal/model"
)

type PreferenceChecker struct {
}

// NewPreferenceChecker creates a new PreferenceChecker
func NewPreferenceChecker() Checker {
	return &PreferenceChecker{}
}

// Check checks if the given request can be matched with the offer
func (pc *PreferenceChecker) Check(offer *model.Offer, request *model.Request) (bool, error) {
	if offer == nil || request == nil {
		return false, fmt.Errorf("offer or request is nil")
	}
	// Check if the request and offer have matching preferences
	matched, err := pc.checkPreferenceMatch(request.Preferences(), offer.Preferences())
	if err != nil {
		return false, fmt.Errorf("failed to check preference match: %w", err)
	}
	if !matched {
		return false, nil
	}
	// Check if the request preferences match with the matched requests in the offer
	for _, matchedRequest := range offer.MatchedRequests() {
		if matchedRequest == nil {
			continue
		}
		matched, err := pc.checkPreferenceMatch(request.Preferences(), matchedRequest.Preferences())
		if err != nil {
			return false, fmt.Errorf("failed to check preference match: %w", err)
		}
		if !matched {
			return false, nil
		}
	}
	return true, nil
}

func (pc *PreferenceChecker) checkPreferenceMatch(requestPref, otherPref *model.Preference) (bool, error) {
	if requestPref == nil || otherPref == nil {
		return false, fmt.Errorf("either request preference or the other preference to compare is nil")
	}
	if requestPref.Gender() != otherPref.Gender() {
		if requestPref.SameGender() || otherPref.SameGender() {
			return false, nil
		}
	}
	return true, nil
}
