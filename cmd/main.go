package main

import (
	"fmt"
	"how-to-server/internal/database"
	"how-to-server/internal/initializers"
	"how-to-server/internal/server"
	"log"
	"net/http"
)

func init() {
	initializers.LoadEnvVariables()
}

func main() {
	// Create a new database service
	dbService := database.New()
	dbService.GetDB()
	defer dbService.Close()

	// Create a new server
	srv := server.NewServer()
	port := srv.Addr[len(":"):]

	// Start the server
	fmt.Printf("Server is running on port %s\n", port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Error starting server: %v", err)
	}

	// Gracefully close the database connection when shutting down
	if err := dbService.Close(); err != nil {
		log.Fatalf("Error closing database: %v", err)
	}
}
