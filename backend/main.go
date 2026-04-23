package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/tmthyjw23/KontrakMatKul_FILKOM/backend/middleware"
	"github.com/tmthyjw23/KontrakMatKul_FILKOM/backend/routes"
)

func main() {
	// ── Muat variabel lingkungan dari .env ───────────────────────────────────
	if err := godotenv.Load(); err != nil {
		log.Println("[WARN] .env tidak ditemukan, menggunakan variabel sistem")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// ── Daftarkan semua route ────────────────────────────────────────────────
	mux := http.NewServeMux()
	routes.Register(mux)

	// ── Terapkan middleware CORS ─────────────────────────────────────────────
	handler := middleware.CORS(mux)

	addr := fmt.Sprintf(":%s", port)
	log.Printf("[INFO] Server berjalan di http://localhost%s", addr)
	log.Printf("[INFO] Endpoint tersedia:")
	log.Printf("         POST   http://localhost%s/api/login", addr)
	log.Printf("         POST   http://localhost%s/api/logout", addr)
	log.Printf("         GET    http://localhost%s/api/health", addr)

	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("[FATAL] Gagal menjalankan server: %v", err)
	}
}
