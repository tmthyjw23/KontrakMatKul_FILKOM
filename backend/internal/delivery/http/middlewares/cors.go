package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const allowedOrigin = "http://localhost:3000"

// CORS strictly allows the local Next.js frontend and blocks unexpected browser origins.
func CORS() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		origin := ctx.GetHeader("Origin")

		if origin != "" && origin != allowedOrigin {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code":    http.StatusForbidden,
				"status":  "error",
				"message": "origin is not allowed",
				"data":    nil,
			})
			return
		}

		if origin == allowedOrigin {
			headers := ctx.Writer.Header()
			headers.Set("Access-Control-Allow-Origin", allowedOrigin)
			headers.Set("Vary", "Origin")
			headers.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			headers.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		}

		if ctx.Request.Method == http.MethodOptions {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}

		ctx.Next()
	}
}
