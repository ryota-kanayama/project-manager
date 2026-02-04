package handlers

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"project-manager/helper"
)

type Handler struct {
	Pool *pgxpool.Pool
}

// ハンドラー: ルート
func (h *Handler) HandleIndex(w http.ResponseWriter, r *http.Request) {
	helper.JsonResponse(w, http.StatusOK, map[string]string{
		"message": "Project Manager API",
	})
}

// ハンドラー: ヘルスチェック
func (h *Handler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	if err := h.Pool.Ping(r.Context()); err != nil {
		helper.ErrorResponse(w, http.StatusServiceUnavailable, "unhealthy")
		return
	}
	helper.JsonResponse(w, http.StatusOK, map[string]string{"status": "healthy"})
}

// ハンドラー: プロジェクト一覧
func (h *Handler) HandleListProjects(w http.ResponseWriter, r *http.Request) {
	helper.JsonResponse(w, http.StatusOK, map[string]any{"projects": []any{}})
}

// ハンドラー: プロジェクト作成
func (h *Handler) HandleCreateProject(w http.ResponseWriter, r *http.Request) {
	helper.JsonResponse(w, http.StatusCreated, map[string]string{"message": "created"})
}

// ハンドラー: プロジェクト取得
func (h *Handler) HandleGetProject(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	helper.JsonResponse(w, http.StatusOK, map[string]string{"id": id})
}

// ハンドラー: プロジェクト更新
func (h *Handler) HandleUpdateProject(w http.ResponseWriter, r *http.Request) {
	helper.JsonResponse(w, http.StatusOK, map[string]string{"message": "updated"})
}

// ハンドラー: プロジェクト削除
func (h *Handler) HandleDeleteProject(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
