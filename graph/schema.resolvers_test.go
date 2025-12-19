package graph

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jxpress/graphql-sampleapp/graph/model"
	"github.com/jxpress/graphql-sampleapp/internal/domain"
	"github.com/jxpress/graphql-sampleapp/internal/repository"
	"github.com/stretchr/testify/assert"
)

// --- Mock Implementations ---

type mockUserRepository struct {
	users   []*domain.User
	user    *domain.User
	err     error
	listErr error
	getErr  error
}

func (m *mockUserRepository) List(ctx context.Context) ([]*domain.User, error) {
	if m.listErr != nil {
		return nil, m.listErr
	}
	if m.err != nil {
		return nil, m.err
	}
	return m.users, nil
}

func (m *mockUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	if m.err != nil {
		return nil, m.err
	}
	if m.user != nil {
		if m.user.ID == id {
			return m.user, nil
		}
	}
	for _, u := range m.users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, errors.New("user not found")
}

type mockMessageRepository struct {
	messages []*domain.Message
	message  *domain.Message
	err      error
}

func (m *mockMessageRepository) List(ctx context.Context) ([]*domain.Message, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.messages, nil
}

func (m *mockMessageRepository) GetByID(ctx context.Context, id string) (*domain.Message, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.message != nil && m.message.ID == id {
		return m.message, nil
	}
	for _, msg := range m.messages {
		if msg.ID == id {
			return msg, nil
		}
	}
	return nil, errors.New("message not found")
}

type mockWeatherAlertMetadataRepository struct {
	metadata  []*domain.WeatherAlertMetadata
	searchIDs []string
	err       error
}

func (m *mockWeatherAlertMetadataRepository) Search(ctx context.Context, filter repository.MetadataFilter) ([]*domain.WeatherAlertMetadata, error) {
	if m.err != nil {
		return nil, m.err
	}
	// Simple filter implementation strictly checking if mock setup matches expectation or just return all for simplicity
	// Checks if filtered region matches
	var result []*domain.WeatherAlertMetadata
	for _, meta := range m.metadata {
		if filter.Region != nil && meta.Region != *filter.Region {
			continue
		}
		if filter.IssuedAfter != nil && !meta.IssuedAt.After(*filter.IssuedAfter) {
			continue
		}
		result = append(result, meta)
	}
	return result, nil
}

func (m *mockWeatherAlertMetadataRepository) SearchIDs(ctx context.Context, filter repository.MetadataFilter) ([]string, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.searchIDs, nil
}

type mockWeatherAlertRepository struct {
	alerts []*domain.WeatherAlert
	alert  *domain.WeatherAlert
	err    error
}

func (m *mockWeatherAlertRepository) GetByIDs(ctx context.Context, ids []string) ([]*domain.WeatherAlert, error) {
	if m.err != nil {
		return nil, m.err
	}
	var result []*domain.WeatherAlert
	for _, id := range ids {
		for _, a := range m.alerts {
			if a.ID == id {
				result = append(result, a)
				break
			}
		}
	}
	return result, nil
}

func (m *mockWeatherAlertRepository) GetByID(ctx context.Context, id string) (*domain.WeatherAlert, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.alert != nil && m.alert.ID == id {
		return m.alert, nil
	}
	return nil, errors.New("not found")
}

// --- Tests ---

func TestQueryResolver_Hello(t *testing.T) {
	r := &queryResolver{}
	got, err := r.Hello(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, "Hello World", got)
}

func TestQueryResolver_Users(t *testing.T) {
	fixedTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		mock    *mockUserRepository
		want    []*model.User
		wantErr bool
	}{
		{
			name: "正常系: ユーザーリスト取得",
			mock: &mockUserRepository{
				users: []*domain.User{
					{ID: "1", Name: "User 1", Email: "u1@example.com", CreatedAt: fixedTime},
				},
			},
			want: []*model.User{
				{ID: "1", Name: "User 1", Email: "u1@example.com", CreatedAt: fixedTime.Format(time.RFC3339)},
			},
			wantErr: false,
		},
		{
			name: "異常系: Repositoryエラー",
			mock: &mockUserRepository{
				err: errors.New("db error"),
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolver := NewResolver(nil, tt.mock, nil, nil)
			q := resolver.Query()
			got, err := q.Users(context.Background())

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.want), len(got))
				if len(got) > 0 {
					assert.Equal(t, tt.want[0].ID, got[0].ID)
				}
			}
		})
	}
}

func TestQueryResolver_User(t *testing.T) {
	fixedTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		id      string
		mock    *mockUserRepository
		want    *model.User
		wantErr bool
	}{
		{
			name: "正常系: ユーザー取得",
			id:   "1",
			mock: &mockUserRepository{
				users: []*domain.User{
					{ID: "1", Name: "User 1", Email: "u1@example.com", CreatedAt: fixedTime},
				},
			},
			want: &model.User{
				ID: "1", Name: "User 1", Email: "u1@example.com", CreatedAt: fixedTime.Format(time.RFC3339),
			},
			wantErr: false,
		},
		{
			name: "異常系: ユーザーが見つからない",
			id:   "99",
			mock: &mockUserRepository{
				users: []*domain.User{},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolver := NewResolver(nil, tt.mock, nil, nil)
			q := resolver.Query()
			got, err := q.User(context.Background(), tt.id)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.ID, got.ID)
			}
		})
	}
}

