package message

import (
	"github.com/gin-gonic/gin"
	"github.com/ricardoalcantara/web-notify/internal/middlewares"
)

func RegisterRoutes(r *gin.Engine) {
	routes := r.Group("/api")
	routes.Use(middlewares.AuthMiddleware())
	routes.POST("/message/notify", notify)
}
