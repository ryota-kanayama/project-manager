package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"project-manager/database"
	"project-manager/helper"
)

type Handler struct{}

// NewHandler ハンドラの共通機能（Index / Health）を提供する。
// エンティティごとのハンドラはサブパッケージで定義する。
func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) HandleIndex(w http.ResponseWriter, _ *http.Request) {
	helper.JsonResponse(w, http.StatusOK, map[string]string{
		"message": "Project Manager API",
	})
}

func (h *Handler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	if err := database.Ping(r.Context()); err != nil {
		slog.ErrorContext(r.Context(), "health check failed", "error", err)
		helper.ErrorResponse(w, http.StatusServiceUnavailable, "unhealthy")
		return
	}
	helper.JsonResponse(w, http.StatusOK, map[string]string{"status": "healthy"})
}

// ParseID URL パスの {id} を UUID として取り出す共通ヘルパー。
func ParseID(r *http.Request) (uuid.UUID, error) {
	return uuid.Parse(r.PathValue("id"))
}

// DecodeJSON リクエストボディを JSON としてデコードする共通ヘルパー。
// 未知フィールドは拒否する。
func DecodeJSON(r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(dst)
}
