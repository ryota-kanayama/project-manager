# ADR-004: エラーログに PII を含めない

## ステータス
メンバー間で合意済み

## コンテキスト
メールアドレス・氏名・電話番号などの PII をエラーログに含めない。`pgx` / `pgconn` の生エラーはクエリパラメータを含むことがあるため、そのまま出さない。

```go
// NG: 生のエラーをそのまま出力（PII が混入する可能性）
slog.ErrorContext(ctx, err.Error())

// OK: 抽象化したメッセージにする
slog.ErrorContext(ctx, "failed to find active worklog by user")
```

```go
// OK: pgconn.PgError のコード・制約名は PII を含まないので出してよい
var pgErr *pgconn.PgError
if errors.As(err, &pgErr) {
    slog.ErrorContext(ctx, "failed to insert into DB",
        "code", pgErr.Code,
        "constraint", pgErr.ConstraintName,
    )
}
```

ID（UUID）は許容するが、不必要に出さない。

ADR-002 に従い key=value 形式で出す（`fmt.Sprintf` でメッセージに埋め込まない）。

### 適用範囲・例外
- DB 操作層を中心に適用（将来サービス層・リポジトリ層が分かれた場合も同様）
- PII を含まないことが明確なメッセージはそのまま出してよい（例: `"connection refused"`）

## 理由
本番ログはモニタリングや外部サービスへ転送されることがあり、PII を含めると個人情報保護上のリスクになる。
