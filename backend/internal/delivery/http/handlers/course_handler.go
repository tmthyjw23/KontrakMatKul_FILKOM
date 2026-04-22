package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"sistemkontrakmatkul/backend/internal/domain/models"
	"sistemkontrakmatkul/backend/internal/domain/services"
)

type CourseHandler struct {
	service services.CourseService
	logger  *zap.Logger
}

func NewCourseHandler(service services.CourseService, logger *zap.Logger) *CourseHandler {
	if logger == nil {
		logger = zap.NewNop()
	}

	return &CourseHandler{
		service: service,
		logger:  logger,
	}
}

func (h *CourseHandler) ListCourses(ctx *gin.Context) {
	courses, err := h.service.ListCourses(ctx.Request.Context())
	if err != nil {
		h.logger.Error("failed to handle list courses", zap.Error(err))

		message := "internal server error"
		if err == models.ErrFailedToFetchCourses {
			message = err.Error()
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"status":  "error",
			"message": message,
			"data":    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"status":  "success",
		"message": "courses fetched successfully",
		"data":    courses,
	})
}
