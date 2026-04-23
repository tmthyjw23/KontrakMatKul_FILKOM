package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// ─── Request / Response shapes ──────────────────────────────────────────────

type LoginRequest struct {
	Role     string `json:"role"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Message     string `json:"message"`
	Role        string `json:"role"`
	Username    string `json:"username"`
	BearerToken string `json:"bearer_token"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type ResetPasswordRequest struct {
	Email       string `json:"email"`
	NewPassword string `json:"new_password"`
}

// ─── Dummy credential store ──────────────────────────────────────────────────

type credential struct {
	Username string
	Password string
	Role     string
}

var credentials = []credential{
	{
		Username: "22010001@student.unklab.ac.id",
		Password: "password123",
		Role:     "student",
	},
	{
		Username: "admin@unklab.ac.id",
		Password: "admin123",
		Role:     "admin",
	},
}

// ─── In-memory password override store ───────────────────────────────────────
// Menyimpan password yang sudah di-reset selama server berjalan.

var (
	passwordOverrides = make(map[string]string)
	overrideMu        sync.RWMutex
)

// getPassword mengembalikan password saat ini untuk user tertentu.
// Jika user pernah reset, gunakan password dari override map.
func getPassword(username string) string {
	overrideMu.RLock()
	defer overrideMu.RUnlock()

	if pw, ok := passwordOverrides[username]; ok {
		return pw
	}
	// Fallback ke credential asli
	for _, c := range credentials {
		if c.Username == username {
			return c.Password
		}
	}
	return ""
}

// ─── JWT Claims ──────────────────────────────────────────────────────────────

type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// ─── Helper: write JSON response ─────────────────────────────────────────────

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

// ─── Login Handler ────────────────────────────────────────────────────────────

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, ErrorResponse{Message: "Method not allowed"})
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Message: "Payload tidak valid"})
		return
	}

	// Normalize inputs
	req.Role = strings.ToLower(strings.TrimSpace(req.Role))
	req.Username = strings.TrimSpace(req.Username)
	req.Password = strings.TrimSpace(req.Password)

	// ── WAJIB: Validasi domain email untuk role Student ──────────────────────
	if req.Role == "student" {
		if !strings.HasSuffix(req.Username, "@student.unklab.ac.id") {
			writeJSON(w, http.StatusBadRequest, ErrorResponse{
				Message: "Email Student harus menggunakan domain @student.unklab.ac.id",
			})
			return
		}
	}

	// ── Cari kredensial yang cocok (mendukung password yang sudah di-reset) ──
	var matched *credential
	for i, c := range credentials {
		currentPw := getPassword(c.Username)
		if c.Username == req.Username && currentPw == req.Password && c.Role == req.Role {
			matched = &credentials[i]
			break
		}
	}

	if matched == nil {
		writeJSON(w, http.StatusUnauthorized, ErrorResponse{
			Message: "Login Gagal: Periksa Email/Password Anda",
		})
		return
	}

	// ── Generate JWT ─────────────────────────────────────────────────────────
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "fallback-secret-change-in-production"
	}

	claims := Claims{
		Username: matched.Username,
		Role:     matched.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "kontrak-matkul-filkom",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Message: "Gagal membuat token"})
		return
	}

	// ── Simpan ke HttpOnly Cookie ─────────────────────────────────────────────
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // set true jika pakai HTTPS
		SameSite: http.SameSiteLaxMode,
		MaxAge:   86400, // 24 jam
	})

	// ── Kembalikan Bearer Token di body JSON ──────────────────────────────────
	writeJSON(w, http.StatusOK, LoginResponse{
		Message:     "Login berhasil",
		Role:        matched.Role,
		Username:    matched.Username,
		BearerToken: tokenString,
	})
}

// ─── Logout Handler ───────────────────────────────────────────────────────────

func Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})
	writeJSON(w, http.StatusOK, map[string]string{"message": "Logout berhasil"})
}

// ─── Reset Password Handler ──────────────────────────────────────────────────

func ResetPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, ErrorResponse{Message: "Method not allowed"})
		return
	}

	var req ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Message: "Payload tidak valid"})
		return
	}

	req.Email = strings.TrimSpace(req.Email)
	req.NewPassword = strings.TrimSpace(req.NewPassword)

	// Validasi input tidak boleh kosong
	if req.Email == "" || req.NewPassword == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Message: "Email dan password baru wajib diisi",
		})
		return
	}

	// Validasi panjang password minimal 6 karakter
	if len(req.NewPassword) < 6 {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Message: "Password baru minimal 6 karakter",
		})
		return
	}

	// Validasi domain email — harus @student.unklab.ac.id atau @unklab.ac.id
	isValidDomain := strings.HasSuffix(req.Email, "@student.unklab.ac.id") ||
		strings.HasSuffix(req.Email, "@unklab.ac.id")

	if !isValidDomain {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{
			Message: "Email harus menggunakan domain @student.unklab.ac.id atau @unklab.ac.id",
		})
		return
	}

	// Periksa apakah email ada di daftar credential
	found := false
	for _, c := range credentials {
		if c.Username == req.Email {
			found = true
			break
		}
	}

	if !found {
		writeJSON(w, http.StatusNotFound, ErrorResponse{
			Message: "Email tidak ditemukan dalam sistem",
		})
		return
	}

	// Simpan password baru ke in-memory override
	overrideMu.Lock()
	passwordOverrides[req.Email] = req.NewPassword
	overrideMu.Unlock()

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "Password berhasil diubah",
	})
}
