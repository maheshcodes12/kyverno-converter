package main

import (
	"kyverno-converter-backend/internal/handlers"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Configure CORS to allow requests from the frontend development server
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173"}
	config.AllowMethods = []string{"POST", "OPTIONS"}
	r.Use(cors.New(config))

	// Setup API routes
	api := r.Group("/api")
	{
		api.POST("/convert", handlers.ConvertPolicyHandler)
	}

	// Start the server
	port := "8080"
	log.Printf("Backend server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
