package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"project-manager/helper"
	"project-manager/model"
)

// WriteServiceError サービス層から返ったドメインエラーを HTTP レスポンスへ変換する。
// 各エンティティのハンドラ（サブパッケージ）から呼ぶ共通ヘルパー。
func WriteServiceError(w http.ResponseWriter, r *http.Request, err error, op string) {
	ctx := r.Context()
	switch {
	case errors.Is(err, model.ErrProjectNotFound):
		slog.DebugContext(ctx, op, "result", "not_found")
		helper.ErrorResponse(w, http.StatusNotFound, "not found")
	case errors.Is(err, model.ErrInvalidInput):
		slog.DebugContext(ctx, op, "result", "invalid_input", "error", err)
		helper.ErrorResponse(w, http.StatusBadRequest, err.Error())
	default:
		slog.ErrorContext(ctx, op, "error", err)
		helper.ErrorResponse(w, http.StatusInternalServerError, "internal error")
	}
}
