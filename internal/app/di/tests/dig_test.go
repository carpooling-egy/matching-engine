package tests

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/mock"
	"go.uber.org/dig"
	"matching-engine/internal/adapter/routing"
	"matching-engine/internal/app/di"
	"matching-engine/internal/app/di/utils"
	"matching-engine/internal/app/starter"
	"matching-engine/internal/model"
	"matching-engine/internal/publisher"
	"matching-engine/internal/repository/postgres"
	"testing"
	"time"
)

// This method is intended to test the DI using the real database and adapters.
// It should be run with a test database and the necessary environment variables set.
func TestDI(t *testing.T) {

	// Path to .env file containing database connection, channel, routing engine configurations
	envPath := ""

	// This is just temp, ideally we should use a test database
	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatal().Msg(".env file not found or failed to load")
	}

	c := di.BuildContainer()
	// Validate container
	if err = validateContainer(c); err != nil {
		log.Fatal().Err(err).Msg("Container validation failed")
	}

}

// validateContainer checks if the container is configured correctly
func validateContainer(c *dig.Container) error {
	return c.Invoke(func(starterService *starter.StarterService) {
		log.Info().Msg("Container validation successful")
	})
}

func TestDIWithMocks(t *testing.T) {
	c := dig.New()

	// Register mock database and adapters
	utils.Must(c.Provide(NewMockDatabase))
	utils.Must(c.Provide(NewMockChannel))
	utils.Must(c.Provide(NewMockRoutingEngine))
	// Register other real services as needed, or more mocks

	// Register the rest of the modules, skipping the real DB/adapters
	di.RegisterGeoServices(c)
	di.RegisterPickupDropoffServices(c)
	di.RegisterTimeMatrixServices(c)
	di.RegisterPathServices(c)
	di.RegisterCheckers(c)
	di.RegisterMatchingServices(c)
	di.RegisterDatabaseRepositoriesAndServices(c)
	// Do not call registerDatabaseServices or registerAdapters if you want only mocks

	di.RegisterStarterService(c)

	// Validate container
	if err := validateContainer(c); err != nil {
		log.Fatal().Err(err).Msg("Container validation failed")
	}
}

type MockPublisher struct {
	mock.Mock
}

type MockRoutingEngine struct {
	mock.Mock
}

func (m *MockRoutingEngine) PlanDrivingRoute(ctx context.Context, routeParams *model.RouteParams) (*model.Route, error) {
	panic("shouldn't be called")
}

func (m *MockRoutingEngine) ComputeDrivingTime(ctx context.Context, routeParams *model.RouteParams) ([]time.Duration, error) {
	panic("shouldn't be called")
}

func (m *MockRoutingEngine) ComputeWalkingTime(ctx context.Context, walkParams *model.WalkParams) (time.Duration, error) {
	panic("shouldn't be called")
}

func (m *MockRoutingEngine) ComputeIsochrone(ctx context.Context, req *model.IsochroneParams) (*model.Isochrone, error) {
	panic("shouldn't be called")
}

func (m *MockRoutingEngine) ComputeDistanceTimeMatrix(ctx context.Context, req *model.DistanceTimeMatrixParams) (*model.DistanceTimeMatrix, error) {
	panic("shouldn't be called")
}

func (m *MockRoutingEngine) SnapPointToRoad(ctx context.Context, point *model.Coordinate) (*model.Coordinate, error) {
	panic("shouldn't be called")
}

func (m *MockPublisher) Publish(results []*model.MatchingResult) error {
	panic("shouldn't be called")
}

func (m *MockPublisher) Close() error {
	panic("shouldn't be called")
}

// Example mock constructors
func NewMockDatabase() *postgres.Database {
	// Return a mock or fake database instance
	return &postgres.Database{} // Replace with actual mock
}

func NewMockChannel() publisher.Publisher {
	// Return a mock adapter
	return &MockPublisher{} // Replace with actual mock
}

func NewMockRoutingEngine() routing.Engine {
	// Return a mock adapter
	return &MockRoutingEngine{} // Replace with actual mock
}
