package main

import (
	"fmt"
	"hr-app-diploma-backend/auth-service/config"
	"hr-app-diploma-backend/auth-service/db"
	"hr-app-diploma-backend/auth-service/grpc"
	"hr-app-diploma-backend/auth-service/handler"
	"hr-app-diploma-backend/auth-service/repository"
	"hr-app-diploma-backend/auth-service/router"
	"hr-app-diploma-backend/auth-service/service"
	"hr-app-diploma-backend/pkg/logger"

	_ "github.com/lib/pq"
)

func main() {
	config.LoadConfig("auth-service/config/config.yaml")
	logger.InitLogger()

	db.InitDB()
	db.RunMigrations()

	userRepo := repository.NewAuthRepository(db.DB)
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	go grpc.StartGRPCServer(authService)
	r := router.SetupRouter(authHandler)

	serverPort := config.AppConfig.Server.Port
	addr := fmt.Sprintf(":%d", serverPort)
	logger.Log.Info("Starting server", "port", serverPort)

	if err := r.Run(addr); err != nil {
		logger.Log.Error("Failed to start server", "error", err.Error())
	}
}
