# Project Documentation: KontrakMatKul_FILKOM

## 1. Project Overview

- **Project Name:** KontrakMatKul_FILKOM
- **Description:** A robust backend service designed to handle student course registration (Kontrak Mata Kuliah). It supports both Student and Admin roles with distinct permissions. It enables administrators to manage course offerings and view all student registrations, while allowing students to view courses, register for them, and manage their individual registration statuses safely.
- **Tech Stack:** Go (Golang 1.25.6), MySQL (`go-sql-driver/mysql`), JWT (`golang-jwt/jwt/v5`), Godotenv (`joho/godotenv`), bcrypt (`golang.org/x/crypto`).
- **Architecture Pattern:** Clean Architecture (Domain, Repository, Usecase, Delivery/HTTP layers).
- **Entry Point:** `main.go`

---

## 2. Project Structure

```text
/backend
├── .env.example
├── go.mod
├── main.go
├── /delivery
│   └── /http
│       ├── admin_handler.go
│       ├── auth_handler.go
│       ├── middleware.go
│       └── student_handler.go
├── /domain
│   ├── auth.go
│   ├── course.go
│   ├── registration.go
│   └── user.go
├── /repository
│   ├── course_repository.go
│   ├── registration_repository.go
│   └── user_repository.go
└── /usecase
    ├── auth_usecase.go
    ├── course_usecase.go
    ├── registration_usecase.go
    └── user_usecase.go
```

---

## 3. What Was Built (Summary of All Components)

### Main Entrypoint
- **File(s):** `main.go`
- **Purpose:** Initializes the database connection, resolves dependencies, wires Clean Architecture layers together, configures HTTP routes using Go 1.22+ `ServeMux`, and starts the server.
- **Key Functions / Methods:**
  - `init()` — Loads `.env` file upon startup.
  - `initDB() *sql.DB` — Validates environment configuration and creates the database connection pool.
  - `main()` — The executable entry point for the REST API server.
- **Dependencies:** `database/sql`, standard library `net/http`, `joho/godotenv`.
- **Notes:** Leverages the method-based routing recently introduced in the Go standard library, eliminating the need for a third-party router like Chi or Gorilla.

### Domain Layer: Auth
- **File(s):** `domain/auth.go`
- **Purpose:** Defines requests, responses, and interfaces related to authentication and session generation.
- **Key Functions / Methods:** 
  - `AuthUsecase` interface containing `Login` and `ValidateClaims`.
- **Dependencies:** `golang-jwt/jwt/v5`.
- **Notes:** Claims structurally combine custom user fields (NIM, Role) with registered JWT claims for expiration logic.

### Domain Layer: Course
- **File(s):** `domain/course.go`
- **Purpose:** Declares the `Course` model and standardizes data/business logic contracts for course-related operations.
- **Key Functions / Methods:** Includes `CourseRepository` and `CourseUsecase` interfaces (e.g., `FetchAll`, `GetByCode`, `Create`, `Delete`).
- **Dependencies:** Standard library `context`.
- **Notes:** Supports metadata like credit hours (`credits`) and allowed quota (`cohort_target`).

### Domain Layer: Registration
- **File(s):** `domain/registration.go`
- **Purpose:** Outlines the concept of a student tying themselves to an available course code with an adjustable status state.
- **Key Functions / Methods:** Includes `RegistrationRepository` and `RegistrationUsecase` interfaces (e.g., `RegisterCourse`, `CancelRegistration`).
- **Dependencies:** Standard library `context`.
- **Notes:** Maintains status tracking ("registered", "cancelled") as part of its core model definition.

### Domain Layer: User (Student)
- **File(s):** `domain/user.go`
- **Purpose:** Governs the `Student` identity definitions and interfaces within the system.
- **Key Functions / Methods:** Includes `UserRepository` and `UserUsecase` interfaces.
- **Dependencies:** Standard library `context`.
- **Notes:** Clearly decouples returning user data from returning sensitive password hashes via separate interface methods.

### Repository Layer: Course Repository
- **File(s):** `repository/course_repository.go`
- **Purpose:** Implements `domain.CourseRepository` abstract interface using MySQL storage infrastructure.
- **Key Functions / Methods:**
  - `FetchAll(ctx context.Context) ([]domain.Course, error)` — Retrieves all available courses.
  - `Create(ctx context.Context, course *domain.Course) error` — Executes raw SQL INSERT.
- **Dependencies:** `database/sql`, `kontrak-matkul/domain`.
- **Notes:** Includes standard context-aware query mechanisms.

### Repository Layer: Registration Repository
- **File(s):** `repository/registration_repository.go`
- **Purpose:** Implements data insertion/retrieval specific to student course registries with MySQL.
- **Key Functions / Methods:**
  - `Create(ctx context.Context, reg *domain.Registration) error` — Creates a registry entry utilizing standard Insert with dynamic `LastInsertId()`.
  - `Cancel(ctx context.Context, id int) error` — Soft deletes using an `UPDATE` status modifier.
- **Dependencies:** `database/sql`, `kontrak-matkul/domain`.
- **Notes:** Returns IDs immediately generated by MySQL during the creation process without complex sequence queries.

