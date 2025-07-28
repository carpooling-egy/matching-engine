package main

import (
	"bufio"
	"fmt"
	"github.com/rs/zerolog/log"
	"matching-engine/cmd/correcteness_test"
	"matching-engine/internal/adapter/valhalla"
	"matching-engine/internal/enums"
	"matching-engine/internal/model"
	"os"
	"strconv"
	"strings"
	"time"
)

func readNextCoordinate(scanner *bufio.Scanner) (*model.Coordinate, error) {
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		var lat, lon float64
		if _, err := fmt.Sscanf(line, "%f, %f", &lat, &lon); err == nil {
			return model.NewCoordinate(lat, lon)
		}
		return nil, fmt.Errorf("invalid coordinate line: %s", line)
	}
	return nil, fmt.Errorf("unexpected end of input while reading coordinate")
}

func readNextDuration(scanner *bufio.Scanner) (time.Duration, error) {
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		mins, err := strconv.Atoi(line)
		if err != nil {
			return 0, fmt.Errorf("invalid duration line: %s", line)
		}
		return time.Duration(mins) * time.Minute, nil
	}
	return 0, fmt.Errorf("unexpected end of input while reading duration")
}

func buildOffer(offerPath []model.PathPoint) *model.Offer {
	preference := model.NewPreference(enums.Female, true)
	return model.NewOffer(
		"offer-id",
		"user-id",
		*offerPath[0].Coordinate(), // source
		*offerPath[len(offerPath)-1].Coordinate(), // destination
		time.Now().Add(1*time.Hour),               // departure
		10*time.Minute,                            // waiting time
		2,                                         // capacity
		*preference,
		time.Now().Add(2*time.Hour), // arrival
		0,
		offerPath,
		nil,
	)
}

func main() {
	// Open input
	inputFile, err := os.Open("cmd/correcteness_test/pickup_dropoff_generator/input.txt")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open input.txt")
	}
	defer inputFile.Close()

	scanner := bufio.NewScanner(inputFile)

	// Parse offer path points
	var offerPath []model.PathPoint
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			break
		}
		if strings.HasPrefix(line, "#") {
			continue
		}
		var lat, lon float64
		if _, err := fmt.Sscanf(line, "%f, %f", &lat, &lon); err != nil {
			log.Warn().Msgf("Skipping invalid coordinate: %s", line)
			continue
		}
		coord, err := model.NewCoordinate(lat, lon)
		if err != nil {
			log.Warn().Err(err).Msg("Invalid coordinate creation")
			continue
		}
		pp := model.NewPathPoint(*coord, enums.Source, time.Now().Add(10*time.Minute), nil, 0)
		offerPath = append(offerPath, *pp)
	}

	// Read pickup/dropoff coordinates and walking duration
	source, err := readNextCoordinate(scanner)
	if err != nil {
		log.Fatal().Err(err).Msg("Error reading source coordinate")
	}
	destination, err := readNextCoordinate(scanner)
	if err != nil {
		log.Fatal().Err(err).Msg("Error reading destination coordinate")
	}
	walkingDuration, err := readNextDuration(scanner)
	if err != nil {
		log.Fatal().Err(err).Msg("Error reading walking duration")
	}

	// Create offer
	offer := buildOffer(offerPath)

	// Setup processor
	engine, err := valhalla.NewValhalla()
	if err != nil {
		log.Fatal().Err(err).Msg("Error creating Valhalla engine")
	}
	pickup, pickupDuration, dropoff, dropoffDuration := correcteness_test.GetPickupDropoffPointsAndDurations(engine, offer, source, walkingDuration, destination)

	// Write output
	outputFile, err := os.Create("cmd/correcteness_test/pickup_dropoff_generator/output.txt")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create output.txt")
	}
	defer func() {
		if err := outputFile.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close output file")
		}
	}()

	formatCoord := func(c model.Coordinate) string {
		return fmt.Sprintf("(%.6f, %.6f)", c.Lat(), c.Lng())
	}

	_, err = fmt.Fprintf(outputFile,
		"Pickup Point: %s\nPickup Duration: %d minutes\n\nDropoff Point: %s\nDropoff Duration: %d minutes\n",
		formatCoord(*pickup), int(pickupDuration.Minutes()),
		formatCoord(*dropoff), int(dropoffDuration.Minutes()),
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to write to output file")
	}
}
