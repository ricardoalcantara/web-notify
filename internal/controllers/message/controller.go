package message

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	routes := r.Group("/api")
	// routes.Use(middlewares.AuthMiddleware())
	routes.POST("/message/notify", notify)
	routes.POST("/message/acknowledge", acknowledge)
}
