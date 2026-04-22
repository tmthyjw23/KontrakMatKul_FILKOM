# Project Progress: Sistem Kontrak Mata Kuliah

> Changelog for the `jordy-backend` branch.

---

## ✅ Step 4 — Concrete Repository, Usecase, DB Schema & Full DI Wiring
*Date: 2026-04-22*

### Overview
This step implements the two innermost layers of Clean Architecture — **Repository** (data access) and **Usecase** (business logic) — with concrete Go structs. It also introduces the `Registration` domain, the PostgreSQL schema migration, and fully wires the dependency injection chain in `main.go`. All `// TODO` stubs in the handlers are now resolved.

---

### New Files

#### `domain/registration.go`
- Added `Registration` struct with fields: `ID`, `StudentNIM`, `CourseCode`, `Status`, `CreatedAt`.
- Defined `RegistrationRepository` interface: `Create`, `GetByNIM`, `GetAll`, `Cancel`.
- Defined `RegistrationUsecase` interface: `RegisterCourse`, `GetRegistrationsByNIM`, `GetAllRegistrations`, `CancelRegistration`.

#### `repository/user_repository.go`
- Concrete `userRepository` backed by `*sql.DB`.
- Implements `GetByNIM` (single row scan) and `GetAll` (multi-row scan).
- Uses safe `$1` PostgreSQL placeholder syntax.

#### `repository/course_repository.go`
- Concrete `courseRepository` backed by `*sql.DB`.
- Implements `FetchAll`, `GetByCode`, `Create`, `Delete`.
- `Delete` checks `RowsAffected` to return a meaningful not-found error.

#### `repository/registration_repository.go`
- Concrete `registrationRepository` backed by `*sql.DB`.
- `Create` uses `RETURNING id, created_at` to populate the struct after insert.
- Implements `GetByNIM`, `GetAll`, `Cancel`.

#### `usecase/user_usecase.go`
- Concrete `userUsecase`; wraps `UserRepository`.
- Validates that NIM is non-empty before delegating.

#### `usecase/course_usecase.go`
- Concrete `courseUsecase`; wraps `CourseRepository`.
- Business rules: code/name required, credits must be positive.
- `DeleteCourse` fetches the course first to surface a clear not-found error.

#### `usecase/registration_usecase.go`
- Concrete `registrationUsecase`; wraps both `RegistrationRepository` and `CourseRepository`.
- **Business Rule 1:** Course must exist before registering.
- **Business Rule 2:** Prevents duplicate active registrations for the same student/course pair.

#### `docs/migrations/001_initial_schema.sql`
- PostgreSQL DDL for all 4 tables: `students`, `courses`, `schedules`, `registrations`.
- Foreign key constraints with `ON DELETE CASCADE`.
- `UNIQUE (student_nim, course_code)` constraint on `registrations` to enforce DB-level deduplication.

---

### Modified Files

#### `delivery/http/admin_handler.go`
- `NewAdminHandler` now accepts `RegistrationUsecase` as a 3rd argument.
- All `// TODO` stubs replaced with real usecase calls.
- Added `GetAllRegistrationsHandler` → `GET /api/v1/admin/registrations`.

#### `delivery/http/student_handler.go`
- `NewStudentHandler` now accepts `RegistrationUsecase` as a 3rd argument.
- All `// TODO` stubs replaced with real usecase calls.
- `RegisterCourseHandler` now expects `{ "nim": "...", "course_code": "..." }` body.
- Added `GetMyRegistrationsHandler` → `GET /api/v1/student/registrations/{nim}`.

#### `main.go`
- Added `initDB()` helper; reads `DATABASE_URL` env var and pings the connection.
- Full DI chain: `Repository → Usecase → Handler`.
- Added `github.com/lib/pq` as blank import for the PostgreSQL driver.
- Two new routes wired: `GET /admin/registrations`, `GET /student/registrations/{nim}`.

---

### Updated Route Table

| Method   | Path                                      | Handler                            | Role    |
|----------|-------------------------------------------|------------------------------------|---------|
| `GET`    | `/api/v1/admin/dashboard`                 | `AdminHandler.Dashboard`           | Admin   |
| `GET`    | `/api/v1/admin/students`                  | `AdminHandler.GetAllStudents`      | Admin   |
| `POST`   | `/api/v1/admin/courses`                   | `AdminHandler.AddCourse`           | Admin   |
| `DELETE` | `/api/v1/admin/courses/{code}`            | `AdminHandler.DeleteCourse`        | Admin   |
| `GET`    | `/api/v1/admin/registrations`             | `AdminHandler.GetAllRegistrations` | Admin   |
| `GET`    | `/api/v1/student/dashboard`               | `StudentHandler.Dashboard`         | Student |
| `GET`    | `/api/v1/student/profile/{nim}`           | `StudentHandler.GetProfile`        | Student |
| `POST`   | `/api/v1/student/courses/register`        | `StudentHandler.RegisterCourse`    | Student |
| `GET`    | `/api/v1/student/registrations/{nim}`     | `StudentHandler.GetMyRegistrations`| Student |
| `GET`    | `/api/v1/courses`                         | `StudentHandler.GetAllCourses`     | Public  |
| `GET`    | `/api/v1/courses/{code}`                  | `StudentHandler.GetCourseDetail`   | Public  |

---

### Next Steps
- [ ] Replace `X-User-Role` header simulation in `middleware.go` with real **JWT validation**.
- [ ] Add an `Admin` domain model (separate from `Student`) for admin authentication.
- [ ] Write a `GET /api/v1/student/cancels/{id}` endpoint for `CancelRegistration`.
- [ ] Add environment variable loading from a `.env` file (consider `os.LookupEnv` + a startup check).
- [ ] Write integration tests against a test database.
