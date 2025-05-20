package tests

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"matching-engine/internal/adapter/valhalla"
	"matching-engine/internal/geo/downsampling"
	"matching-engine/internal/model"
	"math"
	"os"
	"testing"
	"time"
)

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func generateHugePolyline(numPoints int) model.LineString {
	coords := make([]model.Coordinate, numPoints)
	startLat, startLng := 42.43, 1.42
	endLat, endLng := 42.6, 1.7
	deltaLat := (endLat - startLat) / float64(numPoints-1)
	deltaLng := (endLng - startLng) / float64(numPoints-1)

	for i := 0; i < numPoints; i++ {
		t := float64(i) / float64(numPoints-1)
		lat := startLat + float64(i)*deltaLat
		lng := startLng + float64(i)*deltaLng + 0.01*math.Sin(20*math.Pi*t)
		coords[i] = *must(model.NewCoordinate(lat, lng))
	}

	return coords
}

func generateRoute(coords model.LineString) model.LineString {
	v := must(valhalla.NewValhalla())
	routeParam := must(model.NewRouteParams(
		[]model.Coordinate{
			*must(model.NewCoordinate(42.43, 1.42)),
			*must(model.NewCoordinate(42.6, 1.7)),
		},
		time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
	))
	route := must(v.PlanDrivingRoute(context.Background(), routeParam))
	return must(route.Polyline().Coordinates())
}

func TestGenerateGeoJSON(t *testing.T) {
	coords := generateHugePolyline(1_000_000)

	samplers := []downsampling.RouteDownSampler{
		downsampling.NewTimeThresholdDownSampler(downsampling.WithInterval(30 * time.Minute)),
		downsampling.NewRDPDownSampler(downsampling.WithEpsilonMeters(10)),
	}

	names := []string{
		"Original",
		"TimeThreshold",
		"RDP",
	}

	err := GenerateGeoJSON(coords, samplers, names, "simplified_routes.json")
	if err != nil {
		t.Fatalf("failed to generate GeoJSON: %v", err)
	}

	log.Info().Msg("GeoJSON written to simplified_routes.json")
}

type feature struct {
	coords []model.Coordinate
	name   string
	color  string
}

func writeGeoJSON(filename string, features []feature) error {
	geoJSONFeatures := make([]map[string]interface{}, len(features))
	for i, f := range features {
		pts := make([][2]float64, len(f.coords))
		for j, c := range f.coords {
			pts[j] = [2]float64{c.Lng(), c.Lat()}
		}
		geoJSONFeatures[i] = map[string]interface{}{
			"type": "Feature",
			"geometry": map[string]interface{}{
				"type":        "LineString",
				"coordinates": pts,
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

func GenerateGeoJSON(
	route model.LineString,
	samplers []downsampling.RouteDownSampler,
	names []string,
	outputFile string,
) error {
	if len(names) != len(samplers)+1 {
		return fmt.Errorf("expected %d names (1 for original + %d for samplers), got %d", len(samplers)+1, len(samplers), len(names))
	}

	colors := []string{
		"#FF0000", "#00FF00", "#0000FF", "#000000", "#800080",
		"#008080", "#FF00FF", "#A52A2A", "#808080", "#FFFF00",
	}
	if len(samplers)+1 > len(colors) {
		return errors.New("not enough colors for all features")
	}

	features := []feature{
		{
			coords: route,
			name:   names[0],
			color:  colors[0],
		},
	}

	for i, sampler := range samplers {
		start := time.Now()
		simplifiedRoute, err := sampler.DownSample(route)
		elapsed := time.Since(start)

		log.Info().Msgf(
			"Sampler %s took %s to downsample the route",
			names[i+1], elapsed,
		)

		if err != nil {
			return fmt.Errorf("sampler %s failed: %v", names[i+1], err)
		}
		log.Info().Msgf(
			"Sampler %s: Original route size: %d, Simplified route size: %d",
			names[i+1], len(route), len(simplifiedRoute),
		)
		features = append(features, feature{
			coords: simplifiedRoute,
			name:   names[i+1],
			color:  colors[i+1],
		})
	}

	return writeGeoJSON(outputFile, features)
}
