package database

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Service interface {
	Health() map[string]string
	Close() error
	GetDB() *gorm.DB
}

type service struct {
	db *gorm.DB
}

var (
	dbUrl      = os.Getenv("DB_URL")
	dbInstance *service
	once       sync.Once
)

// New creates a new database service.
func New() Service {
	once.Do(func() {
		if dbUrl == "" {
			log.Fatal("DB_URL environment variable is not set")
		}

		newLogger := logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  logger.Info,
				IgnoreRecordNotFoundError: true,
				Colorful:                  true,
			},
		)

		db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{
			Logger: newLogger,
		})
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}

		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("Failed to get database: %v", err)
		}

		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Hour)

		dbInstance = &service{
			db: db,
		}
	})

	return dbInstance
}

// Health returns the health status of the database.
func (s *service) Health() map[string]string {
	stats := make(map[string]string)

	sqlDB, err := s.db.DB()
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

// Close closes the database connection.
func (s *service) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database: %w", err)
	}
	log.Printf("Disconnecting from database")
	return sqlDB.Close()
}

// GetDB returns the underlying *gorm.DB instance
func (s *service) GetDB() *gorm.DB {
	return s.db
}
