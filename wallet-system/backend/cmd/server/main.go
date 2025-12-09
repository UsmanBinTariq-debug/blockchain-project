package main

import (
	"fmt"
	"log"

	"crypto-wallet-backend/internal/api"
	"crypto-wallet-backend/internal/blockchain"
	"crypto-wallet-backend/internal/database"
	"crypto-wallet-backend/pkg/config"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	db, err := database.NewDatabase(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize blockchain
	bc := blockchain.NewBlockchain()

	// Create handler
	handler := api.NewHandler(db, bc, cfg.JWTSecret)

	// Set Gin mode
	if cfg.NodeEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router
	router := gin.Default()

	// Apply middleware
	router.Use(api.CORSMiddleware(cfg.CORSAllowedOrigins))

	// Setup routes
	api.SetupRoutes(router, handler)

	// Start server: bind explicitly to 0.0.0.0 so PaaS like Render can detect the open port
	port := cfg.Port
	addr := fmt.Sprintf("0.0.0.0:%s", port)
	fmt.Printf("Starting server on %s\n", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
