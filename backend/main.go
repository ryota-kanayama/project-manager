package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"project-manager/database"
	"project-manager/handlers"
	"project-manager/schema"

	"github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool

func main() {
	ctx := context.Background()

	// データベース接続
	var err error
	pool, err = database.Connect(ctx)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer pool.Close()
	log.Println("Connected to database")

	// マイグレーション実行
	if err := schema.Migrate(ctx, pool); err != nil {
		log.Fatalf("Migration failed: %v", err)
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
		log.Printf("Server starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exited")
}

// ミドルウェア: ログ出力
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
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