### Repository Layer: User Repository
- **File(s):** `repository/user_repository.go`
- **Purpose:** Implements storage fetching methods associated with students/users against a MySQL backend.
- **Key Functions / Methods:**
  - `GetByNIM(ctx context.Context, nim string) (*domain.Student, error)` — Grabs student general columns.
  - `GetPasswordHashByNIM(ctx context.Context, nim string) (string, error)` — Used solely by the auth logic stack.
- **Dependencies:** `database/sql`, `kontrak-matkul/domain`.

### Usecase Layer: Auth Usecase
- **File(s):** `usecase/auth_usecase.go`
- **Purpose:** Executes concrete business logic encompassing login processing, hash validation, and JWT issuing protocols.
- **Key Functions / Methods:**
  - `Login(ctx context.Context, req *domain.LoginRequest) (*domain.LoginResponse, error)` — Retrieves DB data, runs bcrypt checks, outputs a generated signed token.
  - `ValidateClaims(tokenString string) (*domain.Claims, error)` — Cryptographically verifies returning access tokens.
- **Dependencies:** `golang.org/x/crypto/bcrypt`, `golang-jwt/jwt/v5`, `os`.
- **Notes:** Handles cryptographic validation ensuring DB hashes correctly correlate with plaintext inputs securely.

### Usecase Layer: Course Usecase
- **File(s):** `usecase/course_usecase.go`
- **Purpose:** Bridges course repository calls while inserting pre-verification validation strategies.
- **Key Functions / Methods:**
  - `CreateCourse(ctx context.Context, course *domain.Course) error` — Fails early if code or credits are incorrectly sized.
- **Dependencies:** `kontrak-matkul/domain`.

### Usecase Layer: Registration Usecase
- **File(s):** `usecase/registration_usecase.go`
- **Purpose:** Orchestrates multi-step processes concerning student course enrollments.
- **Key Functions / Methods:**
  - `RegisterCourse(ctx context.Context, nim, courseCode string) (*domain.Registration, error)` — Combines registry checks guaranteeing no duplicate registries and confirming courses fundamentally exist.
- **Dependencies:** `domain.CourseRepository`, `domain.RegistrationRepository`.
- **Notes:** Critical junction enforcing pure application constraints before contacting the database.

### Usecase Layer: User Usecase
- **File(s):** `usecase/user_usecase.go`
- **Purpose:** Executes business actions associated to manipulating student profiles.
- **Key Functions / Methods:**
  - `GetProfile(ctx context.Context, nim string) (*domain.Student, error)`
- **Dependencies:** `kontrak-matkul/domain`.

### Delivery Layer: Authentication Handler
- **File(s):** `delivery/http/auth_handler.go`
- **Purpose:** Web interface dealing with client JSON payload translation targeting login events.
- **Key Functions / Methods:**
  - `LoginHandler(w http.ResponseWriter, r *http.Request)` — Yields JWT payload format structure correctly.
- **Dependencies:** `encoding/json`, `net/http`.

### Delivery Layer: Admin Handler
- **File(s):** `delivery/http/admin_handler.go`
- **Purpose:** API controller restricting access interactions to administrative profiles exclusively.
- **Key Functions / Methods:**
  - `DashboardHandler`, `AddCourseHandler`, `DeleteCourseHandler`, `GetAllStudentsHandler`
- **Dependencies:** `encoding/json`, `net/http`.

### Delivery Layer: Student Handler
- **File(s):** `delivery/http/student_handler.go`
- **Purpose:** Exposes actions catered exactly around student portal experiences.
- **Key Functions / Methods:**
  - `RegisterCourseHandler(w http.ResponseWriter, r *http.Request)` — Triggers the orchestration for adding a class item.
  - `GetMyRegistrationsHandler` and `CancelRegistrationHandler`.
- **Dependencies:** `encoding/json`, `net/http`.

### Delivery Layer: HTTP Middleware
- **File(s):** `delivery/http/middleware.go`
- **Purpose:** Protects API endpoints directly validating HTTP headers and resolving authorization contexts.
- **Key Functions / Methods:**
  - `AuthMiddleware(authUC domain.AuthUsecase, requiredRole string) func(http.HandlerFunc) http.HandlerFunc`
- **Dependencies:** `net/http`, `context`, `strings`.
- **Notes:** Generates a HTTP `Context` append associating the request strictly with an established User identity matching required system rules.

---

## 4. Full Context: Source Code

### `main.go`
```go
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
```

### `domain/auth.go`
```go
package domain

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// LoginRequest is the payload expected by the login endpoint.
type LoginRequest struct {
	NIM      string `json:"nim"`
	Password string `json:"password"`
}

// LoginResponse is returned on successful authentication.
type LoginResponse struct {
	Token     string `json:"token"`
	Role      string `json:"role"`
	ExpiresAt int64  `json:"expires_at"` // Unix timestamp
}

// Claims embeds the standard JWT registered claims and adds our custom fields.
// This struct is used both for signing tokens and for parsing them in middleware.
type Claims struct {
	NIM  string `json:"nim"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

// AuthUsecase defines the contract for authentication business logic.
type AuthUsecase interface {
	// Login validates credentials and returns a signed JWT token on success.
	Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)

	// ValidateClaims parses a raw JWT string and returns the embedded claims.
	ValidateClaims(tokenString string) (*Claims, error)
}

