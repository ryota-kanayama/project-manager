package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"project-manager/ent"
	"project-manager/model"
	"project-manager/repositories"
)

type ProjectService struct {
	repo repositories.ProjectRepository
}

// NewProjectService サービス層が自身の依存（リポジトリ）を組み立てる。
// ent 依存はリポジトリ層に閉じ込められており、ここでは見えない。
func NewProjectService(client *ent.Client) *ProjectService {
	return &ProjectService{repo: repositories.NewProjectRepository(client)}
}

type ProjectInput struct {
	Name        string
	Description *string
	Status      string
	StartDate   *time.Time
	EndDate     *time.Time
}

func (s *ProjectService) List(ctx context.Context) ([]model.ProjectResponse, error) {
	ps, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]model.ProjectResponse, len(ps))
	for i, p := range ps {
		out[i] = model.ProjectResponse{
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
	return out, nil
}

func (s *ProjectService) Get(ctx context.Context, id uuid.UUID) (*model.ProjectResponse, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &model.ProjectResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Status:      string(p.Status),
		StartDate:   p.StartDate,
		EndDate:     p.EndDate,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}, nil
}

func (s *ProjectService) Create(ctx context.Context, in ProjectInput) (*model.ProjectResponse, error) {
	if err := validateInput(in); err != nil {
		return nil, err
	}
	p, err := s.repo.Create(ctx, repositories.ProjectCreateInput{
		Name:        strings.TrimSpace(in.Name),
		Description: in.Description,
		Status:      in.Status,
		StartDate:   in.StartDate,
		EndDate:     in.EndDate,
	})
	if err != nil {
		return nil, err
	}
	return &model.ProjectResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Status:      string(p.Status),
		StartDate:   p.StartDate,
		EndDate:     p.EndDate,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}, nil
}

func (s *ProjectService) Update(ctx context.Context, id uuid.UUID, in ProjectInput) (*model.ProjectResponse, error) {
	if err := validateInput(in); err != nil {
		return nil, err
	}
	p, err := s.repo.Update(ctx, id, repositories.ProjectUpdateInput{
		Name:        strings.TrimSpace(in.Name),
		Description: in.Description,
		Status:      in.Status,
		StartDate:   in.StartDate,
		EndDate:     in.EndDate,
	})
	if err != nil {
		return nil, err
	}
	return &model.ProjectResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Status:      string(p.Status),
		StartDate:   p.StartDate,
		EndDate:     p.EndDate,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}, nil
}

func (s *ProjectService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func validateInput(in ProjectInput) error {
	if strings.TrimSpace(in.Name) == "" {
		return fmt.Errorf("%w: name is required", model.ErrInvalidInput)
	}
	if in.StartDate != nil && in.EndDate != nil && in.EndDate.Before(*in.StartDate) {
		return fmt.Errorf("%w: end_date must be after start_date", model.ErrInvalidInput)
	}
	return nil
}
