# ADR-002: ログ出力には `log/slog`、文字列・エラー組立には `fmt` を使う

## ステータス
メンバー間で合意済み

## コンテキスト
- ログ出力は `log/slog` を使う
- エラーラップ・文字列組立は `fmt` を使う
- `log` パッケージは新規コードでは使わない

ルール:

- 出力: `slog.Info` / `Warn` / `Error` / `Debug`（ctx があれば `*Context` 系）
- ログは key=value の構造化形式で書く
- エラーラップ: `fmt.Errorf("...: %w", err)`
- 文字列組立: `fmt.Sprintf(...)`
- ロガーは `main` で初期化、以降は `slog` のパッケージ関数を使う（独自 Logger を引き回さない）

```go
// OK
slog.Info("connected to database", "host", host)
return fmt.Errorf("connect %s: %w", host, err)
dsn := fmt.Sprintf("postgres://%s@%s/%s", user, host, name)

// NG
fmt.Printf("connected to %s\n", host)        // 診断ログを fmt で出している
log.Println("server starting")               // log パッケージの新規利用
slog.Info(fmt.Sprintf("port=%s", port))      // メッセージに値を埋め込んでいる（key で分離する）
```

## 理由
- `slog` は構造化・レベル・JSON 出力・コンテキスト連携を標準で備える
- 「出力」と「文字列・エラー組立」を別パッケージに分けることで責務が明確になる
- 標準ライブラリのみで完結し、依存追加が不要

## トレードオフ
- 既存の `log.Printf` / `log.Println` は順次 `slog` に移行が必要
- ログメッセージ整形の自由度を犠牲にする（key=value 強制）
- `zap` / `zerolog` は採用しない。性能要件が顕在化したら再検討する
