# ADR-005: エラーは context を持つ下層で wrap する

## ステータス
メンバー間で合意済み

## コンテキスト
| 層 | wrap するか | 内容 |
|---|---|---|
| 下層（DB 操作など） | する | `fmt.Errorf("...: %w", err)` で操作名・引数を含める（PII は ADR-004 参照） |
| 中継のみの層 | しない | `return err` で素通し |
| ハンドラ層 | しない | ログ出力 + HTTP レスポンス変換のみ |

- 必ず `%w` を使う（`%v` / `%s` だと `errors.Is` / `errors.As` が効かない）
- 中継層で機械的に wrap しない（context が増えない prefix が積み重なるため）

```go
// OK: 下層で wrap
func findProjectByID(ctx context.Context, pool *pgxpool.Pool, id string) (*Project, error) {
    var p Project
    if err := pool.QueryRow(ctx, "SELECT ... WHERE id = $1", id).Scan(&p.ID, &p.Name); err != nil {
        return nil, fmt.Errorf("find project by id %s: %w", id, err)
    }
    return &p, nil
}

// OK: 中継のみ
func (s *ProjectService) Get(ctx context.Context, id string) (*Project, error) {
    return findProjectByID(ctx, s.pool, id)
}

// OK: ハンドラ層は wrap せずログ出力 + レスポンス変換
project, err := h.svc.Get(r.Context(), id)
if err != nil {
    slog.ErrorContext(r.Context(), "failed to get project", "error", err)
    helper.ErrorResponse(w, http.StatusInternalServerError, "internal error")
    return
}
```

```go
// NG: %v で埋め込み（errors.Is / As が効かなくなる）
return nil, fmt.Errorf("find project by id %s: %v", id, err)

// NG: 中継層で意味のない wrap
return nil, fmt.Errorf("service get: %w", err)
```

## 理由
- 操作内容と引数を最も詳しく知る下層で context を付けるのが自然
- `%w` で `errors.Is` / `errors.As` を伝播できる
- ADR-003 の各層ログにより、err 単独に context を積み重ねなくても追跡できる

## トレードオフ
- 「自分しか知らない context は何か」の判断が下層の責任になり、レビューで揃える運用が必要
- err 単独で全 context は読めない（ログ集約基盤が前提）

下層 1 回の wrap で context が不足する場面は、例外的に上位層でも wrap してよい。
