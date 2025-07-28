package matchevaluator

import (
	"matching-engine/internal/model"
)

// Evaluator defines the behavior for matching an offer to a request.
type Evaluator interface {
	// Evaluate takes an offer node and a request node, runs any necessary
	// preference checks and path planning, and returns the first feasible path
	// (as a slice of PathPoint) or an error if no valid path is found.
	Evaluate(
		offerNode *model.OfferNode,
		requestNode *model.RequestNode,
	) ([]model.PathPoint, bool, error)
}
