package main

import (
	"go-gin-postgres/internal/api"
	"go-gin-postgres/internal/config"
	db "go-gin-postgres/internal/db"
	"go-gin-postgres/internal/initializers"
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

	db.NewDB(cfg)
	defer db.CloseDB()

	// Create and run server
	srv := api.NewServer(cfg)
	if err := srv.Run(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
