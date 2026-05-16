package repositories

import (
	"context"
	"fmt"
	"log/slog"

	"project-manager/ent"
)

// MapError ent のエラーをドメインエラーに変換し、ログも出力する共通ヘルパー。
// 各リポジトリのサブパッケージから呼び出す。
// NotFound 時はエンティティごとに異なるドメインエラー（例: model.ErrProjectNotFound）を返したいため、
// notFound を引数で受け取る。
func MapError(ctx context.Context, op string, err error, notFound error) error {
	if err == nil {
		return nil
	}
	if ent.IsNotFound(err) {
		slog.DebugContext(ctx, op, "result", "not_found")
		return notFound
	}
	slog.ErrorContext(ctx, op, "error", err)
	return fmt.Errorf("%s: %w", op, err)
}
