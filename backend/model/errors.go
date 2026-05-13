package model

import "errors"

// ドメインエラー
var (
	// ErrProjectNotFound 指定されたプロジェクトが存在しない
	ErrProjectNotFound = errors.New("project not found")
	// ErrInvalidInput 入力値が不正
	ErrInvalidInput = errors.New("invalid input")
)
