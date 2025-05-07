package entity

import (
	"matching-engine/internal/enums"
	"matching-engine/internal/model"
	"time"
)

// PathPointDB is the database model for path points
type PathPointDB struct {
	ID                  string          `gorm:"type:uuid;primaryKey"`
	DriverOfferID       string          `gorm:"type:uuid;not null"`
	PathOrder           int             `gorm:"not null"`
	PointType           enums.PointType `gorm:"column:type;type:point_type;not null"`
	Latitude            float64         `gorm:"type:decimal(10,8);not null"`
	Longitude           float64         `gorm:"type:decimal(11,8);not null"`
	ExpectedArrivalTime time.Time       `gorm:"type:timestamp with time zone;not null"`
	RiderRequestID      string          `gorm:"type:uuid;not null"`        // Foreign key field needs to be explicitly defined
	RiderRequest        *RiderRequestDB `gorm:"foreignKey:RiderRequestID"` // Specify the foreign key field name
}

// TableName specifies the table name for PathPointDB
func (PathPointDB) TableName() string {
	return "path_point"
}

// ToPathPoint converts a PathPointDB to PathPoint domain model
func (p *PathPointDB) ToPathPoint() *model.PathPoint {
	// Use the constructor for Coordinate
	coordinate, _ := model.NewCoordinate(p.Latitude, p.Longitude)

	var riderRequest *model.Request = p.RiderRequest.ToRiderRequest()

	// Fix: Call the constructor function properly
	return model.NewPathPoint(
		*coordinate,
		p.PointType,
		p.ExpectedArrivalTime,
		riderRequest,
	)
}