// TokenExpiry is the lifetime of a generated JWT token.
const TokenExpiry = 24 * time.Hour
```

### `domain/course.go`
```go
package domain

import "context"

// Course represents a subject (mata kuliah) available for registration.
type Course struct {
	Code         string `json:"code"          db:"code"`
	Name         string `json:"name"          db:"name"`
	Class        string `json:"class"         db:"class"`
	LecturerName string `json:"lecturer_name" db:"lecturer_name"`
	Credits      int    `json:"credits"       db:"credits"`
	CohortTarget int    `json:"cohort_target" db:"cohort_target"`
}

// Schedule represents the time and room allocation for a given course.
type Schedule struct {
	CourseCode      string `json:"course_code"      db:"course_code"`
	DayOfWeek       string `json:"day_of_week"      db:"day_of_week"`
	StartTime       string `json:"start_time"       db:"start_time"`
	EndTime         string `json:"end_time"         db:"end_time"`
	Room            string `json:"room"             db:"room"`
	AdditionalNotes string `json:"additional_notes" db:"additional_notes"`
}

// CourseRepository defines the contract for Course data storage operations.
// Concrete implementations will reside in the repository layer.
type CourseRepository interface {
	FetchAll(ctx context.Context) ([]Course, error)
	GetByCode(ctx context.Context, code string) (*Course, error)
	Create(ctx context.Context, course *Course) error
	Delete(ctx context.Context, code string) error
}

// CourseUsecase defines the contract for Course business logic operations.
// Concrete implementations will reside in the usecase layer.
type CourseUsecase interface {
	FetchAllCourses(ctx context.Context) ([]Course, error)
	GetCourseDetails(ctx context.Context, code string) (*Course, error)
	CreateCourse(ctx context.Context, course *Course) error
	DeleteCourse(ctx context.Context, code string) error
}
```

### `domain/registration.go`
```go
package domain

import "context"

// Registration represents a student's course registration record.
type Registration struct {
	ID         int    `json:"id"          db:"id"`
	StudentNIM string `json:"student_nim" db:"student_nim"`
	CourseCode string `json:"course_code" db:"course_code"`
	Status     string `json:"status"      db:"status"` // e.g., "registered", "cancelled"
	CreatedAt  string `json:"created_at"  db:"created_at"`
}

// RegistrationRepository defines the contract for Registration data operations.
type RegistrationRepository interface {
	Create(ctx context.Context, reg *Registration) error
	GetByNIM(ctx context.Context, nim string) ([]Registration, error)
	GetAll(ctx context.Context) ([]Registration, error)
	Cancel(ctx context.Context, id int) error
}

// RegistrationUsecase defines the contract for Registration business logic.
type RegistrationUsecase interface {
	RegisterCourse(ctx context.Context, nim, courseCode string) (*Registration, error)
	GetRegistrationsByNIM(ctx context.Context, nim string) ([]Registration, error)
	GetAllRegistrations(ctx context.Context) ([]Registration, error)
	CancelRegistration(ctx context.Context, id int) error
}
```

### `domain/user.go`
```go
package domain

import "context"

// Student represents a registered student user in the system.
type Student struct {
	NIM          string `json:"nim"           db:"nim"`
	Name         string `json:"name"          db:"name"`
	Faculty      string `json:"faculty"       db:"faculty"`
	StudyProgram string `json:"study_program" db:"study_program"`
	CohortYear   int    `json:"cohort_year"   db:"cohort_year"`
	Role         string `json:"role"          db:"role"`
}

// UserRepository defines the contract for User data storage operations.
// Concrete implementations will reside in the repository layer.
type UserRepository interface {
	GetByNIM(ctx context.Context, nim string) (*Student, error)
	GetAll(ctx context.Context) ([]Student, error)
	// GetPasswordHashByNIM returns only the bcrypt password hash for a given NIM.
	// The hash is intentionally NOT a field on the Student struct to prevent
	// it from being accidentally serialized into JSON API responses.
	GetPasswordHashByNIM(ctx context.Context, nim string) (string, error)
}

// UserUsecase defines the contract for User business logic operations.
// Concrete implementations will reside in the usecase layer.
type UserUsecase interface {
	GetProfile(ctx context.Context, nim string) (*Student, error)
	GetAllStudents(ctx context.Context) ([]Student, error)
}
```

### `repository/course_repository.go`
```go
package repository

import (
	"context"
	"database/sql"
	"fmt"

	"kontrak-matkul/domain"
)

// courseRepository is the concrete implementation of domain.CourseRepository.
type courseRepository struct {
	db *sql.DB
}

// NewCourseRepository creates and returns a new courseRepository.
func NewCourseRepository(db *sql.DB) domain.CourseRepository {
	return &courseRepository{db: db}
}

