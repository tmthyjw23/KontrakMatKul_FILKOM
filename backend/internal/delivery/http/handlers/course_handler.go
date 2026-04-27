package handlers

import (
	"net/http"
	"strconv"

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

func (h *CourseHandler) Create(ctx *gin.Context) {
	var course models.Course
	if err := ctx.ShouldBindJSON(&course); err != nil {
		h.logger.Warn("invalid course payload", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"status":  "error",
			"message": "invalid request payload",
			"data":    nil,
		})
		return
	}

	if err := h.service.Create(ctx.Request.Context(), &course); err != nil {
		h.logger.Error("failed to create course", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"status":  "error",
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"status":  "success",
		"message": "course created successfully",
		"data":    course,
	})
}

func (h *CourseHandler) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"status":  "error",
			"message": "invalid course id",
			"data":    nil,
		})
		return
	}

	var course models.Course
	if err := ctx.ShouldBindJSON(&course); err != nil {
		h.logger.Warn("invalid course payload", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"status":  "error",
			"message": "invalid request payload",
			"data":    nil,
		})
		return
	}
	course.ID = id

	if err := h.service.Update(ctx.Request.Context(), &course); err != nil {
		h.logger.Error("failed to update course", zap.Uint64("course_id", id), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"status":  "error",
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"status":  "success",
		"message": "course updated successfully",
		"data":    course,
	})
}

func (h *CourseHandler) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"status":  "error",
			"message": "invalid course id",
			"data":    nil,
		})
		return
	}

	if err := h.service.Delete(ctx.Request.Context(), id); err != nil {
		h.logger.Error("failed to delete course", zap.Uint64("course_id", id), zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"status":  "error",
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"status":  "success",
		"message": "course deleted successfully",
		"data":    nil,
	})
}
