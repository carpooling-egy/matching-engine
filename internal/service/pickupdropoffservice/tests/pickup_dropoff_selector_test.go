package tests

import (
	"context"
	"errors"
	"matching-engine/internal/adapter/routing"
	"matching-engine/internal/enums"
	"matching-engine/internal/model"
	"matching-engine/internal/service/pickupdropoffservice"
	"matching-engine/internal/service/pickupdropoffservice/pickupdropoffcache"
	"testing"
	"time"
)

// MockPickupDropoffGenerator is a mock implementation of the PickupDropoffGenerator interface
type MockPickupDropoffGenerator struct {
	pickup  *model.PathPoint
	dropoff *model.PathPoint
	err     error
	// Add a field to track calls
	callCount int
}

func NewMockPickupDropoffGenerator(pickup, dropoff *model.PathPoint, err error) *MockPickupDropoffGenerator {
	return &MockPickupDropoffGenerator{
		pickup:    pickup,
		dropoff:   dropoff,
		err:       err,
		callCount: 0,
	}
}

func (m *MockPickupDropoffGenerator) GeneratePickupDropoffPoints(request *model.Request, offer *model.Offer) (*model.PathPoint, *model.PathPoint, error) {
	m.callCount++
	return m.pickup, m.dropoff, m.err
}

// MockRoutingEngine_pickup_dropoff_selector is a mock implementation of the routing.Engine interface
type MockRoutingEngine_pickup_dropoff_selector struct {
	walkingDuration time.Duration
	err             error
	// Add a field to track calls
	callCount int
}

func NewMockRoutingEngine_pickup_dropoff_selector(walkingDuration time.Duration, err error) *MockRoutingEngine_pickup_dropoff_selector {
	return &MockRoutingEngine_pickup_dropoff_selector{
		walkingDuration: walkingDuration,
		err:             err,
		callCount:       0,
	}
}

func (m *MockRoutingEngine_pickup_dropoff_selector) ComputeWalkingTime(ctx context.Context, walkParams *model.WalkParams) (time.Duration, error) {
	m.callCount++
	return m.walkingDuration, m.err
}

// Implement other methods of the routing.Engine interface with empty implementations
func (m *MockRoutingEngine_pickup_dropoff_selector) ComputeDrivingTime(ctx context.Context, routeParams *model.RouteParams) ([]time.Duration, error) {
	return nil, errors.New("not implemented")
}

func (m *MockRoutingEngine_pickup_dropoff_selector) PlanDrivingRoute(ctx context.Context, routeParams *model.RouteParams) (*model.Route, error) {
	return nil, errors.New("not implemented")
}

func (m *MockRoutingEngine_pickup_dropoff_selector) ComputeIsochrone(ctx context.Context, req *model.IsochroneParams) (*model.Isochrone, error) {
	return nil, errors.New("not implemented")
}

func (m *MockRoutingEngine_pickup_dropoff_selector) ComputeDistanceTimeMatrix(ctx context.Context, req *model.DistanceTimeMatrixParams) (*model.DistanceTimeMatrix, error) {
	return nil, errors.New("not implemented")
}

func (m *MockRoutingEngine_pickup_dropoff_selector) SnapPointToRoad(ctx context.Context, point *model.Coordinate) (*model.Coordinate, error) {
	return nil, errors.New("not implemented")
}

