package main

import (
	"bufio"
	"context"
	"fmt"
	"matching-engine/internal/adapter/valhalla"
	"matching-engine/internal/model"
	"os"

	"github.com/rs/zerolog/log"
)

func main() {
	// Create a new Valhalla client
	engine, err := valhalla.NewValhalla()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create Valhalla engine client")
	}

	// Open points.txt file
	pointsFile, err := os.Open("cmd/correcteness_test/snap_point/points.txt")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open points.txt")
	}
	defer func() {
		if err := pointsFile.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close points.txt")
		}
	}()

	scanner := bufio.NewScanner(pointsFile)

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
	var snappedPoints []model.Coordinate
	for _, coord := range coords {
		snappedPoint, err := engine.SnapPointToRoad(context.Background(), &coord)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to snap point: %v", coord)
			continue
		}
		snappedPoints = append(snappedPoints, *snappedPoint)
	}

	// Write snapped points to file, one per line
	outputFile, err := os.Create("cmd/correcteness_test/snap_point/snapped_points.txt")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create snapped_points.txt")
	}
	defer func() {
		if err := outputFile.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close snapped_points.txt")
		}
	}()

	writer := bufio.NewWriter(outputFile)
	for _, pt := range snappedPoints {
		_, err := fmt.Fprintf(writer, "%.8f, %.8f\n", pt.Lat(), pt.Lng())
		if err != nil {
			log.Error().Err(err).Msg("Failed to write to snapped_points.txt")
		}
	}
	if err := writer.Flush(); err != nil {
		log.Error().Err(err).Msg("Failed to flush to snapped_points.txt")
	}
}
