package domain

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// LoginRequest is the payload expected by the login endpoint.
type LoginRequest struct {
	NIM      string `json:"nim"`
	Password string `json:"password"`
}

// LoginResponse is returned on successful authentication.
type LoginResponse struct {
	Token     string `json:"token"`
	Role      string `json:"role"`
	ExpiresAt int64  `json:"expires_at"` // Unix timestamp
}

// Claims embeds the standard JWT registered claims and adds our custom fields.
// This struct is used both for signing tokens and for parsing them in middleware.
type Claims struct {
	NIM  string `json:"nim"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

// AuthUsecase defines the contract for authentication business logic.
type AuthUsecase interface {
	// Login validates credentials and returns a signed JWT token on success.
	Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)

	// ValidateClaims parses a raw JWT string and returns the embedded claims.
	ValidateClaims(tokenString string) (*Claims, error)
}

// TokenExpiry is the lifetime of a generated JWT token.
const TokenExpiry = 24 * time.Hour