// FetchAll retrieves all courses from the database.
func (r *courseRepository) FetchAll(ctx context.Context) ([]domain.Course, error) {
	query := `
		SELECT code, name, class, lecturer_name, credits, cohort_target
		FROM courses
		ORDER BY code ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error fetching all courses: %w", err)
	}
	defer rows.Close()

	var courses []domain.Course
	for rows.Next() {
		var c domain.Course
		if err := rows.Scan(
			&c.Code,
			&c.Name,
			&c.Class,
			&c.LecturerName,
			&c.Credits,
			&c.CohortTarget,
		); err != nil {
			return nil, fmt.Errorf("error scanning course row: %w", err)
		}
		courses = append(courses, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return courses, nil
}

// GetByCode retrieves a single course by its code.
func (r *courseRepository) GetByCode(ctx context.Context, code string) (*domain.Course, error) {
	query := `
		SELECT code, name, class, lecturer_name, credits, cohort_target
		FROM courses
		WHERE code = ?
	`

	row := r.db.QueryRowContext(ctx, query, code)

	var c domain.Course
	err := row.Scan(
		&c.Code,
		&c.Name,
		&c.Class,
		&c.LecturerName,
		&c.Credits,
		&c.CohortTarget,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("course with code %s not found", code)
		}
		return nil, fmt.Errorf("error fetching course: %w", err)
	}

	return &c, nil
}

// Create inserts a new course record into the database.
func (r *courseRepository) Create(ctx context.Context, course *domain.Course) error {
	query := `
		INSERT INTO courses (code, name, class, lecturer_name, credits, cohort_target)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query,
		course.Code,
		course.Name,
		course.Class,
		course.LecturerName,
		course.Credits,
		course.CohortTarget,
	)
	if err != nil {
		return fmt.Errorf("error creating course: %w", err)
	}

	return nil
}

// Delete removes a course record from the database by its code.
func (r *courseRepository) Delete(ctx context.Context, code string) error {
	query := `DELETE FROM courses WHERE code = ?`

	result, err := r.db.ExecContext(ctx, query, code)
	if err != nil {
		return fmt.Errorf("error deleting course: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("course with code %s not found", code)
	}

	return nil
}
```

### `repository/registration_repository.go`
```go
package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"kontrak-matkul/domain"
)

// registrationRepository is the concrete implementation of domain.RegistrationRepository.
type registrationRepository struct {
	db *sql.DB
}

// NewRegistrationRepository creates and returns a new registrationRepository.
func NewRegistrationRepository(db *sql.DB) domain.RegistrationRepository {
	return &registrationRepository{db: db}
}

// Create inserts a new registration record into the database.
// MySQL does not support RETURNING, so we use LastInsertId() to retrieve the new ID
// and set CreatedAt manually.
func (r *registrationRepository) Create(ctx context.Context, reg *domain.Registration) error {
	query := `
		INSERT INTO registrations (student_nim, course_code, status)
		VALUES (?, ?, 'registered')
	`

	result, err := r.db.ExecContext(ctx, query, reg.StudentNIM, reg.CourseCode)
	if err != nil {
		return fmt.Errorf("error creating registration: %w", err)
	}

	// MySQL equivalent of PostgreSQL's RETURNING id
	lastID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error retrieving last insert ID: %w", err)
	}

	reg.ID = int(lastID)
	reg.Status = "registered"
	reg.CreatedAt = time.Now().Format(time.DateTime)
	return nil
}

// GetByNIM retrieves all registrations for a given student NIM.
func (r *registrationRepository) GetByNIM(ctx context.Context, nim string) ([]domain.Registration, error) {
	query := `
		SELECT id, student_nim, course_code, status, created_at
		FROM registrations
		WHERE student_nim = ?
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, nim)
	if err != nil {
		return nil, fmt.Errorf("error fetching registrations for NIM %s: %w", nim, err)
	}
	defer rows.Close()

	var registrations []domain.Registration
	for rows.Next() {
		var reg domain.Registration
		if err := rows.Scan(
			&reg.ID,
			&reg.StudentNIM,
			&reg.CourseCode,
			&reg.Status,
			&reg.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning registration row: %w", err)
		}
		registrations = append(registrations, reg)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return registrations, nil
}

// GetAll retrieves every registration record (admin use).
func (r *registrationRepository) GetAll(ctx context.Context) ([]domain.Registration, error) {
	query := `
		SELECT id, student_nim, course_code, status, created_at
		FROM registrations
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error fetching all registrations: %w", err)
	}
	defer rows.Close()

	var registrations []domain.Registration
	for rows.Next() {
		var reg domain.Registration
		if err := rows.Scan(
			&reg.ID,
			&reg.StudentNIM,
			&reg.CourseCode,
			&reg.Status,
			&reg.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("error scanning registration row: %w", err)
		}
		registrations = append(registrations, reg)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return registrations, nil
}

