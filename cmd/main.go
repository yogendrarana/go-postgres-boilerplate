package main

import (
	"fmt"
	"how-to-server/internal/server"
	"log"
	"net/http"
)

func main() {
	srv := server.NewServer()

	// Extract the port from the Addr field
	port := srv.Addr[len(":"):]

	// Start the server
	fmt.Printf("Server is running on http://localhost:%s\n", port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Error starting server: %v", err)
	}
}