func TestPickupDropoffSelector_GetPickupDropoffPointsAndDurations(t *testing.T) {
	// Create test coordinates
	sourceCoord, _ := model.NewCoordinate(1.0, 1.0)
	destCoord, _ := model.NewCoordinate(2.0, 2.0)
	pickupCoord, _ := model.NewCoordinate(1.1, 1.1)
	dropoffCoord, _ := model.NewCoordinate(1.9, 1.9)

	// Create test times
	now := time.Now()
	later := now.Add(1 * time.Hour)

	// Create test request and offer
	request := model.NewRequest(
		"request1",
		"user1",
		*sourceCoord,
		*destCoord,
		now,
		later,
		15*time.Minute,
		1,
		model.Preference{},
	)

	offer := model.NewOffer(
		"offer1",
		"user2",
		*sourceCoord,
		*destCoord,
		now,
		30*time.Minute,
		4,
		model.Preference{},
		later,
		0,
		nil,
		nil,
	)

	// Create test path points
	pickup := model.NewPathPoint(*pickupCoord, enums.Pickup, now, nil, 0)
	dropoff := model.NewPathPoint(*dropoffCoord, enums.Dropoff, later, nil, 0)

	// Define test cases
	tests := []struct {
		name                   string
		generator              pickupdropoffservice.PickupDropoffGenerator
		routingEngine          routing.Engine
		cache                  *pickupdropoffcache.PickupDropoffCache
		request                *model.Request
		offer                  *model.Offer
		expectedPickup         *model.PathPoint
		expectedDropoff        *model.PathPoint
		expectedPickupWalking  time.Duration
		expectedDropoffWalking time.Duration
		expectError            bool
		expectCacheHit         bool
		prePopulateCache       bool
	}{
		{
			name:                   "Success - Generate new points",
			generator:              NewMockPickupDropoffGenerator(pickup, dropoff, nil),
			routingEngine:          NewMockRoutingEngine_pickup_dropoff_selector(5*time.Minute, nil),
			cache:                  pickupdropoffcache.NewPickupDropoffCache(),
			request:                request,
			offer:                  offer,
			expectedPickup:         pickup,
			expectedDropoff:        dropoff,
			expectedPickupWalking:  5 * time.Minute,
			expectedDropoffWalking: 5 * time.Minute,
			expectError:            false,
			expectCacheHit:         false,
			prePopulateCache:       false,
		},
		{
			name:                   "Success - Cache hit",
			generator:              NewMockPickupDropoffGenerator(pickup, dropoff, nil),
			routingEngine:          NewMockRoutingEngine_pickup_dropoff_selector(5*time.Minute, nil),
			cache:                  pickupdropoffcache.NewPickupDropoffCache(),
			request:                request,
			offer:                  offer,
			expectedPickup:         pickup,
			expectedDropoff:        dropoff,
			expectedPickupWalking:  5 * time.Minute,
			expectedDropoffWalking: 5 * time.Minute,
			expectError:            false,
			expectCacheHit:         true,
			prePopulateCache:       true,
		},
		{
			name:                   "Error - Generator fails",
			generator:              NewMockPickupDropoffGenerator(nil, nil, errors.New("generator error")),
			routingEngine:          NewMockRoutingEngine_pickup_dropoff_selector(5*time.Minute, nil),
			cache:                  pickupdropoffcache.NewPickupDropoffCache(),
			request:                request,
			offer:                  offer,
			expectedPickup:         nil,
			expectedDropoff:        nil,
			expectedPickupWalking:  0,
			expectedDropoffWalking: 0,
			expectError:            true,
			expectCacheHit:         false,
			prePopulateCache:       false,
		},
		{
			name:                   "Error - Walking time calculator fails",
			generator:              NewMockPickupDropoffGenerator(pickup, dropoff, nil),
			routingEngine:          NewMockRoutingEngine_pickup_dropoff_selector(0, errors.New("walking time calculator error")),
			cache:                  pickupdropoffcache.NewPickupDropoffCache(),
			request:                request,
			offer:                  offer,
			expectedPickup:         nil,
			expectedDropoff:        nil,
			expectedPickupWalking:  0,
			expectedDropoffWalking: 0,
			expectError:            true,
			expectCacheHit:         false,
			prePopulateCache:       false,
		},
	}

	// Run test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Pre-populate cache if needed
			if tc.prePopulateCache {
				cacheKey := model.NewOfferRequestKey(tc.offer.ID(), tc.request.ID())
				tc.expectedPickup.SetWalkingDuration(tc.expectedPickupWalking)
				tc.expectedDropoff.SetWalkingDuration(tc.expectedDropoffWalking)
				cacheValue := pickupdropoffcache.NewValue(tc.expectedPickup, tc.expectedDropoff)
				tc.cache.Set(cacheKey, cacheValue)
			}

			// Create walking time calculator with the mock routing engine
			walkingTimeCalculator := pickupdropoffservice.NewWalkingTimeCalculator(tc.routingEngine)

			// Create selector
			selector := pickupdropoffservice.NewPickupDropoffSelector(tc.generator, walkingTimeCalculator, tc.cache)

			// Call the method
			result, err := selector.GetPickupDropoffPointsAndDurations(tc.request, tc.offer)

			// Check error
			if tc.expectError && err == nil {
				t.Errorf("Expected error but got nil")
			}
			if !tc.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			// If we expect an error, don't check the result
			if tc.expectError {
				return
			}

			// Check result
			if result == nil {
				t.Errorf("Expected result but got nil")
				return
			}

			// Check pickup point
			if result.Pickup() != tc.expectedPickup {
				t.Errorf("Expected pickup %v but got %v", tc.expectedPickup, result.Pickup())
			}

			// Check dropoff point
			if result.Dropoff() != tc.expectedDropoff {
				t.Errorf("Expected dropoff %v but got %v", tc.expectedDropoff, result.Dropoff())
			}

			// Check pickup walking duration
			if result.Pickup().WalkingDuration() != tc.expectedPickupWalking {
				t.Errorf("Expected pickup walking duration %v but got %v", tc.expectedPickupWalking, result.Pickup().WalkingDuration())
			}

			// Check dropoff walking duration
			if result.Dropoff().WalkingDuration() != tc.expectedDropoffWalking {
				t.Errorf("Expected dropoff walking duration %v but got %v", tc.expectedDropoffWalking, result.Dropoff().WalkingDuration())
			}
		})
	}
}

