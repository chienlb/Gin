package database

import (
	"fmt"

	"gin-demo/internal/domain"
	"gin-demo/pkg/logger"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	log := logger.Get()

	// AutoMigrate will create the table if it doesn't exist
	if err := db.AutoMigrate(&domain.User{}); err != nil {
		log.Error("Failed to run migrations", err)
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Info("Migrations executed successfully")
	return nil
}
