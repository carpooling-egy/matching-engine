package tests

import (
	"context"
	"fmt"
	"matching-engine/internal/adapter/routing"
	"matching-engine/internal/enums"
	"matching-engine/internal/model"
	"matching-engine/internal/service/earlypruning/prechecker"
	"matching-engine/internal/service/pickupdropoffservice/pickupdropoffcache"
	"testing"
	"time"
)

// MockPickupDropoffSelector is a mock implementation of the PickupDropoffSelectorInterface
type MockPickupDropoffSelector struct {
	value *pickupdropoffcache.Value
	err   error
}

// NewMockPickupDropoffSelector creates a new MockPickupDropoffSelector
func NewMockPickupDropoffSelector(value *pickupdropoffcache.Value, err error) *MockPickupDropoffSelector {
	return &MockPickupDropoffSelector{
		value: value,
		err:   err,
	}
}

// GetPickupDropoffPointsAndDurations implements the PickupDropoffSelectorInterface
func (m *MockPickupDropoffSelector) GetPickupDropoffPointsAndDurations(request *model.Request, offer *model.Offer) (*pickupdropoffcache.Value, error) {
	return m.value, m.err
}

// MockEngine implements the routing.Engine interface for testing
type MockEngine struct {
	// For the route
	durations []time.Duration
	err       error
}

// NewMockEngine creates a new MockEngine
func NewMockEngine(
	durations []time.Duration,
	err error,
) routing.Engine {
	return &MockEngine{
		durations: durations,
		err:       err,
	}
}

// ComputeDrivingTime implements the routing.Engine interface
func (m *MockEngine) ComputeDrivingTime(ctx context.Context, routeParams *model.RouteParams) ([]time.Duration, error) {
	return m.durations, m.err
}

// PlanDrivingRoute implements the routing.Engine interface
func (m *MockEngine) PlanDrivingRoute(ctx context.Context, routeParams *model.RouteParams) (*model.Route, error) {
	return nil, fmt.Errorf("not implemented")
}

// ComputeWalkingTime implements the routing.Engine interface
func (m *MockEngine) ComputeWalkingTime(ctx context.Context, walkParams *model.WalkParams) (time.Duration, error) {
	return 0, fmt.Errorf("not implemented")
}

// ComputeIsochrone implements the routing.Engine interface
func (m *MockEngine) ComputeIsochrone(ctx context.Context, req *model.IsochroneParams) (*model.Isochrone, error) {
	return nil, fmt.Errorf("not implemented")
}

// ComputeDistanceTimeMatrix implements the routing.Engine interface
func (m *MockEngine) ComputeDistanceTimeMatrix(ctx context.Context, req *model.DistanceTimeMatrixParams) (*model.DistanceTimeMatrix, error) {
	return nil, fmt.Errorf("not implemented")
}

// SnapPointToRoad implements the routing.Engine interface
func (m *MockEngine) SnapPointToRoad(ctx context.Context, point *model.Coordinate) (*model.Coordinate, error) {
	return nil, fmt.Errorf("not implemented")
}

// directDuration is a constant representing the direct duration of the trip
const directDuration = 40 * time.Minute

