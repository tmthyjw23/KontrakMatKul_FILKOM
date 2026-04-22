package usecase

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"kontrak-matkul/domain"
)

// authUsecase is the concrete implementation of domain.AuthUsecase.
type authUsecase struct {
	userRepo domain.UserRepository
}

// NewAuthUsecase creates and returns a new authUsecase.
func NewAuthUsecase(ur domain.UserRepository) domain.AuthUsecase {
	return &authUsecase{userRepo: ur}
}

func (u *authUsecase) getSecretKey() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// Fallback for development, should be strictly enforced in production
		return []byte("super-secret-key")
	}
	return []byte(secret)
}

// Login validates credentials and generates a JWT.
func (u *authUsecase) Login(ctx context.Context, req *domain.LoginRequest) (*domain.LoginResponse, error) {
	if req.NIM == "" || req.Password == "" {
		return nil, fmt.Errorf("nim and password are required")
	}

	// 1. Get the user's profile to extract the role
	user, err := u.userRepo.GetByNIM(ctx, req.NIM)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// 2. Fetch the stored password hash
	hash, err := u.userRepo.GetPasswordHashByNIM(ctx, req.NIM)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// 3. Compare the plaintext password with the hash
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// 4. Generate the JWT token
	expirationTime := time.Now().Add(domain.TokenExpiry)
	claims := &domain.Claims{
		NIM:  user.NIM,
		Role: user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(u.getSecretKey())
	if err != nil {
		return nil, fmt.Errorf("could not generate token: %w", err)
	}

	return &domain.LoginResponse{
		Token:     tokenString,
		Role:      user.Role,
		ExpiresAt: expirationTime.Unix(),
	}, nil
}

// ValidateClaims parses the raw JWT string and returns the claims.
func (u *authUsecase) ValidateClaims(tokenString string) (*domain.Claims, error) {
	claims := &domain.Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		// Ensure the signing method is what we expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return u.getSecretKey(), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