// Cancel updates a registration's status to "cancelled".
func (r *registrationRepository) Cancel(ctx context.Context, id int) error {
	query := `UPDATE registrations SET status = 'cancelled' WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error cancelling registration: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("registration with id %d not found", id)
	}

	return nil
}
```

### `repository/user_repository.go`
```go
package repository

import (
	"context"
	"database/sql"
	"fmt"

	"kontrak-matkul/domain"
)

// userRepository is the concrete implementation of domain.UserRepository.
type userRepository struct {
	db *sql.DB
}

// NewUserRepository creates and returns a new userRepository.
func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &userRepository{db: db}
}

// GetByNIM retrieves a single student record by their NIM.
func (r *userRepository) GetByNIM(ctx context.Context, nim string) (*domain.Student, error) {
	query := `
		SELECT nim, name, faculty, study_program, cohort_year, role
		FROM students
		WHERE nim = ?
	`

	row := r.db.QueryRowContext(ctx, query, nim)

	var s domain.Student
	err := row.Scan(
		&s.NIM,
		&s.Name,
		&s.Faculty,
		&s.StudyProgram,
		&s.CohortYear,
		&s.Role,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("student with NIM %s not found", nim)
		}
		return nil, fmt.Errorf("error fetching student: %w", err)
	}

	return &s, nil
}

// GetPasswordHashByNIM retrieves only the password hash for a given student.
func (r *userRepository) GetPasswordHashByNIM(ctx context.Context, nim string) (string, error) {
	query := `SELECT password_hash FROM students WHERE nim = ?`
	
	row := r.db.QueryRowContext(ctx, query, nim)
	var hash string
	if err := row.Scan(&hash); err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("student with NIM %s not found", nim)
		}
		return "", fmt.Errorf("error fetching password hash: %w", err)
	}
	
	return hash, nil
}

// GetAll retrieves all student records from the database.
func (r *userRepository) GetAll(ctx context.Context) ([]domain.Student, error) {
	query := `
		SELECT nim, name, faculty, study_program, cohort_year, role
		FROM students
		ORDER BY cohort_year DESC, name ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error fetching all students: %w", err)
	}
	defer rows.Close()

	var students []domain.Student
	for rows.Next() {
		var s domain.Student
		if err := rows.Scan(
			&s.NIM,
			&s.Name,
			&s.Faculty,
			&s.StudyProgram,
			&s.CohortYear,
			&s.Role,
		); err != nil {
			return nil, fmt.Errorf("error scanning student row: %w", err)
		}
		students = append(students, s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return students, nil
}
```

### `usecase/auth_usecase.go`
```go
package usecase

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"kontrak-matkul/domain"
)

// authUsecase is the concrete implementation of domain.AuthUsecase.
type authUsecase struct {
	userRepo domain.UserRepository
}

// NewAuthUsecase creates and returns a new authUsecase.
func NewAuthUsecase(ur domain.UserRepository) domain.AuthUsecase {
	return &authUsecase{userRepo: ur}
}

func (u *authUsecase) getSecretKey() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// Fallback for development, should be strictly enforced in production
		return []byte("super-secret-key")
	}
	return []byte(secret)
}

// Login validates credentials and generates a JWT.
func (u *authUsecase) Login(ctx context.Context, req *domain.LoginRequest) (*domain.LoginResponse, error) {
	if req.NIM == "" || req.Password == "" {
		return nil, fmt.Errorf("nim and password are required")
	}

	// 1. Get the user's profile to extract the role
	user, err := u.userRepo.GetByNIM(ctx, req.NIM)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// 2. Fetch the stored password hash
	hash, err := u.userRepo.GetPasswordHashByNIM(ctx, req.NIM)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// 3. Compare the plaintext password with the hash
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// 4. Generate the JWT token
	expirationTime := time.Now().Add(domain.TokenExpiry)
	claims := &domain.Claims{
		NIM:  user.NIM,
		Role: user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(u.getSecretKey()) // ERROR_LINE
	if err != nil {
		return nil, fmt.Errorf("could not generate token: %w", err)
	}

	return &domain.LoginResponse{
		Token:     tokenString,
		Role:      user.Role,
		ExpiresAt: expirationTime.Unix(),
	}, nil
}

// ValidateClaims parses the raw JWT string and returns the claims.
func (u *authUsecase) ValidateClaims(tokenString string) (*domain.Claims, error) {
	claims := &domain.Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		// Ensure the signing method is what we expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return u.getSecretKey(), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
```

### `usecase/course_usecase.go`
```go
package usecase

import (
	"context"
	"fmt"

	"kontrak-matkul/domain"
)

// courseUsecase is the concrete implementation of domain.CourseUsecase.
type courseUsecase struct {
	courseRepo domain.CourseRepository
}

// NewCourseUsecase creates and returns a new courseUsecase.
func NewCourseUsecase(cr domain.CourseRepository) domain.CourseUsecase {
	return &courseUsecase{courseRepo: cr}
}

// FetchAllCourses retrieves all available courses.
func (u *courseUsecase) FetchAllCourses(ctx context.Context) ([]domain.Course, error) {
	courses, err := u.courseRepo.FetchAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("FetchAllCourses: %w", err)
	}

	return courses, nil
}

// GetCourseDetails retrieves a single course by code.
func (u *courseUsecase) GetCourseDetails(ctx context.Context, code string) (*domain.Course, error) {
	if code == "" {
		return nil, fmt.Errorf("course code cannot be empty")
	}

	course, err := u.courseRepo.GetByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("GetCourseDetails: %w", err)
	}

	return course, nil
}

// CreateCourse validates and persists a new course.
func (u *courseUsecase) CreateCourse(ctx context.Context, course *domain.Course) error {
	if course.Code == "" {
		return fmt.Errorf("course code is required")
	}
	if course.Name == "" {
		return fmt.Errorf("course name is required")
	}
	if course.Credits <= 0 {
		return fmt.Errorf("credits must be a positive number")
	}

	if err := u.courseRepo.Create(ctx, course); err != nil {
		return fmt.Errorf("CreateCourse: %w", err)
	}

	return nil
}

