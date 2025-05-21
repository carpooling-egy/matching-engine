package postgres

import (
	"context"
	"time"

	"gorm.io/gorm"

	"matching-engine/internal/constants"
	"matching-engine/internal/errors"
	"matching-engine/internal/model"
	"matching-engine/internal/repository"
	"matching-engine/internal/repository/entity"
)

// PostgresDriverOfferRepo implements repository.PostgresDriverOfferRepo
type PostgresDriverOfferRepo struct {
	db *gorm.DB
}

// NewDriverOfferRepository creates a new driver offer repository
func NewPostgresDriverOfferRepository(db *Database) repository.DriverOfferRepo {
	if db == nil {
		panic("db cannot be nil")
	}
	return &PostgresDriverOfferRepo{db: db.DB}
}

// GetByID fetches a driver offer by ID
func (r *PostgresDriverOfferRepo) GetByID(ctx context.Context, id string) (*model.Offer, error) {
	if id == "" {
		return nil, errors.EmptyID("driver offer")
	}
	var driverOfferDB entity.DriverOfferDB
	err := r.db.WithContext(ctx).
		Preload("PathPoints", orderPathPointsByPathOrder).
		Preload("PathPoints.RiderRequest").
		First(&driverOfferDB, "id = ?", id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFound("driver offer", id)
		}
		return nil, errors.DatabaseError("get_driver_offer", err)
	}

	return driverOfferDB.ToDriverOffer(), nil
}

// GetAvailable fetches available driver offers with their paths and associated rider requests
func (r *PostgresDriverOfferRepo) GetAvailable(ctx context.Context, start, end time.Time) ([]*model.Offer, error) {
	if end.Before(start) {
		return nil, errors.InvalidTimeRange()
	}

	var driverOfferDB []entity.DriverOfferDB

	err := r.db.WithContext(ctx).
		Preload("PathPoints", orderPathPointsByPathOrder).
		Preload("PathPoints.RiderRequest").
		Where("departure_time BETWEEN ? AND ?", start, end).
		Where("current_number_of_requests < ?", constants.MaxDriverCapacity).
		Find(&driverOfferDB).Error

	if err != nil {
		return nil, errors.DatabaseError("fetch_available_driver_offers", err)
	}

	return convertToDriverOffers(driverOfferDB), nil
}

// orderPathPointsByPathOrder returns a function that orders path points by their path_order
func orderPathPointsByPathOrder(db *gorm.DB) *gorm.DB {
	return db.Order("path_order ASC")
}

// convertToDriverOffers converts database model to domain model
func convertToDriverOffers(driverOfferDB []entity.DriverOfferDB) []*model.Offer {
	offers := make([]*model.Offer, 0, len(driverOfferDB))

	for i := range driverOfferDB {
		if domainModel := driverOfferDB[i].ToDriverOffer(); domainModel != nil {
			offers = append(offers, domainModel)
		}
	}

	return offers
}
