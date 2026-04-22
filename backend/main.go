package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/go-sql-driver/mysql" // MySQL driver

	deliveryhttp "kontrak-matkul/delivery/http"
	"kontrak-matkul/repository"
	"kontrak-matkul/usecase"
)

func init() {
	// Load variables from .env file into the environment if it exists.
	// We ignore the error because in production environments, variables
	// are typically injected directly without a .env file.
	_ = godotenv.Load()
}

func initDB() *sql.DB {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("❌ DATABASE_URL environment variable is not set")
	}

	db, err := sql.Open("mysql", dsn)
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
	authUC := usecase.NewAuthUsecase(userRepo) // Auth handles login & JWT

	// Handlers (delivery layer)
	adminHandler := deliveryhttp.NewAdminHandler(courseUC, userUC, regUC)
	studentHandler := deliveryhttp.NewStudentHandler(courseUC, userUC, regUC)
	authHandler := deliveryhttp.NewAuthHandler(authUC)

	// ---------------------------------------------------------
	// ROUTER — Go 1.22+ ServeMux with method-based routing
	// ---------------------------------------------------------
	mux := http.NewServeMux()

	// Auth Route
	mux.HandleFunc("POST /api/v1/auth/login", authHandler.LoginHandler)

	// Admin Routes (Protected by AuthMiddleware requiring "Admin" role)
	mux.HandleFunc("GET /api/v1/admin/dashboard",
		deliveryhttp.AuthMiddleware(authUC, "Admin")(adminHandler.DashboardHandler))
	mux.HandleFunc("GET /api/v1/admin/students",
		deliveryhttp.AuthMiddleware(authUC, "Admin")(adminHandler.GetAllStudentsHandler))
	mux.HandleFunc("POST /api/v1/admin/courses",
		deliveryhttp.AuthMiddleware(authUC, "Admin")(adminHandler.AddCourseHandler))
	mux.HandleFunc("DELETE /api/v1/admin/courses/{code}",
		deliveryhttp.AuthMiddleware(authUC, "Admin")(adminHandler.DeleteCourseHandler))
	mux.HandleFunc("GET /api/v1/admin/registrations",
		deliveryhttp.AuthMiddleware(authUC, "Admin")(adminHandler.GetAllRegistrationsHandler))

	// Student Routes (Protected by AuthMiddleware requiring "Student" role)
	mux.HandleFunc("GET /api/v1/student/dashboard",
		deliveryhttp.AuthMiddleware(authUC, "Student")(studentHandler.DashboardHandler))
	mux.HandleFunc("GET /api/v1/student/profile/{nim}",
		deliveryhttp.AuthMiddleware(authUC, "Student")(studentHandler.GetProfileHandler))
	mux.HandleFunc("POST /api/v1/student/courses/register",
		deliveryhttp.AuthMiddleware(authUC, "Student")(studentHandler.RegisterCourseHandler))
	mux.HandleFunc("GET /api/v1/student/registrations/{nim}",
		deliveryhttp.AuthMiddleware(authUC, "Student")(studentHandler.GetMyRegistrationsHandler))
	mux.HandleFunc("DELETE /api/v1/student/registrations/{id}",
		deliveryhttp.AuthMiddleware(authUC, "Student")(studentHandler.CancelRegistrationHandler))

	// Public Routes
	mux.HandleFunc("GET /api/v1/courses", studentHandler.GetAllCoursesHandler)
	mux.HandleFunc("GET /api/v1/courses/{code}", studentHandler.GetCourseDetailHandler)

	// ---------------------------------------------------------
	// SERVER
	// ---------------------------------------------------------
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("🚀 Server listening on http://localhost:%s", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("❌ Server error: %v\n", err)
	}
}
