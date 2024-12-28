package database

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"go-gin-postgres/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db   *gorm.DB
	once sync.Once
)

func NewDB(cfg *config.Config) *gorm.DB {
	once.Do(func() {
		logLevel := logger.Info
		if cfg.Environment == "production" {
			logLevel = logger.Error
		}

		newLogger := logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  logLevel,
				IgnoreRecordNotFoundError: true,
				Colorful:                  true,
			},
		)

		db, err := gorm.Open(postgres.Open(cfg.DBUrl), &gorm.Config{
			Logger:      newLogger,
			PrepareStmt: true,
		})
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}

		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("Failed to get database: %v", err)
		}

		// Configure connection pool
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Hour)

		log.Println("Database connection established successfully")
	})

	return db
}

// Close closes the database connection.
func CloseDB() {
	if db == nil {
		return
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("Error disconnecting from database: %v", err)
		return
	}

	log.Printf("Disconnecting from database")
	if err := sqlDB.Close(); err != nil {
		log.Printf("Error closing database connection: %v", err)
	}
}

// GetDB returns the underlying *gorm.DB instance
func GetDB() *gorm.DB {
	return db
}

// Health returns the health status of the database.
func DBHealth() map[string]string {
	stats := make(map[string]string)

	sqlDB, err := db.DB()
	if err != nil {
		stats["status"] = "error"
		stats["message"] = fmt.Sprintf("Failed to get database: %v", err)
		return stats
	}

	err = sqlDB.Ping()
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("Database down: %v", err)
		return stats
	}

	stats["status"] = "up"
	stats["message"] = "Database is healthy"

	dbStats := sqlDB.Stats()
	stats["open_connections"] = fmt.Sprintf("%d", dbStats.OpenConnections)
	stats["in_use"] = fmt.Sprintf("%d", dbStats.InUse)
	stats["idle"] = fmt.Sprintf("%d", dbStats.Idle)
	stats["wait_count"] = fmt.Sprintf("%d", dbStats.WaitCount)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = fmt.Sprintf("%d", dbStats.MaxIdleClosed)
	stats["max_lifetime_closed"] = fmt.Sprintf("%d", dbStats.MaxLifetimeClosed)

	return stats
}
