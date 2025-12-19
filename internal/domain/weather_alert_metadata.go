package domain

import "time"

const (
	SeverityInfo     = "info"
	SeverityWarning  = "warning"
	SeverityCritical = "critical"
)

type WeatherAlertMetadata struct {
	ID        string
	Region    string
	Severity  string
	IssuedAt  time.Time
	CreatedAt time.Time
}
