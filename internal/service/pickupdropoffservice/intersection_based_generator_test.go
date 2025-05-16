package pickupdropoffservice

import (
	"fmt"
	"matching-engine/internal/geo/processor"
	"matching-engine/internal/model"
	"strings"
	"testing"
	"time"
)

// MockGeospatialProcessor is a mock implementation of the GeospatialProcessor interface
type MockGeospatialProcessor struct {
	pickupCoord          *model.Coordinate
	pickupDuration       time.Duration
	pickupErr            error
	dropoffCoord         *model.Coordinate
	dropoffDuration      time.Duration
	dropoffErr           error
	pickupCallCount      int
	dropoffCallCount     int
	lastPointRequested   *model.Coordinate
	lastWalkingRequested time.Duration
}

func NewMockGeospatialProcessor(
	pickupCoord *model.Coordinate,
	pickupDuration time.Duration,
	pickupErr error,
	dropoffCoord *model.Coordinate,
	dropoffDuration time.Duration,
	dropoffErr error,
) *MockGeospatialProcessor {
	return &MockGeospatialProcessor{
		pickupCoord:     pickupCoord,
		pickupDuration:  pickupDuration,
		pickupErr:       pickupErr,
		dropoffCoord:    dropoffCoord,
		dropoffDuration: dropoffDuration,
		dropoffErr:      dropoffErr,
	}
}

func (m *MockGeospatialProcessor) ComputeClosestRoutePoint(
	point *model.Coordinate,
	walkingTime time.Duration,
) (*model.Coordinate, time.Duration, error) {
	m.lastPointRequested = point
	m.lastWalkingRequested = walkingTime

	// Determine if this is a pickup or dropoff call based on the point
	// This is a simple heuristic for the test - in real code we'd need a more robust way
	// to distinguish between pickup and dropoff calls
	if m.pickupCallCount == 0 {
		m.pickupCallCount++
		return m.pickupCoord, m.pickupDuration, m.pickupErr
	} else {
		m.dropoffCallCount++
		return m.dropoffCoord, m.dropoffDuration, m.dropoffErr
	}
}

// MockProcessorFactory is a mock implementation of processor.ProcessorFactory for testing
type MockProcessorFactory struct {
	processor processor.GeospatialProcessor
	err       error
}

func NewMockProcessorFactory(processor processor.GeospatialProcessor, err error) *MockProcessorFactory {
	return &MockProcessorFactory{
		processor: processor,
		err:       err,
	}
}

func (m *MockProcessorFactory) CreateProcessor(offer *model.Offer) (processor.GeospatialProcessor, error) {
	return m.processor, m.err
}

// We're now using the real IntersectionBasedGenerator instead of a test-specific implementation

