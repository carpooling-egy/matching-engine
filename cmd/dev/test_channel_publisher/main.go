package main

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"matching-engine/internal/adapter/messaging"
	"matching-engine/internal/adapter/messaging/natsjetstream"
	"matching-engine/internal/model"

	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Create a publisher using the messaging interface
	var publisher messaging.Publisher
	var err error

	// Use the existing NATS publisher implementation
	publisher, err = natsjetstream.NewNATSPublisher()
	if err != nil {
		fmt.Printf("Error connecting to NATS: %v\n", err)
		return
	}
	defer func() {
		if err := publisher.Close(); err != nil {
			log.Printf("Error closing publisher: %v", err)
		}
	}()

	// Set up signal handling for graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalChan
		fmt.Println("\nReceived shutdown signal. Closing publisher...")
		if err := publisher.Close(); err != nil {
			log.Printf("Error closing publisher: %v", err)
		}
		os.Exit(0)
	}()

	// Generate and publish 100 matching results
	results := generateMatchingResults(40)
	fmt.Printf("Publishing %d matching results...\n", len(results))

	startTime := time.Now()
	err = publisher.PublishMatchingResults(results)

	if err != nil {
		log.Warn().Err(err).Msg("Failed to publish matching results")
	}
	duration := time.Since(startTime)

	fmt.Printf("Successfully published %d matching results in %v\n", len(results), duration)

	time.Sleep(30 * time.Second) // Wait for a second to ensure messages are sent

	// Generate and publish 100 matching results
	results = generateMatchingResults(40)
	fmt.Printf("Publishing %d matching results...\n", len(results))

	startTime = time.Now()
	err = publisher.PublishMatchingResults(results)

	if err != nil {
		log.Warn().Err(err).Msg("Failed to publish matching results")
	}
	duration = time.Since(startTime)

	fmt.Printf("Returned from publish fn")

	fmt.Println("Press Ctrl+C to exit")

	// Keep the program running until interrupted
	select {}
}

// generateMatchingResults creates a specified number of sample matching results
func generateMatchingResults(count int) []*model.MatchingResult {
	results := make([]*model.MatchingResult, count)

	for i := 0; i < count; i++ {
		// Create a matching result with a unique offer ID
		offerId := fmt.Sprintf("offer-%d", i+1)
		result := model.NewMatchingResult(offerId)

		// Add some random assigned requests
		numRequests := 1 + rand.Intn(3) // 1-3 requests per matching result
		for j := 0; j < numRequests; j++ {
			requestId := fmt.Sprintf("request-%d-%d", i+1, j+1)

			result.AddAssignedMatchedRequest(generateMatchedRequest(requestId))

			// Add points to the path
			addPointsToResult(result, requestId)
		}

		results[i] = result
	}

	return results
}

// generateMatchedRequest creates a random matched request with pickup and dropoff points
func generateMatchedRequest(requestId string) model.MatchedRequest {
	// Generate random coordinates for pickup
	pickupLat := randomCoordinate(-90, 90)
	pickupLng := randomCoordinate(-180, 180)
	pickupCoord, err := model.NewCoordinate(pickupLat, pickupLng)
	if err != nil {
		// In case of invalid coordinates, use default values
		pickupCoord, _ = model.NewCoordinate(40.7128, -74.0060) // NYC coordinates
	}

	// Generate random coordinates for dropoff
	dropoffLat := randomCoordinate(-90, 90)
	dropoffLng := randomCoordinate(-180, 180)
	dropoffCoord, err := model.NewCoordinate(dropoffLat, dropoffLng)
	if err != nil {
		// In case of invalid coordinates, use default values
		dropoffCoord, _ = model.NewCoordinate(34.0522, -118.2437) // LA coordinates
	}

	// Create pickup and dropoff points
	now := time.Now()
	pickupTime := now.Add(time.Duration(rand.Intn(60)) * time.Minute)
	dropoffTime := pickupTime.Add(time.Duration(rand.Intn(120)+30) * time.Minute)

	pickupPoint := model.NewPoint(requestId, pickupCoord, pickupTime, model.PickupPoint)
	dropoffPoint := model.NewPoint(requestId, dropoffCoord, dropoffTime, model.DropoffPoint)

	// Create a new matched request and add the points
	matchedRequest := model.NewMatchedRequestWithRequestId(requestId)
	matchedRequest.AddPickupPoint(pickupPoint)
	matchedRequest.AddDropoffPoint(dropoffPoint)

	return *matchedRequest
}

// addPointsToResult adds pickup and dropoff points to a matching result
func addPointsToResult(result *model.MatchingResult, requestId string) {
	// Generate random coordinates for pickup
	pickupLat := randomCoordinate(-90, 90)
	pickupLng := randomCoordinate(-180, 180)
	pickupCoord, err := model.NewCoordinate(pickupLat, pickupLng)
	if err != nil {
		// In case of invalid coordinates, use default values
		pickupCoord, _ = model.NewCoordinate(40.7128, -74.0060) // NYC coordinates
	}

	// Generate random coordinates for dropoff
	dropoffLat := randomCoordinate(-90, 90)
	dropoffLng := randomCoordinate(-180, 180)
	dropoffCoord, err := model.NewCoordinate(dropoffLat, dropoffLng)
	if err != nil {
		// In case of invalid coordinates, use default values
		dropoffCoord, _ = model.NewCoordinate(34.0522, -118.2437) // LA coordinates
	}

	// Create pickup and dropoff points
	now := time.Now()
	pickupTime := now.Add(time.Duration(rand.Intn(60)) * time.Minute)
	dropoffTime := pickupTime.Add(time.Duration(rand.Intn(120)+30) * time.Minute)

	pickupPoint := model.NewPoint(requestId, pickupCoord, pickupTime, model.PickupPoint)
	dropoffPoint := model.NewPoint(requestId, dropoffCoord, dropoffTime, model.DropoffPoint)

	// Add points to the path
	result.AddPoint(pickupPoint)
	result.AddPoint(dropoffPoint)
}

// randomCoordinate generates a random coordinate within the given range
func randomCoordinate(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}
