package main

import (
	"bufio"
	"fmt"
	"matching-engine/cmd/correcteness_test"
	"os"
	"time"

	"matching-engine/internal/adapter/valhalla"
	"matching-engine/internal/model"

	"github.com/rs/zerolog/log"
)

func main() {
	// Create a new Valhalla client
	engine, err := valhalla.NewValhalla()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create Valhalla engine client")
	}

	// Open points.txt file
	pointsFile, err := os.Open("cmd/correcteness_test/cumulative_travel_time_generator/points.txt")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open points.txt")
	}
	defer func() {
		if err := pointsFile.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close points.txt")
		}
	}()

	scanner := bufio.NewScanner(pointsFile)

	// Read the first line: timestamp
	if !scanner.Scan() {
		log.Fatal().Msg("points.txt is empty or missing timestamp")
	}
	var ts int64
	_, err = fmt.Sscanf(scanner.Text(), "%d", &ts)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse timestamp from first line")
	}
	timestamp := time.Unix(ts, 0)

	// Read coordinates
	var coords []model.Coordinate
	for scanner.Scan() {
		var lat, lon float64
		_, err := fmt.Sscanf(scanner.Text(), "%f, %f", &lat, &lon)
		if err != nil {
			log.Error().Err(err).Msg("Failed to parse coordinates")
			continue
		}
		coord, err := model.NewCoordinate(lat, lon)
		if err != nil {
			log.Error().Err(err).Msg("Failed to create coordinate")
			continue
		}
		coords = append(coords, *coord)
	}
	if err := scanner.Err(); err != nil {
		log.Error().Err(err).Msg("Error reading points.txt")
		return
	}
	drivingTimes := correcteness_test.GetCumulativeTimes(coords, timestamp, engine)
	log.Info().Msgf("Computed driving times: %v", drivingTimes)

	// Write output to file
	outputFile, err := os.Create("cmd/correcteness_test/cumulative_travel_time_generator/computed_times.txt")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create computed_times.txt")
	}
	defer func() {
		if err := outputFile.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close computed_times.txt")
		}
	}()

	_, err = fmt.Fprintf(outputFile, "%v\n", drivingTimes)
	if err != nil {
		log.Error().Err(err).Msg("Failed to write to computed_times.txt")
	}
}
