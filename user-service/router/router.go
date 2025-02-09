package router

import (
	"github.com/gin-gonic/gin"
	"hr-app-diploma-backend/user-service/handler"
	"hr-app-diploma-backend/user-service/middleware"
)

func SetupRouter(userHandler *handler.UserHandler, authMiddleware *middleware.AuthMiddleware) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	{
		users := api.Group("/users")
		users.Use(authMiddleware.VerifyTokenMiddleware())
		{
			users.GET("/me", userHandler.GetUser)
			users.PATCH("/me", userHandler.UpdateUser)
			//users.DELETE("/me", userHandler.DeleteUser)
		}
	}

	return r
}
