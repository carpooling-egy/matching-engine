package entity

// RideMatchDB is the database model for ride matches
type RideMatchDB struct {
	DriverOfferID  string `gorm:"type:uuid;not null;primaryKey"`
	RiderRequestID string `gorm:"type:uuid;not null;primaryKey"`
}

// TableName specifies the table name for RideMatchDB
func (RideMatchDB) TableName() string {
	return "ride_matches"
}
