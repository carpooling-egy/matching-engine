package main

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"matching-engine/internal/adapter/messaging"
	"matching-engine/internal/adapter/messaging/natsjetstream"
	"matching-engine/internal/enums"
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
	var err2 error

	// Use the existing NATS publisher implementation
	publisher, err2 = natsjetstream.NewNATSPublisher()
	if err2 != nil {
		fmt.Printf("Error connecting to NATS: %v\n", err2)
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

	// Generate and publish matching results
	results := generateMatchingResults(40)
	fmt.Printf("Publishing %d matching results...\n", len(results))

	startTime := time.Now()
	err2 = publisher.PublishMatchingResults(results)

	if err2 != nil {
		log.Warn().Err(err2).Msg("Failed to publish matching results")
	}
	duration := time.Since(startTime)

	fmt.Printf("Successfully published %d matching results in %v\n", len(results), duration)

	time.Sleep(30 * time.Second) // Wait to ensure messages are sent

	// Generate and publish more matching results
	results = generateMatchingResults(40)
	fmt.Printf("Publishing %d matching results...\n", len(results))

	startTime = time.Now()
	err2 = publisher.PublishMatchingResults(results)

	if err2 != nil {
		log.Warn().Err(err2).Msg("Failed to publish matching results")
	}
	duration = time.Since(startTime)

	fmt.Printf("Returned from publish fn\n")

	fmt.Println("Press Ctrl+C to exit")

	// Keep the program running until interrupted
	select {}
}

// generateMatchingResults creates a specified number of sample matching results
func generateMatchingResults(count int) []*model.MatchingResult {
	results := make([]*model.MatchingResult, count)

	for i := 0; i < count; i++ {
		// Create a matching result with unique IDs
		userID := fmt.Sprintf("user-%d", i+1)
		offerID := fmt.Sprintf("offer-%d", i+1)

		// Create an empty path slice - will add points later
		newPath := make([]*model.Point, 0)

		// Create an empty assigned requests slice - will add requests later
		assignedRequests := make([]*model.MatchedRequest, 0)

		// Create the matching result using the model constructor
		result := model.NewMatchingResult(userID, offerID, assignedRequests, newPath)

		// Create an offer for this matching result
		offerCoord, _ := model.NewCoordinate(randomCoordinate(-90, 90), randomCoordinate(-180, 180))
		offerDetourTime := time.Duration(rand.Intn(60)) * time.Minute
		offerDepartureTime := time.Now().Add(time.Duration(rand.Intn(120)) * time.Minute)
		offerPreference := model.Preference{}

		offer := model.NewOffer(
			offerID,
			userID,
			*offerCoord,
			*offerCoord,
			offerDetourTime,
			offerDepartureTime,
			nil,
			offerPreference,
			nil,
		)

		// Add some random assigned requests
		numRequests := 1 + rand.Intn(3) // 1-3 requests per matching result
		for j := 0; j < numRequests; j++ {
			requestID := fmt.Sprintf("request-%d-%d", i+1, j+1)

			// Create a request to use as owner for the points
			request := createRequest(requestID)

			// Create pickup and dropoff points for this request
			pickup, dropoff := generatePoints(request)

			// Create matched request with the offer and request
			matchedRequest := model.NewMatchedRequest(offer, request, *pickup, *dropoff)

			// Add to the assigned matched requests
			currentRequests := result.AssignedMatchedRequests()
			currentRequests = append(currentRequests, matchedRequest)
			result.SetAssignedMatchedRequests(currentRequests)

			// Add points to the path
			currentPath := result.NewPath()
			currentPath = append(currentPath, pickup, dropoff)
			result.SetNewPath(currentPath)
		}

		results[i] = result
	}

	return results
}

// createRequest creates a sample request for testing
func createRequest(requestID string) *model.Request {
	// Generate random valid coordinates
	sourceLat := randomCoordinate(-90, 90)
	sourceLng := randomCoordinate(-180, 180)
	sourceCoord, _ := model.NewCoordinate(sourceLat, sourceLng)

	destLat := randomCoordinate(-90, 90)
	destLng := randomCoordinate(-180, 180)
	destCoord, _ := model.NewCoordinate(destLat, destLng)

	// Create time constraints
	now := time.Now()
	earliestDeparture := now.Add(time.Duration(rand.Intn(60)) * time.Minute)
	latestArrival := earliestDeparture.Add(time.Duration(rand.Intn(120)+30) * time.Minute)
	maxWalkingTime := time.Duration(rand.Intn(15)+5) * time.Minute

	// Create the request
	userID := fmt.Sprintf("user-%s", requestID)
	preference := model.Preference{} // Using your Preference type
	numberOfRiders := 1 + rand.Intn(3)

	return model.NewRequest(
		requestID,
		userID,
		*sourceCoord,
		*destCoord,
		earliestDeparture,
		latestArrival,
		maxWalkingTime,
		preference,
		numberOfRiders,
	)
}

// generatePoints creates random pickup and dropoff points for a request
func generatePoints(owner model.Role) (*model.Point, *model.Point) {
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

	pickupPoint := model.NewPoint(owner, *pickupCoord, pickupTime, enums.Pickup)
	dropoffPoint := model.NewPoint(owner, *dropoffCoord, dropoffTime, enums.Dropoff)

	return pickupPoint, dropoffPoint
}

// randomCoordinate generates a random coordinate within the given range
func randomCoordinate(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}
