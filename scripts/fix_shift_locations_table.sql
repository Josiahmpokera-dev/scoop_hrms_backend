-- Fix shift_locations join table structure
-- This script drops and recreates the shift_locations table with correct column names

-- Drop the existing table if it has wrong structure
DROP TABLE IF EXISTS shift_locations;

-- Create the correct join table
CREATE TABLE shift_locations (
    shift_id INTEGER NOT NULL,
    location_id INTEGER NOT NULL,
    PRIMARY KEY (shift_id, location_id),
    CONSTRAINT fk_shift_locations_shift FOREIGN KEY (shift_id) REFERENCES shifts(id) ON DELETE CASCADE,
    CONSTRAINT fk_shift_locations_location FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE CASCADE
);

-- Create indexes for better query performance
CREATE INDEX idx_shift_locations_shift_id ON shift_locations(shift_id);
CREATE INDEX idx_shift_locations_location_id ON shift_locations(location_id);
