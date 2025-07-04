package starter

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"matching-engine/cmd/appmetrics"
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
	original_passenger_number := 0
	for _, request := range requests {
		original_passenger_number += request.NumberOfRiders()
	}
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
		Int("matches", len(matchingResults)).
		Msg("Matching process completed")

	log.Info().
		Str("duration", elapsed.String()).
		Int("matched_requests", totalMatchedRequests).
		Int("matched_drivers", matchedDrivers).
		Int("total_number_of_passengers", totalNumberOfRiders).
		Float64("overall_matching", overall_matching).
		Int("max_matched_requests", max_matched_requests).
		Msg("Matching process completed")

	log.Info().
		Float64("Requests Matching Percentage %", float64(totalMatchedRequests)/float64(len(requests))*100).
		Float64("Drivers Matching Percentage %", float64(matchedDrivers)/float64(len(offers))*100).
		Float64("Overall Matching Percentage %", overall_matching/float64(len(requests)+len(offers))*100).
		Float64("Total Number of Passengers Percentage %", float64(totalNumberOfRiders)/float64(original_passenger_number)*100).
		Msg("Matching statistics")

	appmetrics.IncrementCount("Average matches count", float64(len(matchingResults)))
	appmetrics.TrackTime("Average matching duration", elapsed)
	appmetrics.IncrementCount("Average total matched requests", float64(totalMatchedRequests))
	appmetrics.IncrementCount("Average total matched drivers", float64(matchedDrivers))
	appmetrics.IncrementCount("Average max matched requests", float64(max_matched_requests))
	appmetrics.IncrementCount("Average total number of passengers", float64(totalNumberOfRiders))
	appmetrics.IncrementCount("Average overall matching", overall_matching)
	appmetrics.IncrementCount("Average matching percentage requests", float64(totalMatchedRequests)/float64(len(requests))*100)
	appmetrics.IncrementCount("Average matching percentage drivers", float64(matchedDrivers)/float64(len(offers))*100)
	appmetrics.IncrementCount("Average overall matching percentage", overall_matching/float64(len(requests)+len(offers))*100)

	// Publish results
	// if err = s.publisher.Publish(matchingResults); err != nil {
	//     return fmt.Errorf("failed to publish matching results: %w", err)
	// }
	one_edge_duration := appmetrics.GetTime("one_edge") / time.Duration(appmetrics.GetCount("one_edge"))
	log.Info().
		Str("metric", "one_edge").
		Str("average_duration", one_edge_duration.String()).
		Msg("Average duration for metric")

	log.Info().
		Int("offers", len(offers)).
		Int("requests", len(requests)).
		Int("matches", len(matchingResults)).
		Msg("Matching process completed successfully")

	detour_time_checker_duration := appmetrics.GetTime("detour_time_checker_duration") / time.Duration(appmetrics.GetCount("detour_time_checker_count"))
	log.Info().
		Str("metric", "detour_time_checker_duration").
		Str("average_duration", detour_time_checker_duration.String()).
		Msg("Average duration for metric")
	return nil
}
