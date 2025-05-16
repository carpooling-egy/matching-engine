package tests

import (
	"fmt"
	"matching-engine/internal/model"
	"matching-engine/internal/service/earlypruning/prechecker"
	"testing"
)

// MockChecker is a mock implementation of the Checker interface for testing
type MockChecker struct {
	result bool
	err    error
}

func NewMockChecker(result bool, err error) *MockChecker {
	return &MockChecker{
		result: result,
		err:    err,
	}
}

func (m *MockChecker) Check(offer *model.Offer, request *model.Request) (bool, error) {
	return m.result, m.err
}

func TestCompositeChecker_Check(t *testing.T) {
	// Define test cases
	tests := []struct {
		name        string
		checkers    []prechecker.Checker
		expected    bool
		expectError bool
	}{
		{
			name: "All checkers return true",
			checkers: []prechecker.Checker{
				NewMockChecker(true, nil),
				NewMockChecker(true, nil),
				NewMockChecker(true, nil),
			},
			expected:    true,
			expectError: false,
		},
		{
			name: "One checker returns false",
			checkers: []prechecker.Checker{
				NewMockChecker(true, nil),
				NewMockChecker(false, nil),
				NewMockChecker(true, nil),
			},
			expected:    false,
			expectError: false,
		},
		{
			name: "One checker returns error",
			checkers: []prechecker.Checker{
				NewMockChecker(true, nil),
				NewMockChecker(true, fmt.Errorf("test error")),
				NewMockChecker(true, nil),
			},
			expected:    false,
			expectError: true,
		},
		{
			name:        "Empty list of checkers",
			checkers:    []prechecker.Checker{},
			expected:    true,
			expectError: false,
		},
	}

	// Run test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a composite checker with the test checkers
			checker := prechecker.NewCompositePreChecker(tc.checkers...)

			// Run the check with nil offer and request (not used by mock checkers)
			result, err := checker.Check(nil, nil)

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
