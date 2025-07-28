package tests

import (
	"matching-engine/internal/enums"
	"matching-engine/internal/model"
	"matching-engine/internal/service/checker"
	"testing"
	"time"
)

func TestOverlapChecker_Check(t *testing.T) {
	// Create a new OverlapChecker
	checker := checker.NewOverlapChecker()

	// Define test times for clarity
	now := time.Now()
	oneHourLater := now.Add(1 * time.Hour)
	twoHoursLater := now.Add(2 * time.Hour)
	threeHoursLater := now.Add(3 * time.Hour)

	// Define test cases
	tests := []struct {
		name        string
		offer       *model.Offer
		request     *model.Request
		expected    bool
		expectError bool
	}{
		{
			name: "Overlapping time slots",
			offer: model.NewOffer(
				"offer1", "user1",
				model.Coordinate{}, model.Coordinate{},
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
				model.Coordinate{}, model.Coordinate{},
				oneHourLater,    // earliestDepartureTime
				threeHoursLater, // latestArrivalTime
				10*time.Minute,
				2,
				*model.NewPreference(enums.Female, false),
			),
			expected:    true,
			expectError: false,
		},
		{
			name: "Non-overlapping time slots - request too early",
			offer: model.NewOffer(
				"offer1", "user1",
				model.Coordinate{}, model.Coordinate{},
				twoHoursLater, // departureTime
				30*time.Minute,
				3,
				*model.NewPreference(enums.Male, false),
				threeHoursLater, // maxEstimatedArrivalTime
				0,
				nil,
				nil,
			),
			request: model.NewRequest(
				"request1", "user2",
				model.Coordinate{}, model.Coordinate{},
				now,          // earliestDepartureTime
				oneHourLater, // latestArrivalTime
				10*time.Minute,
				2,
				*model.NewPreference(enums.Female, false),
			),
			expected:    false,
			expectError: false,
		},
		{
			name: "Non-overlapping time slots - request too late",
			offer: model.NewOffer(
				"offer1", "user1",
				model.Coordinate{}, model.Coordinate{},
				now, // departureTime
				30*time.Minute,
				3,
				*model.NewPreference(enums.Male, false),
				oneHourLater, // maxEstimatedArrivalTime
				0,
				nil,
				nil,
			),
			request: model.NewRequest(
				"request1", "user2",
				model.Coordinate{}, model.Coordinate{},
				twoHoursLater,   // earliestDepartureTime
				threeHoursLater, // latestArrivalTime
				10*time.Minute,
				2,
				*model.NewPreference(enums.Female, false),
			),
			expected:    false,
			expectError: false,
		},
		{
			name:        "Nil offer",
			offer:       nil,
			request:     &model.Request{},
			expected:    false,
			expectError: true,
		},
		{
			name:        "Nil request",
			offer:       &model.Offer{},
			request:     nil,
			expected:    false,
			expectError: true,
		},
	}

	// Run test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
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
