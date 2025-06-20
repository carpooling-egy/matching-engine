-- NOTE: Users ,their gender and Cars tables are in a separate database
-- The schema below assumes user UUIDs will be provided from that external source
CREATE TYPE point_type AS ENUM ('pickup', 'dropoff');
CREATE TYPE gender_type AS ENUM ('male', 'female');


-- Rider requests table
CREATE TABLE rider_requests (
    id VARCHAR(50) PRIMARY KEY,
    user_id VARCHAR(50) NOT NULL,

    -- source
    source_latitude DECIMAL(10, 8) NOT NULL,
    source_longitude DECIMAL(11, 8) NOT NULL,
    source_address TEXT,

    -- destination
    destination_latitude DECIMAL(10, 8) NOT NULL,
    destination_longitude DECIMAL(11, 8) NOT NULL,
    destination_address TEXT,

    earliest_departure_time TIMESTAMP WITH TIME ZONE NOT NULL,
    latest_arrival_time TIMESTAMP WITH TIME ZONE NOT NULL,
    max_walking_duration_minutes INTEGER DEFAULT 5,
    number_of_riders INTEGER NOT NULL DEFAULT 1 CHECK (number_of_riders > 0),

    -- Boolean preferences
    same_gender BOOLEAN NOT NULL DEFAULT FALSE,
    user_gender gender_type NOT NULL,

    is_matched BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Driver offers table
CREATE TABLE driver_offers (
    id VARCHAR(50) PRIMARY KEY,
    user_id VARCHAR(50) NOT NULL,
    --source
    source_latitude DECIMAL(10, 8) NOT NULL,
    source_longitude DECIMAL(11, 8) NOT NULL,
    source_address TEXT,

    --destination
    destination_latitude DECIMAL(10, 8) NOT NULL,
    destination_longitude DECIMAL(11, 8) NOT NULL,
    destination_address TEXT,

    departure_time TIMESTAMP WITH TIME ZONE NOT NULL,
    max_estimated_arrival_time TIMESTAMP WITH TIME ZONE,
    estimated_arrival_time TIMESTAMP WITH TIME ZONE, 

    detour_duration_minutes INTEGER DEFAULT 0,
    capacity INTEGER NOT NULL CHECK (capacity > 0),

    current_number_of_requests INTEGER NOT NULL DEFAULT 0,

    -- Boolean preferences
    same_gender BOOLEAN NOT NULL DEFAULT FALSE,
    user_gender gender_type NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE path_point (
    id VARCHAR(50) PRIMARY KEY,
    driver_offer_id VARCHAR(50) NOT NULL REFERENCES driver_offers(id) ON DELETE CASCADE,
    path_order INTEGER NOT NULL,
    type point_type NOT NULL,
    
    latitude DECIMAL(10, 8) NOT NULL,
    longitude DECIMAL(11, 8) NOT NULL,
    walking_duration_minutes INTEGER DEFAULT 0,
    address TEXT,

    -- pickup time or dropoff time depending on the type
    expected_arrival_time TIMESTAMP WITH TIME ZONE NOT NULL,
    rider_request_id VARCHAR(50) NOT NULL REFERENCES rider_requests(id) ON DELETE CASCADE,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (driver_offer_id, path_order)
);

-- Table to store matching information between driver offers and rider requests
CREATE TABLE ride_matches (
    driver_offer_id VARCHAR(50) NOT NULL REFERENCES driver_offers(id) ON DELETE CASCADE,
    rider_request_id VARCHAR(50) NOT NULL REFERENCES rider_requests(id) ON DELETE CASCADE,

    pickup_point_id VARCHAR(50) NOT NULL REFERENCES path_point(id) ON DELETE CASCADE,
    dropoff_point_id VARCHAR(50) NOT NULL REFERENCES path_point(id) ON DELETE CASCADE,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (driver_offer_id, rider_request_id)
);

-- Indexes for performance optimization
CREATE INDEX idx_rider_requests_matching ON rider_requests(is_matched, earliest_departure_time);
CREATE INDEX idx_path_point_driver_path ON path_point(driver_offer_id, path_order);
CREATE INDEX idx_driver_offers_availability ON driver_offers(departure_time, current_number_of_requests);
CREATE INDEX idx_path_point_rider_request ON path_point(rider_request_id);

-- Create a function to update the updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers for each table
CREATE TRIGGER update_rider_requests_updated_at
BEFORE UPDATE ON rider_requests
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_driver_offers_updated_at
BEFORE UPDATE ON driver_offers
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_ride_matches_updated_at
BEFORE UPDATE ON ride_matches
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_path_point_updated_at
BEFORE UPDATE ON path_point
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();