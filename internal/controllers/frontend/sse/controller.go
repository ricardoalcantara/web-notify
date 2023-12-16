package sse

import (
	"github.com/gin-gonic/gin"
	"github.com/ricardoalcantara/web-notify/internal/middlewares"
)

func RegisterRoutes(r *gin.Engine) {
	routes := r.Group("/")
	routes.Use(middlewares.SseAuthMiddleware())
	routes.GET("/sse", get)
}
