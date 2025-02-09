package router

import (
	"hr-app-diploma-backend/auth-service/handler"
	"hr-app-diploma-backend/auth-service/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(authHandler *handler.AuthHandler) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())
	api := r.Group("/api/auth")
	{
		api.POST("/register", authHandler.Register)
		api.POST("/login", authHandler.Login)
	}

	protected := r.Group("/api/protected")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/profile", authHandler.Profile)
	}

	return r
}
