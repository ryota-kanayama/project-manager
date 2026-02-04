package schema

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Schema テーブル作成用SQL
// Goの構造体と対応するテーブル定義を一元管理
var schemas = []string{
	// 拡張機能
	`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`,

	// projects テーブル
	`CREATE TABLE IF NOT EXISTS projects (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		name VARCHAR(255) NOT NULL,
		description TEXT,
		status VARCHAR(50) NOT NULL DEFAULT 'planning',
		start_date DATE,
		end_date DATE,
		created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`,

	// milestones テーブル
	`CREATE TABLE IF NOT EXISTS milestones (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
		name VARCHAR(255) NOT NULL,
		due_date DATE,
		status VARCHAR(50) NOT NULL DEFAULT 'pending',
		created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`,

	// tasks テーブル
	`CREATE TABLE IF NOT EXISTS tasks (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
		parent_id UUID REFERENCES tasks(id) ON DELETE CASCADE,
		milestone_id UUID REFERENCES milestones(id) ON DELETE SET NULL,
		wbs_code VARCHAR(50),
		name VARCHAR(255) NOT NULL,
		description TEXT,
		status VARCHAR(50) NOT NULL DEFAULT 'not_started',
		priority VARCHAR(20) NOT NULL DEFAULT 'medium',
		assignee VARCHAR(255),
		estimated_hours DECIMAL(10, 2),
		actual_hours DECIMAL(10, 2),
		start_date DATE,
		end_date DATE,
		sort_order INTEGER NOT NULL DEFAULT 0,
		created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`,

	// worklogs テーブル
	`CREATE TABLE IF NOT EXISTS worklogs (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		task_id UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
		user_name VARCHAR(255) NOT NULL,
		hours DECIMAL(10, 2) NOT NULL,
		work_date DATE NOT NULL,
		description TEXT,
		created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`,

	// インデックス
	`CREATE INDEX IF NOT EXISTS idx_milestones_project_id ON milestones(project_id)`,
	`CREATE INDEX IF NOT EXISTS idx_tasks_project_id ON tasks(project_id)`,
	`CREATE INDEX IF NOT EXISTS idx_tasks_parent_id ON tasks(parent_id)`,
	`CREATE INDEX IF NOT EXISTS idx_tasks_milestone_id ON tasks(milestone_id)`,
	`CREATE INDEX IF NOT EXISTS idx_worklogs_task_id ON worklogs(task_id)`,
	`CREATE INDEX IF NOT EXISTS idx_worklogs_work_date ON worklogs(work_date)`,
}

// Migrate スキーマをデータベースに適用
func Migrate(ctx context.Context, pool *pgxpool.Pool) error {
	log.Println("Running migrations...")

	for i, schema := range schemas {
		if _, err := pool.Exec(ctx, schema); err != nil {
			return fmt.Errorf("migration %d failed: %w", i, err)
		}
	}

	log.Println("Migrations completed successfully")
	return nil
}

// Drop 全テーブルを削除（開発用）
func Drop(ctx context.Context, pool *pgxpool.Pool) error {
	log.Println("Dropping all tables...")

	dropStatements := []string{
		`DROP TABLE IF EXISTS worklogs CASCADE`,
		`DROP TABLE IF EXISTS tasks CASCADE`,
		`DROP TABLE IF EXISTS milestones CASCADE`,
		`DROP TABLE IF EXISTS projects CASCADE`,
	}

	for _, stmt := range dropStatements {
		if _, err := pool.Exec(ctx, stmt); err != nil {
			return fmt.Errorf("drop failed: %w", err)
		}
	}

	log.Println("All tables dropped")
	return nil
}
