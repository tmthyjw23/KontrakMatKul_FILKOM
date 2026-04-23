package routes

import (
	"net/http"

	"github.com/tmthyjw23/KontrakMatKul_FILKOM/backend/handlers"
)

// Register mendaftarkan semua route API ke mux yang diberikan.
func Register(mux *http.ServeMux) {
	mux.HandleFunc("/api/login", handlers.Login)
	mux.HandleFunc("/api/logout", handlers.Logout)
	mux.HandleFunc("/api/reset-password", handlers.ResetPassword)

	// Health-check sederhana
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok","service":"kontrak-matkul-backend"}`))
	})
}
