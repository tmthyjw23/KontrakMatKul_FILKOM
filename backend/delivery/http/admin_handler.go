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
