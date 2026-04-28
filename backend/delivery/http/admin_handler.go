package http

import (
	"encoding/json"
	"net/http"
	"strconv"

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

// ----------------------------------------------------------------------
// STUB HANDLERS FOR FRONTEND INTEGRATION
// TODO: Implement full usecase/repository logic for these operations
// ----------------------------------------------------------------------

// UpdateCourseHandler handles PUT /api/v1/admin/courses/{code}
func (h *AdminHandler) UpdateCourseHandler(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	var course domain.Course
	if err := json.NewDecoder(r.Body).Decode(&course); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}
	course.Code = code
	// Stub response
	writeJSON(w, http.StatusOK, map[string]any{"message": "Course updated", "course": course})
}

// ApproveRegistrationHandler handles POST /api/v1/admin/registrations/{id}/approve
func (h *AdminHandler) ApproveRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid registration ID"})
		return
	}

	if err := h.RegistrationUsecase.ApproveRegistration(r.Context(), id); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Registration approved successfully"})
}

// RejectRegistrationHandler handles POST /api/v1/admin/registrations/{id}/reject
func (h *AdminHandler) RejectRegistrationHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid registration ID"})
		return
	}

	if err := h.RegistrationUsecase.RejectRegistration(r.Context(), id); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Registration rejected successfully"})
}

// CreateStudentHandler handles POST /api/v1/admin/students
func (h *AdminHandler) CreateStudentHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		NIM          string `json:"nim"`
		Name         string `json:"name"`
		Password     string `json:"password"`
		Faculty      string `json:"faculty"`
		StudyProgram string `json:"study_program"`
		CohortYear   int    `json:"cohort_year"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	student := &domain.Student{
		NIM:          body.NIM,
		Name:         body.Name,
		Faculty:      body.Faculty,
		StudyProgram: body.StudyProgram,
		CohortYear:   body.CohortYear,
	}

	if student.Faculty == "" {
		student.Faculty = "FILKOM"
	}
	if student.StudyProgram == "" {
		student.StudyProgram = "Sistem Informasi"
	}
	if student.CohortYear == 0 {
		student.CohortYear = 2026
	}

	if err := h.UserUsecase.CreateStudent(r.Context(), student, body.Password); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{"message": "Student created", "student": student})
}

// DeleteStudentHandler handles DELETE /api/v1/admin/students/{id}
func (h *AdminHandler) DeleteStudentHandler(w http.ResponseWriter, r *http.Request) {
	nim := r.PathValue("id")
	if err := h.UserUsecase.DeleteStudent(r.Context(), nim); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Student " + nim + " deleted successfully"})
}

// ResetStudentPasswordHandler handles POST /api/v1/admin/students/{nim}/reset-password
func (h *AdminHandler) ResetStudentPasswordHandler(w http.ResponseWriter, r *http.Request) {
	nim := r.PathValue("nim")
	var body struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	if err := h.UserUsecase.ResetPassword(r.Context(), nim, body.Password); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Password reset successful"})
}

// GetContractPeriodHandler handles GET /api/v1/admin/contract-period
func (h *AdminHandler) GetContractPeriodHandler(w http.ResponseWriter, r *http.Request) {
	// Stub response - returning a mocked contract period
	writeJSON(w, http.StatusOK, map[string]any{
		"is_open":    true,
		"start_date": "2026-01-01T00:00:00Z",
		"end_date":   "2026-12-31T23:59:59Z",
	})
}

// UpdateContractPeriodHandler handles PUT /api/v1/admin/contract-period
func (h *AdminHandler) UpdateContractPeriodHandler(w http.ResponseWriter, r *http.Request) {
	var period map[string]any
	if err := json.NewDecoder(r.Body).Decode(&period); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}
	// Stub response
	writeJSON(w, http.StatusOK, map[string]any{"message": "Contract period updated", "period": period})
}
