package main

import (
	"go-gin-postgres/internal/config"
	"go-gin-postgres/internal/database"
	"go-gin-postgres/internal/initializers"
	"go-gin-postgres/internal/server"
	"log"
)

func init() {
	initializers.LoadEnvVariables()
}

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	database.NewDB(cfg)
	defer database.CloseDB()

	// Create and run server
	srv := server.NewServer(cfg)
	if err := srv.Run(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
