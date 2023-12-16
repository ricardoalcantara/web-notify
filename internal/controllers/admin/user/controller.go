package user

import (
	"github.com/gin-gonic/gin"
	"github.com/ricardoalcantara/web-notify/internal/middlewares"
)

func RegisterRoutes(r *gin.Engine) {
	routes := r.Group("/api")
	routes.Use(middlewares.AuthMiddleware())
	routes.GET("/user", get)
	routes.POST("/user", post)
	routes.GET("/user/:userId/token", token)
}
