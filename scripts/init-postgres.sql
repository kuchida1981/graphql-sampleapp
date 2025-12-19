-- Initialize PostgreSQL database for GraphQL sample app
-- This script creates the users table and related indexes

CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes for common query patterns
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at);

-- Weather Alert Metadata table for hybrid PostgreSQL + Firestore storage
CREATE TABLE IF NOT EXISTS weather_alert_metadata (
    id VARCHAR(255) PRIMARY KEY,
    region VARCHAR(100) NOT NULL,
    severity VARCHAR(50) NOT NULL,
    issued_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes for weather alert search queries
CREATE INDEX IF NOT EXISTS idx_weather_alert_metadata_region ON weather_alert_metadata(region);
CREATE INDEX IF NOT EXISTS idx_weather_alert_metadata_issued_at ON weather_alert_metadata(issued_at);
CREATE INDEX IF NOT EXISTS idx_weather_alert_metadata_severity ON weather_alert_metadata(severity);
