package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"project-manager/config"
)

// Init HTTP サーバを起動し、シグナル受信でグレースフルにシャットダウンする。
func Init() error {
	// ミドルウェアチェーン: requestID → accessLog → cors → router
	corsOrigin := config.Get("CORS_ORIGIN", "http://localhost:3000")
	handler := requestIDMiddleware(
		accessLogMiddleware(
			corsMiddleware(corsOrigin)(newRouter()),
		),
	)

	port := config.Get("SERVER_PORT", "8080")
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		slog.Info("server starting", "url", fmt.Sprintf("http://localhost:%s", port))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
	case err := <-errCh:
		return fmt.Errorf("server listen: %w", err)
	}

	slog.Info("shutting down server")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server shutdown: %w", err)
	}
	slog.Info("server exited")
	return nil
}
