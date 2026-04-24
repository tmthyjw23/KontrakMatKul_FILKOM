package main

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"sistemkontrakmatkul/backend/internal/config"
	"sistemkontrakmatkul/backend/internal/delivery/http/handlers"
	"sistemkontrakmatkul/backend/internal/delivery/http/middlewares"
	"sistemkontrakmatkul/backend/internal/delivery/http/routes"
	"sistemkontrakmatkul/backend/internal/repository/mysql"
	"sistemkontrakmatkul/backend/internal/usecase"
	"sistemkontrakmatkul/backend/pkg/database"
	appLogger "sistemkontrakmatkul/backend/pkg/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	logger, err := appLogger.NewZap(cfg.AppEnv)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = logger.Sync()
	}()

	db, err := database.NewMySQL(cfg, logger)
	if err != nil {
		logger.Fatal("failed to initialize mysql", zap.Error(err))
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error("failed to close mysql connection", zap.Error(err))
		}
	}()

	// Dependency Injection flow:
	// infrastructure (db, logger, config)
	//   -> repository
	//   -> usecase/service
	//   -> handler
	//   -> routes
	courseRepository := mysql.NewCourseRepository(db, logger)
	courseUsecase := usecase.NewCourseUsecase(courseRepository, logger)
	courseHandler := handlers.NewCourseHandler(courseUsecase, logger)

	passedCourseRepository := mysql.NewPassedCourseRepository(db)
	prereqRepository := mysql.NewCoursePrerequisiteRepository(db)

	enrollmentRepository := mysql.NewEnrollmentRepository(db, logger)
	enrollmentUsecase := usecase.NewEnrollmentUsecase(enrollmentRepository, passedCourseRepository, prereqRepository, logger)
	enrollmentHandler := handlers.NewEnrollmentHandler(enrollmentUsecase, logger)

	userRepository := mysql.NewUserRepository(db, logger)
	authUsecase := usecase.NewAuthUsecase(userRepository, cfg.JWTSecret, logger)
	authHandler := handlers.NewAuthHandler(authUsecase)

	settingRepository := mysql.NewSystemSettingsRepository(db, logger)
	periodUsecase := usecase.NewPeriodUsecase(settingRepository, logger)
	periodHandler := handlers.NewPeriodHandler(periodUsecase)

	dashboardUsecase := usecase.NewStudentDashboardUsecase(enrollmentRepository, courseRepository, passedCourseRepository)
	dashboardHandler := handlers.NewStudentDashboardHandler(dashboardUsecase)

	prereqUsecase := usecase.NewCoursePrerequisiteUsecase(prereqRepository)
	prereqHandler := handlers.NewCoursePrerequisiteHandler(prereqUsecase)

	jwtMiddleware := middlewares.JWT(cfg.JWTSecret, logger)

	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())

	routes.SetupRoutes(
		router,
		courseHandler,
		enrollmentHandler,
		authHandler,
		periodHandler,
		dashboardHandler,
		prereqHandler,
		jwtMiddleware,
	)

	server := &http.Server{
		Addr:         cfg.Address(),
		Handler:      router,
		ReadTimeout:  cfg.ServerReadTimeout,
		WriteTimeout: cfg.ServerWriteTimeout,
		IdleTimeout:  cfg.ServerIdleTimeout,
	}

	shutdownSignalCtx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	go func() {
		logger.Info("server started", zap.String("address", cfg.Address()))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("failed to start http server", zap.Error(err))
		}
	}()

	<-shutdownSignalCtx.Done()
	logger.Info("shutdown signal received")

	gracefulShutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		cfg.ShutdownGracePeriod,
	)
	defer cancel()

	if err := server.Shutdown(gracefulShutdownCtx); err != nil {
		logger.Error("graceful shutdown failed", zap.Error(err))
		if closeErr := server.Close(); closeErr != nil {
			logger.Error("forced server close failed", zap.Error(closeErr))
		}
	}

	logger.Info("server stopped gracefully")
}
