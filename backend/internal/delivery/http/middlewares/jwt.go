package middlewares

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

const ContextUserIDKey = "user_id"

type jwtClaims struct {
	UserID uint64 `json:"user_id"`
	jwt.RegisteredClaims
}

func JWT(secret string, logger *zap.Logger) gin.HandlerFunc {
	if logger == nil {
		logger = zap.NewNop()
	}

	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			abortUnauthorized(ctx, "missing authorization header")
			return
		}

		const prefix = "Bearer "
		if !strings.HasPrefix(authHeader, prefix) {
			abortUnauthorized(ctx, "invalid authorization scheme")
			return
		}

		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, prefix))
		if tokenString == "" {
			abortUnauthorized(ctx, "missing bearer token")
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &jwtClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}

			return []byte(secret), nil
		})
		if err != nil {
			logger.Warn("failed to validate jwt", zap.Error(err))
			abortUnauthorized(ctx, "invalid token")
			return
		}

		claims, ok := token.Claims.(*jwtClaims)
		if !ok || !token.Valid || claims.UserID == 0 {
			abortUnauthorized(ctx, "invalid token claims")
			return
		}

		ctx.Set(ContextUserIDKey, claims.UserID)
		ctx.Next()
	}
}

func abortUnauthorized(ctx *gin.Context, message string) {
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"code":    http.StatusUnauthorized,
		"status":  "error",
		"message": message,
		"data":    nil,
	})
}
