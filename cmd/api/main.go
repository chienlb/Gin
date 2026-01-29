package main

import (
	"fmt"
	"os"

	"gin-demo/internal/app"
	"gin-demo/internal/config"
	"gin-demo/pkg/logger"
)

func main() {
	// Initialize logger
	log := logger.Init()

	// Load configuration
	cfg := config.Load()
	log.Info("Configuration loaded successfully")

	// Create and initialize server
	server := app.NewServer(cfg)
	if err := server.Initialize(); err != nil {
		log.Error("Failed to initialize server", err)
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Start server
	if err := server.Start(); err != nil {
		log.Error("Server error", err)
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
