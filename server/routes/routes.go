package routes

import (
	"runtime"
	"time"

	"github.com/binduni/bun-golang-react-monorepo/server/config"
	"github.com/binduni/bun-golang-react-monorepo/server/database"
	"github.com/binduni/bun-golang-react-monorepo/server/models"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, cfg *config.Config, db *database.DB) {
	// Root endpoint - API information
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"name":    "Monorepo API",
			"version": "2.0.0",
			"stack":   "Go + Fiber",
			"features": fiber.Map{
				"authentication": "JWT (email/password)",
				"oauth":          "Google, Facebook, Twitter",
				"ids":            "TypeID (type-safe, K-sortable)",
				"roles":          "admin, user, moderator",
			},
			"endpoints": fiber.Map{
				"auth": fiber.Map{
					"register": "POST /api/auth/register",
					"login":    "POST /api/auth/login",
					"refresh":  "POST /api/auth/refresh",
					"logout":   "POST /api/auth/logout",
					"me":       "GET /api/auth/me",
					"sessions": "GET /api/auth/sessions",
				},
				"oauth": fiber.Map{
					"providers": "GET /api/auth/oauth/providers",
					"google":    "GET /api/auth/oauth/google",
					"facebook":  "GET /api/auth/oauth/facebook",
					"twitter":   "GET /api/auth/oauth/twitter",
				},
				"items": fiber.Map{
					"list":   "GET /api/items",
					"get":    "GET /api/items/:id",
					"create": "POST /api/items",
					"update": "PUT /api/items/:id",
					"delete": "DELETE /api/items/:id",
				},
			},
			"docs": "https://github.com/your-repo/docs",
		})
	})

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		health := fiber.Map{
			"status":    "healthy",
			"timestamp": time.Now().Format(time.RFC3339),
			"uptime":    time.Since(startTime).Seconds(),
			"memory": fiber.Map{
				"alloc":      m.Alloc / 1024 / 1024,      // MB
				"totalAlloc": m.TotalAlloc / 1024 / 1024, // MB
				"sys":        m.Sys / 1024 / 1024,        // MB
			},
		}

		// Check database connection
		if db != nil {
			if err := db.Health(); err != nil {
				health["database"] = "unhealthy"
				health["status"] = "degraded"
			} else {
				health["database"] = "healthy"
			}
		} else {
			health["database"] = "not_configured"
		}

		return c.JSON(health)
	})

	// API group
	api := app.Group("/api")

	// Mount auth routes
	auth := api.Group("/auth")
	if db != nil {
		SetupAuthRoutes(auth, cfg, db)
	}

	// Mount items routes
	items := api.Group("/items")
	if db != nil {
		SetupItemsRoutes(items, cfg, db)
	}

	// 404 handler
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(
			models.ErrorResponse("Route not found: " + c.Method() + " " + c.Path()),
		)
	})
}

var startTime = time.Now()
