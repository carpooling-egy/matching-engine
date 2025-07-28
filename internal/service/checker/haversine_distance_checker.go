package checker

import (
	"github.com/umahmood/haversine"
	"matching-engine/internal/app/config"
	"matching-engine/internal/model"
	"time"
)

type HaversineDistanceChecker struct {
}

func NewHaversineDistanceChecker() Checker {
	return &HaversineDistanceChecker{}
}

func (e HaversineDistanceChecker) Check(offer *model.Offer, request *model.Request) (bool, error) {

	driverSource := haversine.Coord{Lat: offer.Source().Lat(), Lon: offer.Source().Lng()}
	driverDestination := haversine.Coord{Lat: offer.Destination().Lat(), Lon: offer.Destination().Lng()}
	requestSource := haversine.Coord{Lat: request.Source().Lat(), Lon: request.Source().Lng()}
	requestDestination := haversine.Coord{Lat: request.Destination().Lat(), Lon: request.Destination().Lng()}

	// Check if the offer can reach the request source before the latest arrival time
	_, driverToRequestSource := haversine.Distance(driverSource, requestSource)
	_, requestSourceToRequestDestination := haversine.Distance(requestSource, requestDestination)
	driverToRequestDestination := driverToRequestSource + requestSourceToRequestDestination
	timeToRequestDestination := e.convertDistanceToTime(driverToRequestDestination)
	if offer.DepartureTime().Add(timeToRequestDestination).After(request.LatestArrivalTime()) {
		return false, nil
	}

	// convert the detour time to distance
	detourDistance := e.convertTimeToDistance(offer.DetourDurationMinutes())

	// Calculate the total distance including detour
	_, requestDestinationToDriverDestination := haversine.Distance(requestDestination, driverDestination)
	totalDistance := driverToRequestSource + requestSourceToRequestDestination + requestDestinationToDriverDestination
	driverDirectDistance, _ := haversine.Distance(driverSource, driverDestination)
	return totalDistance <= driverDirectDistance+detourDistance, nil
}

func (e HaversineDistanceChecker) convertTimeToDistance(t time.Duration) float64 {
	speedKmh := config.GetEnvFloat("FIXED_SPEED_KMH", 27) // Use a fixed speed from config or default to 27 km/h
	hours := t.Hours()
	distance := hours * speedKmh // Convert time to distance
	return distance
}

func (e HaversineDistanceChecker) convertDistanceToTime(d float64) time.Duration {
	speedKmh := config.GetEnvFloat("FIXED_SPEED_KMH", 27) // Use a fixed speed from config or default to 27 km/h
	hours := d / speedKmh                                 // Convert distance to time in hours
	return time.Duration(hours * float64(time.Hour))      // Convert hours to time.Duration
}
