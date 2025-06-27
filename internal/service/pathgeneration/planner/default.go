package planner

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"matching-engine/internal/model"
	"matching-engine/internal/service/pathgeneration/generator"
	"matching-engine/internal/service/pathgeneration/validator"
	"matching-engine/internal/service/pickupdropoffservice"
)

// TODO: Revisit error handling
type DefaultPathPlanner struct {
	pathGenerator         generator.PathGenerator
	pathValidator         validator.PathValidator
	pickupDropoffSelector pickupdropoffservice.PickupDropoffSelectorInterface
}

func NewDefaultPathPlanner(pathGenerator generator.PathGenerator, pathValidator validator.PathValidator, selector pickupdropoffservice.PickupDropoffSelectorInterface) PathPlanner {
	return &DefaultPathPlanner{
		pathGenerator:         pathGenerator,
		pathValidator:         pathValidator,
		pickupDropoffSelector: selector,
	}
}
func (planner *DefaultPathPlanner) FindFirstFeasiblePath(offerNode *model.OfferNode, requestNode *model.RequestNode) ([]model.PathPoint, bool, error) {

	pickupAndDropOffs, err := planner.pickupDropoffSelector.GetPickupDropoffPointsAndDurations(requestNode.Request(), offerNode.Offer())
	if err != nil {
		return nil, false, fmt.Errorf("FindFirstFeasiblePath: error getting pickup & dropoff points: %w", err)
	}

	// NOTE:
	// The path generator is expected to yield a *new slice* for each path iteration.
	// This ensures that modifying the returned slice (e.g., setting ExpectedArrivalTime)
	// inside validation logic is safe and does not affect other iterations.
	pathIter, err := planner.pathGenerator.GeneratePaths(
		offerNode.Offer().Path(),
		pickupAndDropOffs.Pickup(),
		pickupAndDropOffs.Dropoff(),
	)

	if err != nil {
		return nil, false, fmt.Errorf("FindFirstFeasiblePath: error getting path iterator: %w", err)
	}

	// Iterate through candidate paths
	for candidatePath, pathErr := range pathIter {
		if pathErr != nil {
			return nil, false, fmt.Errorf("FindFirstFeasiblePath: error generating path: %w", pathErr)
		}

		// Validate the candidate path
		// NOTE THAT THE FOLLOWING FUNCTION UPDATES THE POINTS IN THE CANDIDATE PATH ITSELF!!
		// (it updates the points with the expected arrival times)
		isValidPath, validateErr := planner.pathValidator.ValidatePath(offerNode, requestNode, candidatePath)
		if validateErr != nil {
			return nil, false, fmt.Errorf("failed to validate path: %w", validateErr)
		}
		if isValidPath {
			// Found a valid path, return it immediately
			return candidatePath, true, nil
		}
	}

	log.Debug().
		Str("offer_id", offerNode.Offer().ID()).
		Str("request_id", requestNode.Request().ID()).
		Msg("No valid paths found for offer and request")
	// No valid paths found
	return nil, false, nil
}
