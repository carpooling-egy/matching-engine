package tests

import (
	"matching-engine/internal/adapter/routing"
	"matching-engine/internal/adapter/valhalla"
	"matching-engine/internal/app/config"
	"matching-engine/internal/app/di"
	"matching-engine/internal/app/di/utils"
	"matching-engine/internal/model"
	matcher2 "matching-engine/internal/service/matcher"
	"testing"
	"time"

	// "github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"go.uber.org/dig"
)

func setupTestingEnvironment() (routing.Engine, error) {
	config.ConfigureLogging()
	// Load environment variables
	if err := config.LoadEnv(); err != nil {
		log.Fatal().Err(err).Msg("Failed to load environment variables")
		return nil, err
	}
	// Create a mock routing engine
	engine, err := valhalla.NewValhalla()
	if err != nil {
		return nil, err
	}
	return engine, nil
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
			name:     "Test No Match Exceeds Requests Limit",
			testFunc: TestNoMatchExceedsRequestsLimit,
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

	// TODO set the env variable with correct variable to test different approaches
	// TODO use a better mechanism to test all variations
	// envPath := ""

	// err := godotenv.Load(envPath)
	// if err != nil {
	// 	log.Error().Msgf("Error loading .env file: %v", err)
	// 	return nil, err
	// }

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
