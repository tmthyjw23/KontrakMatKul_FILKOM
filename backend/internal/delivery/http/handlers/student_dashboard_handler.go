package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"sistemkontrakmatkul/backend/internal/usecase"
)

type StudentDashboardHandler struct {
	usecase usecase.StudentDashboardUsecase
}

func NewStudentDashboardHandler(uc usecase.StudentDashboardUsecase) *StudentDashboardHandler {
	return &StudentDashboardHandler{usecase: uc}
}

func (h *StudentDashboardHandler) GetSchedule(c *gin.Context) {
	userID, _ := c.Get("userId") // Assuming JWT middleware sets this
	uID, _ := strconv.ParseUint(userID.(string), 10, 64)
	
	schedule, err := h.usecase.GetMySchedule(c.Request.Context(), uID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": schedule})
}

func (h *StudentDashboardHandler) GetHistory(c *gin.Context) {
	userID, _ := c.Get("userId")
	uID, _ := strconv.ParseUint(userID.(string), 10, 64)
	
	history, err := h.usecase.GetMyHistory(c.Request.Context(), uID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": history})
}