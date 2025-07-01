package tests

import (
	"context"
	"fmt"
	"matching-engine/internal/adapter/valhalla"
	"matching-engine/internal/geo"
	"matching-engine/internal/model"
	"math/rand"
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

func emust(err error) {
	if err != nil {
		panic(err)
	}
}

func writeToFile(filename, content string) error {
	return os.WriteFile(filename, []byte(content), 0644) // 0644 is rw-r--r--
}

/*
 These are some points that we can use to test the engine
 within the current tiles dataset (Andorra - andorra-latest.osm.pbf):

 Lat: 42.5078, Lon: 1.5211
 Lat: 42.5057, Lon: 1.5265
 Lat: 42.5036, Lon: 1.5148
 Lat: 42.5083, Lon: 1.5353
 Lat: 42.4636, Lon: 1.4912
 Lat: 42.5347, Lon: 1.5830
 Lat: 42.5562, Lon: 1.5339
 Lat: 42.5676, Lon: 1.5980
 Lat: 42.5669, Lon: 1.4846
 Lat: 42.5440, Lon: 1.5148
*/

func TestValhalla_PlanDrivingRoute(t *testing.T) {
	v, err := valhalla.NewValhalla()
	if err != nil {
		t.Fatalf("failed to create Valhalla engine: %v", err)
	}

	testCases := []struct {
		name       string
		routeParam *model.RouteParams
		wantErr    bool
	}{
		{
			name: "valid route",
			routeParam: must(model.NewRouteParams(
				[]model.Coordinate{
					*must(model.NewCoordinate(31.261772319501347, 29.99262985474232)),
					*must(model.NewCoordinate(31.20765210728787, 29.92392636464055)),
				},
				time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			)),
			wantErr: false,
		},
		{
			name:       "nil params",
			routeParam: nil,
			wantErr:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := v.PlanDrivingRoute(context.Background(), tc.routeParam)
			if (err != nil) != tc.wantErr {
				t.Errorf("unexpected error status: got %v, wantErr %v", err, tc.wantErr)
				return
			}

			json := geo.NewGeoJSON().AddRoute(result, "blue")
			for _, point := range tc.routeParam.Waypoints() {
				json.AddPoint(point, "red")
			}

			emust(writeToFile(
				tc.name+".json",
				must(json.Build()),
			))

			fmt.Println(result.Distance())
			fmt.Println(result.Time())
		})
	}
}

func TestValhalla_ComputeDrivingTime(t *testing.T) {
	v, err := valhalla.NewValhalla()
	if err != nil {
		t.Fatalf("failed to create Valhalla engine: %v", err)
	}

	testCases := []struct {
		name       string
		routeParam *model.RouteParams
		wantErr    bool
	}{
		{
			name: "valid time",
			routeParam: must(model.NewRouteParams(
				[]model.Coordinate{
					*must(model.NewCoordinate(29.977462461368575, 31.249469996140675)),
					*must(model.NewCoordinate(29.9811224983645, 31.250405678626862)),
					*must(model.NewCoordinate(29.97828744926288, 31.251670041058134)),
					*must(model.NewCoordinate(29.97376, 31.254408)),
				},
				time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			)),
			wantErr: false,
		},
		{
			name:       "nil params",
			routeParam: nil,
			wantErr:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := v.ComputeDrivingTime(context.Background(), tc.routeParam)
			if (err != nil) != tc.wantErr {
				t.Errorf("unexpected error status: got %v, wantErr %v", err, tc.wantErr)
				return
			}
			fmt.Println(result)
		})
	}
}

func TestValhalla_ComputeWalkingTime(t *testing.T) {
	v, err := valhalla.NewValhalla()
	if err != nil {
		t.Fatalf("failed to create Valhalla engine: %v", err)
	}

	testCases := []struct {
		name      string
		walkParam *model.WalkParams
		wantErr   bool
	}{
		{
			name: "valid walking time",
			walkParam: must(model.NewWalkParams(
				must(model.NewCoordinate(29.977462461368575, 31.249469996140675)),
				must(model.NewCoordinate(29.9811224983645, 31.250405678626862)),
			)),
			wantErr: false,
		},
		{
			name:      "nil params",
			walkParam: nil,
			wantErr:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := v.ComputeWalkingTime(context.Background(), tc.walkParam)
			if (err != nil) != tc.wantErr {
				t.Errorf("unexpected error status: got %v, wantErr %v", err, tc.wantErr)
				return
			}
			fmt.Println(result)
		})
	}
}

