package router

import "github.com/gin-gonic/gin"

func RegesiterUserRoutes(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "good",
		})
	})
}
