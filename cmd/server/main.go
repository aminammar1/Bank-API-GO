package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bank-api/internal/api/routes"
	"github.com/bank-api/internal/config"
	"github.com/bank-api/internal/repository"
)

func main() {
	// Load configuration
	cfg := config.Load()
	
	// Initialize database
	db, err := repository.NewPostgresDB(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	
	// Initialize database tables
	if err := repository.InitializeDatabase(db); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	
	// Setup routes
	router := routes.NewRouter(db, cfg)
	handler := router.SetupRoutes()
	
	// Create server
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler:      handler,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  60 * time.Second,
	}
	
	log.Printf("🏦 Bank API server starting on %s", server.Addr)
	log.Printf("📊 Health check available at: http://%s/api/v1/health", server.Addr)
	
	// Start server
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}
