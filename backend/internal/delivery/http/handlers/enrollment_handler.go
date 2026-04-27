package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"sistemkontrakmatkul/backend/internal/delivery/http/middlewares"
	"sistemkontrakmatkul/backend/internal/domain/models"
	"sistemkontrakmatkul/backend/internal/domain/services"
)

type EnrollmentHandler struct {
	service services.EnrollmentService
	logger  *zap.Logger
}

type apiResponse struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewEnrollmentHandler(service services.EnrollmentService, logger *zap.Logger) *EnrollmentHandler {
	if logger == nil {
		logger = zap.NewNop()
	}

	return &EnrollmentHandler{
		service: service,
		logger:  logger,
	}
}

func (h *EnrollmentHandler) Enroll(ctx *gin.Context) {
	var request models.EnrollmentRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		h.logger.Warn("invalid enrollment payload", zap.Error(err))
		writeJSON(ctx, http.StatusBadRequest, "error", "invalid request payload", nil)
		return
	}

	userIDValue, exists := ctx.Get(middlewares.ContextUserIDKey)
	if !exists {
		h.logger.Warn("missing user id in request context")
		writeJSON(ctx, http.StatusUnauthorized, "error", "unauthorized", nil)
		return
	}

	userID, ok := userIDValue.(uint64)
	if !ok || userID == 0 {
		h.logger.Warn("invalid user id type in request context")
		writeJSON(ctx, http.StatusUnauthorized, "error", "unauthorized", nil)
		return
	}

	request.UserID = userID

	result, err := h.service.Enroll(ctx.Request.Context(), request)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrInvalidEnrollmentRequest),
			errors.Is(err, models.ErrUserNotFound),
			errors.Is(err, models.ErrCourseNotFound):
			h.logger.Warn("enrollment validation failed", zap.Error(err))
			writeJSON(ctx, http.StatusBadRequest, "error", err.Error(), nil)
			return

		case errors.Is(err, models.ErrScheduleConflict),
			errors.Is(err, models.ErrQuotaExceeded),
			errors.Is(err, models.ErrCreditLimitExceeded),
			errors.Is(err, models.ErrAlreadyEnrolled):
			h.logger.Warn("enrollment conflict", zap.Error(err))
			writeJSON(ctx, http.StatusConflict, "error", err.Error(), nil)
			return

		default:
			h.logger.Error("enrollment failed", zap.Error(err))
			writeJSON(ctx, http.StatusInternalServerError, "error", "internal server error", nil)
			return
		}
	}

	writeJSON(ctx, http.StatusOK, "success", "enrollment saved successfully", result)
}

func (h *EnrollmentHandler) ListAll(ctx *gin.Context) {
	enrollments, err := h.service.ListAll(ctx.Request.Context())
	if err != nil {
		h.logger.Error("failed to handle list all enrollments", zap.Error(err))
		writeJSON(ctx, http.StatusInternalServerError, "error", "failed to fetch enrollments", nil)
		return
	}

	writeJSON(ctx, http.StatusOK, "success", "enrollments fetched successfully", enrollments)
}

func (h *EnrollmentHandler) Approve(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		writeJSON(ctx, http.StatusBadRequest, "error", "invalid enrollment id", nil)
		return
	}

	if err := h.service.Approve(ctx.Request.Context(), id); err != nil {
		h.logger.Error("failed to approve enrollment", zap.Uint64("id", id), zap.Error(err))
		writeJSON(ctx, http.StatusInternalServerError, "error", "failed to approve enrollment", nil)
		return
	}

	writeJSON(ctx, http.StatusOK, "success", "enrollment approved successfully", nil)
}

func (h *EnrollmentHandler) Reject(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		writeJSON(ctx, http.StatusBadRequest, "error", "invalid enrollment id", nil)
		return
	}

	if err := h.service.Reject(ctx.Request.Context(), id); err != nil {
		h.logger.Error("failed to reject enrollment", zap.Uint64("id", id), zap.Error(err))
		writeJSON(ctx, http.StatusInternalServerError, "error", "failed to reject enrollment", nil)
		return
	}

	writeJSON(ctx, http.StatusOK, "success", "enrollment rejected successfully", nil)
}

func writeJSON(
	ctx *gin.Context,
	code int,
	status string,
	message string,
	data interface{},
) {
	ctx.JSON(code, apiResponse{
		Code:    code,
		Status:  status,
		Message: message,
		Data:    data,
	})
}
