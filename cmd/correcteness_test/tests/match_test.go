package tests

import (
	"matching-engine/internal/adapter/routing"
	"matching-engine/internal/model"
	"testing"
)

func TestMatch(t *testing.T) {
	engine, err := setupTestingEnvironment()
	if err != nil {
		t.Fatalf("Failed to create Valhalla engine: %v", err)
	}
	tests := []struct {
		name     string
		testFunc func(engine routing.Engine) ([]*model.Offer, []*model.Request, map[string]*model.MatchingResult)
	}{
		{
			name:     "Basic match",
			testFunc: getTest2aData,
		},
		{
			name:     "Match offer with multiple requests",
			testFunc: getTest3aData,
		},
		{
			name:     "Match offer with multiple requests checking dynamic capacity",
			testFunc: getTest3bData,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			offers, requests, expectedResults := tt.testFunc(engine)
			if len(offers) == 0 || len(requests) == 0 {
				t.Fatalf("No offers or requests generated for test %s", tt.name)
			}
			results, err := runMatcher(offers, requests)
			if err != nil {
				t.Fatalf("Matcher failed for test %s: %v", tt.name, err)
			}
			if results == nil && expectedResults == nil {
				return // Both results are nil, which is acceptable
			}
			if len(results) != len(expectedResults) {
				t.Fatalf("Expected %d results, got %d for test %s", len(expectedResults), len(results), tt.name)
			}
			if !compareResults(results, expectedResults) {
				t.Fatalf("Results do not match expected for test %s", tt.name)
			}
		})
	}
}

func TestMatchMultipleOffersMultipleRequests(t *testing.T) {
	testName := "Match Multiple Offers Multiple Requests"

	engine, err := setupTestingEnvironment()
	if err != nil {
		t.Fatalf("[%s] Failed to create Valhalla engine: %v", testName, err)
	}

	offers, requests, expectedResults_a, expectedResults_b := getTest4(engine)
	if len(offers) == 0 || len(requests) == 0 {
		t.Fatalf("[%s] No offers or requests generated", testName)
	}
	results, err := runMatcher(offers, requests)
	if err != nil {
		t.Fatalf("[%s] Matcher failed: %v", testName, err)
	}
	if results == nil && (expectedResults_a == nil || expectedResults_b == nil) {
		return // Both results are nil, which is acceptable
	}
	if len(results) != len(expectedResults_a) && len(results) != len(expectedResults_b) {
		t.Fatalf("[%s] Expected %d or %d results, got %d", testName, len(expectedResults_a), len(expectedResults_b), len(results))
	}
	if !compareResults(results, expectedResults_a) && !compareResults(results, expectedResults_b) {
		t.Fatalf("[%s] Results do not match expected", testName)
	}
}
