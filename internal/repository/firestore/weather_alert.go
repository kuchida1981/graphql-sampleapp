package firestore

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	"github.com/kuchida1981/graphql-sampleapp/internal/domain"
)

type FirestoreWeatherAlertRepository struct {
	client *firestore.Client
}

func NewFirestoreWeatherAlertRepository(client *firestore.Client) *FirestoreWeatherAlertRepository {
	return &FirestoreWeatherAlertRepository{
		client: client,
	}
}

func (r *FirestoreWeatherAlertRepository) GetByID(ctx context.Context, id string) (*domain.WeatherAlert, error) {
	log.Printf("FirestoreWeatherAlertRepository: Fetching weather alert with ID: %s", id)

	doc, err := r.client.Collection("weatherAlerts").Doc(id).Get(ctx)
	if err != nil {
		log.Printf("FirestoreWeatherAlertRepository: Error fetching weather alert %s: %v", id, err)
		return nil, fmt.Errorf("weather alert not found: %w", err)
	}

	var alert domain.WeatherAlert
	if err := doc.DataTo(&alert); err != nil {
		log.Printf("FirestoreWeatherAlertRepository: Error converting document to WeatherAlert: %v", err)
		return nil, err
	}

	log.Printf("FirestoreWeatherAlertRepository: Successfully fetched weather alert: %s", id)
	return &alert, nil
}

func (r *FirestoreWeatherAlertRepository) GetByIDs(ctx context.Context, ids []string) ([]*domain.WeatherAlert, error) {
	log.Printf("FirestoreWeatherAlertRepository: Fetching %d weather alerts", len(ids))

	if len(ids) == 0 {
		return []*domain.WeatherAlert{}, nil
	}

	var alerts []*domain.WeatherAlert
	for _, id := range ids {
		doc, err := r.client.Collection("weatherAlerts").Doc(id).Get(ctx)
		if err != nil {
			log.Printf("FirestoreWeatherAlertRepository: Warning - failed to fetch weather alert %s: %v (skipping)", id, err)
			continue
		}

		var alert domain.WeatherAlert
		if err := doc.DataTo(&alert); err != nil {
			log.Printf("FirestoreWeatherAlertRepository: Warning - failed to convert document %s: %v (skipping)", id, err)
			continue
		}

		alerts = append(alerts, &alert)
	}

	log.Printf("FirestoreWeatherAlertRepository: Successfully fetched %d out of %d weather alerts", len(alerts), len(ids))
	return alerts, nil
}
