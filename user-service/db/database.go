package db

import (
	"database/sql"
	"fmt"
	"hr-app-diploma-backend/pkg/logger"
	"hr-app-diploma-backend/user-service/config"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	cfg := config.AppConfig.Database
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}

	logger.Log.Info("Connected to database")
}
