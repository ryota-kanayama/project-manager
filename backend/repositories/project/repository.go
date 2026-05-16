package project

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"project-manager/database"
	"project-manager/ent"
	entproject "project-manager/ent/project"
	"project-manager/model"
	"project-manager/repositories"
	"project-manager/response"
)

// Repository プロジェクトの永続化を担う。
// ent 依存はこのパッケージに閉じ込め、外側からは見えないようにする。
type Repository struct {
	client *ent.Client
}

// NewRepository Repository を返す。
// 依存（ent クライアント）は database パッケージのシングルトンから取得する。
func NewRepository() *Repository {
	return &Repository{client: database.Client()}
}

// CreateInput 作成時の入力。Status が空文字なら ent のデフォルトを使う。
type CreateInput struct {
	Name        string
	Description *string
	Status      string
	StartDate   *time.Time
	EndDate     *time.Time
}

// UpdateInput 更新時の入力。nil のフィールドは Clear 扱い。
type UpdateInput struct {
	Name        string
	Description *string
	Status      string
	StartDate   *time.Time
	EndDate     *time.Time
}

func (r *Repository) List(ctx context.Context) ([]response.Project, error) {
	ps, err := r.client.Project.Query().
		Order(ent.Desc(entproject.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, repositories.MapError(ctx, "list projects", err, model.ErrProjectNotFound)
	}
	out := make([]response.Project, len(ps))
	for i, p := range ps {
		out[i] = *toResponse(p)
	}
	return out, nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*response.Project, error) {
	p, err := r.client.Project.Get(ctx, id)
	if err != nil {
		return nil, repositories.MapError(ctx, "get project", err, model.ErrProjectNotFound)
	}
	return toResponse(p), nil
}

func (r *Repository) Create(ctx context.Context, in CreateInput) (*response.Project, error) {
	status, err := parseStatus(in.Status, true)
	if err != nil {
		return nil, err
	}
	p, err := r.client.Project.Create().
		SetName(in.Name).
		SetStatus(status).
		SetNillableDescription(in.Description).
		SetNillableStartDate(in.StartDate).
		SetNillableEndDate(in.EndDate).
		Save(ctx)
	if err != nil {
		return nil, repositories.MapError(ctx, "create project", err, model.ErrProjectNotFound)
	}
	return toResponse(p), nil
}

func (r *Repository) Update(ctx context.Context, id uuid.UUID, in UpdateInput) (*response.Project, error) {
	status, err := parseStatus(in.Status, false)
	if err != nil {
		return nil, err
	}
	update := r.client.Project.UpdateOneID(id).
		SetName(in.Name).
		SetStatus(status)

	if in.Description != nil {
		update.SetNillableDescription(in.Description)
	} else {
		update.ClearDescription()
	}
	if in.StartDate != nil {
		update.SetNillableStartDate(in.StartDate)
	} else {
		update.ClearStartDate()
	}
	if in.EndDate != nil {
		update.SetNillableEndDate(in.EndDate)
	} else {
		update.ClearEndDate()
	}

	p, err := update.Save(ctx)
	if err != nil {
		return nil, repositories.MapError(ctx, "update project", err, model.ErrProjectNotFound)
	}
	return toResponse(p), nil
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.client.Project.DeleteOneID(id).Exec(ctx); err != nil {
		return repositories.MapError(ctx, "delete project", err, model.ErrProjectNotFound)
	}
	return nil
}

// toResponse ent.Project からレスポンス DTO への mapper。
func toResponse(p *ent.Project) *response.Project {
	return &response.Project{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Status:      string(p.Status),
		StartDate:   p.StartDate,
		EndDate:     p.EndDate,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

// parseStatus 文字列を ent enum へ変換。allowEmpty=true なら空文字をデフォルトにフォールバック。
func parseStatus(s string, allowEmpty bool) (entproject.Status, error) {
	if s == "" {
		if allowEmpty {
			return entproject.DefaultStatus, nil
		}
		return "", fmt.Errorf("%w: status is required", model.ErrInvalidInput)
	}
	ps := entproject.Status(s)
	if err := entproject.StatusValidator(ps); err != nil {
		return "", fmt.Errorf("%w: invalid status %q", model.ErrInvalidInput, s)
	}
	return ps, nil
}
