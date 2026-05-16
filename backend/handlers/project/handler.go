package project

import (
	"net/http"

	"project-manager/handlers"
	"project-manager/helper"
	"project-manager/request"
	projectsvc "project-manager/services/project"
)

type Handler struct {
	project *projectsvc.Service
}

// NewHandler Handler を組み立てる。
// 引数なしで生成でき、内部で services のシングルトンを参照する。
func NewHandler() *Handler {
	return &Handler{
		project: projectsvc.NewService(),
	}
}

func (p *Handler) HandleListProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := p.project.List(r.Context())
	if err != nil {
		handlers.WriteServiceError(w, r, err, "list projects")
		return
	}
	helper.JsonResponse(w, http.StatusOK, map[string]any{"projects": projects})
}

func (p *Handler) HandleCreateProject(w http.ResponseWriter, r *http.Request) {
	var req request.Project
	if !handlers.Bind(w, r, &req) {
		return
	}

	proj, err := p.project.Create(r.Context(), toInput(req))
	if err != nil {
		handlers.WriteServiceError(w, r, err, "create project")
		return
	}
	helper.JsonResponse(w, http.StatusCreated, proj)
}

func (p *Handler) HandleGetProject(w http.ResponseWriter, r *http.Request) {
	id, err := handlers.ParseID(r)
	if err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, "invalid id")
		return
	}
	proj, err := p.project.Get(r.Context(), id)
	if err != nil {
		handlers.WriteServiceError(w, r, err, "get project")
		return
	}
	helper.JsonResponse(w, http.StatusOK, proj)
}

func (p *Handler) HandleUpdateProject(w http.ResponseWriter, r *http.Request) {
	id, err := handlers.ParseID(r)
	if err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req request.Project
	if !handlers.Bind(w, r, &req) {
		return
	}
	proj, err := p.project.Update(r.Context(), id, toInput(req))
	if err != nil {
		handlers.WriteServiceError(w, r, err, "update project")
		return
	}
	helper.JsonResponse(w, http.StatusOK, proj)
}

func (p *Handler) HandleDeleteProject(w http.ResponseWriter, r *http.Request) {
	id, err := handlers.ParseID(r)
	if err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := p.project.Delete(r.Context(), id); err != nil {
		handlers.WriteServiceError(w, r, err, "delete project")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func toInput(req request.Project) projectsvc.Input {
	return projectsvc.Input{
		Name:        req.Name,
		Description: req.Description,
		Status:      req.Status,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	}
}
