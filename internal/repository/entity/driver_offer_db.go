package entity

import (
	"time"

	"matching-engine/internal/enums"
	"matching-engine/internal/model"
)

// DriverOfferDB is the database model for driver offers
type DriverOfferDB struct {
	ID     string `gorm:"type:varchar(50);primaryKey"`
	UserID string `gorm:"type:varchar(50);not null"`

	SourceLatitude  float64 `gorm:"type:decimal(10,8);not null"`
	SourceLongitude float64 `gorm:"type:decimal(11,8);not null"`

	DestinationLatitude  float64 `gorm:"type:decimal(10,8);not null"`
	DestinationLongitude float64 `gorm:"type:decimal(11,8);not null"`

	DepartureTime           time.Time `gorm:"type:timestamp with time zone;not null"`
	EstimatedArrivalTime    time.Time `gorm:"type:timestamp with time zone"`
	MaxEstimatedArrivalTime time.Time `gorm:"type:timestamp with time zone"`

	DetourDurationMinutes   int `gorm:"default:0"`
	Capacity                int `gorm:"not null;check:capacity > 0"`
	CurrentNumberOfRequests int `gorm:"not null;default:0"`

	SameGender    bool          `gorm:"not null;default:false"`
	AllowsSmoking bool          `gorm:"not null;default:true"`
	AllowsPets    bool          `gorm:"not null;default:true"`
	UserGender    enums.Gender  `gorm:"type:gender_type;not null"`
	PathPoints    []PathPointDB `gorm:"foreignKey:DriverOfferID"`
}

// TableName specifies the table name for DriverOfferDB
func (DriverOfferDB) TableName() string {
	return "driver_offers"
}

// ToDriverOffer converts a DriverOfferDB to DriverOffer domain model
func (d *DriverOfferDB) ToDriverOffer() *model.Offer {
	// Create source coordinate
	sourceCoord, _ := model.NewCoordinate(d.SourceLatitude, d.SourceLongitude)

	// Create destination coordinate
	destCoord, _ := model.NewCoordinate(d.DestinationLatitude, d.DestinationLongitude)

	// Create preference using the constructor
	preferences := model.NewPreference(d.UserGender, d.SameGender, d.AllowsSmoking, d.AllowsPets)

	// Pre-allocate pathPoints slice
	pathPoints := make([]model.PathPoint, 0, len(d.PathPoints)+2) // +2 for source and destination

	// Create a map to store requests by ID
	requestsMap := make(map[string]*model.Request)

	// Add source point
	sourcePoint := model.NewPathPoint(*sourceCoord, enums.Source, d.DepartureTime, nil, 0)
	pathPoints = append(pathPoints, *sourcePoint) // Use the value, not the pointer

	// Process path points from database
	for _, pp := range d.PathPoints {
		pathPoint := pp.ToPathPoint()
		pathPoints = append(pathPoints, *pathPoint) // Use the value, not the pointer

		// If the path point has a request associated with it
		if pathPoint.Owner() != nil {
			request, ok := pathPoint.Owner().AsRequest()
			if ok && request != nil {
				requestID := request.ID()
				requestsMap[requestID] = request
			}
		}
	}

	// Add destination point
	arrivalTime := d.MaxEstimatedArrivalTime
	if arrivalTime.IsZero() {
		arrivalTime = d.EstimatedArrivalTime
	}
	destPoint := model.NewPathPoint(*destCoord, enums.Destination, arrivalTime, nil, 0)
	pathPoints = append(pathPoints, *destPoint) // Use the value, not the pointer

	// Convert requests map to slice
	requests := make([]*model.Request, 0, len(requestsMap))
	for _, request := range requestsMap {
		requests = append(requests, request)
	}

	// Create driver offer with the already converted path points
	driverOffer := model.NewOffer(
		d.ID,
		d.UserID,
		*sourceCoord,
		*destCoord,
		d.DepartureTime,
		time.Duration(d.DetourDurationMinutes)*time.Minute,
		d.Capacity,
		*preferences,
		d.MaxEstimatedArrivalTime,
		d.CurrentNumberOfRequests,
		pathPoints,
		requests,
	)

	// Set the driver as owner of the first and last path points
	if len(driverOffer.Path()) > 0 {
		driverOffer.Path()[0].SetOwner(driverOffer)
		driverOffer.Path()[len(driverOffer.Path())-1].SetOwner(driverOffer)
	}

	return driverOffer
}