func TestValhalla_ComputeIsochrone(t *testing.T) {
	v, err := valhalla.NewValhalla()
	if err != nil {
		t.Fatalf("failed to create Valhalla engine: %v", err)
	}

	testCases := []struct {
		name    string
		req     *model.IsochroneParams
		wantErr bool
	}{
		{
			name: "valid isochrone",
			req: must(model.NewIsochroneParams(
				must(model.NewCoordinate(31.249469996140675, 29.977462461368575)),
				must(model.NewContour(1, model.ContourMetricDistanceInKilometers)),
				model.ProfilePedestrian,
			)),
			wantErr: false,
		},
		{
			name:    "nil params",
			req:     nil,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := v.ComputeIsochrone(context.Background(), tc.req)
			if (err != nil) != tc.wantErr {
				t.Errorf("unexpected error status: got %v, wantErr %v", err, tc.wantErr)
				return
			}

			json := geo.
				NewGeoJSON().
				AddIsochrone(result, "blue").
				AddPoint(*tc.req.Origin(), "red")

			emust(writeToFile(
				tc.name+".json",
				must(json.Build()),
			))
		})
	}
}

func TestValhalla_ComputeDistanceTimeMatrix(t *testing.T) {
	v, err := valhalla.NewValhalla()
	if err != nil {
		t.Fatalf("failed to create Valhalla engine: %v", err)
	}

	testCases := []struct {
		name    string
		req     *model.DistanceTimeMatrixParams
		wantErr bool
	}{
		{
			name: "valid matrix",
			req: must(model.NewDistanceTimeMatrixParams(
				[]model.Coordinate{
					*must(model.NewCoordinate(29.977462461368575, 31.249469996140675)),
				},
				model.ProfilePedestrian,
				model.WithDepartureTime(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)),
				model.WithTargets([]model.Coordinate{
					*must(model.NewCoordinate(29.97828744926288, 31.251670041058134)),
					*must(model.NewCoordinate(29.97376, 31.254408)),
					*must(model.NewCoordinate(29.9811224983645, 31.250405678626862)),
				}),
			)),
			wantErr: false,
		},
		{
			name:    "nil params",
			req:     nil,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := v.ComputeDistanceTimeMatrix(context.Background(), tc.req)
			if (err != nil) != tc.wantErr {
				t.Errorf("unexpected error status: got %v, wantErr %v", err, tc.wantErr)
				return
			}
			fmt.Println(result)
		})
	}
}

func TestValhalla_SnapPointToRoad(t *testing.T) {
	fmt.Println("TestValhalla_SnapPointToRoad")
	v, err := valhalla.NewValhalla()
	if err != nil {
		t.Fatalf("failed to create Valhalla engine: %v", err)
	}

	testCases := []struct {
		name    string
		point   *model.Coordinate
		wantErr bool
	}{
		{
			name:    "valid point",
			point:   must(model.NewCoordinate(29.9811, 31.2504)),
			wantErr: false,
		},
		{
			name:    "nil point",
			point:   nil,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := v.SnapPointToRoad(context.Background(), tc.point)
			if (err != nil) != tc.wantErr {
				t.Errorf("unexpected error status: got %v, wantErr %v", err, tc.wantErr)
				return
			}
			fmt.Println(result)
		})
	}
}

func TestValhalla_StressTest_ComputeDistanceTimeMatrix(t *testing.T) {
	v, err := valhalla.NewValhalla()
	if err != nil {
		t.Fatalf("failed to create Valhalla engine: %v", err)
	}

	numSources := 50
	numTargets := 50
	numRuns := 10 // Number of times to repeat the test

	// Generate sources and targets in Alexandria, Egypt bounding box
	// Approx: lat 31.18 - 31.25, lon 29.85 - 29.98
	genCoord := func() model.Coordinate {
		lat := 31.18 + (0.07 * rand.Float64()) // 31.18 - 31.25
		lon := 29.85 + (0.13 * rand.Float64()) // 29.85 - 29.98
		c, _ := model.NewCoordinate(lat, lon)
		return *c
	}

	var sources []model.Coordinate
	var targets []model.Coordinate
	for i := 0; i < numSources; i++ {
		sources = append(sources, genCoord())
	}
	for i := 0; i < numTargets; i++ {
		targets = append(targets, genCoord())
	}

	params := must(model.NewDistanceTimeMatrixParams(
		sources,
		model.ProfileCar,
		model.WithTargets(targets),
		model.WithDepartureTime(time.Now().Add(time.Minute)),
	))

	var total time.Duration
	var min, max time.Duration

	for i := 0; i < numRuns; i++ {
		start := time.Now()
		result, err := v.ComputeDistanceTimeMatrix(context.Background(), params)
		elapsed := time.Since(start)

		if err != nil {
			t.Fatalf("matrix computation failed on run %d: %v", i+1, err)
		}

		if i == 0 || elapsed < min {
			min = elapsed
		}
		if i == 0 || elapsed > max {
			max = elapsed
		}
		total += elapsed

		fmt.Printf("Run %d: %v (matrix: %dx%d)\n", i+1, elapsed, len(result.Times()), len(result.Times()[0]))
	}

	avg := total / time.Duration(numRuns)
	fmt.Printf("Average: %v, Min: %v, Max: %v (over %d runs)\n", avg, min, max, numRuns)
}
