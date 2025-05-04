package entity

import (
	"matching-engine/internal/model"
	"time"
)

// RiderRequestDB is the database model for rider requests
type RiderRequestDB struct {
	ID                        string        `gorm:"type:uuid;primaryKey"`
	UserID                    string        `gorm:"type:uuid;not null"`
	SourceLatitude            float64       `gorm:"type:decimal(10,8);not null"`
	SourceLongitude           float64       `gorm:"type:decimal(11,8);not null"`
	DestinationLatitude       float64       `gorm:"type:decimal(10,8);not null"`
	DestinationLongitude      float64       `gorm:"type:decimal(11,8);not null"`
	EarliestDepartureTime     time.Time     `gorm:"type:timestamp with time zone;not null"`
	LatestArrivalTime         time.Time     `gorm:"type:timestamp with time zone;not null"`
	MaxWalkingDurationMinutes time.Duration `gorm:"type:interval;default:'5 minutes'"`
	NumberOfRiders            int           `gorm:"not null;default:1;check:number_of_riders > 0"`
	SameGender                bool          `gorm:"not null;default:false"`
	AllowsSmoking             bool          `gorm:"not null;default:true"`
	AllowsPets                bool          `gorm:"not null;default:true"`
	IsMatched                 bool          `gorm:"default:false"`
}

// TableName specifies the table name for RiderRequestDB
func (RiderRequestDB) TableName() string {
	return "rider_requests"
}

// ToRiderRequest converts a RiderRequestDB to RiderRequest domain model
func (r *RiderRequestDB) ToRiderRequest() *models.RiderRequest {
    sourceCoord, _ := models.NewCoordinate(r.SourceLatitude, r.SourceLongitude)
    
    destCoord, _ := models.NewCoordinate(r.DestinationLatitude, r.DestinationLongitude)
    
    preferences := models.NewPreference(r.SameGender, r.AllowsSmoking, r.AllowsPets)
    
    // Call the constructor function properly and handle any potential errors
    riderRequest := models.NewRiderRequest(
        r.ID,
        r.UserID,
        *sourceCoord,
        *destCoord,
        r.EarliestDepartureTime,
        r.LatestArrivalTime,
        r.MaxWalkingDurationMinutes,
        r.NumberOfRiders,
        *preferences,
		r.IsMatched,
    )
    return riderRequest
}