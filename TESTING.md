# テストガイドライン

このプロジェクトのテスト方針と手順について記述します。

## テスト作成の基本方針
- **ユニットテスト優先**: `internal/repository` および `graph` パッケージに対し、外部依存（DB等）をモックしたユニットテストを作成します。
- **カバレッジ**: 具体的な数値を目標にはしませんが、主要な正常系および異常系パスを網羅することを目指します。
- **命名規則**: `Test<関数名>_<シナリオ>` の形式（例: `TestPostgresUserRepository_List_Success`）を推奨します。

## モッキング戦略
- Repository層のテストには [go-sqlmock](https://github.com/DATA-DOG/go-sqlmock) を使用してデータベース接続をモックします。
- Resolver層のテストには、Repositoryインターフェースの手動モック実装（または `gomock`）を使用して、データ層への依存を切り離します。

## テスト実行方法

### すべてのテストを実行
```bash
go test ./...
```

### カバレッジを確認
```bash
go test -cover ./...
```

### 特定のパッケージのみ実行
```bash
go test -v ./graph/...
```

## CI/CD
GitHub Actionsにより、Pull Request作成時に自動的にテストが実行されます。
詳細は `.github/workflows/test.yml` を参照してください。