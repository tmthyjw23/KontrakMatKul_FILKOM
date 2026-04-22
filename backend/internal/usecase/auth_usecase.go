package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"sistemkontrakmatkul/backend/internal/domain/models"
	"sistemkontrakmatkul/backend/internal/domain/repositories"
)

type AuthUsecase struct {
	userRepo repositories.UserRepository
	jwtSecret string
	logger    *zap.Logger
}

func NewAuthUsecase(userRepo repositories.UserRepository, jwtSecret string, logger *zap.Logger) *AuthUsecase {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &AuthUsecase{
		userRepo:   userRepo,
		jwtSecret: jwtSecret,
		logger:    logger,
	}
}

func (u *AuthUsecase) Login(ctx context.Context, req models.UserLoginRequest) (*models.UserResponse, string, error) {
	user, err := u.userRepo.GetByStudentNumber(ctx, req.StudentNumber)
	if err != nil {
		u.logger.Error("failed to fetch user during login", zap.String("student_number", req.StudentNumber), zap.Error(err))
		return nil, "", err
	}

	if user == nil {
		return nil, "", errors.New("invalid student number or password")
	}

	// In a real app, use bcrypt.CompareHashAndPassword
	if user.Password != req.Password {
		return nil, "", errors.New("invalid student number or password")
	}

	// Generate JWT Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
		"iat":     time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(u.jwtSecret))
	if err != nil {
		u.logger.Error("failed to sign jwt token", zap.Error(err))
		return nil, "", err
	}

	response := &models.UserResponse{
		ID:            user.ID,
		StudentNumber: user.StudentNumber,
		FullName:      user.FullName,
		Email:         user.Email,
		Role:          user.Role,
	}

	return response, tokenString, nil
}
