package entity

import (
	"time"

	"matching-engine/internal/model"
)

// DriverOfferDB is the database model for driver offers
type DriverOfferDB struct {
	ID                      string        `gorm:"type:uuid;primaryKey"`
	UserID                  string        `gorm:"type:uuid;not null"`
	SourceLatitude          float64       `gorm:"type:decimal(10,8);not null"`
	SourceLongitude         float64       `gorm:"type:decimal(11,8);not null"`
	DestinationLatitude     float64       `gorm:"type:decimal(10,8);not null"`
	DestinationLongitude    float64       `gorm:"type:decimal(11,8);not null"`
	DepartureTime           time.Time     `gorm:"type:timestamp with time zone;not null"`
	DetourDurationMinutes   time.Duration `gorm:"type:interval;default:'0 minutes'"`
	Capacity                int           `gorm:"not null;check:capacity > 0"`
	CurrentNumberOfRequests int           `gorm:"not null;default:0"`
	SameGender              bool          `gorm:"not null;default:false"`
	AllowsSmoking           bool          `gorm:"not null;default:true"`
	AllowsPets              bool          `gorm:"not null;default:true"`
	PathPoints              []PathPointDB `gorm:"foreignKey:DriverOfferID"`
}

// TableName specifies the table name for DriverOfferDB
func (DriverOfferDB) TableName() string {
	return "driver_offers"
}

// ToDriverOffer converts a DriverOfferDB to DriverOffer domain model
func (d *DriverOfferDB) ToDriverOffer() (*models.DriverOffer) {
    // Create source coordinate
    sourceCoord, _ := models.NewCoordinate(d.SourceLatitude, d.SourceLongitude)


    // Create destination coordinate
    destCoord, _ := models.NewCoordinate(d.DestinationLatitude, d.DestinationLongitude)


    // Create preference using the constructor
    preferences := models.NewPreference(d.SameGender, d.AllowsSmoking, d.AllowsPets)

    // Pre-allocate pathPoints slice
    pathPoints := make([]models.PathPoint, 0, len(d.PathPoints))
    
    // Convert path points
    for _, pp := range d.PathPoints {
        pathPoint := pp.ToPathPoint()
        pathPoints = append(pathPoints, *pathPoint)
    }

    // Create driver offer with the already converted path points
    driverOffer := models.NewDriverOffer(
        d.ID,
        d.UserID,
        *sourceCoord,
        *destCoord,
        d.DepartureTime,
        d.DetourDurationMinutes,
        d.Capacity,
        *preferences,
        d.CurrentNumberOfRequests,
        pathPoints,
    )

    return driverOffer
}
