package helper

import (
	"encoding/json"
	"net/http"
)

// ヘルパー: JSONレスポンス
func JsonResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ヘルパー: エラーレスポンス
func ErrorResponse(w http.ResponseWriter, status int, message string) {
	JsonResponse(w, status, map[string]string{"error": message})
}
