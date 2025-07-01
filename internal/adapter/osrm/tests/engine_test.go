package tests

import (
	"context"
	"fmt"
	"matching-engine/internal/adapter/osrm"
	"matching-engine/internal/geo"
	"matching-engine/internal/model"
	"math/rand"
	"os"
	"sync"
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
 within the current tiles dataset (New York - new-york-latest.osm.pbf):

 Lat: 40.7128, Lon: -74.0060 (New York City)
 Lat: 40.7589, Lon: -73.9851 (Times Square)
 Lat: 40.7505, Lon: -73.9934 (Penn Station)
 Lat: 40.7484, Lon: -73.9857 (Madison Square Garden)
 Lat: 40.7527, Lon: -73.9772 (Grand Central Terminal)
 Lat: 40.7589, Lon: -73.9851 (Empire State Building)
 Lat: 40.7484, Lon: -73.9857 (Bryant Park)
 Lat: 40.7505, Lon: -73.9934 (Herald Square)
 Lat: 40.7527, Lon: -73.9772 (Rockefeller Center)
 Lat: 40.7589, Lon: -73.9851 (Central Park)
*/

func TestOSRM_PlanDrivingRoute(t *testing.T) {
	o, err := osrm.NewOSRM()
	if err != nil {
		t.Fatalf("failed to create OSRM engine: %v", err)
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
					*must(model.NewCoordinate(31.21869167193043, 29.942667902383192)),
					*must(model.NewCoordinate(31.20186826886348, 29.901443738997273)),
				},
				time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			)),
			wantErr: false,
		},
		{
			name: "multi-point route",
			routeParam: must(model.NewRouteParams(
				[]model.Coordinate{
					*must(model.NewCoordinate(40.7128, -74.0060)), // New York City
					*must(model.NewCoordinate(40.7589, -73.9851)), // Times Square
					*must(model.NewCoordinate(40.7505, -73.9934)), // Penn Station
					*must(model.NewCoordinate(40.7484, -73.9857)), // Madison Square Garden
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
			result, err := o.PlanDrivingRoute(context.Background(), tc.routeParam)
			if (err != nil) != tc.wantErr {
				t.Errorf("unexpected error status: got %v, wantErr %v", err, tc.wantErr)
				return
			}

			if !tc.wantErr && result != nil {
				json := geo.NewGeoJSON().AddRoute(result, "blue")
				for _, point := range tc.routeParam.Waypoints() {
					json.AddPoint(point, "red")
				}

				emust(writeToFile(
					tc.name+".json",
					must(json.Build()),
				))

				fmt.Printf("Route distance: %.2f km\n", result.Distance().Value()/1000)
				fmt.Printf("Route time: %v\n", result.Time())
			}
		})
	}
}

