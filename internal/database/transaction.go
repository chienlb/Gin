package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"gorm.io/gorm"
)

// TransactionFunc defines a function that runs within a transaction
type TransactionFunc func(*gorm.DB) error

// WithTransaction executes a function within a database transaction
// If the function returns an error, the transaction is rolled back
// Otherwise, the transaction is committed
func WithTransaction(db *gorm.DB, fn TransactionFunc) error {
	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("Transaction panicked and rolled back: %v", r)
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback().Error; rbErr != nil {
			return fmt.Errorf("failed to rollback transaction: %w (original error: %v)", rbErr, err)
		}
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// WithContext executes a function with a context-aware DB instance
func WithContext(db *gorm.DB, ctx context.Context) *gorm.DB {
	return db.WithContext(ctx)
}

// CreateIndexes creates database indexes for better query performance
func CreateIndexes(db *gorm.DB) error {
	indexes := []string{
		// User indexes
		"CREATE INDEX IF NOT EXISTS idx_users_email ON users(email) WHERE deleted_at IS NULL",
		"CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at DESC)",
		"CREATE INDEX IF NOT EXISTS idx_users_name ON users(name) WHERE deleted_at IS NULL",

		// Composite indexes for common queries
		"CREATE INDEX IF NOT EXISTS idx_users_email_created ON users(email, created_at) WHERE deleted_at IS NULL",
	}

	for _, query := range indexes {
		if err := db.Exec(query).Error; err != nil {
			log.Printf("Warning: Failed to create index: %v", err)
		}
	}

	return nil
}

// Backup creates a logical backup of the database (PostgreSQL specific)
func Backup(db *gorm.DB, outputPath string) error {
	// Get database connection info
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// This is a simplified example
	// In production, use pg_dump or appropriate backup tool
	log.Printf("Backup initiated to %s", outputPath)
	log.Println("Note: Use pg_dump for production backups")

	_ = sqlDB
	return nil
}

// GetDBStats returns database connection statistics
func GetDBStats(db *gorm.DB) sql.DBStats {
	sqlDB, err := db.DB()
	if err != nil {
		return sql.DBStats{}
	}
	return sqlDB.Stats()
}

// HealthCheck verifies database connectivity
func HealthCheck(db *gorm.DB, ctx context.Context) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	return nil
}
