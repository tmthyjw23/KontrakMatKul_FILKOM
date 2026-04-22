package main

import (
	"fmt"
	"log"
	"net/http"

	deliveryhttp "kontrak-matkul/delivery/http"
)

// --- MIDDLEWARE ---

// RoleMiddleware is a basic middleware to check user roles based on Go standard library.
// In a real application, you would extract and validate the user's role from a JWT token or session.
func RoleMiddleware(requiredRole string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Placeholder: Simulate extracting a role from a header for demonstration purposes
		userRole := r.Header.Get("X-User-Role")

		if userRole != requiredRole {
			http.Error(w, fmt.Sprintf("Forbidden: Requires %s role", requiredRole), http.StatusForbidden)
			return
		}

		// Proceed to the next handler if the role matches
		next(w, r)
	}
}

func main() {
	// Initialize Go 1.22+ ServeMux
	mux := http.NewServeMux()

	// ---------------------------------------------------------
	// DEPENDENCY INJECTION (Stubs for now)
	// ---------------------------------------------------------
	
	// We inject nil for the usecases as they are not implemented yet.
	adminHandler := deliveryhttp.NewAdminHandler(nil, nil)
	studentHandler := deliveryhttp.NewStudentHandler(nil, nil)

	// ---------------------------------------------------------
	// API V1: ADMIN ROUTES
	// ---------------------------------------------------------
	
	// Wrapping handlers with Admin Role Middleware
	mux.HandleFunc("GET /api/v1/admin/dashboard", RoleMiddleware("Admin", adminHandler.DashboardHandler))
	
	// Example of another admin route (e.g., adding a course)
	mux.HandleFunc("POST /api/v1/admin/courses", RoleMiddleware("Admin", adminHandler.AddCourseHandler))

	// ---------------------------------------------------------
	// API V1: STUDENT ROUTES
	// ---------------------------------------------------------
	
	// Wrapping handlers with Student Role Middleware
	mux.HandleFunc("GET /api/v1/student/dashboard", RoleMiddleware("Student", studentHandler.DashboardHandler))

	// Example of another student route (e.g., registering for a course)
	mux.HandleFunc("POST /api/v1/student/courses/register", RoleMiddleware("Student", studentHandler.RegisterCourseHandler))

	// ---------------------------------------------------------
	// SHARED / PUBLIC ROUTES
	// ---------------------------------------------------------
	
	// A route accessible to both. For simplicity, attached to the student handler.
	mux.HandleFunc("GET /api/v1/courses", studentHandler.CourseListHandler)

	// Server Configuration
	port := ":8080"
	server := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	log.Printf("Server starting on port %s", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on %s: %v\n", port, err)
	}
}
