package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port            int
	DBUrl           string
	Environment     string
	ShutdownTimeout time.Duration
}

func LoadConfig() (*Config, error) {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		port = 8000
	}

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		return nil, fmt.Errorf("DB_URL environment variable is not set")
	}

	// Default shutdown timeout is 15 seconds if not specified
	shutdownTimeout := 15 * time.Second
	if timeout := os.Getenv("SHUTDOWN_TIMEOUT_SECONDS"); timeout != "" {
		if timeoutInt, err := strconv.Atoi(timeout); err == nil {
			shutdownTimeout = time.Duration(timeoutInt) * time.Second
		}
	}

	return &Config{
		Port:            port,
		DBUrl:           dbUrl,
		Environment:     os.Getenv("APP_ENV"),
		ShutdownTimeout: shutdownTimeout,
	}, nil
}
