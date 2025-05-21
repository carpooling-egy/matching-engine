package tests

import (
	"matching-engine/internal/enums"
	"matching-engine/internal/model"
	"matching-engine/internal/service/checker"
	"testing"
	"time"
)

func TestPreferenceChecker_Check(t *testing.T) {
	// Create a new PreferenceChecker
	checker := checker.NewPreferenceChecker()

	// Define test cases
	tests := []struct {
		name        string
		offer       *model.Offer
		request     *model.Request
		expected    bool
		expectError bool
	}{
		{
			name: "Matching preferences",
			offer: model.NewOffer(
				"offer1", "user1",
				model.Coordinate{}, model.Coordinate{},
				time.Now(), 30*time.Minute,
				3,
				*model.NewPreference(enums.Male, false),
				time.Now().Add(1*time.Hour),
				0,
				nil,
				nil,
			),
			request: model.NewRequest(
				"request1", "user2",
				model.Coordinate{}, model.Coordinate{},
				time.Now(), time.Now().Add(1*time.Hour),
				10*time.Minute,
				2,
				*model.NewPreference(enums.Female, false),
			),
			expected:    true,
			expectError: false,
		},
		{
			name: "Same gender preference not satisfied",
			offer: model.NewOffer(
				"offer1", "user1",
				model.Coordinate{}, model.Coordinate{},
				time.Now(), 30*time.Minute,
				3,
				*model.NewPreference(enums.Male, true),
				time.Now().Add(1*time.Hour),
				0,
				nil,
				nil,
			),
			request: model.NewRequest(
				"request1", "user2",
				model.Coordinate{}, model.Coordinate{},
				time.Now(), time.Now().Add(1*time.Hour),
				10*time.Minute,
				2,
				*model.NewPreference(enums.Female, false),
			),
			expected:    false,
			expectError: false,
		},
		{
			name: "Request with same gender preference not satisfied",
			offer: model.NewOffer(
				"offer1", "user1",
				model.Coordinate{}, model.Coordinate{},
				time.Now(), 30*time.Minute,
				3,
				*model.NewPreference(enums.Male, false),
				time.Now().Add(1*time.Hour),
				0,
				nil,
				nil,
			),
			request: model.NewRequest(
				"request1", "user2",
				model.Coordinate{}, model.Coordinate{},
				time.Now(), time.Now().Add(1*time.Hour),
				10*time.Minute,
				2,
				*model.NewPreference(enums.Female, true),
			),
			expected:    false,
			expectError: false,
		},
		{
			name: "Offer with matched request having incompatible preferences",
			offer: func() *model.Offer {
				// Create an offer with a matched request that has incompatible preferences
				incompatibleRequest := model.NewRequest(
					"request2", "user3",
					model.Coordinate{}, model.Coordinate{},
					time.Now(), time.Now().Add(1*time.Hour),
					10*time.Minute,
					1,
					*model.NewPreference(enums.Male, false),
				)

				offer := model.NewOffer(
					"offer1", "user1",
					model.Coordinate{}, model.Coordinate{},
					time.Now(), 30*time.Minute,
					3,
					*model.NewPreference(enums.Male, false),
					time.Now().Add(1*time.Hour),
					1,
					nil,
					[]*model.Request{incompatibleRequest},
				)

				return offer
			}(),
			request: model.NewRequest(
				"request1", "user2",
				model.Coordinate{}, model.Coordinate{},
				time.Now(), time.Now().Add(1*time.Hour),
				10*time.Minute,
				2,
				*model.NewPreference(enums.Female, true), // requires the same gender (female only)
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
