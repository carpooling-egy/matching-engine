package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"matching-engine/internal/model"
	"matching-engine/internal/repository"
	"matching-engine/internal/repository/postgres"
)

func main() {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Connect to the database
	db, err := postgres.NewDatabase(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Successfully connected to database")

	// Create repositories
	riderRepo := postgres.NewPostgresRiderRequestRepo(db)
	driverRepo := postgres.NewPostgresDriverOfferRepository(db)

	fmt.Println("=== TESTING RIDER REQUEST REPOSITORY ===")
	testRiderRequestRepo(ctx, riderRepo)

	fmt.Println("\n=== TESTING DRIVER OFFER REPOSITORY ===")
	testDriverOfferRepo(ctx, driverRepo)
}

func testRiderRequestRepo(ctx context.Context, repo repository.RiderRequestRepo) {
	// Test GetByID
	fmt.Println("Testing GetByID with 'req-001'...")
	rider, err := repo.GetByID(ctx, "req-001")
	if err != nil {
		fmt.Printf("Error getting rider request: %v\n", err)
	} else {
		fmt.Printf("Found rider request: ID=%s, UserID=%s\n", rider.ID(), rider.UserID())
		fmt.Printf("  Source: (%.6f, %.6f)\n", rider.Source().Lat(), rider.Source().Lng())
		fmt.Printf("  Destination: (%.6f, %.6f)\n", rider.Destination().Lat(), rider.Destination().Lng())
		fmt.Printf("  Earliest Departure: %v\n", rider.EarliestDepartureTime())
		fmt.Printf("  walking: %v\n", rider.MaxWalkingDurationMinutes())
		fmt.Print("  Preferences: ")
		fmt.Printf("gender %s\n", rider.Preferences().Gender().String())
		fmt.Printf("  %v\n", rider.Preferences().SameGender())
	

		fmt.Printf("  Latest Arrival: %v\n", rider.LatestArrivalTime())
		fmt.Printf("  Number of Riders: %d\n", rider.NumberOfRiders())
	}

	// Test FindUnmatched
	fmt.Println("\nTesting FindUnmatched...")
	start := time.Now().Add(-24 * time.Hour) // 1 day ago
	end := time.Now().Add(24 * time.Hour)    // 1 day from now

	unmatched, err := repo.FindUnmatched(ctx, start, end)
	if err != nil {
		fmt.Printf("Error finding unmatched requests: %v\n", err)
	} else {
		fmt.Printf("Found %d unmatched rider requests\n", len(unmatched))
		for i, req := range unmatched {
			fmt.Printf("  %d. ID=%s, UserID=%s\n", i+1, req.ID(), req.UserID())
		}
	}
}

func testDriverOfferRepo(ctx context.Context, repo repository.DriverOfferRepo) {
	// Test GetByID
	fmt.Println("Testing GetByID with 'drv-001'...")
	driver, err := repo.GetByID(ctx, "drv-001")
	if err != nil {
		fmt.Printf("Error getting driver offer: %v\n", err)
		return
	}

	printOfferDetails(driver)

	// Test GetAvailable
	fmt.Println("\nTesting GetAvailable...")
	start := time.Now().Add(-24 * time.Hour) // 1 day ago
	end := time.Now().Add(24 * time.Hour)    // 1 day from now

	available, err := repo.GetAvailable(ctx, start, end)
	if err != nil {
		fmt.Printf("Error finding available driver offers: %v\n", err)
	} else {
		fmt.Printf("Found %d available driver offers\n", len(available))
		for i, offer := range available {
			fmt.Printf("\n--- Available Offer %d ---\n", i+1)
			printOfferDetails(offer)
		}
	}
}

func printOfferDetails(driver *model.Offer) {
	fmt.Println("=== Driver Offer Details ===")
	fmt.Printf("ID: %s\n", driver.ID())
	fmt.Printf("User ID: %s\n", driver.UserID())

	// Basic details
	if driver.Source() != nil {
		fmt.Printf("Source: (%.6f, %.6f)\n", driver.Source().Lat(), driver.Source().Lng())
	}

	if driver.Destination() != nil {
		fmt.Printf("Destination: (%.6f, %.6f)\n", driver.Destination().Lat(), driver.Destination().Lng())
	}

	fmt.Printf("Departure Time: %v\n", driver.DepartureTime())
	fmt.Printf("Capacity: %d\n", driver.Capacity())
	fmt.Printf("Current Number of Requests: %d\n", driver.CurrentNumberOfRequests())
	fmt.Printf("Detour Duration: %v\n", driver.DetourDurationMinutes())
	fmt.Printf("Max Estimated Arrival Time: %v\n", driver.MaxEstimatedArrivalTime())

	// Preferences
	pref := driver.Preferences()
	fmt.Println("\n--- Preferences ---")
	fmt.Printf("Same Gender: %t\n", pref.SameGender())
	fmt.Printf("Gender: %v\n", pref.Gender())

	// Path points
	fmt.Printf("\n--- Path Points (%d) ---\n", len(driver.Path()))
	for i, point := range driver.Path() {
		fmt.Printf("%d. Point Type: %v\n", i+1, point.PointType())

		if point.Coordinate() != nil {
			fmt.Printf("   Position: (%.6f, %.6f)\n", point.Coordinate().Lat(), point.Coordinate().Lng())
		} else {
			fmt.Printf("   Position: <nil>\n")
		}

		fmt.Printf("   Expected Arrival: %v\n", point.ExpectedArrivalTime())
		fmt.Printf("	  Walking Duration: %v\n", point.WalkingDuration())

		// Owner information
		if point.Owner() != nil {
			switch {
			case point.Owner() == driver:
				fmt.Printf("   Owner: This Driver Offer\n")
			default:
				// Try to determine the concrete type
				if req, ok := point.Owner().AsRequest(); ok {
					fmt.Printf("   Owner: Rider Request (ID: %s)\n", req.ID())
				} else if offer, ok := point.Owner().AsOffer(); ok {
					fmt.Printf("   Owner: Driver Offer (ID: %s)\n", offer.ID())
				} else {
					fmt.Printf("   Owner: Unknown Type %T\n", point.Owner())
				}
			}
		} else {
			fmt.Printf("   Owner: <nil>\n")
		}
	}

	// Matched requests
	fmt.Printf("\n--- Matched Requests (%d) ---\n", len(driver.MatchedRequests()))
	for i, matched := range driver.MatchedRequests() {
		fmt.Printf("%d. ", i+1)

		if matched == nil {
			fmt.Println("Matched Request: <nil>")
			continue
		}

		// Request details

		fmt.Printf("Request ID: %s\n", matched.ID())
		fmt.Printf("   User ID: %s\n", matched.UserID())

		fmt.Printf("   # of Riders: %d\n", matched.NumberOfRiders())
		fmt.Printf("   Earliest Departure: %v\n", matched.EarliestDepartureTime())
		fmt.Printf("   Latest Arrival: %v\n", matched.LatestArrivalTime())
		fmt.Printf("   Max Walking Duration: %v\n", matched.MaxWalkingDurationMinutes())

	}
}

func printPathPointDetails(point model.PathPoint) {
	if point.Coordinate() != nil {
		fmt.Printf("      Position: (%.6f, %.6f)\n", point.Coordinate().Lat(), point.Coordinate().Lng())
	} else {
		fmt.Printf("      Position: <nil>\n")
	}

	fmt.Printf("      Point Type: %v\n", point.PointType())
	fmt.Printf("      Expected Arrival: %v\n", point.ExpectedArrivalTime())

	if point.Owner() != nil {
		if req, ok := point.Owner().AsRequest(); ok {
			fmt.Printf("      Owner: Rider Request (ID: %s)\n", req.ID())
		} else if offer, ok := point.Owner().AsOffer(); ok {
			fmt.Printf("      Owner: Driver Offer (ID: %s)\n", offer.ID())
		} else {
			fmt.Printf("      Owner: Unknown Type %T\n", point.Owner())
		}
	} else {
		fmt.Printf("      Owner: <nil>\n")
	}
}
