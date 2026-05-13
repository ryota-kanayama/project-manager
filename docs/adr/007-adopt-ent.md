# ADR-007: ent 採用とレイヤー責務の再定義

## ステータス

メンバー間で合意済み（[ADR-001](001-layered-architecture-responsibilities.md) を置き換え）

## コンテキスト

DB アクセス層に [ent](https://entgo.io/) を採用する。スキーマ・モデル・クエリビルダを Go コード（`ent/schema/*.go`）から一元生成し、コードファーストで実装する。

### レイヤー構成

| 層 | やること | やらないこと |
| --- | --- | --- |
| ハンドラ | リクエストのバインド・バリデーション、サービス呼び出し、サービスが返した DTO の JSON シリアライズ、ドメインエラーの HTTP ステータスへの変換 | ビジネスロジック、DB の直接操作、ent 構造体の参照 |
| サービス | ビジネスロジック（必須・整合性チェック等）、ドメイン入力の組み立て、`repositories` の interface 呼び出し、ent 生成型 → `model` の DTO への詰め替え | HTTP の概念、ent のクエリ Builder |
| リポジトリ | ent クライアントのラップ、`ent.IsNotFound` 等のドメインエラーへの変換、ent enum のバリデーション | ビジネスロジック |
| ent | スキーマ定義、CRUD/クエリ Builder（コード生成） | リポジトリ層の外には漏らさない |

リポジトリ層は **interface として復活させる**。サービス層から ent の Builder API・`ent.IsNotFound`・`project.Status` 等を直接見えないようにすることで、サービス層の可読性とテスタビリティを確保する。

### ディレクトリ構成

```text
backend/
├── handlers/       # ハンドラ層
├── services/       # サービス層（repositories を呼び、DTO に詰め替える）
├── repositories/   # リポジトリ層（ent ラッパー、ent 依存をここに閉じ込める）
├── ent/            # ent コード生成（自動生成）
│   └── schema/     # source of truth（Go でスキーマ定義）
├── model/          # ドメインエラー + API レスポンス DTO
├── helper/         # 共通ヘルパー
├── config/         # 環境変数・.env 読み込み
├── database/       # ent.Client 初期化・保持・Ping
└── server/         # HTTP サーバ起動・ルーター
```

### マイグレーション

- 開発期は ent の自動マイグレーション（`client.Schema.Create(ctx)`）を採用
- 本番運用フェーズ移行時に Versioned Migration（[atlas](https://atlasgo.io/) 連携で SQL ファイル管理）へ切り替える

### エラー変換

- `ent.IsNotFound(err)` を `model.ErrProjectNotFound` 等のドメインエラーへ変換するのはリポジトリ層の責務
- ent enum の妥当性検証もリポジトリ層で行う（`project.StatusValidator` 等）
- ADR-005 の wrap 方針はリポジトリ層に適用する（ent クライアント呼び出しの直下で wrap）

### DTO（API レスポンス）

- レスポンス用 DTO は `model/` パッケージに定義する（例: `model.ProjectResponse`）
- ent 生成型は `Edges` フィールドを含むため JSON シリアライズで `"edges":{}` が混入する。これをサービス層で DTO に詰め替えることで防ぐ
- 詰め替えは各サービスメソッド内に **inline** で書く（補助関数 / mapper を介さない）。Get / Create / Update / List で同じ詰め替えが繰り返されるが、読み下しの簡潔さを優先する
- ハンドラ層は DTO をそのまま `helper.JsonResponse` に渡す（再変換しない）

### 戻り値の型

| 層 | 戻り値の型 |
| --- | --- |
| リポジトリ | `*ent.Project` 等の ent 生成型（コード量との妥協） |
| サービス | `*model.ProjectResponse` 等の DTO |
| ハンドラ | DTO を JSON シリアライズ |

完全に ent を閉じ込めたい場合（リポジトリも独自ドメインモデルを返す）は、`model.Project` + mapper を導入する。サービスの単体テストで repository をモックしたい等の実需が出た時点で再検討する。

## 理由

- スキーマ・構造体・クエリの 3 重管理を解消（ent スキーマが単一の source of truth）
- 型安全な Builder API でクエリを書ける（reflection 不要）
- リレーション（Edges）が一級市民として表現される
- 構造体・JSON タグ・マイグレーション SQL がすべて自動生成され、カラム追加時の変更点が `ent/schema/*.go` 一箇所に収まる

## トレードオフ

- 複雑な集計クエリは ent Builder では書きにくいことがあり、その場合は `entsql` の Modify / Raw SQL 機能を使う
- 学習コストは生 SQL より高い（Edges, Hook, Mixin 等の独自概念）
- 自動マイグレーションは破壊的変更を意図せず反映するリスクがあるため、本番フェーズでは atlas 連携が必須
- ent 生成コードがリポジトリのサイズを増やす（バージョン管理に含めるが、レビュー対象外）
