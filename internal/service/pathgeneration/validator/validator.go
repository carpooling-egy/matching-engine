package validator

import "matching-engine/internal/model"

// PathValidator defines the interface for validating paths in the matching engine
type PathValidator interface {
	// ValidatePath checks if the given path satisfies all constraints.
	// It returns true if the path is valid, false otherwise.
	// An error is returned only for system errors, not for validation failures.
	//
	// Note: This method may modify the provided path by setting expected arrival times.
	ValidatePath(offerNode *model.OfferNode, requestNode *model.RequestNode, path []model.PathPoint) (bool, error)
}
