package tests

import (
	"github.com/rs/zerolog/log"
	"go.uber.org/dig"
	"matching-engine/internal/adapter/routing"
	"matching-engine/internal/adapter/valhalla"
	"matching-engine/internal/app/config"
	"matching-engine/internal/app/di"
	"matching-engine/internal/app/di/utils"
	"matching-engine/internal/model"
	matcher2 "matching-engine/internal/service/matcher"
	"testing"
	"time"
)

func setupTestingEnvironment() (routing.Engine, error) {
	config.ConfigureLogging()

	// Create a mock routing engine
	engine, err := valhalla.NewValhalla()
	if err != nil {
		return nil, err
	}
	return engine, nil
}

func TestNoMatchTimeOverlap(t *testing.T) {
	engine, err := setupTestingEnvironment()
	if err != nil {
		t.Fatalf("Failed to create Valhalla engine: %v", err)
	}
	tests := []struct {
		name     string
		testFunc func(engine routing.Engine) ([]*model.Offer, []*model.Request, map[string]*model.MatchingResult)
	}{
		{
			name:     "Latest arrival before driver departure",
			testFunc: getTest1aiData,
		},
		{
			name:     "Earliest departure after driver max estimated arrival",
			testFunc: getTest1aiiData,
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

func TestNoMatchCapacity(t *testing.T) {
	engine, err := setupTestingEnvironment()
	if err != nil {
		t.Fatalf("Failed to create Valhalla engine: %v", err)
	}
	tests := []struct {
		name     string
		testFunc func(engine routing.Engine) ([]*model.Offer, []*model.Request, map[string]*model.MatchingResult)
	}{
		{
			name:     "Offer capacity less than request riders",
			testFunc: getTest1bData,
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

func TestNoMatchPreferenceMismatch(t *testing.T) {
	engine, err := setupTestingEnvironment()
	if err != nil {
		t.Fatalf("Failed to create Valhalla engine: %v", err)
	}
	tests := []struct {
		name     string
		testFunc func(engine routing.Engine) ([]*model.Offer, []*model.Request, map[string]*model.MatchingResult)
	}{
		{
			name:     "Offer preference does not match request preference",
			testFunc: getTest1ciData,
		},
		{
			name:     "Matched request preference does not match request preference",
			testFunc: getTest1ciiData,
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

func TestNoMatchPreDepartureArrival(t *testing.T) {
	engine, err := setupTestingEnvironment()
	if err != nil {
		t.Fatalf("Failed to create Valhalla engine: %v", err)
	}
	tests := []struct {
		name     string
		testFunc func(engine routing.Engine) ([]*model.Offer, []*model.Request, map[string]*model.MatchingResult)
	}{
		{
			name:     "Driver arrives at pickup before request earliest departure",
			testFunc: getTest1diData,
		},
		{
			name:     "Driver arrives at dropoff after request latest arrival",
			testFunc: getTest1diiData,
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

func TestNoMatchPreDetour(t *testing.T) {
	engine, err := setupTestingEnvironment()
	if err != nil {
		t.Fatalf("Failed to create Valhalla engine: %v", err)
	}
	tests := []struct {
		name     string
		testFunc func(engine routing.Engine) ([]*model.Offer, []*model.Request, map[string]*model.MatchingResult)
	}{
		{
			name:     "Driver's trip with rider exceeds the driver's direct time with detour",
			testFunc: getTest1eData,
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

func TestNoMatchFeasibiltyConstraints(t *testing.T) {
	engine, err := setupTestingEnvironment()
	if err != nil {
		t.Fatalf("Failed to create Valhalla engine: %v", err)
	}
	tests := []struct {
		name     string
		testFunc func(engine routing.Engine) ([]*model.Offer, []*model.Request, map[string]*model.MatchingResult)
	}{
		{
			name:     "Detour passes but Dynamic Capacity constraint fails",
			testFunc: getTest1fiData,
		},
		{
			name:     "Detour fails but Dynamic Capacity constraint passes",
			testFunc: getTest1fiiData,
		},
		{
			name:     "Detour passes but rider doesn't reach before his latest arrival time",
			testFunc: getTest1fiiiData,
		},
		{
			name:     "The matched request's latest arrival time constraint is violated",
			testFunc: getTest1fivData,
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

func TestCorrecteness(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(t2 *testing.T)
	}{
		{
			name:     "Test No Match due to Overlap",
			testFunc: TestNoMatchTimeOverlap,
		},
		{
			name:     "Test No Match due to Capacity",
			testFunc: TestNoMatchCapacity,
		},
		{
			name:     "Test No Match due to Preference Mismatch",
			testFunc: TestNoMatchPreferenceMismatch,
		},
		{
			name:     "Test No Match due to Pre-Departure Arrival",
			testFunc: TestNoMatchPreDepartureArrival,
		},
		{
			name:     "Test No Match due to Pre-Detour",
			testFunc: TestNoMatchPreDetour,
		},
		{
			name:     "Test No Match due to Feasibility Constraints",
			testFunc: TestNoMatchFeasibiltyConstraints,
		},
		{
			name:     "Test Match",
			testFunc: TestMatch,
		},
		{
			name:     "Test Multiple Offers and Requests",
			testFunc: TestMatchMultipleOffersMultipleRequests,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t)
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

func compareResults(results []*model.MatchingResult, expectedResults map[string]*model.MatchingResult) bool {
	if len(results) != len(expectedResults) {
		log.Debug().Msgf("Results length mismatch: got %d, expected %d\n", len(results), len(expectedResults))
		return false
	}
	for _, result := range results {
		expectedResult := expectedResults[result.OfferID()]
		if result.UserID() != expectedResult.UserID() || result.OfferID() != expectedResult.OfferID() {
			log.Debug().Msgf("UserID or OfferID mismatch: got %s/%s, expected %s/%s\n",
				result.UserID(), result.OfferID(), expectedResult.UserID(), expectedResult.OfferID())
			return false
		}
		if len(result.AssignedMatchedRequests()) != len(expectedResult.AssignedMatchedRequests()) {
			log.Debug().Msgf("Assigned matched requests length mismatch: got %d, expected %d\n",
				len(result.AssignedMatchedRequests()), len(expectedResult.AssignedMatchedRequests()))
			return false
		}
		if len(result.NewPath()) != len(expectedResult.NewPath()) {
			log.Debug().Msgf("Path length mismatch: got %d, expected %d\n",
				len(result.NewPath()), len(expectedResult.NewPath()))
			return false
		}
		if result.CurrentNumberOfRequests() != expectedResult.CurrentNumberOfRequests() {
			log.Debug().Msgf("Number of requests mismatch: got %d, expected %d\n",
				result.CurrentNumberOfRequests(), expectedResult.CurrentNumberOfRequests())
			return false
		}
		for _, req := range result.AssignedMatchedRequests() {
			matchedRequests := false
			for _, expectedReq := range expectedResult.AssignedMatchedRequests() {
				if req.ID() == expectedReq.ID() &&
					req.Source().Equal(expectedReq.Source()) &&
					req.Destination().Equal(expectedReq.Destination()) &&
					req.EarliestDepartureTime().Equal(expectedReq.EarliestDepartureTime()) &&
					req.LatestArrivalTime().Equal(expectedReq.LatestArrivalTime()) &&
					req.MaxWalkingDurationMinutes() == expectedReq.MaxWalkingDurationMinutes() &&
					req.NumberOfRiders() == expectedReq.NumberOfRiders() &&
					req.Preferences() == expectedReq.Preferences() {
					matchedRequests = true
				}
			}
			if !matchedRequests {
				log.Debug().
					Str("requestID", req.ID()).
					Msg("Assigned matched request not found in expected results")
				return false
			}
		}
		for i, point := range result.NewPath() {
			if i >= len(expectedResult.NewPath()) {
				log.Debug().
					Int("point", i).
					Msg("Point index out of bounds in expected result")
				return false
			}
			expectedPoint := expectedResult.NewPath()[i]
			if !point.Coordinate().Equal(expectedPoint.Coordinate()) ||
				point.PointType() != expectedPoint.PointType() ||
				!checkTimeOverlap(point.ExpectedArrivalTime(), expectedPoint.ExpectedArrivalTime(), 10*time.Second) ||
				point.WalkingDuration() != expectedPoint.WalkingDuration() ||
				!checkOwnerMatch(point, expectedPoint) {
				log.Debug().
					Int("point", i).
					Msg("Point mismatch: ")
				return false
			}
		}
	}
	return true
}

func checkOwnerMatch(point model.PathPoint, expectedPoint model.PathPoint) bool {
	_, isRequest := point.Owner().AsRequest()
	_, isExpectedRequest := expectedPoint.Owner().AsRequest()
	if isRequest != isExpectedRequest {
		return false
	}
	if point.GetOwnerID() != expectedPoint.GetOwnerID() {
		return false
	}
	return true
}

func checkTimeOverlap(time1, time2 time.Time, tolerance time.Duration) bool {
	// Check if the two times are within the specified tolerance
	return time1.After(time2.Add(-tolerance)) && time1.Before(time2.Add(tolerance))
}

func runMatcher(offers []*model.Offer, requests []*model.Request) ([]*model.MatchingResult, error) {

	c := dig.New()

	// register all dependencies for matching services
	di.RegisterGeoServices(c)
	di.RegisterPickupDropoffServices(c)
	di.RegisterTimeMatrixServices(c)
	di.RegisterPathServices(c)
	di.RegisterCheckers(c)
	di.RegisterMatchingServices(c)
	utils.Must(c.Provide(valhalla.NewValhalla))

	var matches []*model.MatchingResult
	var matchErr error
	err := c.Invoke(func(matcher *matcher2.Matcher) {
		matches, matchErr = matcher.Match(offers, requests)
	})
	if err != nil {
		panic("Failed to invoke matcher in the container: " + err.Error())
	}
	return matches, matchErr
}
