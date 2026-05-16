package response

import (
	"time"

	"github.com/google/uuid"
)

// Project プロジェクトの API レスポンス。
// ent.Project の Edges を含めずシリアライズするために、リポジトリ層で詰め替える。
type Project struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Description *string    `json:"description,omitempty"`
	Status      string     `json:"status"`
	StartDate   *time.Time `json:"start_date,omitempty"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
