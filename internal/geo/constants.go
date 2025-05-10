package geo

import "math"

const WalkingSpeedMPS = 1.4
const EarthRadiusInMeters = 6371000

func MetersPerSecondToDegreesPerSecond(mps float64) float64 {
	return mps * 180 / (EarthRadiusInMeters * math.Pi)
}
