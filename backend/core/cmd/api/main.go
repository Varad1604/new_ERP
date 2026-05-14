package main

import (
	"log"
	"net/http"
	"os"

	httpdelivery "github.com/enterprise-erp/core/internal/delivery/http"
	"github.com/enterprise-erp/core/internal/repository/postgres"
	"github.com/enterprise-erp/core/internal/usecase"
	"github.com/enterprise-erp/core/pkg/auth"
	"github.com/enterprise-erp/core/pkg/database"
	"github.com/gin-gonic/gin"
)

func main() {
	// Database connection configuration
	dbConfig := database.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "erp_admin"),
		Password: getEnv("DB_PASSWORD", "erp_password_dev"),
		DBName:   getEnv("DB_NAME", "enterprise_erp"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	// Connect to Database
	db, err := database.ConnectDB(dbConfig)
	if err != nil {
		log.Printf("Warning: Failed to connect to database. Starting without DB (stub mode). Error: %v", err)
	} else {
		log.Println("Successfully connected to the database.")

		// Run Database Migrations
		migrationsPath := getEnv("MIGRATIONS_PATH", "db/migrations")
		err = database.RunMigrations(db, migrationsPath)
		if err != nil {
			log.Fatalf("Critical: Database migrations failed: %v", err)
		}
	}

	jwtSecret := getEnv("JWT_SECRET", "super-secret-enterprise-key-do-not-use-in-prod")

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
		if db != nil {
			// Initialize Repositories
			userRepo := postgres.NewUserRepository(db)
			financeRepo := postgres.NewFinanceRepository(db)

			// Initialize UseCases
			authUseCase := usecase.NewAuthUseCase(userRepo, jwtSecret)
			financeUseCase := usecase.NewFinanceUseCase(financeRepo)

			// Initialize Handlers & Routes
			authHandler := httpdelivery.NewAuthHandler(authUseCase)
			authHandler.RegisterRoutes(v1)

			// Finance routes require authentication
			financeGroup := v1.Group("/")
			financeGroup.Use(auth.AuthMiddleware(jwtSecret))
			financeHandler := httpdelivery.NewFinanceHandler(financeUseCase)
			financeHandler.RegisterRoutes(financeGroup)
		} else {
			v1.GET("/warning", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "API running without database connection"})
			})
		}
	}

	// Server port configuration
	port := getEnv("PORT", "8080")

	log.Printf("Starting core service on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