func TestQueryResolver_Messages(t *testing.T) {
	fixedTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		mock    *mockMessageRepository
		want    []*model.Message
		wantErr bool
	}{
		{
			name: "正常系: メッセージリスト取得",
			mock: &mockMessageRepository{
				messages: []*domain.Message{
					{ID: "1", Content: "Hello", Author: "User1", CreatedAt: fixedTime},
				},
			},
			want: []*model.Message{
				{ID: "1", Content: "Hello", Author: "User1", CreatedAt: fixedTime.Format(time.RFC3339)},
			},
			wantErr: false,
		},
		{
			name: "異常系: Repositoryエラー",
			mock: &mockMessageRepository{
				err: errors.New("db error"),
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolver := NewResolver(tt.mock, nil, nil, nil)
			q := resolver.Query()
			got, err := q.Messages(context.Background())

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.want), len(got))
			}
		})
	}
}

func TestQueryResolver_Message(t *testing.T) {
	fixedTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		id      string
		mock    *mockMessageRepository
		want    *model.Message
		wantErr bool
	}{
		{
			name: "正常系: メッセージ取得",
			id:   "1",
			mock: &mockMessageRepository{
				messages: []*domain.Message{
					{ID: "1", Content: "Hello", Author: "User1", CreatedAt: fixedTime},
				},
			},
			want: &model.Message{
				ID: "1", Content: "Hello", Author: "User1", CreatedAt: fixedTime.Format(time.RFC3339),
			},
			wantErr: false,
		},
		{
			name: "異常系: メッセージが見つからない",
			id:   "99",
			mock: &mockMessageRepository{
				messages: []*domain.Message{},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolver := NewResolver(tt.mock, nil, nil, nil)
			q := resolver.Query()
			got, err := q.Message(context.Background(), tt.id)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.ID, got.ID)
			}
		})
	}
}

func TestQueryResolver_WeatherAlerts(t *testing.T) {
	fixedTime := time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC)
	region := "Tokyo"

	meta := &domain.WeatherAlertMetadata{
		ID:        "alert1",
		Region:    region,
		Severity:  "Warning",
		IssuedAt:  fixedTime,
		CreatedAt: fixedTime,
	}
	alert := &domain.WeatherAlert{
		ID:          "alert1",
		Title:       "Typhoon",
		Description: "Big typhoon coming",
		RawData:     map[string]interface{}{"pressure": 900},
	}

	tests := []struct {
		name        string
		region      *string
		issuedAfter *string
		mockMeta    *mockWeatherAlertMetadataRepository
		mockAlert   *mockWeatherAlertRepository
		wantLen     int
		wantErr     bool
	}{
		{
			name:   "正常系: フィルタなし",
			region: nil,
			mockMeta: &mockWeatherAlertMetadataRepository{
				metadata: []*domain.WeatherAlertMetadata{meta},
			},
			mockAlert: &mockWeatherAlertRepository{
				alerts: []*domain.WeatherAlert{alert},
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name:   "正常系: 地域フィルタあり（ヒット）",
			region: &region,
			mockMeta: &mockWeatherAlertMetadataRepository{
				metadata: []*domain.WeatherAlertMetadata{meta},
			},
			mockAlert: &mockWeatherAlertRepository{
				alerts: []*domain.WeatherAlert{alert},
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name:   "正常系: 地域フィルタあり（ヒットしない）",
			region: func() *string { s := "Osaka"; return &s }(),
			mockMeta: &mockWeatherAlertMetadataRepository{
				metadata: []*domain.WeatherAlertMetadata{meta},
			},
			mockAlert: &mockWeatherAlertRepository{
				alerts: []*domain.WeatherAlert{alert},
			},
			wantLen: 0, // Mock search implementation filters this
			wantErr: false,
		},
		{
			name:        "正常系: 日時フィルタエラー",
			issuedAfter: func() *string { s := "invalid-date"; return &s }(),
			mockMeta:    &mockWeatherAlertMetadataRepository{},
			mockAlert:   &mockWeatherAlertRepository{},
			wantLen:     0,
			wantErr:     true,
		},
		{
			name: "正常系: WeatherAlert詳細が見つからない（metadataのみ）",
			mockMeta: &mockWeatherAlertMetadataRepository{
				metadata: []*domain.WeatherAlertMetadata{meta},
			},
			mockAlert: &mockWeatherAlertRepository{
				alerts: []*domain.WeatherAlert{}, // Empty alerts
			},
			wantLen: 0,
			wantErr: false, // implementation skips missing alerts
		},
		{
			name: "異常系: WeatherAlertMetadataRepositoryエラー",
			mockMeta: &mockWeatherAlertMetadataRepository{
				err: errors.New("meta DB error"),
			},
			mockAlert: &mockWeatherAlertRepository{},
			wantLen:   0,
			wantErr:   true,
		},
		{
			name: "異常系: WeatherAlertRepositoryエラー",
			mockMeta: &mockWeatherAlertMetadataRepository{
				metadata: []*domain.WeatherAlertMetadata{meta},
			},
			mockAlert: &mockWeatherAlertRepository{
				err: errors.New("alert DB error"),
			},
			wantLen: 0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolver := NewResolver(nil, nil, tt.mockMeta, tt.mockAlert)
			q := resolver.Query()
			got, err := q.WeatherAlerts(context.Background(), tt.region, tt.issuedAfter)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantLen, len(got))
				if tt.wantLen > 0 {
					assert.Equal(t, "alert1", got[0].ID)
					assert.Equal(t, "Typhoon", got[0].Title)
					assert.Contains(t, got[0].RawData, "pressure")
				}
			}
		})
	}
}
