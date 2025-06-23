package tests

import (
	"fmt"
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

func TestCorrecteness(t *testing.T) {
	config.ConfigureLogging()

	// Create a mock routing engine
	engine, err := valhalla.NewValhalla()
	if err != nil {
		t.Fatalf("Failed to create Valhalla engine: %v", err)
	}
	tests := []struct {
		name     string
		testFunc func(engine routing.Engine) ([]*model.Offer, []*model.Request, map[string]*model.MatchingResult)
	}{
		{
			name:     "Test1ai",
			testFunc: getTest1aiData,
		},
		{
			name:     "Test1aii",
			testFunc: getTest1aiiData,
		},
		{
			name:     "Test1b",
			testFunc: getTest1bData,
		},
		{
			name:     "Test1ci",
			testFunc: getTest1ciData,
		},
		{
			name:     "Test1cii",
			testFunc: getTest1ciiData,
		},
		{
			name:     "Test1di",
			testFunc: getTest1diData,
		},
		{
			name:     "Test1dii",
			testFunc: getTest1diiData,
		},
		{
			name:     "Test1e",
			testFunc: getTest1eData,
		},
		{
			name:     "Test1fi",
			testFunc: getTest1fiData,
		},
		{
			name:     "Test1fii",
			testFunc: getTest1fiiData,
		},
		{
			name:     "Test1fiii",
			testFunc: getTest1fiiiData,
		},
		{
			name:     "Test1fiv",
			testFunc: getTest1fivData,
		},
		{
			name:     "Test2a",
			testFunc: getTest2aData,
		},
		{
			name:     "Test3a",
			testFunc: getTest3aData,
		},
		{
			name:     "Test3b",
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

func TestCorrecteness4(t *testing.T) {
	testName := "TestCorrecteness4"
	config.ConfigureLogging()

	// Create a mock routing engine
	engine, err := valhalla.NewValhalla()
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
	fmt.Println("debugging results")

	fmt.Println(len(results))
	fmt.Println(len(results[0].AssignedMatchedRequests()))
	fmt.Println(len(results[1].AssignedMatchedRequests()))

	if len(results) != len(expectedResults_a) && len(results) != len(expectedResults_b) {
		t.Fatalf("[%s] Expected %d or %d results, got %d", testName, len(expectedResults_a), len(expectedResults_b), len(results))
	}
	if !compareResults(results, expectedResults_a) && !compareResults(results, expectedResults_b) {
		t.Fatalf("[%s] Results do not match expected", testName)
	}
}

func compareResults(results []*model.MatchingResult, expectedResults map[string]*model.MatchingResult) bool {
	fmt.Println("calllledddddddd")
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
