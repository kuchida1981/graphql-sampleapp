# テストガイドライン

このドキュメントは、`graphql-sampleapp` プロジェクトにおけるテスト作成の基本方針とベストプラクティスをまとめたものです。

## 目次

1. [テスト作成の基本方針](#テスト作成の基本方針)
2. [モッキング戦略](#モッキング戦略)
3. [テスト命名規則](#テスト命名規則)
4. [テーブル駆動テスト](#テーブル駆動テスト)
5. [カバレッジ方針](#カバレッジ方針)
6. [テスト実行方法](#テスト実行方法)
7. [CI/CD統合](#cicd統合)
8. [必要な依存ライブラリ](#必要な依存ライブラリ)
9. [テストコード例](#テストコード例)

## テスト作成の基本方針

### テストの種類

このプロジェクトでは、以下の種類のテストを実施しています:

1. **ユニットテスト**: 各関数やメソッドの動作を個別にテストします
   - Repository層のテスト
   - ビジネスロジックのテスト
   - GraphQL Resolverのテスト

2. **統合テスト** (今後実装予定):
   - 実際のDBを使用したテスト
   - E2Eテスト

### テストの原則

- **Fast (高速)**: テストは素早く実行できるべきです。DBへの依存を排除してモックを使用します
- **Independent (独立)**: 各テストは他のテストに依存せず、どの順序で実行しても同じ結果になるべきです
- **Repeatable (再現可能)**: 同じ入力に対して常に同じ結果を返すべきです
- **Self-Validating (自己検証)**: テストは成功/失敗を明確に示すべきです
- **Timely (適時)**: テストは実装と同時、または実装前に書くべきです

## モッキング戦略

### PostgreSQLのモッキング

PostgreSQL（`database/sql`）を使用するRepositoryには、[go-sqlmock](https://github.com/DATA-DOG/go-sqlmock) を使用します。

**利点:**
- 実際のDB接続不要
- テストが高速
- 様々なエラーケースを簡単にシミュレート可能

**基本的な使い方:**

```go
import (
    "github.com/DATA-DOG/go-sqlmock"
)

func TestRepository(t *testing.T) {
    // モックDBを作成
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("failed to create mock: %v", err)
    }
    defer db.Close()

    // 期待するクエリとレスポンスを設定
    rows := sqlmock.NewRows([]string{"id", "name"}).
        AddRow("1", "Alice")
    mock.ExpectQuery("SELECT id, name FROM users").
        WillReturnRows(rows)

    // テスト対象を実行
    repo := NewUserRepository(db)
    users, err := repo.List(context.Background())

    // アサーション
    if err != nil {
        t.Errorf("unexpected error: %v", err)
    }
    // ...

    // すべての期待が満たされたか確認
    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("unfulfilled expectations: %v", err)
    }
}
```

### Firestoreのモッキング

Firestore SDKは直接的なモッキングが困難なため、以下のアプローチを推奨します:

1. **モックリポジトリパターン**: 簡易的なモック実装を作成
2. **Firestore Emulator**: 実際のエミュレータを使った統合テスト（推奨）
3. **testcontainers**: Dockerコンテナを使ったテスト

現在のプロジェクトでは、モックリポジトリパターンを採用しています。

**例:**

```go
type mockFirestoreMessageRepository struct {
    messages map[string]*domain.Message
}

func (m *mockFirestoreMessageRepository) GetByID(ctx context.Context, id string) (*domain.Message, error) {
    if msg, ok := m.messages[id]; ok {
        return msg, nil
    }
    return nil, fmt.Errorf("message not found")
}
```

## テスト命名規則

### テスト関数名

テスト関数名は以下の形式で命名します:

```
Test<対象の構造体>_<対象のメソッド名>_<テストケース>
```

**例:**
- `TestPostgresUserRepository_List_Success` - 正常系
- `TestPostgresUserRepository_GetByID_NotFound` - エラーケース

日本語のテストケース名も使用可能です（テーブル駆動テストの`name`フィールド）:

```go
tests := []struct {
    name string
    // ...
}{
    {
        name: "正常系: ユーザーリスト取得成功",
        // ...
    },
    {
        name: "異常系: クエリエラー",
        // ...
    },
}
```

### ファイル名

テストファイルは `_test.go` サフィックスを付けます:

- `postgres_user.go` → `postgres_user_test.go`
- `firestore_message.go` → `firestore_message_test.go`

## テーブル駆動テスト

Goでは、複数のテストケースを効率的に実行するために **テーブル駆動テスト** が推奨されます。

### 基本構造

```go
func TestSomething(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {
            name:    "正常系: 基本的な入力",
            input:   "hello",
            want:    "HELLO",
            wantErr: false,
        },
        {
            name:    "異常系: 空文字列",
            input:   "",
            want:    "",
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := ToUpper(tt.input)

            if (err != nil) != tt.wantErr {
                t.Errorf("ToUpper() error = %v, wantErr %v", err, tt.wantErr)
                return
            }

            if got != tt.want {
                t.Errorf("ToUpper() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### メリット

- **可読性**: テストケースが一目で分かる
- **保守性**: 新しいケースを簡単に追加できる
- **網羅性**: 正常系・異常系を体系的にカバーできる

## カバレッジ方針

### 目標カバレッジ

- **Repository層**: 80%以上
- **GraphQL Resolver層**: 70%以上
- **ビジネスロジック**: 90%以上

### カバレッジ測定方法

```bash
# カバレッジを計算
go test ./... -cover

# 詳細なカバレッジレポートを生成
go test ./... -coverprofile=coverage.out

# HTMLでカバレッジレポートを表示
go tool cover -html=coverage.out
```

### カバレッジの重視ポイント

カバレッジの数字にこだわりすぎず、以下を重視します:

1. **重要なビジネスロジック**: 必ず100%カバー
2. **エラーハンドリング**: 主要なエラーパスをカバー
3. **境界値テスト**: 空配列、nil、境界値などをカバー

## テスト実行方法

### すべてのテストを実行

```bash
go test ./...
```

### 特定のパッケージのテストを実行

```bash
go test ./internal/repository
go test ./graph
```

### 特定のテスト関数を実行

```bash
go test ./internal/repository -run TestPostgresUserRepository
go test ./internal/repository -run TestPostgresUserRepository_List
```

### Verboseモードで実行

```bash
go test ./... -v
```

### カバレッジ付きで実行

```bash
go test ./... -cover
go test ./... -coverprofile=coverage.out
```

## CI/CD統合

### GitHub Actionsの設定例

`.github/workflows/test.yml`:

```yaml
name: Test

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest

    steps:
    - uses: actions/checkout@v3

    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.24'

    - name: Download dependencies
      run: go mod download

    - name: Run tests
      run: go test ./... -v

    - name: Run tests with coverage
      run: go test ./... -coverprofile=coverage.out

    - name: Upload coverage to Codecov (optional)
      uses: codecov/codecov-action@v3
      with:
        files: ./coverage.out
```

### テスト失敗時のCI動作

- PRマージ前に必ずテストをパス
- テスト失敗時はマージをブロック
- カバレッジが閾値を下回った場合は警告

## 必要な依存ライブラリ

### テスト用ライブラリ

```bash
# go-sqlmock: SQLのモッキング
go get github.com/DATA-DOG/go-sqlmock

# (オプション) testify: アサーションライブラリ
go get github.com/stretchr/testify
```

### インストール

```bash
go mod download
go mod tidy
```

## テストコード例

### PostgreSQL Repositoryのテスト例

```go
package repository

import (
    "context"
    "testing"
    "github.com/DATA-DOG/go-sqlmock"
)

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
                rows := sqlmock.NewRows([]string{"id", "name", "email"}).
                    AddRow("user1", "Alice", "alice@example.com")
                mock.ExpectQuery("SELECT id, name, email FROM users WHERE id = \\$1").
                    WithArgs("user1").
                    WillReturnRows(rows)
            },
            want: &domain.User{
                ID:    "user1",
                Name:  "Alice",
                Email: "alice@example.com",
            },
            wantErr: false,
        },
        {
            name: "異常系: ユーザーが見つからない",
            id:   "nonexistent",
            mockFn: func(mock sqlmock.Sqlmock) {
                mock.ExpectQuery("SELECT id, name, email FROM users WHERE id = \\$1").
                    WithArgs("nonexistent").
                    WillReturnError(sql.ErrNoRows)
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
                if got.ID != tt.want.ID || got.Name != tt.want.Name {
                    t.Errorf("GetByID() = %+v, want %+v", got, tt.want)
                }
            }

            if err := mock.ExpectationsWereMet(); err != nil {
                t.Errorf("unfulfilled expectations: %v", err)
            }
        })
    }
}
```

### Firestore Repositoryのモックテスト例

```go
package repository

import (
    "context"
    "testing"
)

type mockFirestoreMessageRepository struct {
    messages map[string]*domain.Message
}

func TestMockFirestoreMessageRepository_GetByID(t *testing.T) {
    tests := []struct {
        name     string
        messages map[string]*domain.Message
        id       string
        wantErr  bool
    }{
        {
            name: "正常系: メッセージ取得成功",
            messages: map[string]*domain.Message{
                "msg1": {ID: "msg1", Content: "Hello"},
            },
            id:      "msg1",
            wantErr: false,
        },
        {
            name:     "異常系: メッセージが見つからない",
            messages: map[string]*domain.Message{},
            id:       "nonexistent",
            wantErr:  true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            repo := &mockFirestoreMessageRepository{
                messages: tt.messages,
            }

            got, err := repo.GetByID(context.Background(), tt.id)

            if (err != nil) != tt.wantErr {
                t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
                return
            }

            if !tt.wantErr && got == nil {
                t.Error("GetByID() returned nil message")
            }
        })
    }
}
```

## ベストプラクティス

### 1. テストは読みやすく書く

```go
// Good
if got != want {
    t.Errorf("GetByID() = %v, want %v", got, want)
}

// Bad
if got != want {
    t.Error("failed")
}
```

### 2. エラーメッセージは具体的に

```go
// Good
t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)

// Bad
t.Error("error occurred")
```

### 3. テストは独立させる

各テストは他のテストに依存せず、単独で実行できるようにします。

### 4. モックの期待は明確に

```go
mock.ExpectQuery("SELECT id, name FROM users WHERE id = \\$1").
    WithArgs("user1").
    WillReturnRows(rows)
```

### 5. 境界値をテストする

- 空のスライス
- nil値
- 0値
- 最大値/最小値

## まとめ

このドキュメントに従うことで、以下が実現できます:

- **一貫性のあるテストコード**: チーム全体で統一されたスタイル
- **高速なテスト実行**: モックを使用したDBへの依存排除
- **高いカバレッジ**: 重要なロジックの網羅的なテスト
- **CI/CD統合**: 自動化されたテスト実行

テストは継続的に改善していくものです。新しいベストプラクティスやパターンが見つかった場合は、このドキュメントを更新してください。