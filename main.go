package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"calico-go-project/database"
	"calico-go-project/handlers"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  Warning: .env file not found or failed to load")
	}

	// Initialize database and auto-run migrations
	if err := database.InitDB(); err != nil {
		log.Fatalf("❌ Failed to initialize database: %v", err)
	}
	defer database.CloseDB()

	// Set up routes
	router := handlers.SetupRoutes()

	// Start server
	log.Println("🚀 Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
