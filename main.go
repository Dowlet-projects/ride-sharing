// @title Ride-Sharing Backend API
// @version 1.0
// @description API for managing passengers and taxists in a ride-sharing application.
// @host 192.168.55.42:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"log"
	"net/http"

	"ride-sharing/config"
	"ride-sharing/db"
	"ride-sharing/handlers"

	_ "ride-sharing/docs"

	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger" // Add this import
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found: %v", err)
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	database, err := db.Initialize(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Set up router
	router := handlers.SetupRouter(database, cfg)
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	router.PathPrefix("/docs/").Handler(http.StripPrefix("/docs/", http.FileServer(http.Dir("docs"))))

	// Specify IP address and port
	addr := "0.0.0.0:8080" // Bind to all interfaces; replace with specific IP like "192.168.1.100:8080" if needed
	log.Printf("Server starting on http://%s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