func TestPickupDropoffSelector_GetPickupDropoffPointsAndDurations_CacheUsage(t *testing.T) {
	// Create test coordinates
	sourceCoord, _ := model.NewCoordinate(1.0, 1.0)
	destCoord, _ := model.NewCoordinate(2.0, 2.0)
	pickupCoord, _ := model.NewCoordinate(1.1, 1.1)
	dropoffCoord, _ := model.NewCoordinate(1.9, 1.9)

	// Create test times
	now := time.Now()
	later := now.Add(1 * time.Hour)

	// Create test request and offer
	request := model.NewRequest(
		"request1",
		"user1",
		*sourceCoord,
		*destCoord,
		now,
		later,
		15*time.Minute,
		1,
		model.Preference{},
	)

	offer := model.NewOffer(
		"offer1",
		"user2",
		*sourceCoord,
		*destCoord,
		now,
		30*time.Minute,
		4,
		model.Preference{},
		later,
		0,
		nil,
		nil,
	)

	// Create test path points
	pickup := model.NewPathPoint(*pickupCoord, enums.Pickup, now, nil, 0)
	dropoff := model.NewPathPoint(*dropoffCoord, enums.Dropoff, later, nil, 0)

	// Create a generator that counts calls
	generator := NewMockPickupDropoffGenerator(pickup, dropoff, nil)

	// Create a routing engine that counts calls
	routingEngine := NewMockRoutingEngine_pickup_dropoff_selector(5*time.Minute, nil)

	// Create walking time calculator with the mock routing engine
	walkingTimeCalculator := pickupdropoffservice.NewWalkingTimeCalculator(routingEngine)

	// Create cache
	cache := pickupdropoffcache.NewPickupDropoffCache()

	// Create selector
	selector := pickupdropoffservice.NewPickupDropoffSelector(generator, walkingTimeCalculator, cache)

	// First call should generate new points
	result1, err := selector.GetPickupDropoffPointsAndDurations(request, offer)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if generator.callCount != 1 {
		t.Errorf("Expected generator to be called once, but was called %d times", generator.callCount)
	}
	if routingEngine.callCount != 2 { // Called twice, once for pickup and once for dropoff
		t.Errorf("Expected routing engine to be called twice, but was called %d times", routingEngine.callCount)
	}

	// Second call should use cache
	result2, err := selector.GetPickupDropoffPointsAndDurations(request, offer)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if generator.callCount != 1 {
		t.Errorf("Expected generator to still be called once, but was called %d times", generator.callCount)
	}
	if routingEngine.callCount != 2 {
		t.Errorf("Expected routing engine to still be called twice, but was called %d times", routingEngine.callCount)
	}

	// Results should be the same
	if result1 != result2 {
		t.Errorf("Expected same result object from cache, but got different objects")
	}
}
