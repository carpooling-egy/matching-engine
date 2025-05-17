package planner

import (
	"fmt"
	"matching-engine/internal/model"
	"matching-engine/internal/service/pathgeneration/generator"
	"matching-engine/internal/service/pathgeneration/validator"
	"matching-engine/internal/service/pickupdropoffservice"
)

// TODO: think if there might be cases when we can get a generation or validation & still want to proceed in path generation
type DefaultPathPlanner struct {
	pathGenerator         generator.PathGenerator
	pathValidator         validator.PathValidator
	pickupDropoffSelector pickupdropoffservice.PickupDropoffSelectorInterface
}

func NewDefaultPathPlanner(pathGenerator generator.PathGenerator, pathValidator validator.PathValidator, selector pickupdropoffservice.PickupDropoffSelectorInterface) *DefaultPathPlanner {
	return &DefaultPathPlanner{
		pathGenerator:         pathGenerator,
		pathValidator:         pathValidator,
		pickupDropoffSelector: selector,
	}
}
func (planner *DefaultPathPlanner) FindFirstFeasiblePath(offerNode *model.OfferNode, requestNode *model.RequestNode) ([]model.PathPoint, bool, error) {

	pickupAndDropOffs, err := planner.pickupDropoffSelector.GetPickupDropoffPointsAndDurations(requestNode.Request(), offerNode.Offer())
	if err != nil {
		return nil, false, fmt.Errorf("FindFirstFeasiblePath: error getting dropoff points: %w", err)
	}

	pathIter := planner.pathGenerator.GeneratePaths(
		offerNode.Offer().Path(),
		pickupAndDropOffs.Pickup(),
		pickupAndDropOffs.Dropoff(),
	)

	// Iterate through candidate paths
	for candidatePath, pathErr := range pathIter {
		if pathErr != nil {
			return nil, false, fmt.Errorf("FindFirstFeasiblePath: error generating path: %w", pathErr)
		}

		// Validate the candidate path
		// NOTE THAT THE FOLLOWING FUNCTION UPDATES THE POINTS IN THE CANDIDATE PATH ITSELF!!
		// (it updates the points with the expected arrival times)
		isValidPath, validateErr := planner.pathValidator.ValidatePath(offerNode, candidatePath)
		if validateErr != nil {
			return nil, false, fmt.Errorf("failed to validate path: %w", validateErr)
		}

		if isValidPath {
			// Found a valid path, return it immediately
			return candidatePath, true, nil
		}
	}

	// No valid paths found
	return nil, false, nil
}
