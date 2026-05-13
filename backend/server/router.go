package server

import (
	"log/slog"
	"net/http"
	"time"

	"project-manager/ent"
	"project-manager/handlers"
)

// newRouter ルーティングとミドルウェアを組み立てて http.Handler を返す。
// ハンドラ層の組み立てはここで行い、上位（main）はハンドラ層を意識しない。
func newRouter(client *ent.Client, corsOrigin string) http.Handler {
	h := handlers.NewHandler(client)

	mux := http.NewServeMux()

	// ルート
	mux.HandleFunc("GET /", h.HandleIndex)
	mux.HandleFunc("GET /health", h.HandleHealth)

	// API
	mux.HandleFunc("GET /api/projects", h.HandleListProjects)
	mux.HandleFunc("POST /api/projects", h.HandleCreateProject)
	mux.HandleFunc("GET /api/projects/{id}", h.HandleGetProject)
	mux.HandleFunc("PUT /api/projects/{id}", h.HandleUpdateProject)
	mux.HandleFunc("DELETE /api/projects/{id}", h.HandleDeleteProject)

	return loggingMiddleware(corsMiddleware(corsOrigin)(mux))
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		slog.Info("request", "method", r.Method, "path", r.URL.Path, "duration", time.Since(start))
	})
}

func corsMiddleware(origin string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
}
