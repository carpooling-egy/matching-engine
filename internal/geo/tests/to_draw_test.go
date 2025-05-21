package tests

import (
	"encoding/json"
	"matching-engine/internal/model"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type feature struct {
	coords  []model.Coordinate
	geoType string
	name    string
	color   string
}

// writeGeoJSON writes a GeoJSON FeatureCollection with multiple features.
func writeGeoJSON(filename string, features []feature) error {
	geoJSONFeatures := make([]map[string]interface{}, len(features))
	for i, f := range features {
		var coordinates interface{}

		switch f.geoType {
		case "Point":
			// For Point, use a simple [lng, lat] array
			if len(f.coords) > 0 {
				coordinates = []float64{f.coords[0].Lng(), f.coords[0].Lat()}
			}
		case "LineString":
			// For LineString, use array of [lng, lat] arrays
			pts := make([][2]float64, len(f.coords))
			for j, c := range f.coords {
				pts[j] = [2]float64{c.Lng(), c.Lat()}
			}
			coordinates = pts
		case "Polygon":
			// For Polygon, wrap coordinates in an additional array (needed structure)
			pts := make([][2]float64, len(f.coords))
			for j, c := range f.coords {
				pts[j] = [2]float64{c.Lng(), c.Lat()}
			}
			coordinates = [][][2]float64{pts} // Array of linear rings
		}

		geoJSONFeatures[i] = map[string]interface{}{
			"type": "Feature",
			"geometry": map[string]interface{}{
				"type":        f.geoType,
				"coordinates": coordinates,
			},
			"properties": map[string]interface{}{
				"name":  f.name,
				"color": f.color,
			},
		}
	}

	gj := map[string]interface{}{
		"type":     "FeatureCollection",
		"features": geoJSONFeatures,
	}
	data, err := json.MarshalIndent(gj, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func TestGeoJsonn(t *testing.T) {
	// Test cases for GeoJSON parsing and validation
	feats := []feature{
		{
			coords:  []model.Coordinate{*must(model.NewCoordinate(42.43, 1.42))},
			geoType: "Point",
			name:    "Test Point",
			color:   "#000000",
		},
		{
			coords:  []model.Coordinate{*must(model.NewCoordinate(42.43, 1.42)), *must(model.NewCoordinate(42.6, 1.7))},
			geoType: "LineString",
			name:    "Test LineString",
			color:   "#FF0000",
		},
	}

	// Write GeoJSON to file
	err := writeGeoJSON("test.json", feats)
	if err != nil {
	}
}

// DrawPolyline reads an encoded polyline from a file and writes it as GeoJSON
func DrawPolyline(inputFilename, outputFilename string, colorHex string) error {
	// Read the polyline from file
	polylineStr, err := readPolylineFromFile(inputFilename)
	if err != nil {
		return err
	}

	// Create a polyline model from the encoded string
	polyline, err := model.NewPolyline(polylineStr)
	if err != nil {
		return err
	}

	// Get coordinates from the polyline
	coords, err := polyline.Coordinates()
	if err != nil {
		return err
	}

	// Convert coords to []model.Coordinate
	modelCoords := make([]model.Coordinate, len(coords))
	for i, coord := range coords {
		modelCoords[i] = coord
	}

	// Create a feature for the polyline
	polylineFeature := feature{
		coords:  modelCoords,
		geoType: "LineString",
		name:    "Polyline Path",
		color:   colorHex,
	}

	// Write to GeoJSON file
	return writeGeoJSON(outputFilename, []feature{polylineFeature})
}

// readPolylineFromFile reads a polyline string from a file
func readPolylineFromFile(filename string) (string, error) {
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

// Usage example:
func DrawSamplePolyline() error {
	return DrawPolyline("path.txt", "polyline_visualization.json", "#0000FF") // Blue line
}

// DrawCircle creates a GeoJSON circle (as polygon) around a given point
func DrawCircle(centerLat, centerLng, radiusMeters float64, outputFilename, colorHex string) error {
	// Number of points to use to approximate the circle
	const numPoints = 64

	// Earth radius in meters
	const earthRadius = 6371000.0

	// Convert radius from meters to degrees (approximate)
	radiusLat := (radiusMeters / earthRadius) * (180.0 / math.Pi)
	radiusLng := (radiusMeters / earthRadius) * (180.0 / math.Pi) / math.Cos(centerLat*math.Pi/180.0)

	// Generate points around the circle
	coords := make([]model.Coordinate, numPoints+1)
	for i := 0; i < numPoints; i++ {
		angle := 2.0 * math.Pi * float64(i) / float64(numPoints)
		lat := centerLat + radiusLat*math.Sin(angle)
		lng := centerLng + radiusLng*math.Cos(angle)
		coord, _ := model.NewCoordinate(lat, lng)
		coords[i] = *coord
	}
	// Close the polygon by repeating the first point
	coords[numPoints] = coords[0]

	// Create features for circle and center point
	features := []feature{
		{
			coords:  coords,
			geoType: "Polygon",
			name:    "Search Radius",
			color:   colorHex,
		},
		{
			coords:  []model.Coordinate{*must(model.NewCoordinate(centerLat, centerLng))},
			geoType: "Point",
			name:    "Search Center",
			color:   "#000000", // Black for the center point
		},
	}

	return writeGeoJSON(outputFilename, features)
}

func DrawPolylineWithSearchArea(
	polylineFile string,
	outputFile string,
	centerLat, centerLng, radiusDegrees float64,
) error {
	// Read the polyline from file
	polylineStr, err := readPolylineFromFile(polylineFile)
	if err != nil {
		return err
	}

	// Create a polyline model from the encoded string
	polyline, err := model.NewPolyline(polylineStr)
	if err != nil {
		return err
	}

	// Get coordinates from the polyline
	lineCoords, err := polyline.Coordinates()
	if err != nil {
		return err
	}

	// Convert coords to []model.Coordinate
	modelLineCoords := make([]model.Coordinate, len(lineCoords))
	for i, coord := range lineCoords {
		modelLineCoords[i] = coord
	}

	// Generate circle points
	const numPoints = 64

	// Use the radius directly in degrees
	radiusLat := radiusDegrees
	// Adjust longitude radius to account for latitude (circles appear as ellipses in Mercator)
	radiusLng := radiusDegrees / math.Cos(centerLat*math.Pi/180.0)

	circleCoords := make([]model.Coordinate, numPoints+1)
	for i := 0; i < numPoints; i++ {
		angle := 2.0 * math.Pi * float64(i) / float64(numPoints)
		lat := centerLat + radiusLat*math.Sin(angle)
		lng := centerLng + radiusLng*math.Cos(angle)

		coord, _ := model.NewCoordinate(lat, lng)
		circleCoords[i] = *coord
	}
	circleCoords[numPoints] = circleCoords[0]

	// Create features
	features := []feature{
		{
			coords:  modelLineCoords,
			geoType: "LineString",
			name:    "Polyline Path",
			color:   "#0000FF", // Blue
		},
		{
			coords:  circleCoords,
			geoType: "Polygon",
			name:    "Search Radius",
			color:   "#FF000080", // Semi-transparent red
		},
		{
			coords:  []model.Coordinate{*must(model.NewCoordinate(centerLat, centerLng))},
			geoType: "Point",
			name:    "Search Center",
			color:   "#000000", // Black
		},
	}

	return writeGeoJSON(outputFile, features)
}

// Usage example for visualization with search area
func DrawSampleWithSearchArea() error {
	return DrawPolylineWithSearchArea(
		"internal/geo/tests/path.txt",
		"polyline_with_search.json",
		42.578417, 1.660117, 0.005,
	)
}
