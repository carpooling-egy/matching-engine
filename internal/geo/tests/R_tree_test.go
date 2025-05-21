package tests

import (
	"matching-engine/internal/geo"
	"matching-engine/internal/geo/pruning"
	"matching-engine/internal/model"

	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// readPolylineFromFile reads a polyline string from a file
func readPolylineFromFileTest(filename string) (string, error) {
	// Get the absolute path to the file
	absPath, err := filepath.Abs(filename)
	if err != nil {
		return "", err
	}

	// Read the entire file content
	content, err := os.ReadFile(absPath)
	if err != nil {
		return "", err
	}

	// Trim whitespace and return
	return strings.TrimSpace(string(content)), nil
}
func TestBuildTreeAndSearch(t *testing.T) {
	polylineStr, _ := readPolylineFromFileTest("path.txt")

	tests := []struct {
		name          string
		polylineStr   string
		queryLat      float64
		queryLng      float64
		radius        float64
		expectResults bool
	}{
		{
			name:          "Point near polyline",
			polylineStr:   polylineStr, // Sample polyline (SF to NYC)
			queryLat:      42.588417,   // Los Angeles area (far from polyline)
			queryLng:      1.630117,
			radius:        0.01, // Small radius
			expectResults: true,
		},
		{
			name:          "Point far from polyline",
			polylineStr:   polylineStr, // Sample polyline (SF to NYC)
			queryLat:      42.588417,   // Los Angeles area (far from polyline)
			queryLng:      1.630117,
			radius:        0.005, // Very small radius
			expectResults: false,
		},
		{
			name:          "Point with large radius",
			polylineStr:   polylineStr, // Sample polyline (SF to NYC)
			queryLat:      42.578417,
			queryLng:      1.660117,
			radius:        0.01, // Large radius should capture segments
			expectResults: true,
		},
	}
	// Create polyline from encoded string
	polyline, err := model.NewPolyline(polylineStr)
	if err != nil {
		t.Fatalf("Failed to create polyline: %v", err)
	}
	coords, _ := polyline.Coordinates()
	factory := pruning.NewRTreePrunerFactory()
	pruner, err := factory.NewRoutePruner(coords)
	if err != nil {
		t.Fatalf("Failed to create pruner: %v", err)
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			// Create query point
			queryPoint, err := model.NewCoordinate(tc.queryLat, tc.queryLng)
			if err != nil {
				t.Fatalf("Failed to create query point: %v", err)
			}
			// Convert radius in degrees to walking time
			// 1. Convert degrees to meters (using great-circle distance)
			radiusInMeters := tc.radius * geo.EarthRadiusInMeters * math.Pi / 180.0

			// 2. Calculate walking time (distance/speed = time)
			walkingTimeSeconds := radiusInMeters / geo.WalkingSpeedMPS

			// 3. Create time.Duration
			walkingTime := time.Duration(walkingTimeSeconds * float64(time.Second))

			// Call the function to test with time duration
			results, err := pruner.Prune(queryPoint, walkingTime)
			print(len(results), " results found")
			if err != nil {
				t.Fatalf("BuildTreeAndSearch failed: %v", err)
			}

			// Check results
			if tc.expectResults && len(results) == 0 {
				t.Errorf("Expected to find results, but got none")
			}

			if !tc.expectResults && len(results) > 0 {
				t.Errorf("Expected to find no results, but got %d", len(results))
			}

			// Check if the results are within the expected radius and other conditions are out of the radius
			if len(results) > 0 {
				// 1. Verify that all returned coordinates are actually part of the original route
				coordMap := make(map[string]bool)
				for _, c := range coords {
					key := fmt.Sprintf("%.6f:%.6f", c.Lat(), c.Lng())
					coordMap[key] = true
				}

				for i, result := range results {
					key := fmt.Sprintf("%.6f:%.6f", result.Lat(), result.Lng())
					if !coordMap[key] {
						t.Errorf("Result %d contains coordinate not in original route: %v", i, result)
					}
				}

				// 2. Check that all results are within the expected radius
				for i, result := range results {
					// Calculate distance between query point and result
					latDiff := result.Lat() - queryPoint.Lat()
					lngDiff := result.Lng() - queryPoint.Lng()
					distSquared := latDiff*latDiff + lngDiff*lngDiff

					// Convert search radius to squared value for comparison
					// (avoid expensive sqrt operations)
					radiusSquared := tc.radius * tc.radius

					// Allow small margin of error for floating point comparisons
					if distSquared > (radiusSquared * 1.01) { // 1% margin of error
						t.Errorf("Result %d at %v is outside search radius (dist²: %f, radius²: %f)",
							i, result, distSquared, radiusSquared)
					}
				}

				// 3. Check that no coordinates within radius were missed
				// This is optional and more intensive
				for i, c := range coords {
					// Calculate if coordinate is in radius
					latDiff := c.Lat() - queryPoint.Lat()
					lngDiff := c.Lng() - queryPoint.Lng()
					distSquared := latDiff*latDiff + lngDiff*lngDiff
					radiusSquared := tc.radius * tc.radius
					epsilon := 1e-10
					if distSquared < (radiusSquared - epsilon) {
						// This coordinate should be in results
						found := false
						for _, result := range results {
							if c.Lat() == result.Lat() && c.Lng() == result.Lng() {
								found = true
								break
							}
						}
						if !found {
							t.Errorf("Coordinate %d within radius was not included in results: %v, %f %f %t",
								i, c, distSquared, radiusSquared, (distSquared < radiusSquared))
						}
					}
				}
			}
		})
	}
}

// TestNewSegment tests the NewSegment function
func TestNewSegment(t *testing.T) {
	// Create two coordinates
	a, _ := model.NewCoordinate(37.7749, -122.4194) // San Francisco
	b, _ := model.NewCoordinate(40.7128, -74.0060)  // New York

	// Create a segment
	segment := pruning.NewSegment(*a, *b)

	// Test that segment has correct coordinates
	if segment.A.Lat() != a.Lat() || segment.A.Lng() != a.Lng() {
		t.Errorf("Segment point A incorrect: got %v, want %v", segment.A, *a)
	}

	if segment.B.Lat() != b.Lat() || segment.B.Lng() != b.Lng() {
		t.Errorf("Segment point B incorrect: got %v, want %v", segment.B, *b)
	}

	// Test that bounds are correctly calculated
	bounds := segment.Bounds()

	// Check that bounds contain both points
	minLat := math.Min(a.Lat(), b.Lat())
	maxLat := math.Max(a.Lat(), b.Lat())
	minLng := math.Min(a.Lng(), b.Lng())
	maxLng := math.Max(a.Lng(), b.Lng())

	if bounds.PointCoord(0) != minLat ||
		bounds.LengthsCoord(0) != maxLat-minLat ||
		bounds.PointCoord(1) != minLng ||
		bounds.LengthsCoord(1) != maxLng-minLng {
		t.Errorf("Segment bounds incorrect: got %v", bounds)
	}
}

// TestNewPolyline tests the NewPolyline function
func TestNewPolyline(t *testing.T) {
	DrawSampleWithSearchArea()
}
