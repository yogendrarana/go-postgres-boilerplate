package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-gin-postgres/internal/config"
	router "go-gin-postgres/internal/router"
)

type Server struct {
	httpServer *http.Server
	cfg        *config.Config
}

func NewServer(cfg *config.Config) *Server {
	router := router.NewRouter()

	server := &Server{
		httpServer: &http.Server{
			Addr:         fmt.Sprintf(":%d", cfg.Port),
			Handler:      router,
			IdleTimeout:  time.Minute,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 30 * time.Second,
		},
		cfg: cfg,
	}

	return server
}

func (s *Server) Run() error {
	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	log.Printf("Starting server on port %d", s.cfg.Port)

	// Start server
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Error starting server: %v", err)
			serverStopCtx()
		}
	}()

	// Wait for interrupt signal
	<-sig
	log.Printf("Shutdown signal received")

	// Shutdown signal with grace period from config
	shutdownCtx, cancel := context.WithTimeout(serverCtx, s.cfg.ShutdownTimeout)
	defer cancel()

	// Trigger graceful shutdown
	err := s.httpServer.Shutdown(shutdownCtx)
	if err != nil {
		log.Printf("Error during shutdown: %v", err)
	}

	// Wait for server context to be stopped
	serverStopCtx()
	<-serverCtx.Done()
	log.Println("Server exited gracefully")
	return nil
}
