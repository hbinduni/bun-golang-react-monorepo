package handlers

import (
	"github.com/binduni/bun-golang-react-monorepo/server/config"
	"github.com/binduni/bun-golang-react-monorepo/server/database"
	"github.com/binduni/bun-golang-react-monorepo/server/middleware"
	"github.com/binduni/bun-golang-react-monorepo/server/models"
	"github.com/binduni/bun-golang-react-monorepo/server/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

type ItemsHandler struct {
	db     *database.DB
	config *config.Config
}

func NewItemsHandler(db *database.DB, cfg *config.Config) *ItemsHandler {
	return &ItemsHandler{
		db:     db,
		config: cfg,
	}
}

// ListItems returns all items for the authenticated user
func (h *ItemsHandler) ListItems(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse("Unauthorized"))
	}

	items, err := h.db.GetUserItems(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to retrieve items"))
	}

	return c.JSON(models.SuccessResponse(items))
}

// GetItem returns a single item by ID
func (h *ItemsHandler) GetItem(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse("Unauthorized"))
	}

	itemID := c.Params("id")
	if itemID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Item ID is required"))
	}

	item, err := h.db.GetItemByID(c.Context(), itemID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse("Item not found"))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to retrieve item"))
	}

	// Verify ownership
	if item.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(models.ErrorResponse("Access denied"))
	}

	return c.JSON(models.SuccessResponse(item))
}

// CreateItem creates a new item
func (h *ItemsHandler) CreateItem(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse("Unauthorized"))
	}

	var req struct {
		Title       string            `json:"title"`
		Description string            `json:"description"`
		Status      models.ItemStatus `json:"status"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Invalid request body"))
	}

	// Validate input
	if req.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Title is required"))
	}

	// Default status if not provided
	if req.Status == "" {
		req.Status = models.ItemStatusActive
	}

	// Create item
	item := &models.Item{
		ID:          utils.NewItemID(),
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
	}

	if err := h.db.CreateItem(c.Context(), item); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to create item"))
	}

	return c.Status(fiber.StatusCreated).JSON(models.SuccessResponse(item))
}

// UpdateItem updates an existing item
func (h *ItemsHandler) UpdateItem(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse("Unauthorized"))
	}

	itemID := c.Params("id")
	if itemID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Item ID is required"))
	}

	// Get existing item
	item, err := h.db.GetItemByID(c.Context(), itemID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse("Item not found"))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to retrieve item"))
	}

	// Verify ownership
	if item.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(models.ErrorResponse("Access denied"))
	}

	var req struct {
		Title       string            `json:"title"`
		Description string            `json:"description"`
		Status      models.ItemStatus `json:"status"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Invalid request body"))
	}

	// Update fields
	if req.Title != "" {
		item.Title = req.Title
	}
	if req.Description != "" {
		item.Description = req.Description
	}
	if req.Status != "" {
		item.Status = req.Status
	}

	if err := h.db.UpdateItem(c.Context(), item); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to update item"))
	}

	return c.JSON(models.SuccessResponse(item))
}

// DeleteItem deletes an item
func (h *ItemsHandler) DeleteItem(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse("Unauthorized"))
	}

	itemID := c.Params("id")
	if itemID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Item ID is required"))
	}

	// Get existing item to verify ownership
	item, err := h.db.GetItemByID(c.Context(), itemID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse("Item not found"))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to retrieve item"))
	}

	// Verify ownership
	if item.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(models.ErrorResponse("Access denied"))
	}

	if err := h.db.DeleteItem(c.Context(), itemID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to delete item"))
	}

	return c.JSON(models.SuccessResponse(fiber.Map{
		"message": "Item deleted successfully",
	}))
}
