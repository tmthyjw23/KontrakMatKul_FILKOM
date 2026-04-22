package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"sistemkontrakmatkul/backend/internal/delivery/http/handlers"
	"sistemkontrakmatkul/backend/internal/delivery/http/middlewares"
)

func SetupRoutes(
	router *gin.Engine,
	enrollmentHandler *handlers.EnrollmentHandler,
	jwtMiddleware gin.HandlerFunc,
) {
	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"status":  "success",
			"message": "course registration backend is running",
			"data": gin.H{
				"service": "sistem-kontrak-api",
			},
		})
	})

	apiV1 := router.Group("/api/v1")
	apiV1.Use(middlewares.CORS())
	{
		apiV1.POST("/enrollments", jwtMiddleware, enrollmentHandler.Enroll)
		apiV1.OPTIONS("/enrollments", func(ctx *gin.Context) {
			ctx.Status(http.StatusNoContent)
		})
	}
}
