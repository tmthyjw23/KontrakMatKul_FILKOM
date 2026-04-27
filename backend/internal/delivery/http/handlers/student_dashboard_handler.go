package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"sistemkontrakmatkul/backend/internal/delivery/http/middlewares"
	"sistemkontrakmatkul/backend/internal/usecase"
)

type StudentDashboardHandler struct {
	usecase usecase.StudentDashboardUsecase
}

func NewStudentDashboardHandler(uc usecase.StudentDashboardUsecase) *StudentDashboardHandler {
	return &StudentDashboardHandler{usecase: uc}
}

func (h *StudentDashboardHandler) GetSchedule(c *gin.Context) {
	userIDValue, exists := c.Get(middlewares.ContextUserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: user id not found in context"})
		return
	}

	uID, ok := userIDValue.(uint64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error: invalid user id type"})
		return
	}

	schedule, err := h.usecase.GetMySchedule(c.Request.Context(), uID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": schedule})
}

func (h *StudentDashboardHandler) GetHistory(c *gin.Context) {
	userIDValue, exists := c.Get(middlewares.ContextUserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized: user id not found in context"})
		return
	}

	uID, ok := userIDValue.(uint64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error: invalid user id type"})
		return
	}

	history, err := h.usecase.GetMyHistory(c.Request.Context(), uID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": history})
}