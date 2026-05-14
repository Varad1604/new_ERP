package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize router
	router := gin.Default()

	// Middleware
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "UP",
			"service": "core-erp-engine",
		})
	})

	// API Version 1 Group
	v1 := router.Group("/api/v1")
	{
		// Module stubs
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/login", func(c *gin.Context) { c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented"}) })
		}

		financeGroup := v1.Group("/finance")
		{
			financeGroup.POST("/journal", func(c *gin.Context) { c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented"}) })
			financeGroup.GET("/ledger", func(c *gin.Context) { c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented"}) })
		}
	}

	// Server port configuration
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting core service on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
