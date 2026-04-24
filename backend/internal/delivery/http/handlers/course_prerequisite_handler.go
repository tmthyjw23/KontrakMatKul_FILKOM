package handlers

import (
	"net/http"

	"sistemkontrakmatkul/backend/internal/domain/models"
	"sistemkontrakmatkul/backend/internal/usecase"

	"github.com/gin-gonic/gin"
)

type CoursePrerequisiteHandler struct {
	usecase usecase.CoursePrerequisiteUsecase
}

func NewCoursePrerequisiteHandler(uc usecase.CoursePrerequisiteUsecase) *CoursePrerequisiteHandler {
	return &CoursePrerequisiteHandler{usecase: uc}
}

func (h *CoursePrerequisiteHandler) Add(c *gin.Context) {
	var req models.CoursePrerequisiteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.usecase.AddPrerequisite(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "prerequisite added successfully"})
}

func (h *CoursePrerequisiteHandler) Remove(c *gin.Context) {
	// Simplified for now, typically uses path params
	var req models.CoursePrerequisiteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.usecase.RemovePrerequisite(c.Request.Context(), req.CourseID, req.PrerequisiteCourseID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "prerequisite removed successfully"})
}

func (h *CoursePrerequisiteHandler) List(c *gin.Context) {
	// Logic to parse courseID and call usecase...
	c.JSON(http.StatusOK, gin.H{"message": "list of prerequisites"})
}