func TestIntersectionBasedGenerator_GeneratePickupDropoffPoints(t *testing.T) {
	// Create test coordinates
	sourceCoord, _ := model.NewCoordinate(1.0, 1.0)
	destCoord, _ := model.NewCoordinate(2.0, 2.0)
	pickupCoord, _ := model.NewCoordinate(1.1, 1.1)
	dropoffCoord, _ := model.NewCoordinate(1.9, 1.9)

	// Create test times
	now := time.Now()
	later := now.Add(1 * time.Hour)

	// Define test cases
	tests := []struct {
		name                 string
		request              *model.Request
		offer                *model.Offer
		pickupCoord          *model.Coordinate
		pickupDuration       time.Duration
		pickupErr            error
		dropoffCoord         *model.Coordinate
		dropoffDuration      time.Duration
		dropoffErr           error
		factoryErr           error
		expectedPickupCoord  *model.Coordinate
		expectedDropoffCoord *model.Coordinate
		expectError          bool
		errorContains        string
	}{
		{
			name: "Both pickup and dropoff have intersections",
			request: model.NewRequest(
				"request1",
				"user1",
				*sourceCoord,
				*destCoord,
				now,
				later,
				15*time.Minute, // Max walking duration
				1,
				model.Preference{},
			),
			offer: model.NewOffer(
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
			),
			pickupCoord:          pickupCoord,
			pickupDuration:       10 * time.Minute, // Less than max walking duration
			pickupErr:            nil,
			dropoffCoord:         dropoffCoord,
			dropoffDuration:      10 * time.Minute, // Less than max walking duration
			dropoffErr:           nil,
			factoryErr:           nil,
			expectedPickupCoord:  pickupCoord,  // Should use computed route point
			expectedDropoffCoord: dropoffCoord, // Should use computed route point
			expectError:          false,
		},
		{
			name: "Pickup has intersection but dropoff doesn't",
			request: model.NewRequest(
				"request2",
				"user1",
				*sourceCoord,
				*destCoord,
				now,
				later,
				15*time.Minute, // Max walking duration
				1,
				model.Preference{},
			),
			offer: model.NewOffer(
				"offer2",
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
			),
			pickupCoord:          pickupCoord,
			pickupDuration:       10 * time.Minute, // Less than max walking duration
			pickupErr:            nil,
			dropoffCoord:         dropoffCoord,
			dropoffDuration:      20 * time.Minute, // More than max walking duration
			dropoffErr:           nil,
			factoryErr:           nil,
			expectedPickupCoord:  pickupCoord, // Should use computed route point
			expectedDropoffCoord: destCoord,   // Should use original destination
			expectError:          false,
		},
		{
			name: "Pickup doesn't have intersection but dropoff does",
			request: model.NewRequest(
				"request3",
				"user1",
				*sourceCoord,
				*destCoord,
				now,
				later,
				15*time.Minute, // Max walking duration
				1,
				model.Preference{},
			),
			offer: model.NewOffer(
				"offer3",
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
			),
			pickupCoord:          pickupCoord,
			pickupDuration:       20 * time.Minute, // More than max walking duration
			pickupErr:            nil,
			dropoffCoord:         dropoffCoord,
			dropoffDuration:      10 * time.Minute, // Less than max walking duration
			dropoffErr:           nil,
			factoryErr:           nil,
			expectedPickupCoord:  sourceCoord,  // Should use original source
			expectedDropoffCoord: dropoffCoord, // Should use computed route point
			expectError:          false,
		},
		{
			name: "Neither pickup nor dropoff have intersections",
			request: model.NewRequest(
				"request4",
				"user1",
				*sourceCoord,
				*destCoord,
				now,
				later,
				15*time.Minute, // Max walking duration
				1,
				model.Preference{},
			),
			offer: model.NewOffer(
				"offer4",
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
			),
			pickupCoord:          pickupCoord,
			pickupDuration:       20 * time.Minute, // More than max walking duration
			pickupErr:            nil,
			dropoffCoord:         dropoffCoord,
			dropoffDuration:      20 * time.Minute, // More than max walking duration
			dropoffErr:           nil,
			factoryErr:           nil,
			expectedPickupCoord:  sourceCoord, // Should use original source
			expectedDropoffCoord: destCoord,   // Should use original destination
			expectError:          false,
		},
		{
			name:          "Nil request",
			request:       nil,
			offer:         model.NewOffer("offer5", "user2", *sourceCoord, *destCoord, now, 30*time.Minute, 4, model.Preference{}, later, 0, nil, nil),
			expectError:   true,
			errorContains: "request or offer is nil",
		},
		{
			name:          "Nil offer",
			request:       model.NewRequest("request5", "user1", *sourceCoord, *destCoord, now, later, 15*time.Minute, 1, model.Preference{}),
			offer:         nil,
			expectError:   true,
			errorContains: "request or offer is nil",
		},
		{
			name: "Error creating processor",
			request: model.NewRequest(
				"request6",
				"user1",
				*sourceCoord,
				*destCoord,
				now,
				later,
				15*time.Minute,
				1,
				model.Preference{},
			),
			offer: model.NewOffer(
				"offer6",
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
			),
			factoryErr:    fmt.Errorf("factory error"),
			expectError:   true,
			errorContains: "failed to create processor",
		},
		{
			name: "Error computing pickup point",
			request: model.NewRequest(
				"request7",
				"user1",
				*sourceCoord,
				*destCoord,
				now,
				later,
				15*time.Minute,
				1,
				model.Preference{},
			),
			offer: model.NewOffer(
				"offer7",
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
			),
			pickupErr:     fmt.Errorf("pickup error"),
			expectError:   true,
			errorContains: "failed to compute closest route point",
		},
		{
			name: "Error computing dropoff point",
			request: model.NewRequest(
				"request8",
				"user1",
				*sourceCoord,
				*destCoord,
				now,
				later,
				15*time.Minute,
				1,
				model.Preference{},
			),
			offer: model.NewOffer(
				"offer8",
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
			),
			pickupCoord:    pickupCoord,
			pickupDuration: 10 * time.Minute,
			dropoffErr:     fmt.Errorf("dropoff error"),
			expectError:    true,
			errorContains:  "failed to compute closest route point",
		},
	}

	// Run test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create mock processor
			mockProcessor := NewMockGeospatialProcessor(
				tc.pickupCoord,
				tc.pickupDuration,
				tc.pickupErr,
				tc.dropoffCoord,
				tc.dropoffDuration,
				tc.dropoffErr,
			)

			// Create mock factory
			mockFactory := NewMockProcessorFactory(mockProcessor, tc.factoryErr)

			// Create generator using the real implementation
			generator := NewIntersectionBasedGenerator(mockFactory)

			// Call the method
			pickup, dropoff, err := generator.GeneratePickupDropoffPoints(tc.request, tc.offer)

			// Check error
			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error but got nil")
				} else if tc.errorContains != "" && !containsString(err.Error(), tc.errorContains) {
					t.Errorf("Expected error to contain '%s' but got: %v", tc.errorContains, err)
				}
				return
			}

			if err != nil {
				t.Errorf("Expected no error but got: %v", err)
				return
			}

			// Check pickup point
			if pickup == nil {
				t.Errorf("Expected pickup point but got nil")
			} else if !coordinatesEqual(pickup.Coordinate(), tc.expectedPickupCoord) {
				t.Errorf("Expected pickup coordinate %v but got %v", tc.expectedPickupCoord, pickup.Coordinate())
			}

			// Check dropoff point
			if dropoff == nil {
				t.Errorf("Expected dropoff point but got nil")
			} else if !coordinatesEqual(dropoff.Coordinate(), tc.expectedDropoffCoord) {
				t.Errorf("Expected dropoff coordinate %v but got %v", tc.expectedDropoffCoord, dropoff.Coordinate())
			}

			// Check that the processor was called with the correct parameters
			if !tc.expectError && tc.request != nil && tc.offer != nil {
				if mockProcessor.pickupCallCount != 1 {
					t.Errorf("Expected processor to be called once for pickup, but was called %d times", mockProcessor.pickupCallCount)
				}
				if mockProcessor.dropoffCallCount != 1 {
					t.Errorf("Expected processor to be called once for dropoff, but was called %d times", mockProcessor.dropoffCallCount)
				}
			}
		})
	}
}

// Helper function to check if a string contains another string
func containsString(s, substr string) bool {
	return strings.Contains(s, substr)
}

// Helper function to check if two coordinates are equal
func coordinatesEqual(c1, c2 *model.Coordinate) bool {
	if c1 == nil && c2 == nil {
		return true
	}
	if c1 == nil || c2 == nil {
		return false
	}
	return c1.Lat() == c2.Lat() && c1.Lng() == c2.Lng()
}
