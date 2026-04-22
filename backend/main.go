package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq" // PostgreSQL driver — blank import registers the driver

	deliveryhttp "kontrak-matkul/delivery/http"
	"kontrak-matkul/repository"
	"kontrak-matkul/usecase"
)

// initDB opens a connection pool to the PostgreSQL database.
// The connection string is read from the DATABASE_URL environment variable.
// Example value: "postgres://user:password@localhost:5432/kontrak_matkul?sslmode=disable"
func initDB() *sql.DB {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("❌ DATABASE_URL environment variable is not set")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("❌ Failed to open database connection: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("❌ Failed to ping database: %v", err)
	}

	log.Println("✅ Database connection established")
	return db
}

func main() {
	// ---------------------------------------------------------
	// DATABASE
	// ---------------------------------------------------------
	db := initDB()
	defer db.Close()

	// ---------------------------------------------------------
	// DEPENDENCY INJECTION (Repository → Usecase → Handler)
	// ---------------------------------------------------------

	// Repositories (data layer)
	userRepo := repository.NewUserRepository(db)
	courseRepo := repository.NewCourseRepository(db)
	regRepo := repository.NewRegistrationRepository(db)

	// Usecases (business logic layer)
	userUC := usecase.NewUserUsecase(userRepo)
	courseUC := usecase.NewCourseUsecase(courseRepo)
	regUC := usecase.NewRegistrationUsecase(regRepo, courseRepo)

	// Handlers (delivery layer)
	adminHandler := deliveryhttp.NewAdminHandler(courseUC, userUC, regUC)
	studentHandler := deliveryhttp.NewStudentHandler(courseUC, userUC, regUC)

	// ---------------------------------------------------------
	// ROUTER — Go 1.22+ ServeMux with method-based routing
	// ---------------------------------------------------------
	mux := http.NewServeMux()

	// Admin Routes
	mux.HandleFunc("GET /api/v1/admin/dashboard",
		deliveryhttp.RoleMiddleware("Admin", adminHandler.DashboardHandler))
	mux.HandleFunc("GET /api/v1/admin/students",
		deliveryhttp.RoleMiddleware("Admin", adminHandler.GetAllStudentsHandler))
	mux.HandleFunc("POST /api/v1/admin/courses",
		deliveryhttp.RoleMiddleware("Admin", adminHandler.AddCourseHandler))
	mux.HandleFunc("DELETE /api/v1/admin/courses/{code}",
		deliveryhttp.RoleMiddleware("Admin", adminHandler.DeleteCourseHandler))
	mux.HandleFunc("GET /api/v1/admin/registrations",
		deliveryhttp.RoleMiddleware("Admin", adminHandler.GetAllRegistrationsHandler))

	// Student Routes
	mux.HandleFunc("GET /api/v1/student/dashboard",
		deliveryhttp.RoleMiddleware("Student", studentHandler.DashboardHandler))
	mux.HandleFunc("GET /api/v1/student/profile/{nim}",
		deliveryhttp.RoleMiddleware("Student", studentHandler.GetProfileHandler))
	mux.HandleFunc("POST /api/v1/student/courses/register",
		deliveryhttp.RoleMiddleware("Student", studentHandler.RegisterCourseHandler))
	mux.HandleFunc("GET /api/v1/student/registrations/{nim}",
		deliveryhttp.RoleMiddleware("Student", studentHandler.GetMyRegistrationsHandler))

	// Public Routes
	mux.HandleFunc("GET /api/v1/courses",
		studentHandler.GetAllCoursesHandler)
	mux.HandleFunc("GET /api/v1/courses/{code}",
		studentHandler.GetCourseDetailHandler)

	// ---------------------------------------------------------
	// SERVER
	// ---------------------------------------------------------
	port := ":8080"
	server := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	log.Printf("🚀 Server listening on http://localhost%s", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("❌ Server error: %v\n", err)
	}
}
