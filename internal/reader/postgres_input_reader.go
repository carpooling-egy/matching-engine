package reader

import (
	"context"
	"fmt"
	"matching-engine/internal/model"
	"matching-engine/internal/repository"
	"matching-engine/internal/repository/postgres"
	"time"
)

type PostgresInputReader struct {
	requestsRepository repository.RiderRequestRepo
	offersRepository   repository.DriverOfferRepo
	db                 *postgres.Database
}

func NewPostgresInputReader(db *postgres.Database, requestsRepo repository.RiderRequestRepo, offersRepo repository.DriverOfferRepo) MatchInputReader {

	// Connect to the database
	return &PostgresInputReader{
		db:                 db,
		requestsRepository: requestsRepo,
		offersRepository:   offersRepo,
	}
}

func (r *PostgresInputReader) GetOffersAndRequests(ctx context.Context) ([]*model.Request, []*model.Offer, bool, error) {

	// TODO - check if we need to add a timeout to the context
	// TODO - Read start and end time from config
	requests, err := r.requestsRepository.GetUnmatched(ctx, time.Now(), time.Now().Add(24*time.Hour))

	if err != nil {
		return nil, nil, false, fmt.Errorf("failed to get requests %w", err)
	}
	if len(requests) == 0 {
		return nil, nil, false, nil
	}

	offers, err := r.offersRepository.GetAvailable(ctx, time.Now(), time.Now().Add(24*time.Hour))
	if err != nil {
		return nil, nil, false, fmt.Errorf("failed to get offers %w", err)
	}
	if len(offers) == 0 {
		return nil, nil, false, nil
	}
	return requests, offers, true, nil
}

func (r *PostgresInputReader) Close() error {
	err := r.db.Close()
	if err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}
	return nil
}
