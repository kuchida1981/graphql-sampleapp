# docker-compose-setup Specification

## Purpose
TBD - created by archiving change add-firestore-docker-integration. Update Purpose after archive.
## Requirements
### Requirement: Docker Compose Configuration
システムはDocker Compose設定ファイルを提供し、アプリケーションとFirestore Emulatorサービスを定義することをMUSTとする。

#### Scenario: docker-compose.ymlが存在する
- **WHEN** プロジェクトルートディレクトリを確認する
- **THEN** `docker-compose.yml` ファイルが存在する
- **AND** ファイルには `app` サービスと `firestore` サービスが定義されている

#### Scenario: サービス間ネットワーク接続
- **WHEN** `docker-compose up` を実行する
- **THEN** `app` サービスは `firestore` サービスにサービス名で接続できる
- **AND** `app` サービスの環境変数 `FIRESTORE_EMULATOR_HOST` が `firestore:8081` に設定されている

### Requirement: Application Container
システムはGoアプリケーション用のDockerイメージを定義し、GraphQLサーバーを実行することをMUSTとする。

#### Scenario: Dockerfileが存在する
- **WHEN** プロジェクトルートディレクトリを確認する
- **THEN** `Dockerfile` が存在する
- **AND** マルチステージビルドを使用している

#### Scenario: アプリケーションコンテナの起動
- **WHEN** `docker-compose up app` を実行する
- **THEN** Goアプリケーションがビルドされコンテナ内で起動する
- **AND** ポート8080でGraphQLサーバーがリッスンする
- **AND** ホストマシンから `http://localhost:8080` でアクセスできる

#### Scenario: ホットリロード対応（開発モード）
- **WHEN** Docker Composeで起動時にボリュームマウントを使用する
- **THEN** ホストのソースコード変更がコンテナ内に反映される
- **AND** コンテナ再ビルドなしで開発できる

### Requirement: Firestore Emulator Container
システムはFirestore Emulator用のコンテナを定義し、ローカルFirestoreインスタンスを提供することをMUSTとする。

#### Scenario: Firestore Emulatorサービスの定義
- **WHEN** `docker-compose.yml` を確認する
- **THEN** `firestore` サービスが `google/cloud-sdk` イメージを使用している
- **AND** `gcloud emulators firestore start` コマンドでEmulatorを起動する
- **AND** ポート8081でEmulatorがリッスンする

#### Scenario: Firestore Emulatorの起動
- **WHEN** `docker-compose up firestore` を実行する
- **THEN** Firestore Emulatorが正常に起動する
- **AND** `app` サービスからポート8081で接続できる
- **AND** エミュレーターログが標準出力に表示される

### Requirement: Environment Configuration
システムは環境変数を使用してコンテナ設定を管理することをMUSTとする。

#### Scenario: 環境変数の設定
- **WHEN** `docker-compose.yml` の `app` サービス定義を確認する
- **THEN** `FIRESTORE_EMULATOR_HOST` 環境変数が設定されている
- **AND** `GCP_PROJECT_ID` 環境変数がダミー値（例: `demo-project`）に設定されている
- **AND** `PORT` 環境変数が `8080` に設定されている

#### Scenario: .envファイルのサポート（オプション）
- **WHEN** プロジェクトルートに `.env` ファイルを配置する
- **THEN** Docker Composeが `.env` ファイルから環境変数を読み込む
- **AND** `.env.example` がサンプルとして提供されている

### Requirement: Service Orchestration
システムは複数サービスを正しい順序で起動し、依存関係を管理することをMUSTとする。

#### Scenario: サービス起動順序
- **WHEN** `docker-compose up` を実行する
- **THEN** `firestore` サービスが先に起動する
- **AND** `firestore` が ready 状態になった後に `app` サービスが起動する
- **AND** `depends_on` ディレクティブで依存関係が定義されている

#### Scenario: 全サービスの一括起動
- **WHEN** `docker-compose up` をパラメータなしで実行する
- **THEN** すべてのサービス（app, firestore）が起動する
- **AND** 両サービスのログが統合されて表示される

### Requirement: Development Workflow
システムは開発者が簡単に環境を起動・停止・リセットできることをMUSTとする。

#### Scenario: 環境の起動
- **WHEN** `docker-compose up -d` を実行する
- **THEN** すべてのサービスがバックグラウンドで起動する
- **AND** コマンドが即座に完了する

#### Scenario: 環境の停止
- **WHEN** `docker-compose down` を実行する
- **THEN** すべてのコンテナが停止され削除される
- **AND** ネットワークが削除される

#### Scenario: ログの確認
- **WHEN** `docker-compose logs -f app` を実行する
- **THEN** `app` サービスのリアルタイムログが表示される
- **AND** Ctrl+C でログストリームを終了できる

#### Scenario: データのリセット
- **WHEN** `docker-compose down -v` を実行する
- **THEN** コンテナ、ネットワーク、およびボリュームが削除される
- **AND** 次回起動時にFirestoreデータがクリーンな状態から開始する

