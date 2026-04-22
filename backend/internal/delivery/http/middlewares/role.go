package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RoleAuth(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get(ContextUserRoleKey)
		if !exists || role != requiredRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code":    http.StatusForbidden,
				"status":  "error",
				"message": "you do not have permission to access this resource",
			})
			return
		}
		c.Next()
	}
}