// DeleteCourse removes a course by code after validating it exists.
func (u *courseUsecase) DeleteCourse(ctx context.Context, code string) error {
	if code == "" {
		return fmt.Errorf("course code cannot be empty")
	}

	// Confirm course exists before attempting deletion
	if _, err := u.courseRepo.GetByCode(ctx, code); err != nil {
		return fmt.Errorf("DeleteCourse: %w", err)
	}

	if err := u.courseRepo.Delete(ctx, code); err != nil {
		return fmt.Errorf("DeleteCourse: %w", err)
	}

	return nil
}
```

### `usecase/registration_usecase.go`
```go
package usecase

import (
	"context"
	"fmt"

	"kontrak-matkul/domain"
)

// registrationUsecase is the concrete implementation of domain.RegistrationUsecase.
type registrationUsecase struct {
	regRepo    domain.RegistrationRepository
	courseRepo domain.CourseRepository
}

// NewRegistrationUsecase creates and returns a new registrationUsecase.
// It takes both a registration and course repository to enforce business rules
// (e.g., a student can only register for an existing course).
func NewRegistrationUsecase(rr domain.RegistrationRepository, cr domain.CourseRepository) domain.RegistrationUsecase {
	return &registrationUsecase{
		regRepo:    rr,
		courseRepo: cr,
	}
}

// RegisterCourse validates and creates a new registration entry.
func (u *registrationUsecase) RegisterCourse(ctx context.Context, nim, courseCode string) (*domain.Registration, error) {
	if nim == "" {
		return nil, fmt.Errorf("student NIM cannot be empty")
	}
	if courseCode == "" {
		return nil, fmt.Errorf("course code cannot be empty")
	}

	// Business Rule: Ensure the course exists before registering
	if _, err := u.courseRepo.GetByCode(ctx, courseCode); err != nil {
		return nil, fmt.Errorf("RegisterCourse: course not found: %w", err)
	}

	// Business Rule: Prevent duplicate registrations
	existing, err := u.regRepo.GetByNIM(ctx, nim)
	if err != nil {
		return nil, fmt.Errorf("RegisterCourse: error checking existing registrations: %w", err)
	}
	for _, r := range existing {
		if r.CourseCode == courseCode && r.Status == "registered" {
			return nil, fmt.Errorf("RegisterCourse: student %s is already registered for course %s", nim, courseCode)
		}
	}

	reg := &domain.Registration{
		StudentNIM: nim,
		CourseCode: courseCode,
	}

	if err := u.regRepo.Create(ctx, reg); err != nil {
		return nil, fmt.Errorf("RegisterCourse: %w", err)
	}

	return reg, nil
}

// GetRegistrationsByNIM fetches all registrations for a given student.
func (u *registrationUsecase) GetRegistrationsByNIM(ctx context.Context, nim string) ([]domain.Registration, error) {
	if nim == "" {
		return nil, fmt.Errorf("NIM cannot be empty")
	}

	return u.regRepo.GetByNIM(ctx, nim)
}

// GetAllRegistrations fetches every registration (admin only).
func (u *registrationUsecase) GetAllRegistrations(ctx context.Context) ([]domain.Registration, error) {
	return u.regRepo.GetAll(ctx)
}

// CancelRegistration cancels a registration by ID.
func (u *registrationUsecase) CancelRegistration(ctx context.Context, id int) error {
	if id <= 0 {
		return fmt.Errorf("invalid registration ID")
	}

	if err := u.regRepo.Cancel(ctx, id); err != nil {
		return fmt.Errorf("CancelRegistration: %w", err)
	}

	return nil
}
```

### `usecase/user_usecase.go`
```go
package usecase

import (
	"context"
	"fmt"

	"kontrak-matkul/domain"
)

// userUsecase is the concrete implementation of domain.UserUsecase.
type userUsecase struct {
	userRepo domain.UserRepository
}

// NewUserUsecase creates and returns a new userUsecase.
func NewUserUsecase(ur domain.UserRepository) domain.UserUsecase {
	return &userUsecase{userRepo: ur}
}

// GetProfile fetches a student's profile by NIM.
func (u *userUsecase) GetProfile(ctx context.Context, nim string) (*domain.Student, error) {
	if nim == "" {
		return nil, fmt.Errorf("NIM cannot be empty")
	}

	student, err := u.userRepo.GetByNIM(ctx, nim)
	if err != nil {
		return nil, fmt.Errorf("GetProfile: %w", err)
	}

	return student, nil
}

// GetAllStudents fetches all registered students.
func (u *userUsecase) GetAllStudents(ctx context.Context) ([]domain.Student, error) {
	students, err := u.userRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetAllStudents: %w", err)
	}

	return students, nil
}
```

### `delivery/http/admin_handler.go`
```go
package http

import (
	"encoding/json"
	"net/http"

	"kontrak-matkul/domain"
)

