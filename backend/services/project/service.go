package project

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	projectrepo "project-manager/repositories/project"
	"project-manager/response"
)

type Service struct {
	repo *projectrepo.Repository
}

// NewService サービス層が自身の依存（リポジトリ）を組み立てる。
// 引数なしで生成でき、内部で repositories のシングルトンを参照する。
func NewService() *Service {
	return &Service{repo: projectrepo.NewRepository()}
}

// Input サービスへの入力。バリデーションは呼び出し元（handlers.Bind）で実施済みとして扱う。
type Input struct {
	Name        string
	Description *string
	Status      string
	StartDate   *time.Time
	EndDate     *time.Time
}

func (s *Service) List(ctx context.Context) ([]response.Project, error) {
	return s.repo.List(ctx)
}

func (s *Service) Get(ctx context.Context, id uuid.UUID) (*response.Project, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) Create(ctx context.Context, in Input) (*response.Project, error) {
	return s.repo.Create(ctx, projectrepo.CreateInput{
		Name:        strings.TrimSpace(in.Name),
		Description: in.Description,
		Status:      in.Status,
		StartDate:   in.StartDate,
		EndDate:     in.EndDate,
	})
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, in Input) (*response.Project, error) {
	return s.repo.Update(ctx, id, projectrepo.UpdateInput{
		Name:        strings.TrimSpace(in.Name),
		Description: in.Description,
		Status:      in.Status,
		StartDate:   in.StartDate,
		EndDate:     in.EndDate,
	})
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
