package middleware

import "net/http"

// AllowedOrigins adalah daftar origin yang diizinkan.
var AllowedOrigins = []string{
	"http://localhost:3000",
	"http://127.0.0.1:3000",
}

// CORS membungkus handler dengan header CORS yang diperlukan agar
// frontend React (port 3000) dapat mengakses API ini.
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// Periksa apakah origin ada di daftar yang diizinkan
		allowed := false
		for _, o := range AllowedOrigins {
			if o == origin {
				allowed = true
				break
			}
		}

		if allowed {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "86400")

		// Preflight request – langsung balas 204 No Content
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
