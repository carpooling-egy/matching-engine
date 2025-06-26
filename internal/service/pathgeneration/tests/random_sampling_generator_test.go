package tests

import (
	"github.com/stretchr/testify/assert"
	"matching-engine/internal/enums"
	"matching-engine/internal/model"
	"matching-engine/internal/service/pathgeneration/generator"
	"testing"
	"time"
)

func TestRandomGeneratePaths_EmptyPath(t *testing.T) {
	// Arrange
	pathGenerator := generator.NewRandomTopologicalGenerator()

	var emptyPath []model.PathPoint
	pickup := createTestPathPoint(enums.Pickup, time.Now())
	dropoff := createTestPathPoint(enums.Dropoff, time.Now().Add(10*time.Minute))

	// Act
	_, err := pathGenerator.GeneratePaths(emptyPath, &pickup, &dropoff)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "path must contain at least two points")
}

func TestRandomGeneratePaths_SinglePointPath(t *testing.T) {
	// Arrange
	pathGenerator := generator.NewRandomTopologicalGenerator()

	singlePointPath := []model.PathPoint{createTestPathPoint(enums.Pickup, time.Now())}
	pickup := createTestPathPoint(enums.Pickup, time.Now().Add(5*time.Minute))
	dropoff := createTestPathPoint(enums.Dropoff, time.Now().Add(10*time.Minute))

	// Act
	_, err := pathGenerator.GeneratePaths(singlePointPath, &pickup, &dropoff)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "path must contain at least two points")
}

func TestRandomGeneratePaths_BasicPath(t *testing.T) {
	// Arrange
	pathGenerator := generator.NewRandomTopologicalGenerator()

	now := time.Now()

	// create points
	source := createTestPathPoint(enums.Source, now)
	pickup := createTestPathPoint(enums.Pickup, now.Add(5*time.Minute))
	dropoff := createTestPathPoint(enums.Dropoff, now.Add(10*time.Minute))
	destination := createTestPathPoint(enums.Destination, now.Add(15*time.Minute))

	// Create a basic path with two points

	path := []model.PathPoint{
		source,
		destination,
	}

	// Act
	seq, err := pathGenerator.GeneratePaths(path, &pickup, &dropoff)
	assert.NoError(t, err)

	// Collect all generated paths
	var generatedPaths [][]model.PathPoint
	for candidatePath, pathErr := range seq {
		assert.NoError(t, pathErr)
		generatedPaths = append(generatedPaths, candidatePath)
	}

	// Assert
	// Should have exactly one path: [original1, pickup, dropoff, original2]
	assert.Len(t, generatedPaths, 1)
	assert.Len(t, generatedPaths[0], 4)

	// Check the order of points
	assert.Equal(t, source.ID(), generatedPaths[0][0].ID())
	assert.Equal(t, pickup.ID(), generatedPaths[0][1].ID())
	assert.Equal(t, dropoff.ID(), generatedPaths[0][2].ID())
	assert.Equal(t, destination.ID(), generatedPaths[0][3].ID())
}
