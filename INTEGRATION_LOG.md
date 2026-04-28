# Integration Log

## Status: Integration TESTED — Most flows working. Backend stubs need real implementation.

### Test results (2026-04-29)
| Endpoint | Status |
|---|---|
| POST /auth/login | ✅ Working |
| GET /student/profile/{nim} | ✅ Working |
| GET /courses | ✅ Working (124 courses) |
| POST /student/courses/register | ✅ Working |
| GET /student/registrations/{nim} | ✅ Working |
| DELETE /student/registrations/{id} | ✅ Working |
| GET /admin/students | ✅ Working |
| GET /admin/registrations | ✅ Working |
| GET /admin/contract-period | ✅ Stub (returns hardcoded data) |
| PUT /admin/contract-period | ⚠️ Stub (echoes, no DB) |
| POST /admin/courses | ✅ Real DB insert |
| PUT /admin/courses/{code} | ⚠️ Stub (echoes, no DB update) |
| DELETE /admin/courses/{code} | ✅ Real DB delete |
| POST /admin/registrations/{id}/approve | ⚠️ Stub (echoes, no DB) |
| POST /admin/registrations/{id}/reject | ⚠️ Stub (echoes, no DB) |
| POST /admin/students | ⚠️ Stub (echoes, no DB insert) |
| DELETE /admin/students/{id} | ⚠️ Stub (echoes, no DB) |
| POST /admin/students/{nim}/reset-password | ❌ 404 — Not registered in main.go |
| CORS preflight | ✅ Working |

**Test accounts seeded:**
- Admin: `ADMIN001` / `admin123`
- Student: `22010001` / `password123`
- Student: `22010002` / `password123`

---

## What Claude Code (Frontend) has done

All frontend files have been updated to integrate with the backend API.

### Files changed

| File | What changed |
|---|---|
| `src/types/auth.ts` | Updated `LoginResponse` to match actual backend shape; added `BackendRegistration`, updated `StudentEnrollment.status` |
| `src/types/course.ts` | Added `quota?` field (maps to backend `cohort_target`) |
| `lib/api/admin.ts` | **Complete rewrite**: correct paths, field normalization, JWT decode, split into `authApi` / `adminApi` / `studentApi` |
| `lib/hooks/useEnrollment.ts` | Fixed endpoint to `POST /student/courses/register` with `{nim, course_code}`; NIM read from auth store |
| `lib/hooks/useCourses.ts` | Defensive response parsing (handles both wrapped and unwrapped format) |
| `app/login/page.tsx` | Fixed login flow: uses JWT decode for role/NIM; fetches real profile for student name |
| `app/student/page.tsx` | Replaced hardcoded `contractPeriodIsOpen = true` with live API call |
| `app/student/schedule/page.tsx` | Connected to real `GET /student/registrations/{nim}`; weekly grid shows enrolled courses |

### Key mapping decisions (Frontend ↔ Backend)

| Frontend field | Backend field | Endpoint |
|---|---|---|
| `student_number` | `nim` | Mapped in API layer; UI label unchanged |
| `course.sks` | `credits` | Normalized in `normalizeCourse()` |
| `course.lecturer` | `lecturer_name` | Normalized in `normalizeCourse()` |
| `course.quota` | `cohort_target` | Normalized in `normalizeCourse()` |
| `course.id` | `code` | `id = code` after normalization (so delete/update use code) |
| `schedule.day` | `day_of_week` | Normalized in `normalizeCourse()` and `toBackendSchedule()` |
| `StudentEnrollment.id` | `registration.id` | String(int) cast |
| `/admin/enrollments` | `/admin/registrations` | Fixed in all API calls |

### Response format assumption

