package routes

import (
	"github.com/binduni/bun-golang-react-monorepo/server/config"
	"github.com/binduni/bun-golang-react-monorepo/server/database"
	"github.com/binduni/bun-golang-react-monorepo/server/handlers"
	"github.com/binduni/bun-golang-react-monorepo/server/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupItemsRoutes(router fiber.Router, cfg *config.Config, db *database.DB) {
	itemsHandler := handlers.NewItemsHandler(db, cfg)

	// All items routes require authentication
	router.Use(middleware.AuthMiddleware(cfg.JWTSecret))

	router.Get("/", itemsHandler.ListItems)
	router.Get("/:id", itemsHandler.GetItem)
	router.Post("/", itemsHandler.CreateItem)
	router.Put("/:id", itemsHandler.UpdateItem)
	router.Delete("/:id", itemsHandler.DeleteItem)
}
