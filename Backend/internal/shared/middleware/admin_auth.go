package middleware

import (
	"Re_Shop/Backend/internal/modules/user/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleValue, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "login is required",
			})
			c.Abort()
			return
		}

		role, ok := roleValue.(int8)
		if !ok || role != model.UserRoleAdmin {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    http.StatusForbidden,
				"message": "admin permission is required",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
