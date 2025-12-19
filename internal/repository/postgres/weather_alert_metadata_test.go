package postgres

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kuchida1981/graphql-sampleapp/internal/repository"
)

func TestPostgresWeatherAlertMetadataRepository_Search(t *testing.T) {
	now := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	region := "Tokyo"
	issuedAfter := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		filter  repository.MetadataFilter
		mockFn  func(mock sqlmock.Sqlmock)
		want    int
		wantErr bool
	}{
		{
			name:   "正常系: フィルタなしで検索",
			filter: repository.MetadataFilter{},
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "region", "severity", "issued_at", "created_at"}).
					AddRow("alert1", "Tokyo", "warning", now, now).
					AddRow("alert2", "Osaka", "info", now.Add(-24*time.Hour), now.Add(-24*time.Hour))
				mock.ExpectQuery("SELECT id, region, severity, issued_at, created_at FROM weather_alert_metadata ORDER BY issued_at DESC").
					WillReturnRows(rows)
			},
			want:    2,
			wantErr: false,
		},
		{
			name: "正常系: 地域フィルタ付き検索",
			filter: repository.MetadataFilter{
				Region: &region,
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "region", "severity", "issued_at", "created_at"}).
					AddRow("alert1", "Tokyo", "warning", now, now)
				mock.ExpectQuery("SELECT id, region, severity, issued_at, created_at FROM weather_alert_metadata WHERE region = \\$1 ORDER BY issued_at DESC").
					WithArgs("Tokyo").
					WillReturnRows(rows)
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "正常系: 日時フィルタ付き検索",
			filter: repository.MetadataFilter{
				IssuedAfter: &issuedAfter,
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "region", "severity", "issued_at", "created_at"}).
					AddRow("alert1", "Tokyo", "warning", now, now)
				mock.ExpectQuery("SELECT id, region, severity, issued_at, created_at FROM weather_alert_metadata WHERE issued_at >= \\$1 ORDER BY issued_at DESC").
					WithArgs(issuedAfter).
					WillReturnRows(rows)
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "正常系: 複合フィルタ検索",
			filter: repository.MetadataFilter{
				Region:      &region,
				IssuedAfter: &issuedAfter,
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "region", "severity", "issued_at", "created_at"}).
					AddRow("alert1", "Tokyo", "warning", now, now)
				mock.ExpectQuery("SELECT id, region, severity, issued_at, created_at FROM weather_alert_metadata WHERE region = \\$1 AND issued_at >= \\$2 ORDER BY issued_at DESC").
					WithArgs("Tokyo", issuedAfter).
					WillReturnRows(rows)
			},
			want:    1,
			wantErr: false,
		},
		{
			name:   "異常系: クエリエラー",
			filter: repository.MetadataFilter{},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, region, severity, issued_at, created_at FROM weather_alert_metadata ORDER BY issued_at DESC").
					WillReturnError(errors.New("database connection error"))
			},
			want:    0,
			wantErr: true,
		},
		{
			name:   "正常系: 結果が0件",
			filter: repository.MetadataFilter{},
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "region", "severity", "issued_at", "created_at"})
				mock.ExpectQuery("SELECT id, region, severity, issued_at, created_at FROM weather_alert_metadata ORDER BY issued_at DESC").
					WillReturnRows(rows)
			},
			want:    0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create mock: %v", err)
			}
			defer db.Close()

			tt.mockFn(mock)

			repo := NewPostgresWeatherAlertMetadataRepository(db)
			got, err := repo.Search(context.Background(), tt.filter)

			if (err != nil) != tt.wantErr {
				t.Errorf("Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && len(got) != tt.want {
				t.Errorf("Search() got %d records, want %d", len(got), tt.want)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestPostgresWeatherAlertMetadataRepository_SearchIDs(t *testing.T) {
	region := "Tokyo"
	issuedAfter := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		filter  repository.MetadataFilter
		mockFn  func(mock sqlmock.Sqlmock)
		want    []string
		wantErr bool
	}{
		{
			name:   "正常系: IDリスト取得",
			filter: repository.MetadataFilter{},
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow("alert1").
					AddRow("alert2").
					AddRow("alert3")
				mock.ExpectQuery("SELECT id FROM weather_alert_metadata ORDER BY issued_at DESC").
					WillReturnRows(rows)
			},
			want:    []string{"alert1", "alert2", "alert3"},
			wantErr: false,
		},
		{
			name: "正常系: 地域フィルタ付きIDリスト取得",
			filter: repository.MetadataFilter{
				Region: &region,
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow("alert1")
				mock.ExpectQuery("SELECT id FROM weather_alert_metadata WHERE region = \\$1 ORDER BY issued_at DESC").
					WithArgs("Tokyo").
					WillReturnRows(rows)
			},
			want:    []string{"alert1"},
			wantErr: false,
		},
		{
			name: "正常系: 日時フィルタ付きIDリスト取得",
			filter: repository.MetadataFilter{
				IssuedAfter: &issuedAfter,
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id"}).
					AddRow("alert1").
					AddRow("alert2")
				mock.ExpectQuery("SELECT id FROM weather_alert_metadata WHERE issued_at >= \\$1 ORDER BY issued_at DESC").
					WithArgs(issuedAfter).
					WillReturnRows(rows)
			},
			want:    []string{"alert1", "alert2"},
			wantErr: false,
		},
		{
			name:   "異常系: クエリエラー",
			filter: repository.MetadataFilter{},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id FROM weather_alert_metadata ORDER BY issued_at DESC").
					WillReturnError(errors.New("database connection error"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:   "正常系: 結果が0件",
			filter: repository.MetadataFilter{},
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id"})
				mock.ExpectQuery("SELECT id FROM weather_alert_metadata ORDER BY issued_at DESC").
					WillReturnRows(rows)
			},
			want:    []string{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create mock: %v", err)
			}
			defer db.Close()

			tt.mockFn(mock)

			repo := NewPostgresWeatherAlertMetadataRepository(db)
			got, err := repo.SearchIDs(context.Background(), tt.filter)

			if (err != nil) != tt.wantErr {
				t.Errorf("SearchIDs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(got) != len(tt.want) {
					t.Errorf("SearchIDs() got %d IDs, want %d", len(got), len(tt.want))
					return
				}

				for i, id := range got {
					if id != tt.want[i] {
						t.Errorf("SearchIDs() got[%d] = %v, want %v", i, id, tt.want[i])
					}
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestPostgresWeatherAlertMetadataRepository_Search_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	// スキャンエラーを引き起こすために不正な型を返す
	rows := sqlmock.NewRows([]string{"id", "region", "severity", "issued_at", "created_at"}).
		AddRow("alert1", "Tokyo", "warning", "invalid-date", time.Now())
	mock.ExpectQuery("SELECT id, region, severity, issued_at, created_at FROM weather_alert_metadata ORDER BY issued_at DESC").
		WillReturnRows(rows)

	repo := NewPostgresWeatherAlertMetadataRepository(db)
	_, err = repo.Search(context.Background(), repository.MetadataFilter{})

	if err == nil {
		t.Error("Search() expected scan error, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestPostgresWeatherAlertMetadataRepository_Search_RowsError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "region", "severity", "issued_at", "created_at"}).
		AddRow("alert1", "Tokyo", "warning", time.Now(), time.Now()).
		RowError(0, sql.ErrConnDone)
	mock.ExpectQuery("SELECT id, region, severity, issued_at, created_at FROM weather_alert_metadata ORDER BY issued_at DESC").
		WillReturnRows(rows)

	repo := NewPostgresWeatherAlertMetadataRepository(db)
	_, err = repo.Search(context.Background(), repository.MetadataFilter{})

	if err == nil {
		t.Error("Search() expected row error, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}