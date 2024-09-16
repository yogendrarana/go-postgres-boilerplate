package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"

	"how-to-server/internal/database"
	"how-to-server/internal/routers"
)

type Server struct {
	port   int
	router *gin.Engine
	db     database.Service
}

// NewServer initializes and returns an *http.Server with the Gin router
func NewServer() *http.Server {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		port = 8000
	}

	// Initialize the server and database
	srv := &Server{
		port:   port,
		router: routers.NewRouter(),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", srv.port),
		Handler:      srv.router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
