package message

import "github.com/gin-gonic/gin"

func acknowledge(c *gin.Context) {
	// h.DB
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
