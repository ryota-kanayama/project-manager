package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

// Init .env ファイルを読み込む（envFile が空ならスキップ）。
func Init(envFile string) {
	if envFile == "" {
		return
	}
	if err := godotenv.Load(envFile); err != nil {
		panic(err)
	}
	slog.Info("loaded env file", "path", envFile)
}

// Get 環境変数を取得し、未設定なら fallback を返す。
func Get(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}
