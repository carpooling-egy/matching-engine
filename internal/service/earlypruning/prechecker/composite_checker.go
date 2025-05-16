package prechecker

import (
	"fmt"
	"matching-engine/internal/model"
)

type CompositeChecker struct {
	checkers []Checker
}

func NewCompositePreChecker(preCheckers ...Checker) *CompositeChecker {
	return &CompositeChecker{
		checkers: preCheckers,
	}
}

func (c *CompositeChecker) Check(offer *model.Offer, request *model.Request) (bool, error) {
	for _, checker := range c.checkers {
		ok, err := checker.Check(offer, request)
		if err != nil {
			return false, fmt.Errorf("checker %T failed: %w", checker, err)
		}
		if !ok {
			return false, nil
		}
	}
	return true, nil
}
