package starter

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"matching-engine/internal/reader"
	"matching-engine/internal/service/matcher"
	"time"
)

type StarterService struct {
	reader  reader.MatchInputReader
	matcher *matcher.Matcher
	//publisher publisher.Publisher
}

// NewStarterService creates a new starter service
func NewStarterService(reader reader.MatchInputReader, matcher *matcher.Matcher /*, publisher publisher.Publisher*/) *StarterService {
	return &StarterService{
		reader:  reader,
		matcher: matcher,
		//publisher: publisher,
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
	startTime := time.Now()
	matchingResults, err := s.matcher.Match(offers, requests)
	if err != nil {
		return fmt.Errorf("failed to match offers and requests: %w", err)
	}
	elapsed := time.Since(startTime)

	if len(matchingResults) == 0 {
		log.Info().Msg("No matches found")
		return nil
	}
	// Publish results
	//if err = s.publisher.Publish(matchingResults); err != nil {
	//	return fmt.Errorf("failed to publish matching results: %w", err)
	//}

	max_matched_requests := 0
	totalMatchedRequests := 0
	totalNumberOfRiders := 0
	for _, match := range matchingResults {
		totalMatchedRequests += match.CurrentNumberOfRequests()
		max_matched_requests = max(max_matched_requests, match.CurrentNumberOfRequests())
		for _, request := range match.AssignedMatchedRequests() {
			totalNumberOfRiders += request.NumberOfRiders()
		}
	}
	// / Calculate the request fulfillment rate
	// requestFulfillmentRate := float64(totalMatchedRequests) / float64(len(requests)) * 100

	// The number of matched drivers is simply the number of matching results
	matchedDrivers := len(matchingResults)
	// driverUtilizationRate := float64(matchedDrivers) / float64(len(offers)) * 100

	overall_matching := float64(totalMatchedRequests + matchedDrivers)

	log.Info().
		Str("duration", elapsed.String()).
		Int("matches", len(matchingResults)).
		Msg("Matching process completed")

	log.Info().
		Str("duration", elapsed.String()).
		Int("matched_requests", totalMatchedRequests).
		Int("matched_drivers", matchedDrivers).
		Float64("overall_matching", overall_matching).
		Int("max_matched_requests", max_matched_requests).
		Msg("Matching process completed")

	// Publish results
	// if err = s.publisher.Publish(matchingResults); err != nil {
	//     return fmt.Errorf("failed to publish matching results: %w", err)
	// }

	log.Info().
		Int("offers", len(offers)).
		Int("requests", len(requests)).
		Int("matches", len(matchingResults)).
		Msg("Matching process completed successfully")

	return nil
}
