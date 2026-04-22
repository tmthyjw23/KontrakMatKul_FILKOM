package http

import (
	"fmt"
	"net/http"

	"kontrak-matkul/domain"
)

// AdminHandler handles HTTP requests related to Admin operations.
type AdminHandler struct {
	CourseUsecase domain.CourseUsecase
	UserUsecase   domain.UserUsecase
}

// NewAdminHandler creates a new instance of AdminHandler.
func NewAdminHandler(cu domain.CourseUsecase, uu domain.UserUsecase) *AdminHandler {
	return &AdminHandler{
		CourseUsecase: cu,
		UserUsecase:   uu,
	}
}

func (h *AdminHandler) DashboardHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the Admin Dashboard!")
}

func (h *AdminHandler) AddCourseHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Course created successfully (Admin Only)")
}
