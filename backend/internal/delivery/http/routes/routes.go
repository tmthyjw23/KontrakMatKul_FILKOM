package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"sistemkontrakmatkul/backend/internal/delivery/http/handlers"
	"sistemkontrakmatkul/backend/internal/delivery/http/middlewares"
)

func SetupRoutes(
	router *gin.Engine,
	courseHandler *handlers.CourseHandler,
	enrollmentHandler *handlers.EnrollmentHandler,
	authHandler *handlers.AuthHandler,
	periodHandler *handlers.PeriodHandler,
	dashboardHandler *handlers.StudentDashboardHandler,
	prereqHandler *handlers.CoursePrerequisiteHandler,
	jwtMiddleware gin.HandlerFunc,
) {
	router.Use(middlewares.CORS())

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
	{
		// Public routes
		apiV1.POST("/login", authHandler.Login)
		apiV1.GET("/courses", courseHandler.ListCourses)

		// Student/User routes
		apiV1.GET("/period/status", jwtMiddleware, periodHandler.GetStatus)
		apiV1.POST("/enrollments", jwtMiddleware, enrollmentHandler.Enroll)
		apiV1.GET("/dashboard/schedule", jwtMiddleware, dashboardHandler.GetSchedule)
		apiV1.GET("/dashboard/history", jwtMiddleware, dashboardHandler.GetHistory)

		// Admin routes
		adminGroup := apiV1.Group("/admin")
		adminGroup.Use(jwtMiddleware, middlewares.RoleAuth("ADMIN"))
		{
			adminGroup.PUT("/period", periodHandler.UpdateStatus)
			adminGroup.POST("/prerequisites", prereqHandler.Add)
			adminGroup.DELETE("/prerequisites", prereqHandler.Remove)
		}

		apiV1.OPTIONS("/enrollments", func(ctx *gin.Context) {
			ctx.Status(http.StatusNoContent)
		})
	}
}
