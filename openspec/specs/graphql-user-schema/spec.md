# graphql-user-schema Specification

## Purpose
TBD - created by archiving change add-postgresql-user-integration. Update Purpose after archive.
## Requirements
### Requirement: GraphQL Schema Extension
システムはUserデータを取得するためのGraphQLクエリをスキーマに追加することをMUSTとする。

#### Scenario: Userタイプの定義
- **WHEN** GraphQLスキーマファイル（`schema.graphqls`）を確認する
- **THEN** `User` タイプが定義されている
- **AND** `id: ID!` フィールドが含まれる
- **AND** `name: String!` フィールドが含まれる
- **AND** `email: String!` フィールドが含まれる
- **AND** `createdAt: String!` フィールドが含まれる

#### Scenario: usersクエリの定義
- **WHEN** GraphQLスキーマの `Query` タイプを確認する
- **THEN** `users: [User!]!` クエリが定義されている
- **AND** 既存の `hello`, `messages`, `message` クエリは維持されている

#### Scenario: userクエリの定義（ID指定）
- **WHEN** GraphQLスキーマの `Query` タイプを確認する
- **THEN** `user(id: ID!): User` クエリが定義されている
- **AND** 引数 `id` が必須（`ID!`）である

### Requirement: Resolver Implementation
システムはUserRepositoryを使用してGraphQLクエリを解決することをMUSTとする。

#### Scenario: usersリゾルバーの実装
- **WHEN** `users` クエリが実行される
- **THEN** Resolverは `userRepo.List(ctx)` を呼び出す
- **AND** PostgreSQLから全Userレコードを取得する
- **AND** User配列をGraphQLレスポンスとして返す

#### Scenario: userリゾルバーの実装（ID指定）
- **WHEN** `user(id: "1")` クエリが実行される
- **THEN** Resolverは `userRepo.GetByID(ctx, "1")` を呼び出す
- **AND** 指定されたIDのUserレコードを取得する
- **AND** 該当するUserをGraphQLレスポンスとして返す

#### Scenario: ユーザー未存在のエラーハンドリング
- **WHEN** 存在しないIDで `user(id: "nonexistent")` クエリが実行される
- **THEN** Repositoryは `nil` と `error` を返す
- **AND** Resolverは GraphQL エラーレスポンスを返す
- **AND** エラーメッセージは "user not found" を含む

### Requirement: GraphQL Integration Testing
システムはGraphQL Playgroundを通じてUserクエリをテストできることをMUSTとする。

#### Scenario: usersクエリの実行成功
- **WHEN** GraphQL Playgroundで `{ users { id name email createdAt } }` を実行する
- **THEN** PostgreSQLに存在する全Userが返される
- **AND** レスポンスは有効なJSONである
- **AND** HTTPステータスコードは200である

#### Scenario: userクエリの実行成功（ID指定）
- **WHEN** GraphQL Playgroundで `{ user(id: "1") { id name email createdAt } }` を実行する
- **THEN** ID "1" のUserが返される
- **AND** レスポンスはUserの全フィールドを含む

#### Scenario: スキーマエクスプローラーでの確認
- **WHEN** GraphQL Playgroundのスキーマエクスプローラーを開く
- **THEN** `Query.users` が表示される
- **AND** `Query.user(id: ID!)` が表示される
- **AND** `User` タイプとそのフィールドが表示される

#### Scenario: 後方互換性の確認
- **WHEN** 既存の `messages` クエリを実行する
- **THEN** Firestore Message機能が正常に動作する
- **AND** User機能の追加による影響がない

