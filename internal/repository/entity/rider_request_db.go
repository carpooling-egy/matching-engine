package entity

import (
	"time"

	"matching-engine/internal/enums"
	"matching-engine/internal/model"
)

// RiderRequestDB is the database model for rider requests
type RiderRequestDB struct {
	ID                        string        `gorm:"type:varchar(50);primaryKey"`
	UserID                    string        `gorm:"type:varchar(50);not null"`
	SourceLatitude            float64       `gorm:"type:decimal(10,8);not null"`
	SourceLongitude           float64       `gorm:"type:decimal(11,8);not null"`
	DestinationLatitude       float64       `gorm:"type:decimal(10,8);not null"`
	DestinationLongitude      float64       `gorm:"type:decimal(11,8);not null"`
	EarliestDepartureTime     time.Time     `gorm:"type:timestamp with time zone;not null"`
	LatestArrivalTime         time.Time     `gorm:"type:timestamp with time zone;not null"`
	MaxWalkingDurationMinutes time.Duration `gorm:"column:max_walking_duration_minutes;default:10"`
	NumberOfRiders            int           `gorm:"not null;default:1;check:number_of_riders > 0"`
	SameGender                bool          `gorm:"not null;default:false"`
	AllowsSmoking             bool          `gorm:"not null;default:true"`
	AllowsPets                bool          `gorm:"not null;default:true"`
	UserGender                enums.Gender  `gorm:"type:gender_type;not null"`
}

// TableName specifies the table name for RiderRequestDB
func (RiderRequestDB) TableName() string {
	return "rider_requests"
}

// ToRiderRequest converts a RiderRequestDB to RiderRequest domain model
func (r *RiderRequestDB) ToRiderRequest() *model.Request {
	sourceCoord, _ := model.NewCoordinate(r.SourceLatitude, r.SourceLongitude)

	destCoord, _ := model.NewCoordinate(r.DestinationLatitude, r.DestinationLongitude)

	preferences := model.NewPreference(r.UserGender, r.SameGender, r.AllowsSmoking, r.AllowsPets)

	// Call the constructor function properly and handle any potential errors
	riderRequest := model.NewRequest(
		r.ID,
		r.UserID,
		*sourceCoord,
		*destCoord,
		r.EarliestDepartureTime,
		r.LatestArrivalTime,
		r.MaxWalkingDurationMinutes,
		r.NumberOfRiders,
		*preferences,
	)
	return riderRequest
}
