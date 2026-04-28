package http

import (
	"encoding/json"
	"errors"
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
		if errors.Is(err, domain.ErrMaxCreditsExceeded) {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		if errors.Is(err, domain.ErrScheduleConflict) {
			writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
			return
		}
		if errors.Is(err, domain.ErrPendingRegistration) {
			writeJSON(w, http.StatusForbidden, map[string]string{"error": err.Error()})
			return
		}
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
