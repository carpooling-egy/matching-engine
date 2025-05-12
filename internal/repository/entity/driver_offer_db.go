package entity

import (
	"time"

	"matching-engine/internal/enums"
	"matching-engine/internal/model"
)

// DriverOfferDB is the database model for driver offers
type DriverOfferDB struct {
	ID                      string        `gorm:"type:varchar(50);primaryKey"`
	UserID                  string        `gorm:"type:varchar(50);not null"`

	SourceLatitude          float64       `gorm:"type:decimal(10,8);not null"`
	SourceLongitude         float64       `gorm:"type:decimal(11,8);not null"`

	DestinationLatitude     float64       `gorm:"type:decimal(10,8);not null"`
	DestinationLongitude    float64       `gorm:"type:decimal(11,8);not null"`

	DepartureTime           time.Time     `gorm:"type:timestamp with time zone;not null"`
	EstimatedArrivalTime    time.Time     `gorm:"type:timestamp with time zone"`
	MaxEstimatedArrivalTime time.Time     `gorm:"type:timestamp with time zone"`
	
	DetourDurationMinutes   time.Duration `gorm:"default:0"`
	Capacity                int           `gorm:"not null;check:capacity > 0"`
	CurrentNumberOfRequests int           `gorm:"not null;default:0"`

	SameGender              bool          `gorm:"not null;default:false"`
	AllowsSmoking           bool          `gorm:"not null;default:true"`
	AllowsPets              bool          `gorm:"not null;default:true"`
	UserGender              enums.Gender  `gorm:"type:gender_type;not null"`

	PathPoints              []PathPointDB `gorm:"foreignKey:DriverOfferID"`
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
	pathPoints := make([]*model.PathPoint, 0, len(d.PathPoints))

	// Create a map to store pickup points by request ID
	pickupPointsByRequestID := make(map[string]*model.PathPoint)
	// And store path points by request ID for both pickup and dropoff
	dropoffPointsByRequestID := make(map[string]*model.PathPoint)

	matchedRequests := make([]*model.MatchedRequest, 0, len(d.PathPoints)/2)

	sourcePoint := model.NewPathPoint(*sourceCoord, enums.Source, d.DepartureTime, nil)
	pathPoints = append(pathPoints, sourcePoint)

	for _, pp := range d.PathPoints {
		pathPoint := pp.ToPathPoint()
		pathPoints = append(pathPoints, pathPoint)

		// If the path point has a request associated with it
		if pathPoint.Owner() != nil {
			request, _ := pathPoint.Owner().AsRequest()
			requestID := request.ID()

			// If it's a pickup point, store it
			if pathPoint.PointType() == enums.Pickup {
				pickupPointsByRequestID[requestID] = pathPoint
			} else if pathPoint.PointType() == enums.Dropoff {
				dropoffPointsByRequestID[requestID] = pathPoint
			}

		}
	}

	destPoint := model.NewPathPoint(*destCoord, enums.Destination, d.MaxEstimatedArrivalTime, nil)
	pathPoints = append(pathPoints, destPoint)

	// Now create matched requests
	for requestID, pickupPoint := range pickupPointsByRequestID {
		// Check if there's a corresponding dropoff point
		dropoffPoint, exists := dropoffPointsByRequestID[requestID]
		if exists {
			request, _ := pickupPoint.Owner().AsRequest()
			// Create matched request using the constructor
			matchedRequest := model.NewMatchedRequest(request, *pickupPoint, *dropoffPoint)
			matchedRequests = append(matchedRequests, matchedRequest)
		}
	}

	// Create driver offer with the already converted path points
	driverOffer := model.NewOffer(
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
		matchedRequests,
	)

	driverOffer.Path()[0].SetOwner(driverOffer)
	driverOffer.Path()[len(driverOffer.Path())-1].SetOwner(driverOffer)

	return driverOffer
}
