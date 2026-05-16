package helper

import "context"

// ctxKey ctx 値の衝突を避けるための非公開キー型。
type ctxKey int

const (
	requestIDKey ctxKey = iota
)

// WithRequestID リクエスト ID を context に保存する。
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey, id)
}

// RequestIDFromContext context からリクエスト ID を取り出す。未設定なら空文字。
func RequestIDFromContext(ctx context.Context) string {
	v, _ := ctx.Value(requestIDKey).(string)
	return v
}
