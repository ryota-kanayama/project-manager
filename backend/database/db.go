package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/jackc/pgx/v5/stdlib" // pgx ドライバを database/sql に登録
	"project-manager/config"
	"project-manager/ent"
)

var (
	client *ent.Client
	sqlDB  *sql.DB
)

// Init データベースに接続し、ent スキーマを適用する。
func Init() error {
	ctx := context.Background()

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		config.Get("DB_USER", "pmuser"),
		config.Get("DB_PASSWORD", "pmpassword"),
		config.Get("DB_HOST", "localhost"),
		config.Get("DB_PORT", "5432"),
		config.Get("DB_NAME", "project_manager"),
		config.Get("DB_SSLMODE", "disable"),
	)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("open sql: %w", err)
	}
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return fmt.Errorf("ping db: %w", err)
	}
	sqlDB = db

	drv := entsql.OpenDB(dialect.Postgres, db)
	client = ent.NewClient(ent.Driver(drv))
	slog.Info("connected to database")

	if err := client.Schema.Create(ctx); err != nil {
		return fmt.Errorf("ent schema create: %w", err)
	}
	slog.Info("ent schema applied")
	return nil
}

// Close ent クライアントをクローズする。
func Close() {
	if client != nil {
		_ = client.Close()
	}
}

// Client ent クライアントを返す。Init 前に呼ぶと nil。
func Client() *ent.Client {
	return client
}

// Ping データベース接続を確認する。ヘルスチェック用。
func Ping(ctx context.Context) error {
	if sqlDB == nil {
		return errors.New("database not initialized")
	}
	return sqlDB.PingContext(ctx)
}
