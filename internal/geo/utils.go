package geo

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/paulmach/go.geojson"

	"matching-engine/internal/errors"
	"matching-engine/internal/model"
)

// -------------------------------------------------------------------
// GeoJSONBuilder — a fluent, chainable *builder pattern* for composing
// GeoJSON FeatureCollections.  Methods return the same builder so you
// can write:
//
//     gj, err := geo.NewGeoJSON().
//         AddRoute(route, "#FF0").
//         AddPoint(coord, "#0F0").
//         Build()
//
// The builder accumulates operations; the first error encountered is
// stored and returned by Build(), so you don't need to check every
// step.  If you *do* want to inspect errors earlier, call Err().
//
// Thread‑safety: NOT provided.  Wrap your own locking if you need it.
// -------------------------------------------------------------------

type GeoJSONBuilder struct {
	fc  *geojson.FeatureCollection
	err error // first error encountered (sticky)
}

// NewGeoJSON returns an empty builder ready for chaining.
func NewGeoJSON() *GeoJSONBuilder {
	return &GeoJSONBuilder{fc: geojson.NewFeatureCollection()}
}

// ParseGeoJSON initialises a builder from an existing GeoJSON payload.
// The input can be:
//   - an empty string                 → empty collection
//   - a FeatureCollection             → used as‑is
//   - a single Feature                → wrapped in a new collection
//
// Any other input is considered invalid and results in an error.
func ParseGeoJSON(existing string) (*GeoJSONBuilder, error) {
	var fc *geojson.FeatureCollection

	switch {
	case existing == "":
		// Empty input → start fresh.
		fc = geojson.NewFeatureCollection()

	default:
		// Try to unmarshal as a FeatureCollection first.
		var tmp geojson.FeatureCollection
		if err := json.Unmarshal([]byte(existing), &tmp); err == nil && len(tmp.Features) > 0 {
			fc = &tmp
			break
		}

		// Fallback: maybe it's a single Feature → wrap it.
		f, err := geojson.UnmarshalFeature([]byte(existing))
		if err != nil {
			return nil, fmt.Errorf("invalid GeoJSON: %w", err)
		}
		fc = geojson.NewFeatureCollection()
		fc.AddFeature(f)
	}

	return &GeoJSONBuilder{fc: fc}, nil
}

// Err returns the first error that occurred while building.
func (b *GeoJSONBuilder) Err() error { return b.err }

// Build produces the final GeoJSON string (or the stored error).
func (b *GeoJSONBuilder) Build() (string, error) {
	if b.err != nil {
		return "", b.err
	}
	out, err := b.fc.MarshalJSON()
	if err != nil {
		return "", fmt.Errorf("marshal GeoJSON: %w", err)
	}
	return string(out), nil
}

// AddRoute decodes the route's polyline into a LineString and appends
// it.  Colour (stroke) is optional.
func (b *GeoJSONBuilder) AddRoute(r *model.Route, color string) *GeoJSONBuilder {
	if b.err != nil { // short‑circuit if we're already in error state
		return b
	}
	if r == nil {
		b.err = errors.New("route is nil")
		return b
	}
	ls, err := r.Polyline().Coordinates()
	if err != nil {
		b.err = fmt.Errorf("decode polyline: %w", err)
		return b
	}
	return b.AddLineString(ls, color)
}

// AddLineString appends an existing line string; requires ≥2 points.
func (b *GeoJSONBuilder) AddLineString(ls model.LineString, color string) *GeoJSONBuilder {
	if b.err != nil {
		return b
	}
	if len(ls) < 2 {
		b.err = errors.New("line string needs at least 2 coordinates")
		return b
	}

	coords := make([][]float64, len(ls))
	for i, c := range ls {
		coords[i] = []float64{c.Lng(), c.Lat()} // GeoJSON wants [lon, lat]
	}

	feature := geojson.NewLineStringFeature(coords)
	if color != "" {
		feature.SetProperty("stroke", color)
	}
	b.fc.AddFeature(feature)
	return b
}

// AddPoint appends a single coordinate as a Point feature.
func (b *GeoJSONBuilder) AddPoint(c model.Coordinate, color string) *GeoJSONBuilder {
	if b.err != nil {
		return b
	}
	feature := geojson.NewPointFeature([]float64{c.Lng(), c.Lat()})
	if color != "" {
		feature.SetProperty("marker-color", color)
	}
	b.fc.AddFeature(feature)
	return b
}

// AddCircle adds a circle centered at the given coordinate with the specified radius in meters.
// The circle is approximated as a polygon with 64 vertices.
// Color (fill and stroke) is optional.
func (b *GeoJSONBuilder) AddCircle(center model.Coordinate, radiusMeters float64, color string) *GeoJSONBuilder {
	if b.err != nil {
		return b
	}

	const numPoints = 64 // Fixed number of points for the circle

	// Convert radius from meters to degrees (approximate)
	// Earth's circumference is about 40,075,000 meters at the equator
	// So 1 degree is approximately 111,319.5 meters (40,075,000 / 360)
	radiusDegrees := radiusMeters / 111319.5

	// Generate the circle as a polygon
	coords := make([][]float64, numPoints+1) // +1 to close the loop

	for i := 0; i < numPoints; i++ {
		angle := 2 * math.Pi * float64(i) / float64(numPoints)
		lat := center.Lat() + radiusDegrees*math.Sin(angle)
		lng := center.Lng() + radiusDegrees*math.Cos(angle)/math.Cos(center.Lat()*math.Pi/180)
		coords[i] = []float64{lng, lat} // GeoJSON wants [lon, lat]
	}
	// Close the loop
	coords[numPoints] = coords[0]

	// Create a polygon feature
	feature := geojson.NewPolygonFeature([][][]float64{coords})

	if color != "" {
		feature.SetProperty("stroke", color)
		feature.SetProperty("fill", color)
		feature.SetProperty("fill-opacity", 0.2) // Semi-transparent fill
	}

	b.fc.AddFeature(feature)
	return b
}

// AddIsochrone adds an isochrone polygon to the GeoJSON.
// The isochrone should contain a valid polygon geometry.
// Color (fill and stroke) is optional.
func (b *GeoJSONBuilder) AddIsochrone(isochrone *model.Isochrone, color string) *GeoJSONBuilder {
	if b.err != nil {
		return b
	}
	if isochrone == nil {
		b.err = errors.New("isochrone is nil")
		return b
	}

	polygons := isochrone.Polygons()
	if len(polygons) == 0 {
		b.err = errors.New("isochrone has no polygons")
		return b
	}

	// Process each polygon in the isochrone
	for _, polygon := range polygons {
		if len(polygon) == 0 {
			continue
		}

		// Convert coordinates to GeoJSON format (lng, lat)
		geoJsonPolygon := make([][][]float64, len(polygon))
		for i, ring := range polygon {
			geoJsonRing := make([][]float64, len(ring))
			for j, coord := range ring {
				geoJsonRing[j] = []float64{coord.Lng(), coord.Lat()}
			}
			geoJsonPolygon[i] = geoJsonRing
		}

		// Create a polygon feature
		feature := geojson.NewPolygonFeature(geoJsonPolygon)

		if color != "" {
			feature.SetProperty("stroke", color)
			feature.SetProperty("fill", color)
			feature.SetProperty("fill-opacity", 0.2) // Semi-transparent fill
		}

		b.fc.AddFeature(feature)
	}

	return b
}
