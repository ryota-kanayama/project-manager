package server

import (
	"net/http"

	"project-manager/handlers"
	projecthandler "project-manager/handlers/project"
)

// newRouter ハンドラを組み立て、ルートを登録した http.Handler を返す。
// ミドルウェアの適用は server.Init 側で行う。
func newRouter() http.Handler {
	mux := http.NewServeMux()

	// ルート
	{
		h := handlers.NewHandler()
		mux.HandleFunc("GET /", h.HandleIndex)
		mux.HandleFunc("GET /health", h.HandleHealth)
	}

	// プロジェクト
	{
		ph := projecthandler.NewHandler()
		mux.HandleFunc("GET /api/projects", ph.HandleListProjects)
		mux.HandleFunc("POST /api/projects", ph.HandleCreateProject)
		mux.HandleFunc("GET /api/projects/{id}", ph.HandleGetProject)
		mux.HandleFunc("PUT /api/projects/{id}", ph.HandleUpdateProject)
		mux.HandleFunc("DELETE /api/projects/{id}", ph.HandleDeleteProject)
	}

	return mux
}