func computeMaxEstimatedArrivalTime(offerDepartureTime time.Time, detourTime time.Duration) time.Time {
	return offerDepartureTime.Add(directDuration + detourTime)
}
func TestDetourTimeChecker_Check(t *testing.T) {
	// Define test times for clarity - use a fixed time in the future
	now := time.Now().Add(24 * time.Hour) // Use tomorrow to ensure it's in the future
	quarterHourLater := now.Add(15 * time.Minute)
	oneHourLater := now.Add(1 * time.Hour)
	twoHoursLater := now.Add(2 * time.Hour)

	// Create test coordinates
	sourceCoord := model.Coordinate{} // Using zero values for simplicity
	destCoord := model.Coordinate{}
	pickupCoord := model.Coordinate{}
	dropoffCoord := model.Coordinate{}

	// Define test cases
	tests := []struct {
		name                 string
		offer                *model.Offer
		request              *model.Request
		mockSelectorValue    *pickupdropoffcache.Value
		mockSelectorErr      error
		mockDrivingDurations []time.Duration
		mockDrivingErr       error
		expected             bool
		expectError          bool
	}{
		{
			name: "Valid detour within acceptable range",
			offer: model.NewOffer(
				"offer1", "user1",
				sourceCoord, destCoord,
				now, 30*time.Minute, // 30 minutes detour allowed
				3,
				*model.NewPreference(enums.Male, false),
				twoHoursLater, // maxEstimatedArrivalTime
				0,
				nil,
				nil,
			),
			request: model.NewRequest(
				"request1", "user2",
				pickupCoord, dropoffCoord,
				now,           // earliestDepartureTime
				twoHoursLater, // latestArrivalTime
				10*time.Minute,
				2,
				*model.NewPreference(enums.Female, false),
			),
			mockSelectorValue: pickupdropoffcache.NewValue(
				model.NewPathPoint(pickupCoord, enums.Pickup, now, nil, 5*time.Minute),
				model.NewPathPoint(dropoffCoord, enums.Dropoff, oneHourLater, nil, 10*time.Minute),
			),
			mockSelectorErr:      nil,
			mockDrivingDurations: []time.Duration{0, 20 * time.Minute, 40 * time.Minute, 60 * time.Minute}, // Source, Pickup, Dropoff, Destination
			mockDrivingErr:       nil,
			expected:             true,
			expectError:          false,
		},
		{
			name: "Detour exceeds acceptable range",
			offer: model.NewOffer(
				"offer1", "user1",
				sourceCoord, destCoord,
				now, 10*time.Minute, // Only 10 minutes detour allowed
				3,
				*model.NewPreference(enums.Male, false),
				computeMaxEstimatedArrivalTime(now, 10*time.Minute),
				0,
				nil,
				nil,
			),
			request: model.NewRequest(
				"request1", "user2",
				pickupCoord, dropoffCoord,
				now,           // earliestDepartureTime
				twoHoursLater, // latestArrivalTime
				10*time.Minute,
				2,
				*model.NewPreference(enums.Female, false),
			),
			mockSelectorValue: pickupdropoffcache.NewValue(
				model.NewPathPoint(pickupCoord, enums.Pickup, now, nil, 5*time.Minute),
				model.NewPathPoint(dropoffCoord, enums.Dropoff, oneHourLater, nil, 10*time.Minute),
			),
			mockSelectorErr:      nil,
			mockDrivingDurations: []time.Duration{0, 20 * time.Minute, 40 * time.Minute, 60 * time.Minute}, // Source, Pickup, Dropoff, Destination
			mockDrivingErr:       nil,
			expected:             false, // Total trip duration (60 min) exceeds maxEstimatedArrivalTime (50 min)
			expectError:          false,
		},
		{
			name: "Pickup time constraint not met",
			offer: model.NewOffer(
				"offer1", "user1",
				sourceCoord, destCoord,
				now, 30*time.Minute,
				3,
				*model.NewPreference(enums.Male, false),
				twoHoursLater, // maxEstimatedArrivalTime
				0,
				nil,
				nil,
			),
			request: model.NewRequest(
				"request1", "user2",
				pickupCoord, dropoffCoord,
				oneHourLater,  // earliestDepartureTime - rider wants to be picked up after one hour
				twoHoursLater, // latestArrivalTime
				10*time.Minute,
				2,
				*model.NewPreference(enums.Female, false),
			),
			mockSelectorValue: pickupdropoffcache.NewValue(
				model.NewPathPoint(pickupCoord, enums.Pickup, now, nil, 5*time.Minute),
				model.NewPathPoint(dropoffCoord, enums.Dropoff, oneHourLater, nil, 10*time.Minute),
			),
			mockSelectorErr:      nil,
			mockDrivingDurations: []time.Duration{0, 20 * time.Minute, 40 * time.Minute, 60 * time.Minute}, // Source, Pickup, Dropoff, Destination
			mockDrivingErr:       nil,
			expected:             false, // Driver arrives at pickup before rider's earliest departure time
			expectError:          false,
		},
		{
			name: "Driver arrives after rider's earliest departure time but before pickup time",
			offer: model.NewOffer(
				"offer1", "user1",
				sourceCoord, destCoord,
				now, 30*time.Minute,
				3,
				*model.NewPreference(enums.Male, false),
				twoHoursLater, // maxEstimatedArrivalTime
				0,
				nil,
				nil,
			),
			request: model.NewRequest(
				"request1", "user2",
				pickupCoord, dropoffCoord,
				quarterHourLater, // earliestDepartureTime - rider wants to be picked up after quarter-hour
				twoHoursLater,    // latestArrivalTime
				10*time.Minute,
				2,
				*model.NewPreference(enums.Female, false),
			),
			mockSelectorValue: pickupdropoffcache.NewValue(
				model.NewPathPoint(pickupCoord, enums.Pickup, now, nil, 10*time.Minute),
				model.NewPathPoint(dropoffCoord, enums.Dropoff, oneHourLater, nil, 10*time.Minute),
			),
			mockSelectorErr:      nil,
			mockDrivingDurations: []time.Duration{0, 20 * time.Minute, 40 * time.Minute, 60 * time.Minute}, // Source, Pickup, Dropoff, Destination
			mockDrivingErr:       nil,
			expected:             true,
			expectError:          false,
		},
		{
			name: "Dropoff time constraint not met",
			offer: model.NewOffer(
				"offer1", "user1",
				sourceCoord, destCoord,
				now, 30*time.Minute,
				3,
				*model.NewPreference(enums.Male, false),
				computeMaxEstimatedArrivalTime(now, 30*time.Minute), // maxEstimatedArrivalTime
				0,
				nil,
				nil,
			),
			request: model.NewRequest(
				"request1", "user2",
				pickupCoord, dropoffCoord,
				now,          // earliestDepartureTime
				oneHourLater, // latestArrivalTime - rider wants to arrive within one hour
				10*time.Minute,
				2,
				*model.NewPreference(enums.Female, false),
			),
			mockSelectorValue: pickupdropoffcache.NewValue(
				model.NewPathPoint(pickupCoord, enums.Pickup, now, nil, 5*time.Minute),
				model.NewPathPoint(dropoffCoord, enums.Dropoff, oneHourLater, nil, 10*time.Minute),
			),
			mockSelectorErr:      nil,
			mockDrivingDurations: []time.Duration{0, 20 * time.Minute, 70 * time.Minute, 90 * time.Minute}, // Source, Pickup, Dropoff, Destination
			mockDrivingErr:       nil,
			expected:             false, // Driver arrives at dropoff after rider's latest arrival time
			expectError:          false,
		},
		{
			name: "Driver arrives before rider's latest arrival time but after dropoff time",
			offer: model.NewOffer(
				"offer1", "user1",
				sourceCoord, destCoord,
				now, 30*time.Minute,
				3,
				*model.NewPreference(enums.Male, false),
				computeMaxEstimatedArrivalTime(now, 30*time.Minute), // maxEstimatedArrivalTime
				0,
				nil,
				nil,
			),
			request: model.NewRequest(
				"request1", "user2",
				pickupCoord, dropoffCoord,
				now,          // earliestDepartureTime
				oneHourLater, // latestArrivalTime - rider wants to arrive within one hour
				10*time.Minute,
				2,
				*model.NewPreference(enums.Female, false),
			),
			mockSelectorValue: pickupdropoffcache.NewValue(
				model.NewPathPoint(pickupCoord, enums.Pickup, now, nil, 5*time.Minute),
				model.NewPathPoint(dropoffCoord, enums.Dropoff, oneHourLater, nil, 10*time.Minute),
			),
			mockSelectorErr:      nil,
			mockDrivingDurations: []time.Duration{0, 20 * time.Minute, 55 * time.Minute, 60 * time.Minute}, // Source, Pickup, Dropoff, Destination
			mockDrivingErr:       nil,
			expected:             false,
			expectError:          false,
		},
		{
			name: "Exact match of driver arrival with pickup and dropoff times",
			offer: model.NewOffer(
				"offer1", "user1",
				sourceCoord, destCoord,
				now, 30*time.Minute,
				3,
				*model.NewPreference(enums.Male, false),
				computeMaxEstimatedArrivalTime(now, 30*time.Minute), // maxEstimatedArrivalTime
				0,
				nil,
				nil,
			),
			request: model.NewRequest(
				"request1", "user2",
				pickupCoord, dropoffCoord,
				now.Add(10*time.Minute), // earliestDepartureTime
				oneHourLater,            // latestArrivalTime - rider wants to arrive within one hour
				10*time.Minute,
				2,
				*model.NewPreference(enums.Female, false),
			),
			mockSelectorValue: pickupdropoffcache.NewValue(
				model.NewPathPoint(pickupCoord, enums.Pickup, now, nil, 5*time.Minute),
				model.NewPathPoint(dropoffCoord, enums.Dropoff, oneHourLater, nil, 10*time.Minute),
			),
			mockSelectorErr:      nil,
			mockDrivingDurations: []time.Duration{0, 15 * time.Minute, 50 * time.Minute, 70 * time.Minute}, // Source, Pickup, Dropoff, Destination
			mockDrivingErr:       nil,
			expected:             true, // Driver arrives at pickup and dropoff at the exact time of the pickup and dropoff times
			expectError:          false,
		},
		{
			name: "Error from selector",
			offer: model.NewOffer(
				"offer1", "user1",
				sourceCoord, destCoord,
				now, 30*time.Minute,
				3,
				*model.NewPreference(enums.Male, false),
				computeMaxEstimatedArrivalTime(now, 30*time.Minute), // maxEstimatedArrivalTime
				0,
				nil,
				nil,
			),
			request: model.NewRequest(
				"request1", "user2",
				pickupCoord, dropoffCoord,
				now,           // earliestDepartureTime
				twoHoursLater, // latestArrivalTime
				10*time.Minute,
				2,
				*model.NewPreference(enums.Female, false),
			),
			mockSelectorValue:    nil,
			mockSelectorErr:      context.DeadlineExceeded,
			mockDrivingDurations: nil,
			mockDrivingErr:       nil,
			expected:             false,
			expectError:          true,
		},
		{
			name: "Error from routing engine (driving time)",
			offer: model.NewOffer(
				"offer1", "user1",
				sourceCoord, destCoord,
				now, 30*time.Minute,
				3,
				*model.NewPreference(enums.Male, false),
				computeMaxEstimatedArrivalTime(now, 30*time.Minute), // maxEstimatedArrivalTime
				0,
				nil,
				nil,
			),
			request: model.NewRequest(
				"request1", "user2",
				pickupCoord, dropoffCoord,
				now,           // earliestDepartureTime
				twoHoursLater, // latestArrivalTime
				10*time.Minute,
				2,
				*model.NewPreference(enums.Female, false),
			),
			mockSelectorValue: pickupdropoffcache.NewValue(
				model.NewPathPoint(pickupCoord, enums.Pickup, now, nil, 5*time.Minute),
				model.NewPathPoint(dropoffCoord, enums.Dropoff, oneHourLater, nil, 10*time.Minute),
			),
			mockSelectorErr:      nil,
			mockDrivingDurations: nil,
			mockDrivingErr:       context.DeadlineExceeded,
			expected:             false,
			expectError:          true,
		},
		// This test case is no longer relevant as the direct route call is no longer made
		// Keeping it for backward compatibility but updating expectations
		{
			name: "Direct route no longer used",
			offer: model.NewOffer(
				"offer1", "user1",
				sourceCoord, destCoord,
				now, 30*time.Minute,
				3,
				*model.NewPreference(enums.Male, false),
				twoHoursLater, // maxEstimatedArrivalTime
				0,
				nil,
				nil,
			),
			request: model.NewRequest(
				"request1", "user2",
				pickupCoord, dropoffCoord,
				now,           // earliestDepartureTime
				twoHoursLater, // latestArrivalTime
				10*time.Minute,
				2,
				*model.NewPreference(enums.Female, false),
			),
			mockSelectorValue: pickupdropoffcache.NewValue(
				model.NewPathPoint(pickupCoord, enums.Pickup, now, nil, 5*time.Minute),
				model.NewPathPoint(dropoffCoord, enums.Dropoff, oneHourLater, nil, 10*time.Minute),
			),
			mockSelectorErr:      nil,
			mockDrivingDurations: []time.Duration{0, 20 * time.Minute, 40 * time.Minute, 60 * time.Minute},
			mockDrivingErr:       nil,
			expected:             true,  // Should pass as the direct route is no longer checked
			expectError:          false, // No error expected as the direct route is no longer checked
		},
	}

	// Run test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create mock selector
			mockSelector := NewMockPickupDropoffSelector(tc.mockSelectorValue, tc.mockSelectorErr)

			// Create mock routing engine
			mockEngine := NewMockEngine(
				tc.mockDrivingDurations,
				tc.mockDrivingErr,
			)

			// Create an instance of our test-specific DetourTimeChecker with our mock objects
			checker := prechecker.NewDetourTimeChecker(mockSelector, mockEngine)

			// Run the check
			result, err := checker.Check(tc.offer, tc.request)

			// Check error
			if tc.expectError && err == nil {
				t.Errorf("Expected error but got nil")
			}
			if !tc.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			// Check result
			if result != tc.expected {
				t.Errorf("Expected result %v but got %v", tc.expected, result)
			}
		})
	}
}