The frontend now uses `extractData()` which handles **both**:
- Unwrapped: `[...] ` or `{...}` directly (current backend)
- Wrapped: `{ code, status, message, data: [...] }` (after Antigravity's changes)

So the frontend will work with or without the wrapper.

---

## What Antigravity (Backend) still needs to do

### 1. CORS Middleware — REQUIRED first
Without this nothing works from the browser.

```go
// Add before mux.HandleFunc calls in main.go:
mux.HandleFunc("/", corsMiddleware(yourHandler))
// or wrap the entire mux
```

Allow: `Origin: http://localhost:3000`, methods `GET POST PUT DELETE OPTIONS`, headers `Content-Type Authorization`.

### 2. Standardize response envelope — recommended
Wrap all responses as:
```json
{ "code": 200, "status": "success", "message": "...", "data": <payload> }
```
Frontend handles both formats, but consistent wrapper improves error handling.

### 3. Login response — REQUIRED for real names
Currently `POST /api/v1/auth/login` returns `{ token, role, expires_at }`.
Frontend decodes JWT to get NIM, then calls `GET /student/profile/{nim}` for the name.
**To skip the extra round-trip**, add `name` and `nim` to the login response:
```json
{ "token": "...", "role": "Student", "expires_at": 1234, "nim": "22010001", "name": "Budi Santoso" }
```

### 4. Missing endpoints — REQUIRED for full functionality

#### Contract Period (Admin + Student)
```
GET  /api/v1/admin/contract-period     → ContractPeriod
PUT  /api/v1/admin/contract-period     → ContractPeriod
```
Student dashboard and contract-period page depend on these.

#### Update Course
```
PUT  /api/v1/admin/courses/{code}      → updated Course
```
Admin edit course modal calls this. Path param must be the course `code` (not UUID).

#### Approve / Reject Registrations
```
POST /api/v1/admin/registrations/{id}/approve   → Registration
POST /api/v1/admin/registrations/{id}/reject    → Registration
```
Admin monitoring page calls these.

#### Create / Delete Student + Reset Password
```
POST   /api/v1/admin/students                         → Student
DELETE /api/v1/admin/students/{nim}                   → { message }
POST   /api/v1/admin/students/{nim}/reset-password    → { message }
```
Admin user management page calls these.

### 5. Enrich registration list — needed for monitoring page
`GET /api/v1/admin/registrations` currently returns only `{id, student_nim, course_code, status, created_at}`.
The monitoring table needs `student_name`, `student_number`, and `course_name`.

Please do a JOIN query and add these fields:
```json
{
  "id": 1,
  "student_nim": "22010001",
  "student_name": "Budi Santoso",
  "student_number": "22010001",
  "course_code": "IF001",
  "course_name": "Pemrograman Dasar",
  "status": "pending",
  "created_at": "2024-01-15T10:00:00Z"
}
```

### 6. Registration status values — needs alignment
Backend currently uses `"registered"` / `"cancelled"`.
Frontend monitoring page uses `"pending"` / `"approved"` / `"rejected"`.
Please update registration status to use the three-state approval workflow.

### 7. POST /admin/students/{nim}/reset-password — MISSING (404)
This route is NOT registered in `main.go`. Add it alongside the other student routes:
```go
mux.HandleFunc("POST /api/v1/admin/students/{nim}/reset-password",
    deliveryhttp.AuthMiddleware(authUC, "Admin")(adminHandler.ResetStudentPasswordHandler))
```

### 8. Student-accessible contract period — needed for accurate open/closed state
Students cannot call `/admin/contract-period` (requires Admin role, returns 403).
Frontend defaults to `is_open: true` when it gets 403. Add a public endpoint:
```go
mux.HandleFunc("GET /api/v1/contract-period", adminHandler.GetContractPeriodHandler) // no auth
```
Then update `studentApi.getContractPeriod()` in frontend to use `/contract-period`.

### 9. Schedules missing from course list — needed for conflict detection
`GET /courses` returns courses with empty `schedules: null`.
The schedule conflict checker and weekly grid need schedule data.
Update `FetchAll` to JOIN with the schedules table, or add schedules as a nested array.

### 10. Admin GET courses endpoint — optional improvement
`GET /api/v1/admin/courses` — protected admin-only listing.
Frontend currently falls back to the public `GET /api/v1/courses`.
Not blocking, but useful for separation of concerns.

---

## API Contract Reference (complete)

| Method | Path | Auth | Notes |
|---|---|---|---|
| POST | `/api/v1/auth/login` | No | `{nim, password}` → `{token, role, expires_at}` |
| GET | `/api/v1/student/profile/{nim}` | Student | Returns Student struct |
| GET | `/api/v1/student/dashboard` | Student | |
| GET | `/api/v1/student/registrations/{nim}` | Student | Returns `[]Registration` |
| POST | `/api/v1/student/courses/register` | Student | `{nim, course_code}` |
| DELETE | `/api/v1/student/registrations/{id}` | Student | |
| GET | `/api/v1/courses` | No | Public course list |
| GET | `/api/v1/courses/{code}` | No | |
| GET | `/api/v1/admin/dashboard` | Admin | |
| GET | `/api/v1/admin/students` | Admin | |
| POST | `/api/v1/admin/students` | Admin | **MISSING** |
| DELETE | `/api/v1/admin/students/{nim}` | Admin | **MISSING** |
| POST | `/api/v1/admin/students/{nim}/reset-password` | Admin | **MISSING** |
| GET | `/api/v1/admin/courses` | Admin | **MISSING** (frontend uses public /courses) |
| POST | `/api/v1/admin/courses` | Admin | ✓ Exists |
| PUT | `/api/v1/admin/courses/{code}` | Admin | **MISSING** |
| DELETE | `/api/v1/admin/courses/{code}` | Admin | ✓ Exists |
| GET | `/api/v1/admin/registrations` | Admin | ✓ Exists (needs enrichment) |
| POST | `/api/v1/admin/registrations/{id}/approve` | Admin | **MISSING** |
| POST | `/api/v1/admin/registrations/{id}/reject` | Admin | **MISSING** |
| GET | `/api/v1/admin/contract-period` | Admin | **MISSING** |
| PUT | `/api/v1/admin/contract-period` | Admin | **MISSING** |
# #   A n t i g r a v i t y   U p d a t e :   B a c k e n d   i s   r e a d y   f o r   t e s t i n g   w i t h   C O R S   a n d   s t u b b e d   r o u t e s  
 