// AdminHandler handles HTTP requests for Admin-specific operations.
type AdminHandler struct {
	CourseUsecase       domain.CourseUsecase
	UserUsecase         domain.UserUsecase
	RegistrationUsecase domain.RegistrationUsecase
}

// NewAdminHandler creates a new AdminHandler with injected usecases.
func NewAdminHandler(cu domain.CourseUsecase, uu domain.UserUsecase, ru domain.RegistrationUsecase) *AdminHandler {
	return &AdminHandler{
		CourseUsecase:       cu,
		UserUsecase:         uu,
		RegistrationUsecase: ru,
	}
}

// writeJSON is a helper to write a JSON-encoded response with the given status code.
func writeJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}

// DashboardHandler handles GET /api/v1/admin/dashboard
func (h *AdminHandler) DashboardHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"message": "Welcome to the Admin Dashboard!",
		"role":    "Admin",
	})
}

// GetAllStudentsHandler handles GET /api/v1/admin/students
func (h *AdminHandler) GetAllStudentsHandler(w http.ResponseWriter, r *http.Request) {
	students, err := h.UserUsecase.GetAllStudents(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, students)
}

// AddCourseHandler handles POST /api/v1/admin/courses
func (h *AdminHandler) AddCourseHandler(w http.ResponseWriter, r *http.Request) {
	var course domain.Course
	if err := json.NewDecoder(r.Body).Decode(&course); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid request body: " + err.Error(),
		})
		return
	}

	if err := h.CourseUsecase.CreateCourse(r.Context(), &course); err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"message": "Course created successfully",
		"course":  course,
	})
}

// DeleteCourseHandler handles DELETE /api/v1/admin/courses/{code}
func (h *AdminHandler) DeleteCourseHandler(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	if code == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Course code is required"})
		return
	}

	if err := h.CourseUsecase.DeleteCourse(r.Context(), code); err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "Course deleted successfully",
		"code":    code,
	})
}

// GetAllRegistrationsHandler handles GET /api/v1/admin/registrations
func (h *AdminHandler) GetAllRegistrationsHandler(w http.ResponseWriter, r *http.Request) {
	registrations, err := h.RegistrationUsecase.GetAllRegistrations(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, registrations)
}
```

### `delivery/http/auth_handler.go`
```go
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
```

### `delivery/http/middleware.go`
```go
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
```

### `delivery/http/student_handler.go`
```go
package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"kontrak-matkul/domain"
)

// StudentHandler handles HTTP requests for Student-specific operations.
type StudentHandler struct {
	CourseUsecase       domain.CourseUsecase
	UserUsecase         domain.UserUsecase
	RegistrationUsecase domain.RegistrationUsecase
}

// NewStudentHandler creates a new StudentHandler with injected usecases.
func NewStudentHandler(cu domain.CourseUsecase, uu domain.UserUsecase, ru domain.RegistrationUsecase) *StudentHandler {
	return &StudentHandler{
		CourseUsecase:       cu,
		UserUsecase:         uu,
		RegistrationUsecase: ru,
	}
}

// DashboardHandler handles GET /api/v1/student/dashboard
func (h *StudentHandler) DashboardHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"message": "Welcome to the Student Dashboard!",
		"role":    "Student",
	})
}

// GetProfileHandler handles GET /api/v1/student/profile/{nim}
func (h *StudentHandler) GetProfileHandler(w http.ResponseWriter, r *http.Request) {
	nim := r.PathValue("nim")
	if nim == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "NIM is required"})
		return
	}

	student, err := h.UserUsecase.GetProfile(r.Context(), nim)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, student)
}

// GetAllCoursesHandler handles GET /api/v1/courses
func (h *StudentHandler) GetAllCoursesHandler(w http.ResponseWriter, r *http.Request) {
	courses, err := h.CourseUsecase.FetchAllCourses(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, courses)
}

// GetCourseDetailHandler handles GET /api/v1/courses/{code}
func (h *StudentHandler) GetCourseDetailHandler(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	if code == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Course code is required"})
		return
	}

	course, err := h.CourseUsecase.GetCourseDetails(r.Context(), code)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, course)
}

// RegisterCourseHandler handles POST /api/v1/student/courses/register
func (h *StudentHandler) RegisterCourseHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		NIM        string `json:"nim"`
		CourseCode string `json:"course_code"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid request body: " + err.Error(),
		})
		return
	}

	if body.NIM == "" || body.CourseCode == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Both nim and course_code fields are required",
		})
		return
	}

	registration, err := h.RegistrationUsecase.RegisterCourse(r.Context(), body.NIM, body.CourseCode)
	if err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"message":      "Course registration successful",
		"registration": registration,
	})
}

// GetMyRegistrationsHandler handles GET /api/v1/student/registrations/{nim}
func (h *StudentHandler) GetMyRegistrationsHandler(w http.ResponseWriter, r *http.Request) {
	nim := r.PathValue("nim")
	if nim == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "NIM is required"})
		return
	}

	registrations, err := h.RegistrationUsecase.GetRegistrationsByNIM(r.Context(), nim)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, registrations)
}

// CancelRegistrationHandler handles DELETE /api/v1/student/registrations/{id}
func (h *StudentHandler) CancelRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Registration ID is required"})
		return
	}

	var id int
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid registration ID format"})
		return
	}

	if err := h.RegistrationUsecase.CancelRegistration(r.Context(), id); err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "Registration cancelled successfully",
	})
}
```

