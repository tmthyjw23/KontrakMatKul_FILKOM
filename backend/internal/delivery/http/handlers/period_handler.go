package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"sistemkontrakmatkul/backend/internal/domain/models"
	"sistemkontrakmatkul/backend/internal/usecase"
)

type PeriodHandler struct {
	periodUsecase *usecase.PeriodUsecase
}

func NewPeriodHandler(periodUsecase *usecase.PeriodUsecase) *PeriodHandler {
	return &PeriodHandler{
		periodUsecase: periodUsecase,
	}
}

func (h *PeriodHandler) GetStatus(c *gin.Context) {
	res, err := h.periodUsecase.GetPeriodStatus(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"status":  "error",
			"message": "failed to fetch period status",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"status":  "success",
		"data":    res,
	})
}

func (h *PeriodHandler) UpdateStatus(c *gin.Context) {
	var req models.PeriodUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"status":  "error",
			"message": "invalid request payload",
		})
		return
	}

	if err := h.periodUsecase.UpdatePeriodStatus(c.Request.Context(), req.IsOpen); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"status":  "error",
			"message": "failed to update period status",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"status":  "success",
		"message": "period status updated successfully",
	})
}
