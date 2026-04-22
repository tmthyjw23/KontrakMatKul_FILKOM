package http

import (
	"fmt"
	"net/http"
)

// RoleMiddleware is an HTTP middleware that enforces role-based access control.
// It reads the user's role from the "X-User-Role" request header and compares it
// against the required role for the route being accessed.
//
// NOTE: In a production system, the role should be extracted and validated from
// a signed JWT token, not a plain header. This is a placeholder for scaffolding.
func RoleMiddleware(requiredRole string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userRole := r.Header.Get("X-User-Role")

		if userRole != requiredRole {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w,
				fmt.Sprintf(`{"error": "Forbidden: requires '%s' role"}`, requiredRole),
				http.StatusForbidden,
			)
			return
		}

		next(w, r)
	}
}