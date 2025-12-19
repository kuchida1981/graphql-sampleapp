//go:build ignore

package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/kuchida1981/graphql-sampleapp/internal/firestore"
	"github.com/kuchida1981/graphql-sampleapp/internal/postgres"
)

func main() {
	ctx := context.Background()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://graphql_user:graphql_pass@localhost:5432/graphql_db?sslmode=disable"
		log.Printf("DATABASE_URL not set, using default: %s", databaseURL)
	}

	db, err := postgres.NewClient(ctx, databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer db.Close()

	projectID := os.Getenv("GCP_PROJECT_ID")
	if projectID == "" {
		projectID = "demo-project"
	}

	firestoreClient, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to connect to Firestore: %v", err)
	}
	defer firestoreClient.Close()

	metadataRecords := []struct {
		id        string
		region    string
		severity  string
		issuedAt  time.Time
		createdAt time.Time
	}{
		{"alert-tokyo-001", "Tokyo", "warning", time.Now().Add(-48 * time.Hour), time.Now().Add(-48 * time.Hour)},
		{"alert-tokyo-002", "Tokyo", "info", time.Now().Add(-24 * time.Hour), time.Now().Add(-24 * time.Hour)},
		{"alert-tokyo-003", "Tokyo", "critical", time.Now().Add(-12 * time.Hour), time.Now().Add(-12 * time.Hour)},
		{"alert-osaka-001", "Osaka", "warning", time.Now().Add(-36 * time.Hour), time.Now().Add(-36 * time.Hour)},
		{"alert-osaka-002", "Osaka", "info", time.Now().Add(-18 * time.Hour), time.Now().Add(-18 * time.Hour)},
		{"alert-osaka-003", "Osaka", "critical", time.Now().Add(-6 * time.Hour), time.Now().Add(-6 * time.Hour)},
		{"alert-kyoto-001", "Kyoto", "warning", time.Now().Add(-30 * time.Hour), time.Now().Add(-30 * time.Hour)},
		{"alert-kyoto-002", "Kyoto", "info", time.Now().Add(-15 * time.Hour), time.Now().Add(-15 * time.Hour)},
		{"alert-kyoto-003", "Kyoto", "warning", time.Now().Add(-3 * time.Hour), time.Now().Add(-3 * time.Hour)},
		{"alert-kyoto-004", "Kyoto", "critical", time.Now().Add(-1 * time.Hour), time.Now().Add(-1 * time.Hour)},
	}

	pgQuery := `
		INSERT INTO weather_alert_metadata (id, region, severity, issued_at, created_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE
		SET region = EXCLUDED.region,
		    severity = EXCLUDED.severity,
		    issued_at = EXCLUDED.issued_at,
		    created_at = EXCLUDED.created_at
	`

	for _, record := range metadataRecords {
		_, err := db.ExecContext(ctx, pgQuery, record.id, record.region, record.severity, record.issuedAt, record.createdAt)
		if err != nil {
			log.Printf("Failed to insert metadata %s: %v", record.id, err)
			continue
		}
		log.Printf("Inserted metadata: %s", record.id)
	}

	firestoreAlerts := []struct {
		id              string
		title           string
		description     string
		rawData         map[string]interface{}
		affectedAreas   []string
		recommendations []string
	}{
		{
			"alert-tokyo-001",
			"Strong Wind Warning",
			"Strong winds expected in Tokyo area",
			map[string]interface{}{
				"temperature":   map[string]interface{}{"value": 15.2, "unit": "celsius"},
				"windSpeed":     map[string]interface{}{"value": 25.5, "unit": "m/s"},
				"precipitation": map[string]interface{}{"value": 0, "unit": "mm"},
				"pressure":      map[string]interface{}{"value": 1013.2, "unit": "hPa"},
			},
			[]string{"Chiyoda", "Minato", "Shibuya"},
			[]string{"Stay indoors", "Secure loose objects"},
		},
		{
			"alert-tokyo-002",
			"Clear Weather Information",
			"Clear weather expected for the next 24 hours",
			map[string]interface{}{
				"temperature":   map[string]interface{}{"value": 22.5, "unit": "celsius"},
				"windSpeed":     map[string]interface{}{"value": 5.2, "unit": "m/s"},
				"precipitation": map[string]interface{}{"value": 0, "unit": "mm"},
				"pressure":      map[string]interface{}{"value": 1015.8, "unit": "hPa"},
			},
			[]string{"All areas"},
			[]string{"Good day for outdoor activities"},
		},
		{
			"alert-tokyo-003",
			"Severe Thunderstorm Critical Alert",
			"Severe thunderstorm with heavy rainfall imminent",
			map[string]interface{}{
				"temperature":   map[string]interface{}{"value": 18.0, "unit": "celsius"},
				"windSpeed":     map[string]interface{}{"value": 35.0, "unit": "m/s"},
				"precipitation": map[string]interface{}{"value": 80, "unit": "mm"},
				"pressure":      map[string]interface{}{"value": 995.5, "unit": "hPa"},
			},
			[]string{"All areas"},
			[]string{"Seek shelter immediately", "Avoid travel"},
		},
		{
			"alert-osaka-001",
			"Heavy Rain Warning",
			"Heavy rainfall expected in Osaka region",
			map[string]interface{}{
				"temperature":   map[string]interface{}{"value": 19.5, "unit": "celsius"},
				"windSpeed":     map[string]interface{}{"value": 15.0, "unit": "m/s"},
				"precipitation": map[string]interface{}{"value": 50, "unit": "mm"},
				"pressure":      map[string]interface{}{"value": 1008.0, "unit": "hPa"},
			},
			[]string{"Kita", "Chuo", "Naniwa"},
			[]string{"Carry umbrella", "Watch for flooding"},
		},
		{
			"alert-osaka-002",
			"Mild Weather Information",
			"Mild weather conditions throughout the day",
			map[string]interface{}{
				"temperature":   map[string]interface{}{"value": 20.0, "unit": "celsius"},
				"windSpeed":     map[string]interface{}{"value": 8.0, "unit": "m/s"},
				"precipitation": map[string]interface{}{"value": 0, "unit": "mm"},
				"pressure":      map[string]interface{}{"value": 1012.5, "unit": "hPa"},
			},
			[]string{"All areas"},
			[]string{"Enjoy your day"},
		},
		{
			"alert-osaka-003",
			"Typhoon Critical Alert",
			"Typhoon approaching Osaka bay area",
			map[string]interface{}{
				"temperature":   map[string]interface{}{"value": 16.5, "unit": "celsius"},
				"windSpeed":     map[string]interface{}{"value": 45.0, "unit": "m/s"},
				"precipitation": map[string]interface{}{"value": 120, "unit": "mm"},
				"pressure":      map[string]interface{}{"value": 985.0, "unit": "hPa"},
			},
			[]string{"All areas"},
			[]string{"Evacuate if instructed", "Stock emergency supplies"},
		},
		{
			"alert-kyoto-001",
			"Fog Warning",
			"Dense fog reducing visibility",
			map[string]interface{}{
				"temperature":   map[string]interface{}{"value": 12.0, "unit": "celsius"},
				"windSpeed":     map[string]interface{}{"value": 3.0, "unit": "m/s"},
				"precipitation": map[string]interface{}{"value": 0, "unit": "mm"},
				"pressure":      map[string]interface{}{"value": 1016.0, "unit": "hPa"},
			},
			[]string{"Northern districts"},
			[]string{"Drive carefully", "Use fog lights"},
		},
		{
			"alert-kyoto-002",
			"Pleasant Weather Information",
			"Pleasant spring weather expected",
			map[string]interface{}{
				"temperature":   map[string]interface{}{"value": 18.5, "unit": "celsius"},
				"windSpeed":     map[string]interface{}{"value": 6.5, "unit": "m/s"},
				"precipitation": map[string]interface{}{"value": 0, "unit": "mm"},
				"pressure":      map[string]interface{}{"value": 1014.2, "unit": "hPa"},
			},
			[]string{"All areas"},
			[]string{"Perfect for sightseeing"},
		},
		{
			"alert-kyoto-003",
			"Thunderstorm Warning",
			"Isolated thunderstorms possible in the evening",
			map[string]interface{}{
				"temperature":   map[string]interface{}{"value": 21.0, "unit": "celsius"},
				"windSpeed":     map[string]interface{}{"value": 12.0, "unit": "m/s"},
				"precipitation": map[string]interface{}{"value": 25, "unit": "mm"},
				"pressure":      map[string]interface{}{"value": 1010.5, "unit": "hPa"},
			},
			[]string{"Eastern districts"},
			[]string{"Postpone outdoor activities", "Stay informed"},
		},
		{
			"alert-kyoto-004",
			"Flash Flood Critical Alert",
			"Flash flood warning due to heavy upstream rainfall",
			map[string]interface{}{
				"temperature":   map[string]interface{}{"value": 17.5, "unit": "celsius"},
				"windSpeed":     map[string]interface{}{"value": 18.0, "unit": "m/s"},
				"precipitation": map[string]interface{}{"value": 95, "unit": "mm"},
				"pressure":      map[string]interface{}{"value": 1002.0, "unit": "hPa"},
			},
			[]string{"Riverside areas"},
			[]string{"Move to higher ground", "Avoid riverbanks"},
		},
	}

	for _, alert := range firestoreAlerts {
		_, err := firestoreClient.Collection("weatherAlerts").Doc(alert.id).Set(ctx, map[string]interface{}{
			"id":              alert.id,
			"title":           alert.title,
			"description":     alert.description,
			"rawData":         alert.rawData,
			"affectedAreas":   alert.affectedAreas,
			"recommendations": alert.recommendations,
		})
		if err != nil {
			log.Printf("Failed to insert Firestore alert %s: %v", alert.id, err)
			continue
		}
		log.Printf("Inserted Firestore alert: %s", alert.id)
	}

	log.Println("Weather alerts seeding completed successfully!")
}
