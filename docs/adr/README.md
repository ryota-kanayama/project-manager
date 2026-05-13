# ADR (Architecture Decision Records) 一覧

- [ADR-001](001-layered-architecture-responsibilities.md): レイヤードアーキテクチャの責務分担
- [ADR-002](002-logging-and-formatting.md): ログ出力には `log/slog`、文字列・エラー組立には `fmt` を使う
- [ADR-003](003-log-on-errors.md): エラーを出したらログも出す
- [ADR-004](004-no-pii-in-error-logs.md): エラーログに PII を含めない
- [ADR-005](005-wrap-errors-at-lower-layers.md): エラーは context を持つ下層で wrap する
- [ADR-007](007-adopt-ent.md): ent 採用とレイヤー責務の再定義（ADR-001 を置き換え）
