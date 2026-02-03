package routes

import (
	"github.com/binduni/bun-golang-react-monorepo/server/config"
	"github.com/binduni/bun-golang-react-monorepo/server/database"
	"github.com/binduni/bun-golang-react-monorepo/server/handlers"
	"github.com/binduni/bun-golang-react-monorepo/server/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(router fiber.Router, cfg *config.Config, db *database.DB) {
	authHandler := handlers.NewAuthHandler(db, cfg)

	// Public routes (no authentication required)
	router.Post("/register", authHandler.Register)
	router.Post("/login", authHandler.Login)
	router.Post("/refresh", authHandler.RefreshToken)

	// Protected routes (authentication required)
	protected := router.Group("", middleware.AuthMiddleware(cfg.JWTSecret))
	protected.Post("/logout", authHandler.Logout)
	protected.Get("/me", authHandler.GetCurrentUser)
	protected.Get("/sessions", authHandler.GetSessions)
}
