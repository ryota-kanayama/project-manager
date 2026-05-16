package request

import "time"

// Project プロジェクト作成・更新のリクエストボディ。
//
// validate タグでフィールド単位のバリデーションを行う。複雑なバリデーション
// （業務ルール）はサービス層で実施する。
type Project struct {
	Name        string     `json:"name" validate:"required,max=255"`
	Description *string    `json:"description,omitempty"`
	Status      string     `json:"status,omitempty" validate:"omitempty,oneof=planning in_progress completed on_hold"`
	StartDate   *time.Time `json:"start_date,omitempty"`
	EndDate     *time.Time `json:"end_date,omitempty" validate:"omitempty,gtefield=StartDate"`
}