func TestOSRM_ComputeDrivingTime(t *testing.T) {
	o, err := osrm.NewOSRM()
	if err != nil {
		t.Fatalf("failed to create OSRM engine: %v", err)
	}

	testCases := []struct {
		name       string
		routeParam *model.RouteParams
		wantErr    bool
	}{
		{
			name: "valid driving time",
			routeParam: must(model.NewRouteParams(
				[]model.Coordinate{
					*must(model.NewCoordinate(40.7128, -74.0060)), // New York City
					*must(model.NewCoordinate(40.7589, -73.9851)), // Times Square
				},
				time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			)),
			wantErr: false,
		},
		{
			name: "multi-point driving time",
			routeParam: must(model.NewRouteParams(
				[]model.Coordinate{
					*must(model.NewCoordinate(40.7128, -74.0060)), // New York City
					*must(model.NewCoordinate(40.7589, -73.9851)), // Times Square
					*must(model.NewCoordinate(40.7505, -73.9934)), // Penn Station
					*must(model.NewCoordinate(40.7484, -73.9857)), // Madison Square Garden
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
			result, err := o.ComputeDrivingTime(context.Background(), tc.routeParam)
			if (err != nil) != tc.wantErr {
				t.Errorf("unexpected error status: got %v, wantErr %v", err, tc.wantErr)
				return
			}

			if !tc.wantErr && result != nil {
				fmt.Printf("Driving times: %v\n", result)
			}
		})
	}
}

func TestOSRM_ComputeWalkingTime(t *testing.T) {
	o, err := osrm.NewOSRM()
	if err != nil {
		t.Fatalf("failed to create OSRM engine: %v", err)
	}

	testCases := []struct {
		name      string
		walkParam *model.WalkParams
		wantErr   bool
	}{
		{
			name: "valid walking time",
			walkParam: must(model.NewWalkParams(
				must(model.NewCoordinate(31.219824800885945, 29.94221584806283)),
				must(model.NewCoordinate(31.2279084030221, 29.94136620165611)),
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
			result, err := o.ComputeWalkingTime(context.Background(), tc.walkParam)
			if (err != nil) != tc.wantErr {
				t.Errorf("unexpected error status: got %v, wantErr %v", err, tc.wantErr)
				return
			}

			if !tc.wantErr {
				fmt.Printf("Walking time: %v\n", result)
			}
		})
	}
}

func TestOSRM_ComputeDistanceTimeMatrix(t *testing.T) {
	o, err := osrm.NewOSRM()
	if err != nil {
		t.Fatalf("failed to create OSRM engine: %v", err)
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
					*must(model.NewCoordinate(40.7128, -74.0060)), // New York City
				},
				model.ProfileCar,
				model.WithDepartureTime(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)),
				model.WithTargets([]model.Coordinate{
					*must(model.NewCoordinate(40.7589, -73.9851)), // Times Square
					*must(model.NewCoordinate(40.7505, -73.9934)), // Penn Station
					*must(model.NewCoordinate(40.7484, -73.9857)), // Madison Square Garden
				}),
			)),
			wantErr: false,
		},
		{
			name: "car profile matrix",
			req: must(model.NewDistanceTimeMatrixParams(
				[]model.Coordinate{
					*must(model.NewCoordinate(40.7128, -74.0060)), // New York City
				},
				model.ProfileCar,
				model.WithDepartureTime(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)),
				model.WithTargets([]model.Coordinate{
					*must(model.NewCoordinate(40.7589, -73.9851)), // Times Square
					*must(model.NewCoordinate(40.7505, -73.9934)), // Penn Station
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
			result, err := o.ComputeDistanceTimeMatrix(context.Background(), tc.req)
			if (err != nil) != tc.wantErr {
				t.Errorf("unexpected error status: got %v, wantErr %v", err, tc.wantErr)
				return
			}

			if !tc.wantErr && result != nil {
				fmt.Printf("Matrix result: %+v\n", result)
			}
		})
	}
}

func TestOSRM_SnapPointToRoad(t *testing.T) {
	o, err := osrm.NewOSRM()
	if err != nil {
		t.Fatalf("failed to create OSRM engine: %v", err)
	}

	testCases := []struct {
		name    string
		point   *model.Coordinate
		wantErr bool
	}{
		{
			name:    "valid point",
			point:   must(model.NewCoordinate(40.7128, -74.0060)), // New York City
			wantErr: false,
		},
		{
			name:    "another valid point",
			point:   must(model.NewCoordinate(40.7589, -73.9851)), // Times Square
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
			result, err := o.SnapPointToRoad(context.Background(), tc.point)
			if (err != nil) != tc.wantErr {
				t.Errorf("unexpected error status: got %v, wantErr %v", err, tc.wantErr)
				return
			}

			if !tc.wantErr && result != nil {
				json := geo.NewGeoJSON().
					AddPoint(*tc.point, "green").
					AddPoint(*result, "blue")

				emust(writeToFile(
					tc.name+".json",
					must(json.Build()),
				))

				fmt.Printf("Snapped point: %.6f, %.6f\n", result.Lat(), result.Lng())
			}
		})
	}
}

func TestOSRM_StressTest_ComputeDistanceTimeMatrix(t *testing.T) {
	o, err := osrm.NewOSRM()
	if err != nil {
		t.Fatalf("failed to create OSRM engine: %v", err)
	}

	numSources := 100
	numTargets := 100
	numRuns := 100

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

	var wg sync.WaitGroup
	durations := make([]time.Duration, numRuns)
	errs := make([]error, numRuns)
	mu := sync.Mutex{}

	for i := 0; i < numRuns; i++ {
		wg.Add(1)
		go func(runIdx int) {
			defer wg.Done()
			start := time.Now()
			result, err := o.ComputeDistanceTimeMatrix(context.Background(), params)
			elapsed := time.Since(start)
			mu.Lock()
			durations[runIdx] = elapsed
			errs[runIdx] = err
			mu.Unlock()
			if err == nil {
				fmt.Printf("Run %d: %v (matrix: %dx%d)\n", runIdx+1, elapsed, len(result.Times()), len(result.Times()[0]))
			} else {
				fmt.Printf("Run %d: ERROR: %v\n", runIdx+1, err)
			}
		}(i)
	}
	wg.Wait()

	// Calculate stats
	var total time.Duration
	var min, max time.Duration
	first := true
	for i, d := range durations {
		if errs[i] != nil {
			t.Fatalf("matrix computation failed on run %d: %v", i+1, errs[i])
		}
		if first || d < min {
			min = d
		}
		if first || d > max {
			max = d
		}
		total += d
		first = false
	}
	avg := total / time.Duration(numRuns)
	fmt.Printf("Average: %v, Min: %v, Max: %v (over %d runs, parallel)\n", avg, min, max, numRuns)
}