---

## 5. Database / Data Models

- **Student:**
  - `nim` (string) — the primary identifier
  - `name` (string)
  - `faculty` (string)
  - `study_program` (string)
  - `cohort_year` (int)
  - `role` (string) — defines "Admin" or "Student"
  - `password_hash` (string) — (Stored within MySQL only, not natively serialized).

- **Course:**
  - `code` (string)
  - `name` (string)
  - `class` (string)
  - `lecturer_name` (string)
  - `credits` (int)
  - `cohort_target` (int)

- **Registration:**
  - `id` (int) — Auto Increment Primary ID
  - `student_nim` (string)
  - `course_code` (string)
  - `status` (string) — `registered` or `cancelled`
  - `created_at` (string) — Date format mapping to `DATETIME`.

---

## 6. API Endpoints

| Method | Endpoint | Handler | Description | Auth Required |
|--------|----------|---------|-------------|---------------|
| POST   | `/api/v1/auth/login` | `LoginHandler` | Retrieves access tokens for the system by validating credentials | No |
| GET    | `/api/v1/courses` | `GetAllCoursesHandler` | Fetches a generalized list of all courses | No |
| GET    | `/api/v1/courses/{code}` | `GetCourseDetailHandler` | Look up precise course info | No |
| GET    | `/api/v1/admin/dashboard` | `DashboardHandler` | Tests "Admin" validation successfully parsed by router | Yes (Admin) |
| GET    | `/api/v1/admin/students` | `GetAllStudentsHandler` | Pulls entire student directory matrix | Yes (Admin) |
| POST   | `/api/v1/admin/courses` | `AddCourseHandler` | Pushes a single `Course` insert | Yes (Admin) |
| DELETE | `/api/v1/admin/courses/{code}` | `DeleteCourseHandler` | Hard deletes a valid course | Yes (Admin) |
| GET    | `/api/v1/admin/registrations` | `GetAllRegistrationsHandler`| Fetches all universal registrations system-wide | Yes (Admin) |
| GET    | `/api/v1/student/dashboard` | `DashboardHandler` | Tests "Student" validation mechanism | Yes (Student) |
| GET    | `/api/v1/student/profile/{nim}`| `GetProfileHandler` | Looks up individual user data fields | Yes (Student) |
| POST   | `/api/v1/student/courses/register`| `RegisterCourseHandler` | Attempts a new course queue registry against NIM and code | Yes (Student) |
| GET    | `/api/v1/student/registrations/{nim}`| `GetMyRegistrationsHandler`| Obtains array of items queued/registered for specific student | Yes (Student) |
| DELETE | `/api/v1/student/registrations/{id}`| `CancelRegistrationHandler`| Sets previous registry into a "cancelled" database state | Yes (Student) |

---

## 7. Environment & Configuration

Environment Variables:
- `PORT`: (e.g., 8080) Defines HTTP listening port via `net/http` package.
- `JWT_SECRET`: Local or production signing salt cryptographic array.
- `DATABASE_URL`: Connection string.

**`.env.example`**
```env
# Application Configuration
PORT=8080
JWT_SECRET=super-secret-key-development-only

# Database Configuration
# Format: user:password@tcp(address:port)/dbname?params
# Ensure parseTime=true is present for DATETIME variables -> time.Time scanning
DATABASE_URL=root:root@tcp(127.0.0.1:3306)/kontrak_matkul?parseTime=true&charset=utf8mb4
```

---

## 8. How to Run

```bash
# Provide variables
cp .env.example .env

# Verify and download dependencies
go mod tidy

# Start executing backend
go run main.go
```

---

## 9. Dependencies

| Package | Version | Purpose |
|---------|---------|---------|
| `github.com/go-sql-driver/mysql` | v1.9.3 | Connects generic Go `database/sql` definitions to exact MySQL binary paths |
| `github.com/golang-jwt/jwt/v5` | v5.3.1 | Core crypto-suite used to mathematically verify and mint Bearer logins |
| `github.com/joho/godotenv` | v1.5.1 | Rapid string injection parser specifically targeting `.env` syntax |
| `golang.org/x/crypto` | v0.50.0 | Powers `bcrypt` libraries |

---

## 10. Known Issues / TODOs

- [ ] Need a defined database migration format (e.g., `schema.sql`) to immediately spin up tables for testers since `RETURNING` constraints are not standardized outside of setup boundaries.
- [ ] Incorporate pagination offsets directly into `FetchAll` functions within repositories due to the possibility of massive response matrices in a student/university database.
- [ ] Connect custom logging routines for production auditing mechanisms in `AuthMiddleware`.

---

## 11. Session Notes

- **Date:** April 23, 2026
- **Total files created:** 15 files scanned and unified.
- **Total lines of code (approx):** ~1,000 lines of active `.go` implementation (excluding metadata).
- **Summary of decisions made during this session:** Synthesized and exported the Clean Architecture setup into a cohesive, structured README documentation format that includes all layered logic from router delivery down to `net/http` server implementation, preserving `[domain, repository, usecase, delivery]` structure rigorously for developer hand-off.
