package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"project-manager/database"
	"project-manager/ent"
	"project-manager/helper"
	"project-manager/model"
	"project-manager/services"
)

type Handler struct {
	project *services.ProjectService
}

// NewHandler ハンドラ層が自身の依存（サービス）を組み立てる。
func NewHandler(client *ent.Client) *Handler {
	return &Handler{
		project: services.NewProjectService(client),
	}
}

// projectRequest プロジェクト作成・更新のリクエストボディ
type projectRequest struct {
	Name        string     `json:"name"`
	Description *string    `json:"description,omitempty"`
	Status      string     `json:"status,omitempty"`
	StartDate   *time.Time `json:"start_date,omitempty"`
	EndDate     *time.Time `json:"end_date,omitempty"`
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

func (h *Handler) HandleListProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := h.project.List(r.Context())
	if err != nil {
		h.writeServiceError(w, r, err, "list projects")
		return
	}
	helper.JsonResponse(w, http.StatusOK, map[string]any{"projects": projects})
}

func (h *Handler) HandleCreateProject(w http.ResponseWriter, r *http.Request) {
	var req projectRequest
	if err := decodeJSON(r, &req); err != nil {
		slog.DebugContext(r.Context(), "failed to decode create project body", "error", err)
		helper.ErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	p, err := h.project.Create(r.Context(), toProjectInput(req))
	if err != nil {
		h.writeServiceError(w, r, err, "create project")
		return
	}
	helper.JsonResponse(w, http.StatusCreated, p)
}

func (h *Handler) HandleGetProject(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, "invalid id")
		return
	}
	p, err := h.project.Get(r.Context(), id)
	if err != nil {
		h.writeServiceError(w, r, err, "get project")
		return
	}
	helper.JsonResponse(w, http.StatusOK, p)
}

func (h *Handler) HandleUpdateProject(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req projectRequest
	if err := decodeJSON(r, &req); err != nil {
		slog.DebugContext(r.Context(), "failed to decode update project body", "error", err)
		helper.ErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}
	p, err := h.project.Update(r.Context(), id, toProjectInput(req))
	if err != nil {
		h.writeServiceError(w, r, err, "update project")
		return
	}
	helper.JsonResponse(w, http.StatusOK, p)
}

func (h *Handler) HandleDeleteProject(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.project.Delete(r.Context(), id); err != nil {
		h.writeServiceError(w, r, err, "delete project")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func toProjectInput(req projectRequest) services.ProjectInput {
	return services.ProjectInput{
		Name:        req.Name,
		Description: req.Description,
		Status:      req.Status,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	}
}

func parseID(r *http.Request) (uuid.UUID, error) {
	return uuid.Parse(r.PathValue("id"))
}

func decodeJSON(r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(dst)
}

// writeServiceError サービス層から返ったエラーを HTTP レスポンスへ変換する。
func (h *Handler) writeServiceError(w http.ResponseWriter, r *http.Request, err error, op string) {
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
