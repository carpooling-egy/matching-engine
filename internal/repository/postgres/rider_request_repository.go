package postgres

import (
	"context"
	"time"

	"gorm.io/gorm"

	"matching-engine/internal/errors"
	"matching-engine/internal/model"
	"matching-engine/internal/repository"
	"matching-engine/internal/repository/entity"
)

// PostgresRiderRequestRepo implements repository.PostgresRiderRequestRepo
type PostgresRiderRequestRepo struct {
	db *gorm.DB
}

// NewPostgresRiderRequestRepo creates a new rider request repository
func NewPostgresRiderRequestRepo(db *Database) repository.RiderRequestRepo {
	if db == nil {
		panic("db cannot be nil")
	}
	return &PostgresRiderRequestRepo{db: db.DB}
}

// GetByID fetches a rider request by ID
func (r *PostgresRiderRequestRepo) GetByID(ctx context.Context, id string) (*model.Request, error) {
	if id == "" {
		return nil, errors.EmptyID("rider request")
	}

	var riderRequestDB entity.RiderRequestDB

	err := r.db.WithContext(ctx).
		Omit("created_at, updated_at").
		First(&riderRequestDB, "id = ?", id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFound("rider request", id)
		}
		return nil, errors.DatabaseError("get_rider_request", err)
	}

	return riderRequestDB.ToRiderRequest(), nil
}

// GetUnmatched finds rider requests that can be matched within the specified time window
func (r *PostgresRiderRequestRepo) GetUnmatched(ctx context.Context, start, end time.Time) ([]*model.Request, error) {
	if end.Before(start) {
		return nil, errors.InvalidTimeRange()
	}

	var riderRequestDB []entity.RiderRequestDB

	err := r.db.WithContext(ctx).
		Omit("created_at, updated_at").
		Where("is_matched = false").
		Where("earliest_departure_time BETWEEN ? AND ?", start, end).
		Find(&riderRequestDB).Error

	if err != nil {
		return nil, errors.DatabaseError("find_matchable_requests", err)
	}

	return convertToRiderRequests(riderRequestDB), nil
}

// convertToRiderRequests converts database model to domain model
func convertToRiderRequests(riderRequestDB []entity.RiderRequestDB) []*model.Request {
	requests := make([]*model.Request, 0, len(riderRequestDB))

	for i := range riderRequestDB {
		if domainModel := riderRequestDB[i].ToRiderRequest(); domainModel != nil {
			requests = append(requests, domainModel)
		}
	}

	return requests
}
