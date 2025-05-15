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
	// For the regular route
	regularDurations []time.Duration
	regularErr       error

	// For the direct route
	directDurations []time.Duration
	directErr       error
}

// NewMockEngine creates a new MockEngine
func NewMockEngine(
	regularDurations []time.Duration,
	regularErr error,
	directDurations []time.Duration,
	directErr error,
) routing.Engine {
	return &MockEngine{
		regularDurations: regularDurations,
		regularErr:       regularErr,
		directDurations:  directDurations,
		directErr:        directErr,
	}
}

// ComputeDrivingTime implements the routing.Engine interface
func (m *MockEngine) ComputeDrivingTime(ctx context.Context, routeParams *model.RouteParams) ([]time.Duration, error) {
	// Check if this is the direct route call (only 2 waypoints)
	if len(routeParams.Waypoints()) == 2 {
		return m.directDurations, m.directErr
	}
	// Otherwise, it's the full route call
	return m.regularDurations, m.regularErr
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
		mockDirectDurations  []time.Duration
		mockDirectErr        error
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
				*model.NewPreference(enums.Male, false, false, false),
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
				*model.NewPreference(enums.Female, false, false, false),
			),
			mockSelectorValue: pickupdropoffcache.NewValue(
				model.NewPathPoint(pickupCoord, enums.Pickup, now, nil),
				model.NewPathPoint(dropoffCoord, enums.Dropoff, oneHourLater, nil),
				5*time.Minute,  // pickupWalkingDuration
				10*time.Minute, // dropoffWalkingDuration
			),
			mockSelectorErr:      nil,
			mockDrivingDurations: []time.Duration{0, 20 * time.Minute, 40 * time.Minute, 60 * time.Minute}, // Source, Pickup, Dropoff, Destination
			mockDrivingErr:       nil,
			mockDirectDurations:  []time.Duration{0, 40 * time.Minute}, // Source, Destination (direct route)
			mockDirectErr:        nil,
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
				*model.NewPreference(enums.Male, false, false, false),
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
				*model.NewPreference(enums.Female, false, false, false),
			),
			mockSelectorValue: pickupdropoffcache.NewValue(
				model.NewPathPoint(pickupCoord, enums.Pickup, now, nil),
				model.NewPathPoint(dropoffCoord, enums.Dropoff, oneHourLater, nil),
				5*time.Minute,  // pickupWalkingDuration
				10*time.Minute, // dropoffWalkingDuration
			),
			mockSelectorErr:      nil,
			mockDrivingDurations: []time.Duration{0, 20 * time.Minute, 40 * time.Minute, 60 * time.Minute}, // Source, Pickup, Dropoff, Destination
			mockDrivingErr:       nil,
			mockDirectDurations:  []time.Duration{0, 30 * time.Minute}, // Source, Destination (direct route)
			mockDirectErr:        nil,
			expected:             false, // Detour is 30 minutes (60-30), which exceeds the 10 minutes allowed
			expectError:          false,
		},
		{
			name: "Pickup time constraint not met",
			offer: model.NewOffer(
				"offer1", "user1",
				sourceCoord, destCoord,
				now, 30*time.Minute,
				3,
				*model.NewPreference(enums.Male, false, false, false),
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
				*model.NewPreference(enums.Female, false, false, false),
			),
			mockSelectorValue: pickupdropoffcache.NewValue(
				model.NewPathPoint(pickupCoord, enums.Pickup, now, nil),
				model.NewPathPoint(dropoffCoord, enums.Dropoff, oneHourLater, nil),
				5*time.Minute,  // pickupWalkingDuration
				10*time.Minute, // dropoffWalkingDuration
			),
			mockSelectorErr:      nil,
			mockDrivingDurations: []time.Duration{0, 20 * time.Minute, 40 * time.Minute, 60 * time.Minute}, // Source, Pickup, Dropoff, Destination
			mockDrivingErr:       nil,
			mockDirectDurations:  []time.Duration{0, 40 * time.Minute}, // Source, Destination (direct route)
			mockDirectErr:        nil,
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
				*model.NewPreference(enums.Male, false, false, false),
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
				*model.NewPreference(enums.Female, false, false, false),
			),
			mockSelectorValue: pickupdropoffcache.NewValue(
				model.NewPathPoint(pickupCoord, enums.Pickup, now, nil),
				model.NewPathPoint(dropoffCoord, enums.Dropoff, oneHourLater, nil),
				10*time.Minute, // pickupWalkingDuration
				10*time.Minute, // dropoffWalkingDuration
			),
			mockSelectorErr:      nil,
			mockDrivingDurations: []time.Duration{0, 20 * time.Minute, 40 * time.Minute, 60 * time.Minute}, // Source, Pickup, Dropoff, Destination
			mockDrivingErr:       nil,
			mockDirectDurations:  []time.Duration{0, 40 * time.Minute}, // Source, Destination (direct route)
			mockDirectErr:        nil,
			expected:             false,
			expectError:          false,
		},
		{
			name: "Dropoff time constraint not met",
			offer: model.NewOffer(
				"offer1", "user1",
				sourceCoord, destCoord,
				now, 30*time.Minute,
				3,
				*model.NewPreference(enums.Male, false, false, false),
				twoHoursLater, // maxEstimatedArrivalTime
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
				*model.NewPreference(enums.Female, false, false, false),
			),
			mockSelectorValue: pickupdropoffcache.NewValue(
				model.NewPathPoint(pickupCoord, enums.Pickup, now, nil),
				model.NewPathPoint(dropoffCoord, enums.Dropoff, oneHourLater, nil),
				5*time.Minute,  // pickupWalkingDuration
				10*time.Minute, // dropoffWalkingDuration
			),
			mockSelectorErr:      nil,
			mockDrivingDurations: []time.Duration{0, 20 * time.Minute, 70 * time.Minute, 90 * time.Minute}, // Source, Pickup, Dropoff, Destination
			mockDrivingErr:       nil,
			mockDirectDurations:  []time.Duration{0, 40 * time.Minute}, // Source, Destination (direct route)
			mockDirectErr:        nil,
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
				*model.NewPreference(enums.Male, false, false, false),
				twoHoursLater, // maxEstimatedArrivalTime
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
				*model.NewPreference(enums.Female, false, false, false),
			),
			mockSelectorValue: pickupdropoffcache.NewValue(
				model.NewPathPoint(pickupCoord, enums.Pickup, now, nil),
				model.NewPathPoint(dropoffCoord, enums.Dropoff, oneHourLater, nil),
				5*time.Minute,  // pickupWalkingDuration
				10*time.Minute, // dropoffWalkingDuration
			),
			mockSelectorErr:      nil,
			mockDrivingDurations: []time.Duration{0, 20 * time.Minute, 55 * time.Minute, 60 * time.Minute}, // Source, Pickup, Dropoff, Destination
			mockDrivingErr:       nil,
			mockDirectDurations:  []time.Duration{0, 40 * time.Minute}, // Source, Destination (direct route)
			mockDirectErr:        nil,
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
				*model.NewPreference(enums.Male, false, false, false),
				twoHoursLater, // maxEstimatedArrivalTime
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
				*model.NewPreference(enums.Female, false, false, false),
			),
			mockSelectorValue: pickupdropoffcache.NewValue(
				model.NewPathPoint(pickupCoord, enums.Pickup, now, nil),
				model.NewPathPoint(dropoffCoord, enums.Dropoff, oneHourLater, nil),
				5*time.Minute,  // pickupWalkingDuration
				10*time.Minute, // dropoffWalkingDuration
			),
			mockSelectorErr:      nil,
			mockDrivingDurations: []time.Duration{0, 15 * time.Minute, 50 * time.Minute, 70 * time.Minute}, // Source, Pickup, Dropoff, Destination
			mockDrivingErr:       nil,
			mockDirectDurations:  []time.Duration{0, 40 * time.Minute}, // Source, Destination (direct route)
			mockDirectErr:        nil,
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
				*model.NewPreference(enums.Male, false, false, false),
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
				*model.NewPreference(enums.Female, false, false, false),
			),
			mockSelectorValue:    nil,
			mockSelectorErr:      context.DeadlineExceeded,
			mockDrivingDurations: nil,
			mockDrivingErr:       nil,
			mockDirectDurations:  nil,
			mockDirectErr:        nil,
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
				*model.NewPreference(enums.Male, false, false, false),
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
				*model.NewPreference(enums.Female, false, false, false),
			),
			mockSelectorValue: pickupdropoffcache.NewValue(
				model.NewPathPoint(pickupCoord, enums.Pickup, now, nil),
				model.NewPathPoint(dropoffCoord, enums.Dropoff, oneHourLater, nil),
				5*time.Minute,  // pickupWalkingDuration
				10*time.Minute, // dropoffWalkingDuration
			),
			mockSelectorErr:      nil,
			mockDrivingDurations: nil,
			mockDrivingErr:       context.DeadlineExceeded,
			mockDirectDurations:  nil,
			mockDirectErr:        nil,
			expected:             false,
			expectError:          true,
		},
		{
			name: "Error from routing engine (direct route)",
			offer: model.NewOffer(
				"offer1", "user1",
				sourceCoord, destCoord,
				now, 30*time.Minute,
				3,
				*model.NewPreference(enums.Male, false, false, false),
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
				*model.NewPreference(enums.Female, false, false, false),
			),
			mockSelectorValue: pickupdropoffcache.NewValue(
				model.NewPathPoint(pickupCoord, enums.Pickup, now, nil),
				model.NewPathPoint(dropoffCoord, enums.Dropoff, oneHourLater, nil),
				5*time.Minute,  // pickupWalkingDuration
				10*time.Minute, // dropoffWalkingDuration
			),
			mockSelectorErr:      nil,
			mockDrivingDurations: []time.Duration{0, 20 * time.Minute, 40 * time.Minute, 60 * time.Minute},
			mockDrivingErr:       nil,
			mockDirectDurations:  nil,
			mockDirectErr:        context.DeadlineExceeded,
			expected:             false,
			expectError:          true,
		},
	}

	// Run test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create mock selector
			mockSelector := NewMockPickupDropoffSelector(tc.mockSelectorValue, tc.mockSelectorErr)

			// Create mock routing engine with two different behaviors
			mockEngine := NewMockEngine(
				tc.mockDrivingDurations,
				tc.mockDrivingErr,
				tc.mockDirectDurations,
				tc.mockDirectErr,
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
