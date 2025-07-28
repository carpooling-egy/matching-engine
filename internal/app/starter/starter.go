package starter

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"matching-engine/internal/publisher"
	"matching-engine/internal/reader"
	"matching-engine/internal/service/matcher"
)

type StarterService struct {
	reader    reader.MatchInputReader
	matcher   *matcher.Matcher
	publisher publisher.Publisher
}

// NewStarterService creates a new starter service
func NewStarterService(reader reader.MatchInputReader, matcher *matcher.Matcher, publisher publisher.Publisher) *StarterService {
	return &StarterService{
		reader:    reader,
		matcher:   matcher,
		publisher: publisher,
	}
}

// Start initiates the matching process
func (s *StarterService) Start(ctx context.Context) error {
	log.Info().Msg("Starting matching process...")

	// Get offers and requests
	requests, offers, exists, err := s.reader.GetOffersAndRequests(ctx)
	if err != nil {
		s.reader.Close()
		return fmt.Errorf("failed to get offers and requests: %w", err)
	}
	if !exists {
		s.reader.Close()
		log.Info().Msg("No offers or requests found")
		return nil
	}

	// Close reader
	err = s.reader.Close()
	if err != nil {
		log.Error().Err(err).Msg("Failed to close reader")
	}

	// Process matching
	matchingResults, err := s.matcher.Match(offers, requests)
	if err != nil {
		return fmt.Errorf("failed to match offers and requests: %w", err)
	}

	if len(matchingResults) == 0 {
		log.Info().Msg("No matches found")
		return nil
	}
	// Publish results
	if err = s.publisher.Publish(matchingResults); err != nil {
		return fmt.Errorf("failed to publish matching results: %w", err)
	}

	log.Info().
		Int("offers", len(offers)).
		Int("requests", len(requests)).
		Int("matches", len(matchingResults)).
		Msg("Matching process completed successfully")

	return nil
}
