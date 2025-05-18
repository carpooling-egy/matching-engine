package tests

import (
	"matching-engine/internal/enums"
	"matching-engine/internal/model"
	"math/rand"
	"time"
)

func createTestOffer() *model.Offer {
	source, _ := model.NewCoordinate(1.0, 1.0)
	destination, _ := model.NewCoordinate(2.0, 2.0)
	departureTime := time.Now().Add(1 * time.Hour)

	return model.NewOffer(
		"offer-123",
		"user-abc",
		*source,
		*destination,
		departureTime,
		15*time.Minute,
		4,
		model.Preference{},
		departureTime.Add(1*time.Hour),
		0,
		nil,
		nil,
	)
}

func createTestRequest() *model.Request {
	source, _ := model.NewCoordinate(1.5, 1.5)
	destination, _ := model.NewCoordinate(2.5, 2.5)

	return model.NewRequest(
		"request-456",
		"user-xyz",
		*source,
		*destination,
		time.Now(),
		time.Now().Add(2*time.Hour),
		10*time.Minute,
		1,
		model.Preference{},
	)
}

func createPathPoint(coordinate *model.Coordinate, pointType enums.PointType, owner model.Role) *model.PathPoint {
	return model.NewPathPoint(
		*coordinate,
		pointType,
		time.Now(),
		owner,
		5*time.Minute,
	)
}

func generateRandomTimeDistanceMatrices(size int) ([][]time.Duration, [][]model.Distance) {
	timeMatrix := make([][]time.Duration, size)
	distanceMatrix := make([][]model.Distance, size)
	for i := 0; i < size; i++ {
		timeMatrix[i] = make([]time.Duration, size)
		distanceMatrix[i] = make([]model.Distance, size)
		for j := 0; j < size; j++ {
			if i == j {
				timeMatrix[i][j] = 0
				distanceMatrix[i][j] = *must(model.NewDistance(0, model.DistanceUnitKilometer))
			} else {
				// Random duration between 5 and 30 minutes
				minutes := 5 + rand.Intn(26)
				timeMatrix[i][j] = time.Duration(minutes) * time.Minute
				// Random distance between 0.5 and 2.0 kilometers
				dist := 0.5 + rand.Float32()*1.5
				distanceMatrix[i][j] = *must(model.NewDistance(dist, model.DistanceUnitKilometer))
			}
		}
	}
	return timeMatrix, distanceMatrix
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func resetPathPointIDCounter() {
	model.NextPointID = 1
}
