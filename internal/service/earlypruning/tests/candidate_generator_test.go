package tests

import (
	"matching-engine/internal/enums"
	"matching-engine/internal/model"
	"matching-engine/internal/service/earlypruning"
	"testing"
	"time"
)

func TestPreChecksCandidateGenerator_GenerateCandidates(t *testing.T) {
	// Create a simple checker that always returns true
	checker := &testChecker{shouldMatch: true}

	// Create the candidate generator
	generator := earlypruning.NewPreChecksCandidateGenerator(checker)

	// Create test offers and requests
	offers := createTestOffers()
	requests := createTestRequests()

	// Generate candidates
	candidateIterator, err := generator.GenerateCandidates(offers, requests)

	// Check if there was an error
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check if the iterator was created
	if candidateIterator == nil {
		t.Fatal("Expected candidate iterator to be created, got nil")
	}

	// Count the number of candidates
	candidateCount := 0
	for candidate, err := range candidateIterator.Candidates() {
		if err != nil {
			t.Fatalf("Expected no error while iterating, got %v", err)
		}
		if candidate == nil {
			t.Fatal("Expected candidate to be non-nil")
		}
		candidateCount++
	}

	// Since our checker always returns true, we should have offers * requests candidates
	expectedCount := len(offers) * len(requests)
	if candidateCount != expectedCount {
		t.Fatalf("Expected %d candidates, got %d", expectedCount, candidateCount)
	}
}

// Test with a checker that rejects all matches
func TestPreChecksCandidateGenerator_NoMatches(t *testing.T) {
	// Create a checker that always returns false
	checker := &testChecker{shouldMatch: false}

	// Create the candidate generator
	generator := earlypruning.NewPreChecksCandidateGenerator(checker)

	// Create test offers and requests
	offers := createTestOffers()
	requests := createTestRequests()

	// Generate candidates
	candidateIterator, err := generator.GenerateCandidates(offers, requests)

	// Check if there was an error
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check if the iterator was created
	if candidateIterator == nil {
		t.Fatal("Expected candidate iterator to be created, got nil")
	}

	// Count the number of candidates
	candidateCount := 0
	for candidate, err := range candidateIterator.Candidates() {
		if err != nil {
			t.Fatalf("Expected no error while iterating, got %v", err)
		}
		if candidate != nil {
			candidateCount++
		}
	}

	// Since our checker always returns false, we should have 0 candidates
	if candidateCount != 0 {
		t.Fatalf("Expected 0 candidates, got %d", candidateCount)
	}
}

// Test checker implementation for testing
type testChecker struct {
	shouldMatch bool
}

func (c *testChecker) Check(offer *model.Offer, request *model.Request) (bool, error) {
	return c.shouldMatch, nil
}

// Helper function to create test offers
func createTestOffers() []*model.Offer {
	now := time.Now()

	// Create a preference
	preference := model.NewPreference(enums.Male, false, false, false)

	// Create coordinates for source and destination
	source, _ := model.NewCoordinate(40.7128, -74.0060)       // New York
	destination, _ := model.NewCoordinate(34.0522, -118.2437) // Los Angeles

	// Create offers
	offer1 := model.NewOffer(
		"offer1",
		"user1",
		*source,
		*destination,
		now,
		30*time.Minute, // 30 minutes detour
		4,              // capacity
		*preference,
		now.Add(5*time.Hour), // max arrival time
		0,                    // current number of requests
		nil,                  // path
		nil,                  // matched requests
	)

	offer2 := model.NewOffer(
		"offer2",
		"user2",
		*source,
		*destination,
		now.Add(30*time.Minute), // depart 30 minutes later
		45*time.Minute,          // 45 minutes detour
		3,                       // capacity
		*preference,
		now.Add(6*time.Hour), // max arrival time
		0,                    // current number of requests
		nil,                  // path
		nil,                  // matched requests
	)

	return []*model.Offer{offer1, offer2}
}

// Helper function to create test requests
func createTestRequests() []*model.Request {
	now := time.Now()

	// Create a preference
	preference := model.NewPreference(enums.Female, false, false, false)

	// Create coordinates for source and destination
	source, _ := model.NewCoordinate(40.7128, -74.0060)       // New York
	destination, _ := model.NewCoordinate(34.0522, -118.2437) // Los Angeles

	// Create requests
	request1 := model.NewRequest(
		"request1",
		"user3",
		*source,
		*destination,
		now,                  // earliest departure time
		now.Add(5*time.Hour), // latest arrival time
		15*time.Minute,       // max walking duration
		2,                    // number of riders
		*preference,
	)

	request2 := model.NewRequest(
		"request2",
		"user4",
		*source,
		*destination,
		now.Add(15*time.Minute), // earliest departure time 15 minutes later
		now.Add(6*time.Hour),    // latest arrival time
		20*time.Minute,          // max walking duration
		1,                       // number of riders
		*preference,
	)

	return []*model.Request{request1, request2}
}
