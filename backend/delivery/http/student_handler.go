package http

import (
	"fmt"
	"net/http"

	"kontrak-matkul/domain"
)

// StudentHandler handles HTTP requests related to Student operations.
type StudentHandler struct {
	CourseUsecase domain.CourseUsecase
	UserUsecase   domain.UserUsecase
}

// NewStudentHandler creates a new instance of StudentHandler.
func NewStudentHandler(cu domain.CourseUsecase, uu domain.UserUsecase) *StudentHandler {
	return &StudentHandler{
		CourseUsecase: cu,
		UserUsecase:   uu,
	}
}

func (h *StudentHandler) DashboardHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the Student Dashboard!")
}

func (h *StudentHandler) RegisterCourseHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Registered for course successfully (Student Only)")
}

func (h *StudentHandler) CourseListHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "List of Courses (Accessible to authorized users)")
}
