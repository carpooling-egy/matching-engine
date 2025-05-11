package pruning

import (
    "fmt"
    "time"

    "matching-engine/internal/geo"
    "matching-engine/internal/model"

    "github.com/dhconnelly/rtreego"
)

// RTreePruner implements the RoutePruner interface using R-tree spatial indexing
type RTreePruner struct {
    tree *rtreego.Rtree
}

// Prune filters the route to include only segments within the specified threshold
// from the origin point, implementing the RoutePruner interface
func (p *RTreePruner) Prune(origin *model.Coordinate, threshold time.Duration) (model.LineString, error) {
    // Convert query point to rtreego.Point
    query := rtreego.Point{origin.Lat(), origin.Lng()}
    
    // Convert time threshold to distance in degrees
    thresholdDistance := geo.MetersToDegrees(float64(threshold.Seconds()) * geo.WalkingSpeedMPS)

    // Create bounding box (square) around the circle for initial filtering
    searchRect, err := rtreego.NewRect(
        rtreego.Point{query[0] - thresholdDistance, query[1] - thresholdDistance},
        []float64{2 * thresholdDistance, 2 * thresholdDistance},
    )
    if err != nil {
        return model.LineString{}, err
    }

    // Get candidate segments using R-tree spatial index
    results := p.tree.SearchIntersect(searchRect)
    
    // Use a map to collect unique coordinates within the threshold
    coordMap := make(map[string]model.Coordinate)
    thresholdDistanceSquared := thresholdDistance * thresholdDistance
    
    // Process each segment found in the search area
    for _, result := range results {
        segment, ok := result.(*Segment)
        if !ok {
            continue // Skip if not a Segment
        }
        
        // Check if segment endpoints are within the threshold distance
        if isWithinThreshold(origin, &segment.A, thresholdDistanceSquared) {
            addToCoordinateMap(coordMap, segment.A)
        }
        
        if isWithinThreshold(origin, &segment.B, thresholdDistanceSquared) {
            addToCoordinateMap(coordMap, segment.B)
        }
    }
    
    // Convert the map of unique coordinates to a LineString
    return mapToLineString(coordMap), nil
}

// isWithinThreshold checks if a coordinate is within the squared threshold distance
func isWithinThreshold(origin *model.Coordinate, point *model.Coordinate, thresholdDistanceSquared float64) bool {
    return squaredDistance(origin, point) <= thresholdDistanceSquared
}

// addToCoordinateMap adds a coordinate to the map using a formatted string key
func addToCoordinateMap(coordMap map[string]model.Coordinate, coord model.Coordinate) {
    key := fmt.Sprintf("%.6f:%.6f", coord.Lat(), coord.Lng())
    coordMap[key] = coord
}

// mapToLineString converts a map of coordinates to a LineString
func mapToLineString(coordMap map[string]model.Coordinate) model.LineString {
    coords := make([]model.Coordinate, 0, len(coordMap))
    for _, coord := range coordMap {
        coords = append(coords, coord)
    }
    return coords
}

// squaredDistance calculates the squared Euclidean distance between two coordinates
func squaredDistance(a *model.Coordinate, b *model.Coordinate) float64 {
    latDiff := a.Lat() - b.Lat()
    lngDiff := a.Lng() - b.Lng()
    return latDiff*latDiff + lngDiff*lngDiff
}