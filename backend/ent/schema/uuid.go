package schema

import "github.com/google/uuid"

// newUUIDv7 UUIDv7 を生成する。各スキーマの ID デフォルトで使う。
func newUUIDv7() uuid.UUID {
	v, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}
	return v
}
