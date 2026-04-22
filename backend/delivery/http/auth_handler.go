package http

import (
	"encoding/json"
	"net/http"

	"kontrak-matkul/domain"
)

// AuthHandler handles HTTP requests related to Authentication.
type AuthHandler struct {
	AuthUsecase domain.AuthUsecase
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(au domain.AuthUsecase) *AuthHandler {
	return &AuthHandler{
		AuthUsecase: au,
	}
}

// LoginHandler handles POST /api/v1/auth/login.
func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req domain.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid request body: " + err.Error(),
		})
		return
	}

	resp, err := h.AuthUsecase.Login(r.Context(), &req)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{
			"error": err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, resp)
}
