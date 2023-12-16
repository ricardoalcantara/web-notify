package user

import "github.com/gin-gonic/gin"

func post(c *gin.Context) {
	// h.DB
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
