package postgres

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kuchida1981/graphql-sampleapp/internal/domain"
)

func TestPostgresUserRepository_List(t *testing.T) {
	tests := []struct {
		name    string
		mockFn  func(mock sqlmock.Sqlmock)
		want    []*domain.User
		wantErr bool
	}{
		{
			name: "正常系: ユーザーリスト取得成功",
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "created_at"}).
					AddRow("user1", "Alice", "alice@example.com", time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)).
					AddRow("user2", "Bob", "bob@example.com", time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC))
				mock.ExpectQuery("SELECT id, name, email, created_at FROM users ORDER BY created_at DESC").
					WillReturnRows(rows)
			},
			want: []*domain.User{
				{ID: "user1", Name: "Alice", Email: "alice@example.com", CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
				{ID: "user2", Name: "Bob", Email: "bob@example.com", CreatedAt: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)},
			},
			wantErr: false,
		},
		{
			name: "正常系: ユーザーが0件",
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "created_at"})
				mock.ExpectQuery("SELECT id, name, email, created_at FROM users ORDER BY created_at DESC").
					WillReturnRows(rows)
			},
			want:    []*domain.User{},
			wantErr: false,
		},
		{
			name: "異常系: クエリエラー",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, name, email, created_at FROM users ORDER BY created_at DESC").
					WillReturnError(errors.New("database connection error"))
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "異常系: スキャンエラー",
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "created_at"}).
					AddRow("user1", "Alice", "alice@example.com", "invalid-date")
				mock.ExpectQuery("SELECT id, name, email, created_at FROM users ORDER BY created_at DESC").
					WillReturnRows(rows)
			},
			want:    nil,
			wantErr: true,
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

			repo := NewPostgresUserRepository(db)
			got, err := repo.List(context.Background())

			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(got) != len(tt.want) {
					t.Errorf("List() got %d users, want %d users", len(got), len(tt.want))
					return
				}

				for i, user := range got {
					if user.ID != tt.want[i].ID || user.Name != tt.want[i].Name ||
						user.Email != tt.want[i].Email || !user.CreatedAt.Equal(tt.want[i].CreatedAt) {
						t.Errorf("List() got user[%d] = %+v, want %+v", i, user, tt.want[i])
					}
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestPostgresUserRepository_GetByID(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		mockFn  func(mock sqlmock.Sqlmock)
		want    *domain.User
		wantErr bool
	}{
		{
			name: "正常系: ユーザー取得成功",
			id:   "user1",
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "created_at"}).
					AddRow("user1", "Alice", "alice@example.com", time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
				mock.ExpectQuery("SELECT id, name, email, created_at FROM users WHERE id = \\$1").
					WithArgs("user1").
					WillReturnRows(rows)
			},
			want: &domain.User{
				ID:        "user1",
				Name:      "Alice",
				Email:     "alice@example.com",
				CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			wantErr: false,
		},
		{
			name: "異常系: ユーザーが見つからない (sql.ErrNoRows)",
			id:   "nonexistent",
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, name, email, created_at FROM users WHERE id = \\$1").
					WithArgs("nonexistent").
					WillReturnError(sql.ErrNoRows)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "異常系: スキャンエラー",
			id:   "user1",
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "created_at"}).
					AddRow("user1", "Alice", "alice@example.com", "invalid-date")
				mock.ExpectQuery("SELECT id, name, email, created_at FROM users WHERE id = \\$1").
					WithArgs("user1").
					WillReturnRows(rows)
			},
			want:    nil,
			wantErr: true,
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

			repo := NewPostgresUserRepository(db)
			got, err := repo.GetByID(context.Background(), tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got != nil && tt.want != nil {
				if got.ID != tt.want.ID || got.Name != tt.want.Name ||
					got.Email != tt.want.Email || !got.CreatedAt.Equal(tt.want.CreatedAt) {
					t.Errorf("GetByID() = %+v, want %+v", got, tt.want)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unfulfilled expectations: %v", err)
			}
		})
	}
}