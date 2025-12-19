package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/kuchida1981/graphql-sampleapp/internal/domain"
)

type PostgresWeatherAlertMetadataRepository struct {
	db *sql.DB
}

func NewPostgresWeatherAlertMetadataRepository(db *sql.DB) *PostgresWeatherAlertMetadataRepository {
	return &PostgresWeatherAlertMetadataRepository{db: db}
}

func (r *PostgresWeatherAlertMetadataRepository) SearchIDs(ctx context.Context, filter MetadataFilter) ([]string, error) {
	log.Printf("PostgresWeatherAlertMetadataRepository: Searching IDs with filter: %+v", filter)

	query := "SELECT id FROM weather_alert_metadata"
	var conditions []string
	var args []interface{}
	argIndex := 1

	if filter.Region != nil {
		conditions = append(conditions, fmt.Sprintf("region = $%d", argIndex))
		args = append(args, *filter.Region)
		argIndex++
	}

	if filter.IssuedAfter != nil {
		conditions = append(conditions, fmt.Sprintf("issued_at >= $%d", argIndex))
		args = append(args, *filter.IssuedAfter)
		argIndex++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY issued_at DESC"

	log.Printf("PostgresWeatherAlertMetadataRepository: Executing query: %s with args: %v", query, args)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Printf("PostgresWeatherAlertMetadataRepository: Failed to query: %v", err)
		return nil, fmt.Errorf("failed to search weather alert metadata: %w", err)
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			log.Printf("PostgresWeatherAlertMetadataRepository: Failed to scan ID: %v", err)
			return nil, fmt.Errorf("failed to scan ID: %w", err)
		}
		ids = append(ids, id)
	}

	if err := rows.Err(); err != nil {
		log.Printf("PostgresWeatherAlertMetadataRepository: Row iteration error: %v", err)
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	log.Printf("PostgresWeatherAlertMetadataRepository: Found %d IDs", len(ids))
	return ids, nil
}

func (r *PostgresWeatherAlertMetadataRepository) Search(ctx context.Context, filter MetadataFilter) ([]*domain.WeatherAlertMetadata, error) {
	log.Printf("PostgresWeatherAlertMetadataRepository: Searching metadata with filter: %+v", filter)

	query := "SELECT id, region, severity, issued_at, created_at FROM weather_alert_metadata"
	var conditions []string
	var args []interface{}
	argIndex := 1

	if filter.Region != nil {
		conditions = append(conditions, fmt.Sprintf("region = $%d", argIndex))
		args = append(args, *filter.Region)
		argIndex++
	}

	if filter.IssuedAfter != nil {
		conditions = append(conditions, fmt.Sprintf("issued_at >= $%d", argIndex))
		args = append(args, *filter.IssuedAfter)
		argIndex++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY issued_at DESC"

	log.Printf("PostgresWeatherAlertMetadataRepository: Executing query: %s with args: %v", query, args)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Printf("PostgresWeatherAlertMetadataRepository: Failed to query: %v", err)
		return nil, fmt.Errorf("failed to search weather alert metadata: %w", err)
	}
	defer rows.Close()

	var metadataList []*domain.WeatherAlertMetadata
	for rows.Next() {
		var metadata domain.WeatherAlertMetadata
		if err := rows.Scan(&metadata.ID, &metadata.Region, &metadata.Severity, &metadata.IssuedAt, &metadata.CreatedAt); err != nil {
			log.Printf("PostgresWeatherAlertMetadataRepository: Failed to scan metadata: %v", err)
			return nil, fmt.Errorf("failed to scan metadata: %w", err)
		}
		metadataList = append(metadataList, &metadata)
	}

	if err := rows.Err(); err != nil {
		log.Printf("PostgresWeatherAlertMetadataRepository: Row iteration error: %v", err)
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	log.Printf("PostgresWeatherAlertMetadataRepository: Found %d metadata records", len(metadataList))
	return metadataList, nil
}
