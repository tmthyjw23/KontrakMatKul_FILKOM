package http

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"kontrak-matkul/domain"
)

type contextKey string

// UserContextKey is the key used to store the JWT Claims in the request context.
const UserContextKey = contextKey("user_claims")

// AuthMiddleware creates a middleware that enforces role-based access control
// by validating a Bearer JWT token.
func AuthMiddleware(authUC domain.AuthUsecase, requiredRole string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				writeJSON(w, http.StatusUnauthorized, map[string]string{
					"error": "Missing or invalid Authorization header (Bearer token required)",
				})
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := authUC.ValidateClaims(tokenStr)
			if err != nil {
				writeJSON(w, http.StatusUnauthorized, map[string]string{
					"error": "Invalid token: " + err.Error(),
				})
				return
			}

			// Ensure the user role matches the required role (if a required role is specified)
			if requiredRole != "" && claims.Role != requiredRole {
				writeJSON(w, http.StatusForbidden, map[string]string{
					"error": fmt.Sprintf("Forbidden: requires '%s' role", requiredRole),
				})
				return
			}

			// Attach claims to the request Context so handlers can know WHO is making the request
			ctx := context.WithValue(r.Context(), UserContextKey, claims)
			next(w, r.WithContext(ctx))
		}
	}
}