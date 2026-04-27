package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func RoleAuth(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get(ContextUserRoleKey)
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code":    http.StatusForbidden,
				"status":  "error",
				"message": "user role not found in context",
			})
			return
		}

		role, ok := roleVal.(string)
		if !ok || strings.ToUpper(role) != strings.ToUpper(requiredRole) {
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
