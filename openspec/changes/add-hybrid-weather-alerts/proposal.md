# Proposal: add-hybrid-weather-alerts

## Summary
PostgreSQLでメタデータ検索を行い、Firestoreから詳細データを取得するハイブリッドデータストア構成を実装します。気象アラートドメインを題材に、スキーマ変更に強い柔軟なアーキテクチャパターンを示します。

## What Changes
この変更により、以下の新機能が追加されます:

### 追加される機能
1. **PostgreSQL メタデータレイヤー**
   - `weather_alert_metadata` テーブル（id, region, severity, issued_at）
   - 地域・日時による効率的な検索用インデックス
   - `WeatherAlertMetadataRepository` インターフェースと PostgreSQL 実装

2. **Firestore 本体データレイヤー**
   - `weatherAlerts` コレクション（スキーマフリーな詳細データ）
   - `WeatherAlertRepository` インターフェースと Firestore 実装
   - バッチ取得メソッド (`GetByIDs`)

3. **GraphQL ハイブリッドクエリ**
   - `WeatherAlert` GraphQL タイプ（PostgreSQL + Firestore のマージデータ）
   - `weatherAlerts(region: String, issuedAfter: String)` クエリ
   - 2段階データフロー: PostgreSQL 検索 → Firestore 詳細取得

4. **サンプルデータとテスト**
   - シードスクリプト (`scripts/seed-weather-alerts.go`)
   - 10件の気象アラートサンプルデータ（東京、大阪、京都）

### データフロー
```
GraphQL Query (region: "Tokyo", issuedAfter: "2025-12-19T00:00:00Z")
    ↓
PostgreSQL: SELECT id WHERE region='Tokyo' AND issued_at >= '2025-12-19'
    → ["alert-001", "alert-002"]
    ↓
Firestore: BatchGet weatherAlerts/alert-001, weatherAlerts/alert-002
    → [{title: "...", rawData: {...}}, ...]
    ↓
GraphQL Response: Merge metadata + details
    → [{ id, region, severity, title, rawData, ... }]
```

## Why
この変更は、以下の課題を解決します:

### 課題背景（別サービスでの実運用経験より）
- 気象データの元スキーマが頻繁に変更される
- RDB スキーマ定義が厳密なため、スキーマ変更の影響範囲が大きい
- スキーマ変更のたびにマイグレーション・デプロイが必要

### 解決策
- **スキーマ変更に強いデータ**: Firestore の `rawData` フィールドに JSON として柔軟に格納
- **検索・集計が必要なデータ**: PostgreSQL にメタデータとして正規化して格納
- GraphQL で両方のデータストアを透過的に扱い、クライアントはシンプルな API で利用可能

### メリット
- ✅ 気象庁などの外部データソースのスキーマ変更に強い
- ✅ PostgreSQL で効率的な検索・集計（地域・日時インデックス）
- ✅ Firestore でスキーマレスな詳細データ保存
- ✅ 既存の User (PostgreSQL) と Message (Firestore) パターンを踏襲した一貫性

## How
実装は5つのフェーズに分かれます:

### Phase 1: PostgreSQL Metadata Layer
- Migration スクリプトで `weather_alert_metadata` テーブル作成
- `WeatherAlertMetadata` ドメインモデル定義
- `PostgresWeatherAlertMetadataRepository` 実装（動的 WHERE 句）

### Phase 2: Firestore Data Layer
- `WeatherAlert` ドメインモデル定義（`rawData` は `map[string]interface{}`）
- `FirestoreWeatherAlertRepository` 実装（バッチ取得ロジック）

### Phase 3: GraphQL Integration
- `WeatherAlert` GraphQL タイプ定義
- `weatherAlerts` クエリ実装（2段階データ取得 + マージ）
- `server.go` でリポジトリ依存性注入

### Phase 4: Data Seeding and Testing
- PostgreSQL マイグレーション実行
- シードスクリプトで 10 件のサンプルデータ投入
- GraphQL Playground でクエリテスト

### Phase 5: Code Quality
- gofmt/goimports 実行
- README 更新
- テスト実行

詳細は `tasks.md` を参照してください。

## Impact
### 既存機能への影響
- **既存コードへの影響**: なし（新規ドメイン追加のみ）
- **既存 GraphQL スキーマ**: `hello`, `messages`, `users` クエリは変更なし
- **既存データベース**: 新規テーブル追加のみ（既存テーブルは変更なし）

### 新規依存関係
- なし（既存の PostgreSQL, Firestore クライアントを再利用）

### パフォーマンス考慮事項
- **N+1 問題**: Firestore の `GetByIDs` でバッチ取得を実装（最大 500 件）
- **データ整合性**: PostgreSQL にメタデータがあるが Firestore にデータがない場合は警告ログ + スキップ

## Risks and Mitigations
### リスク 1: 2つのデータストア間の同期
- **リスク**: PostgreSQL と Firestore が別々のデータストアのため、トランザクション保証がない
- **軽減策**: シードスクリプトで PostgreSQL → Firestore を順次実行。エラー時は手動ロールバック（初期実装ではシンプルさ優先）

### リスク 2: Firestore データ欠損
- **リスク**: PostgreSQL にメタデータがあるが Firestore にデータがない場合
- **軽減策**: Resolver で Firestore 取得失敗時は警告ログを記録し、該当データをレスポンスから除外（部分的成功を許容）

### リスク 3: パフォーマンス（大量データ時）
- **リスク**: 検索結果が数千件になる場合、Firestore バッチ取得が遅延する可能性
- **軽減策**: 初期実装ではページネーションなし（シンプルさ優先）。将来的に必要に応じてカーソルベースページネーションを追加

## Alternatives Considered
### 代替案 1: 全てを PostgreSQL に格納
- **理由**: スキーマ変更に弱い。気象データの柔軟性が失われる
- **却下理由**: 元データスキーマ変更のたびにマイグレーションが必要

### 代替案 2: 全てを Firestore に格納
- **理由**: 複雑な検索クエリ（範囲検索・複合インデックス）のパフォーマンスが悪い
- **却下理由**: 地域・日時による効率的な検索が困難

### 採用案: PostgreSQL (メタデータ) + Firestore (本体データ)
- **理由**: 検索効率とスキーマ柔軟性の両立を実現

## Open Questions
なし（ユーザー確認済み: 気象アラートドメイン、地域フィルタ + 日時範囲検索）

## References
- 既存実装: `openspec/changes/archive/2025-12-19-add-postgresql-user-integration`
- 既存実装: `openspec/changes/archive/2025-12-19-add-firestore-docker-integration`
- プロジェクトコンテキスト: `openspec/project.md`
- 設計詳細: `design.md`
- 実装タスク: `tasks.md`
- 仕様デルタ:
  - `specs/weather-alert-metadata-postgresql/spec.md`
  - `specs/weather-alert-firestore/spec.md`
  - `specs/weather-alert-graphql-integration/spec.md`