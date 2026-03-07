package main

import (
	"log"
	"os"

	"github.com/binduni/bun-golang-react-monorepo/server/config"
	"github.com/binduni/bun-golang-react-monorepo/server/database"
	"github.com/binduni/bun-golang-react-monorepo/server/routes"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Connect to database
	var db *database.DB
	if cfg.DatabaseURL != "" {
		var err error
		db, err = database.Connect(cfg.DatabaseURL)
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
		defer db.Close()
	} else {
		log.Println("⚠️  No DATABASE_URL provided, running without database")
	}

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "Monorepo API v" + Version,
		ErrorHandler: errorHandler,
	})

	// Global middleware
	app.Use(recover.New())

	// Logger middleware (development only)
	if cfg.IsDevelopment() {
		app.Use(logger.New(logger.Config{
			Format: "[${time}] ${status} - ${method} ${path} (${latency})\n",
		}))
	}

	// CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.GetAllowedOrigins(),
		AllowCredentials: true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "X-Request-Id"},
	}))

	// Setup routes
	routes.SetupRoutes(app, cfg, db)

	// Start server
	port := cfg.Port
	log.Printf("🚀 Server starting on port %s (env: %s)", port, cfg.Environment)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// Global error handler
func errorHandler(c fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal server error"

	// Check for Fiber errors
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	// Development mode - show detailed errors
	if os.Getenv("ENVIRONMENT") == "development" {
		message = err.Error()
	}

	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"error":   message,
	})
}
