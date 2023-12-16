package user

import "github.com/gin-gonic/gin"

func get(c *gin.Context) {
	// h.DB
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