func TestOSRM_StressTest_PlanDrivingRoute(t *testing.T) {
	o, err := osrm.NewOSRM()
	if err != nil {
		t.Fatalf("failed to create OSRM engine: %v", err)
	}

	numRuns := 1000

	genCoord := func() model.Coordinate {
		lat := 22.0 + (9.5 * rand.Float64())  // 22.0 - 31.5
		lon := 25.0 + (10.0 * rand.Float64()) // 25.0 - 35.0
		c, _ := model.NewCoordinate(lat, lon)
		return *c
	}

	var wg sync.WaitGroup
	durations := make([]time.Duration, numRuns)
	errs := make([]error, numRuns)
	mu := sync.Mutex{}

	for i := 0; i < numRuns; i++ {
		wg.Add(1)
		go func(runIdx int) {
			defer wg.Done()
			startCoord := genCoord()
			endCoord := genCoord()
			params := must(model.NewRouteParams(
				[]model.Coordinate{startCoord, endCoord},
				time.Now().Add(time.Minute),
			))
			start := time.Now()
			result, err := o.PlanDrivingRoute(context.Background(), params)
			elapsed := time.Since(start)
			mu.Lock()
			durations[runIdx] = elapsed
			errs[runIdx] = err
			mu.Unlock()
			if err == nil {
				fmt.Printf("Run %d: %v (distance: %.2f km)\n", runIdx+1, elapsed, result.Distance().Value()/1000)
			} else {
				fmt.Printf("Run %d: ERROR: %v\n", runIdx+1, err)
			}
		}(i)
	}
	wg.Wait()

	// Calculate stats
	var total time.Duration
	var min, max time.Duration
	first := true
	for i, d := range durations {
		if errs[i] != nil {
			t.Fatalf("PlanDrivingRoute failed on run %d: %v", i+1, errs[i])
		}
		if first || d < min {
			min = d
		}
		if first || d > max {
			max = d
		}
		total += d
		first = false
	}
	avg := total / time.Duration(numRuns)
	fmt.Printf("PlanDrivingRoute: Average: %v, Min: %v, Max: %v (over %d runs, parallel)\n", avg, min, max, numRuns)
}

func TestOSRM_StressTest_SnapPointToRoad(t *testing.T) {
	o, err := osrm.NewOSRM()
	if err != nil {
		t.Fatalf("failed to create OSRM engine: %v", err)
	}

	numRuns := 10000

	genCoord := func() model.Coordinate {
		lat := 22.0 + (9.5 * rand.Float64())  // 22.0 - 31.5
		lon := 25.0 + (10.0 * rand.Float64()) // 25.0 - 35.0
		c, _ := model.NewCoordinate(lat, lon)
		return *c
	}

	var wg sync.WaitGroup
	durations := make([]time.Duration, numRuns)
	errs := make([]error, numRuns)
	mu := sync.Mutex{}

	for i := 0; i < numRuns; i++ {
		wg.Add(1)
		go func(runIdx int) {
			defer wg.Done()
			point := genCoord()
			start := time.Now()
			result, err := o.SnapPointToRoad(context.Background(), &point)
			elapsed := time.Since(start)
			mu.Lock()
			durations[runIdx] = elapsed
			errs[runIdx] = err
			mu.Unlock()
			if err == nil {
				fmt.Printf("Run %d: %v (snapped: %.6f, %.6f)\n", runIdx+1, elapsed, result.Lat(), result.Lng())
			} else {
				fmt.Printf("Run %d: ERROR: %v\n", runIdx+1, err)
			}
		}(i)
	}
	wg.Wait()

	// Calculate stats
	var total time.Duration
	var min, max time.Duration
	first := true
	for i, d := range durations {
		if errs[i] != nil {
			t.Fatalf("SnapPointToRoad failed on run %d: %v", i+1, errs[i])
		}
		if first || d < min {
			min = d
		}
		if first || d > max {
			max = d
		}
		total += d
		first = false
	}
	avg := total / time.Duration(numRuns)
	fmt.Printf("SnapPointToRoad: Average: %v, Min: %v, Max: %v (over %d runs, parallel)\n", avg, min, max, numRuns)
}
