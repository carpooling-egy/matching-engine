package tests

import (
	"github.com/stretchr/testify/assert"
	"matching-engine/internal/enums"
	"matching-engine/internal/model"
	"matching-engine/internal/service/pathgeneration/generator"
	"testing"
	"time"
)

// Helper function to create a test path point
func createTestPathPoint(pointType enums.PointType, arrivalTime time.Time) model.PathPoint {
	coordinate := must(model.NewCoordinate(0, 0))
	walkingDuration := time.Duration(0)
	pp := model.NewPathPoint(*coordinate, pointType, arrivalTime, nil, walkingDuration)
	return *pp
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func TestGeneratePaths_EmptyPath(t *testing.T) {
	// Arrange
	pathGenerator := generator.NewInsertionPathGenerator()

	var emptyPath []model.PathPoint
	pickup := createTestPathPoint(enums.Pickup, time.Now())
	dropoff := createTestPathPoint(enums.Dropoff, time.Now().Add(10*time.Minute))

	// Act
	_, err := pathGenerator.GeneratePaths(emptyPath, &pickup, &dropoff)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "path must contain at least two points")
}

func TestGeneratePaths_SinglePointPath(t *testing.T) {
	// Arrange
	pathGenerator := generator.NewInsertionPathGenerator()

	singlePointPath := []model.PathPoint{createTestPathPoint(enums.Pickup, time.Now())}
	pickup := createTestPathPoint(enums.Pickup, time.Now().Add(5*time.Minute))
	dropoff := createTestPathPoint(enums.Dropoff, time.Now().Add(10*time.Minute))

	// Act
	_, err := pathGenerator.GeneratePaths(singlePointPath, &pickup, &dropoff)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "path must contain at least two points")
}

func TestGeneratePaths_BasicPath(t *testing.T) {
	// Arrange
	pathGenerator := generator.NewInsertionPathGenerator()

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

func TestGeneratePaths_LongerPath(t *testing.T) {
	// Arrange
	pathGenerator := generator.NewInsertionPathGenerator()

	now := time.Now()

	// create points
	source := createTestPathPoint(enums.Source, now)
	pickup := createTestPathPoint(enums.Pickup, now.Add(5*time.Minute))
	dropoff := createTestPathPoint(enums.Dropoff, now.Add(10*time.Minute))
	p1 := createTestPathPoint(enums.Pickup, now.Add(11*time.Minute))
	p2 := createTestPathPoint(enums.Dropoff, now.Add(12*time.Minute))
	destination := createTestPathPoint(enums.Destination, now.Add(15*time.Minute))
	// Create a basic path with two points

	path := []model.PathPoint{
		source,
		p1,
		p2,
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
	assert.Len(t, generatedPaths[0], 6)

	//// Check the order of points
	assert.Equal(t, source.ID(), generatedPaths[0][0].ID())
	assert.Equal(t, pickup.ID(), generatedPaths[0][1].ID())
	assert.Equal(t, dropoff.ID(), generatedPaths[0][2].ID())
	assert.Equal(t, p1.ID(), generatedPaths[0][3].ID())
	assert.Equal(t, p2.ID(), generatedPaths[0][4].ID())
	assert.Equal(t, destination.ID(), generatedPaths[0][5].ID())
}

func TestGeneratePaths_LongerPath2(t *testing.T) {
	// Arrange
	pathGenerator := generator.NewInsertionPathGenerator()

	now := time.Now()

	// create points
	source := createTestPathPoint(enums.Source, now)
	p1 := createTestPathPoint(enums.Pickup, now.Add(11*time.Minute))
	p2 := createTestPathPoint(enums.Dropoff, now.Add(12*time.Minute))
	pickup := createTestPathPoint(enums.Pickup, now.Add(13*time.Minute))
	dropoff := createTestPathPoint(enums.Dropoff, now.Add(14*time.Minute))
	destination := createTestPathPoint(enums.Destination, now.Add(15*time.Minute))
	// Create a basic path with two points

	path := []model.PathPoint{
		source,
		p1,
		p2,
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
	assert.Len(t, generatedPaths, 6)
	assert.Len(t, generatedPaths[0], 6)

	//// Check the order of points
	assert.Equal(t, pickup.ID(), generatedPaths[0][1].ID())
	assert.Equal(t, dropoff.ID(), generatedPaths[0][2].ID())
	assert.Equal(t, pickup.ID(), generatedPaths[1][1].ID())
	assert.Equal(t, dropoff.ID(), generatedPaths[1][3].ID())
	assert.Equal(t, pickup.ID(), generatedPaths[2][1].ID())
	assert.Equal(t, dropoff.ID(), generatedPaths[2][4].ID())
	assert.Equal(t, pickup.ID(), generatedPaths[3][2].ID())
	assert.Equal(t, dropoff.ID(), generatedPaths[3][3].ID())
	assert.Equal(t, pickup.ID(), generatedPaths[4][2].ID())
	assert.Equal(t, dropoff.ID(), generatedPaths[4][4].ID())
	assert.Equal(t, pickup.ID(), generatedPaths[5][3].ID())
	assert.Equal(t, dropoff.ID(), generatedPaths[5][4].ID())
}

func TestGeneratePaths_LongerPath3(t *testing.T) {
	// Arrange
	pathGenerator := generator.NewInsertionPathGenerator()

	now := time.Now()

	// create points
	source := createTestPathPoint(enums.Source, now)
	pickup := createTestPathPoint(enums.Pickup, now.Add(5*time.Minute))
	dropoff := createTestPathPoint(enums.Dropoff, now.Add(10*time.Minute))
	p1 := createTestPathPoint(enums.Pickup, now.Add(11*time.Minute))
	p2 := createTestPathPoint(enums.Dropoff, now.Add(12*time.Minute))
	p3 := createTestPathPoint(enums.Pickup, now.Add(13*time.Minute))
	p4 := createTestPathPoint(enums.Dropoff, now.Add(14*time.Minute))
	destination := createTestPathPoint(enums.Destination, now.Add(15*time.Minute))
	// Create a basic path with two points

	path := []model.PathPoint{
		source,
		p1,
		p2,
		p3,
		p4,
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
	assert.Len(t, generatedPaths[0], 8)

	//// Check the order of points
	assert.Equal(t, source.ID(), generatedPaths[0][0].ID())
	assert.Equal(t, pickup.ID(), generatedPaths[0][1].ID())
	assert.Equal(t, dropoff.ID(), generatedPaths[0][2].ID())
	assert.Equal(t, p1.ID(), generatedPaths[0][3].ID())
	assert.Equal(t, p2.ID(), generatedPaths[0][4].ID())
	assert.Equal(t, p3.ID(), generatedPaths[0][5].ID())
	assert.Equal(t, p4.ID(), generatedPaths[0][6].ID())
	assert.Equal(t, destination.ID(), generatedPaths[0][7].ID())
}
