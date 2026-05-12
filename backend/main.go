package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"project-manager/database"
	"project-manager/handlers"
	"project-manager/schema"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var pool *pgxpool.Pool

func main() {
	envFile := flag.String("env", "", "path to .env file (optional)")
	flag.Parse()

	if *envFile != "" {
		if err := godotenv.Load(*envFile); err != nil {
			panic(err)
		}
		slog.Info("loaded env file", "path", *envFile)
	}

	ctx := context.Background()

	// データベース接続
	var err error
	pool, err = database.Connect(ctx)
	if err != nil {
		panic(err)
	}
	defer pool.Close()
	slog.Info("connected to database")

	// マイグレーション実行
	if err := schema.Migrate(ctx, pool); err != nil {
		panic(err)
	}

	// ハンドラー初期化
	h := &handlers.Handler{Pool: pool}

	// ルーター設定
	mux := http.NewServeMux()

	// ルート
	mux.HandleFunc("GET /", h.HandleIndex)
	mux.HandleFunc("GET /health", h.HandleHealth)

	// APIルート
	mux.HandleFunc("GET /api/projects", h.HandleListProjects)
	mux.HandleFunc("POST /api/projects", h.HandleCreateProject)
	mux.HandleFunc("GET /api/projects/{id}", h.HandleGetProject)
	mux.HandleFunc("PUT /api/projects/{id}", h.HandleUpdateProject)
	mux.HandleFunc("DELETE /api/projects/{id}", h.HandleDeleteProject)

	// ミドルウェアを適用
	wrappedHandler := loggingMiddleware(corsMiddleware(mux))

	// サーバー起動
	port := database.GetEnv("SERVER_PORT", "8080")
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      wrappedHandler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		slog.Info("server starting", "url", fmt.Sprintf("http://localhost:%s", port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		panic(err)
	}
	slog.Info("server exited")
}

// ミドルウェア: ログ出力
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		slog.Info("request", "method", r.Method, "path", r.URL.Path, "duration", time.Since(start))
	})
}

// ミドルウェア: CORS
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := database.GetEnv("CORS_ORIGIN", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
