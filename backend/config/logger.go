package config

import (
	"context"
	"log/slog"
	"os"

	"project-manager/helper"
)

// InitLogger slog のデフォルトロガーを設定する。
// context から req_id を取り出して各ログレコードに自動付与する。
func InitLogger() {
	base := slog.NewTextHandler(os.Stdout, nil)
	slog.SetDefault(slog.New(&ctxHandler{Handler: base}))
}

// ctxHandler context から req_id を取り出してログに付与する slog.Handler ラッパー。
type ctxHandler struct {
	slog.Handler
}

func (h *ctxHandler) Handle(ctx context.Context, r slog.Record) error {
	if id := helper.RequestIDFromContext(ctx); id != "" {
		r.AddAttrs(slog.String("req_id", id))
	}
	return h.Handler.Handle(ctx, r)
}
