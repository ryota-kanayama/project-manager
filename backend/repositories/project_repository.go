package repositories

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"project-manager/ent"
	"project-manager/ent/project"
	"project-manager/model"
)

// ProjectRepository プロジェクトの永続化インタフェース。
// ent 依存はここに閉じ込め、サービス層からは ent のクエリ Builder が見えないようにする。
type ProjectRepository interface {
	List(ctx context.Context) ([]*ent.Project, error)
	GetByID(ctx context.Context, id uuid.UUID) (*ent.Project, error)
	Create(ctx context.Context, in ProjectCreateInput) (*ent.Project, error)
	Update(ctx context.Context, id uuid.UUID, in ProjectUpdateInput) (*ent.Project, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// ProjectCreateInput 作成時の入力。Status が空文字なら ent のデフォルトを使う。
type ProjectCreateInput struct {
	Name        string
	Description *string
	Status      string
	StartDate   *time.Time
	EndDate     *time.Time
}

// ProjectUpdateInput 更新時の入力。nil のフィールドは Clear 扱い。
type ProjectUpdateInput struct {
	Name        string
	Description *string
	Status      string
	StartDate   *time.Time
	EndDate     *time.Time
}

// projectRepository ent ベースの実装。
type projectRepository struct {
	client *ent.Client
}

// NewProjectRepository ent クライアントから ProjectRepository を構築する。
func NewProjectRepository(client *ent.Client) ProjectRepository {
	return &projectRepository{client: client}
}

func (r *projectRepository) List(ctx context.Context) ([]*ent.Project, error) {
	projects, err := r.client.Project.Query().
		Order(ent.Desc(project.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "failed to list projects", "error", err)
		return nil, fmt.Errorf("list projects: %w", err)
	}
	return projects, nil
}

func (r *projectRepository) GetByID(ctx context.Context, id uuid.UUID) (*ent.Project, error) {
	p, err := r.client.Project.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			slog.DebugContext(ctx, "project not found", "id", id)
			return nil, model.ErrProjectNotFound
		}
		slog.ErrorContext(ctx, "failed to get project", "error", err, "id", id)
		return nil, fmt.Errorf("get project %s: %w", id, err)
	}
	return p, nil
}

func (r *projectRepository) Create(ctx context.Context, in ProjectCreateInput) (*ent.Project, error) {
	status, err := parseProjectStatus(in.Status, true)
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
		slog.ErrorContext(ctx, "failed to create project", "error", err)
		return nil, fmt.Errorf("create project: %w", err)
	}
	return p, nil
}

func (r *projectRepository) Update(ctx context.Context, id uuid.UUID, in ProjectUpdateInput) (*ent.Project, error) {
	status, err := parseProjectStatus(in.Status, false)
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
		if ent.IsNotFound(err) {
			slog.DebugContext(ctx, "project not found for update", "id", id)
			return nil, model.ErrProjectNotFound
		}
		slog.ErrorContext(ctx, "failed to update project", "error", err, "id", id)
		return nil, fmt.Errorf("update project %s: %w", id, err)
	}
	return p, nil
}

func (r *projectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.client.Project.DeleteOneID(id).Exec(ctx); err != nil {
		if ent.IsNotFound(err) {
			slog.DebugContext(ctx, "project not found for delete", "id", id)
			return model.ErrProjectNotFound
		}
		slog.ErrorContext(ctx, "failed to delete project", "error", err, "id", id)
		return fmt.Errorf("delete project %s: %w", id, err)
	}
	return nil
}

// parseProjectStatus 文字列を ent enum へ変換。allowEmpty=true なら空文字をデフォルトにフォールバック。
func parseProjectStatus(s string, allowEmpty bool) (project.Status, error) {
	if s == "" {
		if allowEmpty {
			return project.DefaultStatus, nil
		}
		return "", fmt.Errorf("%w: status is required", model.ErrInvalidInput)
	}
	ps := project.Status(s)
	if err := project.StatusValidator(ps); err != nil {
		return "", fmt.Errorf("%w: invalid status %q", model.ErrInvalidInput, s)
	}
	return ps, nil
}
