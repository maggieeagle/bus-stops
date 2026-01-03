CREATE TABLE IF NOT EXISTS stops (
    stop_id BIGINT PRIMARY KEY,
    stop_code VARCHAR(50),
    stop_name VARCHAR(255) NOT NULL,
    stop_lat DOUBLE PRECISION NOT NULL,
    stop_lon DOUBLE PRECISION NOT NULL,
    zone_id BIGINT,
    alias VARCHAR(255),
    stop_area VARCHAR(255),
    stop_desc TEXT,
    lest_x DOUBLE PRECISION,
    lest_y DOUBLE PRECISION,
    zone_name VARCHAR(100),
    authority VARCHAR(100) 
);

CREATE TABLE IF NOT EXISTS routes (
    route_id CHAR(32) PRIMARY KEY,
    agency_id INT NOT NULL,
    route_short_name VARCHAR(50) NOT NULL,
    route_long_name VARCHAR(255) NOT NULL,
    route_type SMALLINT NOT NULL CHECK (route_type IN (0,1,2,3,4,5,6,7)),  
    route_color CHAR(6) DEFAULT 'FFFFFF',
    competent_authority VARCHAR(100),
    route_desc TEXT
);


CREATE TABLE IF NOT EXISTS trips (
    trip_id BIGINT PRIMARY KEY,
    route_id CHAR(32) NOT NULL,
    service_id BIGINT NOT NULL,
    trip_headsign VARCHAR(255),
    trip_long_name VARCHAR(255),
    direction_code VARCHAR(10),
    shape_id BIGINT,
    wheelchair_accessible SMALLINT DEFAULT 0 CHECK (wheelchair_accessible IN (0,1,2)),

    CONSTRAINT fk_trips_route FOREIGN KEY (route_id) REFERENCES routes (route_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS stop_times (
    trip_id BIGINT NOT NULL,
    arrival_time TEXT NOT NULL,
    departure_time TEXT NOT NULL,
    stop_id BIGINT NOT NULL,
    stop_sequence INT NOT NULL CHECK (stop_sequence > 0),
    pickup_type SMALLINT DEFAULT 0 CHECK (pickup_type IN (0,1,2,3)),
    drop_off_type SMALLINT DEFAULT 0 CHECK (drop_off_type IN (0,1,2,3)),

    PRIMARY KEY (trip_id, stop_sequence),

    CONSTRAINT fk_stop_times_trip FOREIGN KEY (trip_id) REFERENCES trips (trip_id) ON DELETE CASCADE,
    CONSTRAINT fk_stop_times_stop FOREIGN KEY (stop_id) REFERENCES stops (stop_id)ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_trips_route_id ON trips (route_id);
CREATE INDEX IF NOT EXISTS idx_trips_service_id ON trips (service_id);
CREATE INDEX IF NOT EXISTS idx_trips_shape_id ON trips (shape_id);

CREATE INDEX IF NOT EXISTS idx_stop_times_stop_id ON stop_times (stop_id);
CREATE INDEX IF NOT EXISTS idx_stop_times_trip_id ON stop_times (trip_id);
