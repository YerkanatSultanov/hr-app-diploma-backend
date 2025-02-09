package db

import (
	"database/sql"
	"fmt"
	"hr-app-diploma-backend/auth-service/config"
	"hr-app-diploma-backend/pkg/logger"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations() {
	dbConfig := config.AppConfig.Database
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.DBName, dbConfig.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Log.Error("Failed to connect to AuthService database for migrations",
			slog.String("error", err.Error()))
		return
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		logger.Log.Error("Failed to create migration driver",
			slog.String("error", err.Error()))
		return
	}

	wd, _ := os.Getwd()
	migrationsPath := fmt.Sprintf("file://%s/auth-service/db/migrations", filepath.Join(wd))

	m, err := migrate.NewWithDatabaseInstance(migrationsPath, "postgres", driver)
	if err != nil {
		logger.Log.Error("Failed to create migration instance",
			slog.String("error", err.Error()))
		return
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Log.Error("Migration failed", slog.String("error", err.Error()))
		return
	}

	logger.Log.Info("AuthService migrations applied successfully")
}
