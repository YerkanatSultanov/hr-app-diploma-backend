package main

import (
	"fmt"
	"hr-app-diploma-backend/pkg/logger"
	"hr-app-diploma-backend/user-service/config"
	"hr-app-diploma-backend/user-service/db"
	"hr-app-diploma-backend/user-service/grpc"
	"hr-app-diploma-backend/user-service/handler"
	"hr-app-diploma-backend/user-service/middleware"
	"hr-app-diploma-backend/user-service/repository"
	"hr-app-diploma-backend/user-service/router"
	"hr-app-diploma-backend/user-service/service"
)

func main() {
	config.LoadConfig("user-service/config/config.yaml")
	logger.InitLogger()

	db.InitDB()
	db.RunMigrations()

	authClient, err := grpc.NewAuthServiceClient(config.AppConfig.AuthService.Address)
	if err != nil {
		logger.Log.Error("Failed to connect to auth-service", "error", err.Error())
		return
	}

	userRepo := repository.NewUserRepository(db.DB)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)
	authMiddleware := middleware.NewAuthMiddleware(authClient)

	r := router.SetupRouter(userHandler, authMiddleware)

	port := config.AppConfig.Server.Port
	logger.Log.Info("Starting user-service", "port", port)

	err = r.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		logger.Log.Error("Failed to start server", "error", err)
	}
}
