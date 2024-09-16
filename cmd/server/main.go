package main

import (
	"fmt"
	"how-to-server/internal/database"
	"how-to-server/internal/server"
	"log"
	"net/http"

	"github.com/pressly/goose"
)

func main() {
	// database service
	dbService := database.New()
	defer dbService.Close()

	// run migrations
	if err := runMigrations(dbService); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	srv := server.NewServer()
	port := srv.Addr[len(":"):]

	// Start the server
	fmt.Printf("Server is running on port %s\n", port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Error starting server: %v", err)
	}
}

func runMigrations(dbService database.Service) error {
	db := dbService.GetDB()
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get *sql.DB: %w", err)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	// Get the migration directory from environment variable or use default
	migrationDir := "./migration"
	if migrationDir == "" {
		migrationDir = "migrations"
	}

	if err := goose.Up(sqlDB, "migrations"); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
