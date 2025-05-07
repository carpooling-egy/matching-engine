package tests

import (
	"context"
	"fmt"
	"matching-engine/internal/adapter/valhalla"
	"matching-engine/internal/model"
	"testing"
	"time"
)

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
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
					*must(model.NewCoordinate(42.43, 1.42)),
					*must(model.NewCoordinate(42.6, 1.7)),
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
			fmt.Println(result)
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
					*must(model.NewCoordinate(42.5078, 1.5211)),
					*must(model.NewCoordinate(42.5057, 1.5265)),
					*must(model.NewCoordinate(42.5057, 1.5650)),
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
				must(model.NewCoordinate(42.5669, 1.4846)),
				must(model.NewCoordinate(42.5440, 1.5148)),
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
				must(model.NewCoordinate(42.5347, 1.5830)),
				must(model.NewContour(30, model.ContourMetricDistanceInKilometers)),
				model.Pedestrian,
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
			fmt.Println(result)
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
					*must(model.NewCoordinate(42.5078, 1.5211)),
					*must(model.NewCoordinate(42.5440, 1.5148)),
					*must(model.NewCoordinate(42.5057, 1.5265)),
					*must(model.NewCoordinate(42.5440, 1.5148)),
					*must(model.NewCoordinate(42.5057, 1.5265)),
				},
				time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
				model.Auto,
				model.WithTargets([]model.Coordinate{
					*must(model.NewCoordinate(42.5057, 1.5265)),
					*must(model.NewCoordinate(42.5440, 1.5148)),
					*must(model.NewCoordinate(42.5057, 1.5265)),
					*must(model.NewCoordinate(42.5036, 1.5148)),
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
			point:   must(model.NewCoordinate(42.50828, 1.53210)),
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
