package main

import (
	"log"
	"net/http"
	"os"

	"github.com/kifangamukundi/gm/loan/config"
	"github.com/kifangamukundi/gm/loan/database"
	"github.com/kifangamukundi/gm/loan/emails"
	"github.com/kifangamukundi/gm/loan/jobs"
	"github.com/kifangamukundi/gm/loan/migrations"
	"github.com/kifangamukundi/gm/libs/rates"
	"github.com/kifangamukundi/gm/loan/routes"
	"github.com/kifangamukundi/gm/loan/seeds"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file (if present)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, proceeding with system environment variables.")
	}

	// Get the GIN_MODE environment variable and convert it to a boolean
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "true" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Initialize the database connection
	err := database.ConnectDB()
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	// Get the database instance
	db := database.GetDB()

	// Run migrations
	shouldRunMigrations := os.Getenv("RUN_MIGRATIONS") == "true"
	if shouldRunMigrations {
		migrations.RunMigrations()
	}

	// Run all seeders
	shouldSeedData := os.Getenv("SEED_DATA") == "true"
	if shouldSeedData {
		if err := seeds.SeedAll(db); err != nil {
			log.Fatalf("Error running seeders: %v", err)
		}
	}

	r := gin.Default()

	// Initialize CORS configuration
	r.Use(config.InitializeCorsConfig())

	// Initialize Redis connection
	rates.InitRedis()

	// Example route for testing
	r.GET("/api", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"company": "Kifanga Mukundi, Inc",
		})
	})

	// Initialize routes with the database instance
	routes.InitializeRoutes(r, db)

	// Initialize cron jobs
	jobs.InitializeJobs()

	// Load and Initialize Mail Configurations from the config package
	mailConfig := config.LoadMailConfig()

	// Initialize the email configuration
	emails.InitializeMailConfig(mailConfig)

	// Set the port based on the environment (development or production)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on http://localhost:%s", port)
	r.Run(":" + port)
}
