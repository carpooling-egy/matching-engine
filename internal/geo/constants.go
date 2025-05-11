package geo

import "math"

const WalkingSpeedMPS = 1.4
const EarthRadiusInMeters = 6371000

func MetersToDegrees(meters float64) float64 {
    return meters * 180 / (EarthRadiusInMeters * math.Pi)
}
