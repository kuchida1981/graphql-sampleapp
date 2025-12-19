# docker-compose-postgresql Specification Delta

## MODIFIED Requirements

### Requirement: PostgreSQL Service in Docker Compose
システムはdocker-composeでPostgreSQLコンテナを起動することをMUSTとする。

#### Scenario: postgresサービスの定義
- **WHEN** `docker-compose.yml` を確認する
- **THEN** `postgres` サービスが定義されている
- **AND** イメージは `postgres:16-alpine` である
- **AND** ポート `5432:5432` がマッピングされている

#### Scenario: PostgreSQL環境変数の設定
- **WHEN** `postgres` サービスの環境変数を確認する
- **THEN** `POSTGRES_USER=graphql_user` が設定されている
- **AND** `POSTGRES_PASSWORD=graphql_pass` が設定されている
- **AND** `POSTGRES_DB=graphql_db` が設定されている

#### Scenario: データ永続化
- **WHEN** `postgres` サービスのボリューム設定を確認する
- **THEN** `postgres_data:/var/lib/postgresql/data` ボリュームがマウントされている
- **AND** `./scripts/init-postgres.sql:/docker-entrypoint-initdb.d/init.sql` がマウントされている

#### Scenario: 初期化スクリプトの自動実行
- **WHEN** `docker compose up` で初回起動する
- **THEN** `/docker-entrypoint-initdb.d/init.sql` が自動的に実行される
- **AND** `users` テーブルが作成される

### Requirement: Application Service Integration
システムはアプリケーションサービスからPostgreSQLに接続できることをMUSTとする。

#### Scenario: DATABASE_URL環境変数の設定
- **WHEN** `app` サービスの環境変数を確認する
- **THEN** `DATABASE_URL=postgres://graphql_user:graphql_pass@postgres:5432/graphql_db?sslmode=disable` が設定されている

#### Scenario: サービス依存関係
- **WHEN** `app` サービスの `depends_on` を確認する
- **THEN** `postgres` が含まれている
- **AND** `firestore` も維持されている

#### Scenario: 起動順序
- **WHEN** `docker compose up` を実行する
- **THEN** `postgres` サービスが先に起動する
- **AND** `app` サービスはPostgreSQLの起動後に起動する

### Requirement: Volume Management
システムはPostgreSQLデータを永続化するボリュームを定義することをMUSTとする。

#### Scenario: postgres_dataボリュームの定義
- **WHEN** `docker-compose.yml` の `volumes` セクションを確認する
- **THEN** `postgres_data:` が定義されている

#### Scenario: データの永続化確認
- **WHEN** `docker compose down` 後に `docker compose up` する
- **THEN** PostgreSQLデータが保持されている
- **AND** usersテーブルとデータが残っている
