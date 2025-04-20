package service

import (
	"context"
	"github.com/rs/zerolog/log"
	"matching-engine/internal/model"
	"matching-engine/internal/repository"
)

// MatcherService handles the matching of riders to drivers
type MatcherService struct {
	repo repository.Repository
}

// NewMatcherService creates a new instance of MatcherService
func NewMatcherService(repo repository.Repository) *MatcherService {
	return &MatcherService{
		repo: repo,
	}
}

// RunMatchingAlgorithm processes pending rider requests and attempts to match them with available drivers
func (s *MatcherService) RunMatchingAlgorithm(ctx context.Context) error {
	// Get all pending rider requests
	pendingRequests, err := s.repo.GetPendingRiderRequests(ctx)
	if err != nil {
		return err
	}

	log.Info().Msgf("Processing %d pending rider requests", len(pendingRequests))

	// In a real implementation, this would use a sophisticated algorithm to match riders with drivers
	// For now, we'll just simulate matching by updating the status of each request
	for _, request := range pendingRequests {
		// In a real implementation, we would find the best driver for this request
		// For now, we'll just mark it as matched with a dummy driver ID
		driverID := "driver-123" // This would be determined by the matching algorithm

		// Save the match
		err := s.repo.SaveMatch(ctx, request.ID, driverID)
		if err != nil {
			log.Error().Err(err).Str("requestID", request.ID).Msg("Error saving match for request")
			continue
		}

		// Update the request status
		err = s.repo.UpdateRiderRequestStatus(ctx, request.ID, model.StatusMatched)
		if err != nil {
			log.Error().Err(err).Str("requestID", request.ID).Msg("Error updating status for request")
			continue
		}

		log.Info().Str("requestID", request.ID).Str("driverID", driverID).Msg("Matched rider request with driver")
	}

	return nil
}
