package logger

import (
	"context"
	"log/slog"
	"os"

	"project-manager/helper"
)

// Init slog のデフォルトロガーを設定する。
// context にリクエスト ID があれば全ログに自動付与する。
func Init() {
	base := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	slog.SetDefault(slog.New(&ctxHandler{Handler: base}))
}

// ctxHandler context.Context から値を読んでログ属性に変換する slog.Handler。
type ctxHandler struct {
	slog.Handler
}

func (h *ctxHandler) Handle(ctx context.Context, r slog.Record) error {
	if id := helper.RequestIDFromContext(ctx); id != "" {
		r.AddAttrs(slog.String("req_id", id))
	}
	return h.Handler.Handle(ctx, r)
}

func (h *ctxHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ctxHandler{Handler: h.Handler.WithAttrs(attrs)}
}

func (h *ctxHandler) WithGroup(name string) slog.Handler {
	return &ctxHandler{Handler: h.Handler.WithGroup(name)}
}
