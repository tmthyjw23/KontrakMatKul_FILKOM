# Project Progress: Sistem Kontrak Mata Kuliah

> Changelog for the `jordy-backend` branch.

---

## ✅ Step 5 — Real Authentication, JWT Middleware, & Environment Variables
*Date: 2026-04-23*

### Overview
This step replaces the temporary placeholder role authentication with a proper JWT (JSON Web Token) based authentication system. It includes password hashing via `bcrypt`, environment variable usage via `godotenv`, and a new `/api/v1/auth/login` endpoint to securely generate tokens for users.

---

### New Files

#### `domain/auth.go`
- Defined `LoginRequest` (NIM, Password) and `LoginResponse` (Token, Role, ExpiresAt).
- Added Custom `Claims` struct embedding `jwt.RegisteredClaims` to carry `NIM` and `Role` inside the token.
- Defined `AuthUsecase` interface (`Login` and `ValidateClaims`).

#### `usecase/auth_usecase.go`
- Concrete `authUsecase` which utilizes `bcrypt.CompareHashAndPassword` securely without exposing the password hash outside the DB layer.
- Implements `Login` to issue a JWT token signed securely by `HMAC-SHA256` using `JWT_SECRET`.
- Implements `ValidateClaims` to reliably decode and verify incoming Bearer tokens.

#### `delivery/http/auth_handler.go`
- Added the `AuthHandler` module processing the `POST /api/v1/auth/login` payload.
- Injected `AuthUsecase` specifically for handling user sessions and securely returning the JWT payload to the client.

#### `.env.example`
- Created a blueprint of `.env` configuration file to streamline setting up `PORT`, `JWT_SECRET`, and `DATABASE_URL`.

---

### Modified Files

#### `domain/user.go`
- Added `GetPasswordHashByNIM(ctx context.Context, nim string) (string, error)` to strictly abstract the bcrypt hash logic outside of the generic data struct. The hash is intentionally left out of the `Student` struct completely to prevent any accidental leakage through API JSON responses.

#### `repository/user_repository.go`
- Implemented `GetPasswordHashByNIM` fetching specifically the password string to cross-check in the Auth layer. 

#### `delivery/http/middleware.go`
- Entirely overhauled `AuthMiddleware` (now fully a builder `func() func()`).
- Parses, validates, and extracts standard `Authorization: Bearer <token>` data strings.
- Now securely propagates verified user `Claims` directly down into the `http.Request` Context utilizing the `UserContextKey` abstraction.

#### `delivery/http/student_handler.go`
- Resolved the pending Next Step for Cancellation endpoints by creating `CancelRegistrationHandler` mapped to `DELETE /api/v1/student/registrations/{id}` safely.

#### `docs/migrations/001_initial_schema.sql`
- Reconfigured the `students` DML instructions adding `password_hash VARCHAR(255) NOT NULL` strictly.

#### `main.go`
- Integrated `godotenv` to safely load initialization parameters automatically. 
- Integrated the new `AuthUsecase` and `AuthHandler`. 
- Overhauled and explicitly nested protected endpoints to utilize the updated `deliveryhttp.AuthMiddleware(authUC, "<Role>")` handlers structure systematically.

---

### Updated Route Table

| Method   | Path                                      | Handler                               | Required Role |
|----------|-------------------------------------------|---------------------------------------|---------------|
| `POST`   | `/api/v1/auth/login`                      | `AuthHandler.Login`                   | None / Public |
| `GET`    | `/api/v1/courses`                         | `StudentHandler.GetAllCourses`        | None / Public |
| `GET`    | `/api/v1/courses/{code}`                  | `StudentHandler.GetCourseDetail`      | None / Public |
| `GET`    | `/api/v1/admin/dashboard`                 | `AdminHandler.Dashboard`              | Admin         |
| `GET`    | `/api/v1/admin/students`                  | `AdminHandler.GetAllStudents`         | Admin         |
| `POST`   | `/api/v1/admin/courses`                   | `AdminHandler.AddCourse`              | Admin         |
| `DELETE` | `/api/v1/admin/courses/{code}`            | `AdminHandler.DeleteCourse`           | Admin         |
| `GET`    | `/api/v1/admin/registrations`             | `AdminHandler.GetAllRegistrations`    | Admin         |
| `GET`    | `/api/v1/student/dashboard`               | `StudentHandler.Dashboard`            | Student       |
| `GET`    | `/api/v1/student/profile/{nim}`           | `StudentHandler.GetProfile`           | Student       |
| `POST`   | `/api/v1/student/courses/register`        | `StudentHandler.RegisterCourse`       | Student       |
| `GET`    | `/api/v1/student/registrations/{nim}`     | `StudentHandler.GetMyRegistrations`   | Student       |
| `DELETE` | `/api/v1/student/registrations/{id}`      | `StudentHandler.CancelRegistration`   | Student       |

---

### Next Steps 🚀
- [ ] Incorporate comprehensive Unit/Integration Testing (e.g., `_test.go` files).
- [ ] Add an `Admin` domain model (separate from `Student`) for admin system differentiation if required.
- [ ] Build a database Seeder script to initialize test/admin data.
- [ ] Secure CORS routing headers for Frontend connectivity integration.
