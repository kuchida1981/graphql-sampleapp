# Design: Testing Framework and Policy

## Architecture Overview
本プロジェクトのテスト戦略は、以下の3つの層で構成されます:

```
┌─────────────────────────────────────┐
│   GraphQL Resolver Tests            │ ← モックRepositoryを使用
├─────────────────────────────────────┤
│   Repository Unit Tests             │ ← モックDB接続を使用
├─────────────────────────────────────┤
│   Domain Model Tests (必要に応じて) │ ← ビジネスロジックがある場合
└─────────────────────────────────────┘
```

## Testing Philosophy

### 1. テスト作成の基本方針
- **すべての公開メソッドにテストを書く**: Repositoryインターフェース、Resolverなど外部から呼び出される関数は必ずテスト対象
- **DBは常にモック**: ユニットテストではDBへの実際の接続は行わない
- **テーブル駆動テスト**: 複数のケースは`[]struct`でまとめて記述
- **明示的なエラーケース**: 正常系だけでなく、異常系も必ずテスト
- **テスト名の明確化**: `Test<FunctionName>_<Scenario>`形式で命名

### 2. モッキング戦略

#### Database Mocking
- **PostgreSQL**: `github.com/DATA-DOG/go-sqlmock`を使用
  - `sql.DB`の動作をモック
  - クエリの期待値とレスポンスを定義

- **Firestore**: カスタムモックまたは`gomock`で`*firestore.Client`をモック
  - インターフェース化が必要な場合はRepositoryインターフェースでラップ

#### Repository Mocking
- Resolver層のテストでは、Repositoryインターフェースのモックを作成
- `gomock`または手動モック実装を使用

### 3. テストファイル構成
```
internal/repository/
├── postgres_user.go
├── postgres_user_test.go          ← Repositoryのユニットテスト
├── firestore_message.go
└── firestore_message_test.go

graph/
├── schema.resolvers.go
└── schema.resolvers_test.go       ← Resolverのユニットテスト
```

### 4. テストカバレッジ方針
- **カバレッジ目標**: 具体的なパーセンテージは設定しない
- **意味のあるカバレッジ**: 重要なビジネスロジックとエラーハンドリングを優先
- **測定**: `go test -cover ./...`で確認可能にする
- **CI/CD**: すべてのPRでテストが自動実行される

## Component-Specific Design

### Repository Tests
各Repositoryのテストは以下の構造で実装:

```go
func TestPostgresUserRepository_List(t *testing.T) {
    tests := []struct {
        name    string
        mockDB  func() (*sql.DB, sqlmock.Sqlmock)
        want    []*domain.User
        wantErr bool
    }{
        {
            name: "正常系: ユーザーリストを取得",
            mockDB: func() (*sql.DB, sqlmock.Sqlmock) {
                db, mock, _ := sqlmock.New()
                rows := sqlmock.NewRows(...)
                mock.ExpectQuery("SELECT id, name...").WillReturnRows(rows)
                return db, mock
            },
            want:    []*domain.User{...},
            wantErr: false,
        },
        {
            name: "異常系: クエリ失敗",
            mockDB: func() (*sql.DB, sqlmock.Sqlmock) {
                db, mock, _ := sqlmock.New()
                mock.ExpectQuery("SELECT id, name...").WillReturnError(errors.New("db error"))
                return db, mock
            },
            want:    nil,
            wantErr: true,
        },
    }
    // テスト実行ロジック
}
```

### Resolver Tests
Resolverのテストは以下の構造で実装:

```go
type mockUserRepository struct {
    users []*domain.User
    err   error
}

func (m *mockUserRepository) List(ctx context.Context) ([]*domain.User, error) {
    return m.users, m.err
}

func TestQueryResolver_Users(t *testing.T) {
    tests := []struct {
        name     string
        mockRepo UserRepository
        want     []*model.User
        wantErr  bool
    }{
        {
            name: "正常系: ユーザーリストを取得",
            mockRepo: &mockUserRepository{
                users: []*domain.User{...},
                err:   nil,
            },
            want:    []*model.User{...},
            wantErr: false,
        },
        // ...
    }
    // テスト実行ロジック
}
```

## Dependencies

### Testing Libraries
- `testing` (標準ライブラリ)
- `github.com/DATA-DOG/go-sqlmock` - PostgreSQLモック
- `github.com/golang/mock/gomock` (オプション) - インターフェースモック生成

### Why These Libraries?
- **sqlmock**: デファクトスタンダードのSQLモックライブラリ、軽量で高機能
- **gomock**: Googleが提供する公式モックライブラリ、型安全
- **標準testing**: 外部依存を最小限にし、Goのイディオムに従う

## CI/CD Integration
- `.github/workflows/test.yml` (または同等) を追加
- すべてのPRで`go test ./...`を実行
- テスト失敗時はマージをブロック

## Documentation
テストポリシーは`TESTING.md`として文書化し、以下を含む:
- テスト作成のガイドライン
- モッキングの方法
- テーブル駆動テストの書き方
- CI/CDでのテスト実行方法

## Trade-offs

### 採用したアプローチ
- **ユニットテスト優先**: 高速で依存が少ない
- **モックDB**: セットアップ不要、高速実行

### トレードオフ
- **実際のDB挙動との差異**: モックでは再現できないDB固有の挙動がある
  - 対策: 統合テストは将来的に追加可能 (このプロポーザルの範囲外)
- **モックのメンテナンス**: 実装変更時にモックも更新が必要
  - 対策: インターフェースベースの設計で影響を最小化

## Future Considerations
- 統合テスト (実際のDB使用) の追加
- E2Eテスト
- パフォーマンステスト
- mutation/subscription resolverのテスト (現在未実装)