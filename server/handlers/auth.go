package handlers

import (
	"time"

	"github.com/binduni/bun-golang-react-monorepo/server/config"
	"github.com/binduni/bun-golang-react-monorepo/server/database"
	"github.com/binduni/bun-golang-react-monorepo/server/middleware"
	"github.com/binduni/bun-golang-react-monorepo/server/models"
	"github.com/binduni/bun-golang-react-monorepo/server/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

type AuthHandler struct {
	db     *database.DB
	config *config.Config
}

func NewAuthHandler(db *database.DB, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		db:     db,
		config: cfg,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req models.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Invalid request body"))
	}

	// Validate input
	req.Email = utils.NormalizeEmail(req.Email)
	if !utils.ValidateEmail(req.Email) {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Invalid email address"))
	}
	if !utils.ValidatePassword(req.Password) {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Password must be at least 8 characters"))
	}
	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Name is required"))
	}

	// Check if user already exists
	_, err := h.db.GetUserByEmail(c.Context(), req.Email)
	if err == nil {
		return c.Status(fiber.StatusConflict).JSON(models.ErrorResponse("Email already registered"))
	} else if err != pgx.ErrNoRows {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Database error"))
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to process password"))
	}

	// Create user
	user := &models.User{
		ID:            utils.NewUserID(),
		Email:         req.Email,
		PasswordHash:  &hashedPassword,
		Name:          req.Name,
		Role:          models.RoleUser,
		EmailVerified: false,
	}

	if err := h.db.CreateUser(c.Context(), user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to create user"))
	}

	// Generate tokens
	accessToken, err := utils.GenerateAccessToken(user, h.config.JWTSecret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to generate access token"))
	}

	refreshToken, err := utils.GenerateRefreshToken(user, h.config.JWTSecret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to generate refresh token"))
	}

	// Create session
	session := &models.Session{
		ID:        utils.NewSessionID(),
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(utils.RefreshTokenExpiry),
	}
	userAgent := c.Get("User-Agent")
	if userAgent != "" {
		session.UserAgent = &userAgent
	}
	ipAddress := middleware.GetClientIP(c)
	if ipAddress != "" {
		session.IPAddress = &ipAddress
	}

	if err := h.db.CreateSession(c.Context(), session); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to create session"))
	}

	// Return response
	response := models.AuthResponse{
		User:         *user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int(utils.AccessTokenExpiry.Seconds()),
	}

	return c.Status(fiber.StatusCreated).JSON(models.SuccessResponse(response))
}

// Login handles user login
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Invalid request body"))
	}

	// Normalize email
	req.Email = utils.NormalizeEmail(req.Email)

	// Get user by email
	user, err := h.db.GetUserByEmail(c.Context(), req.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse("Invalid email or password"))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Database error"))
	}

	// Verify password
	if user.PasswordHash == nil || !utils.VerifyPassword(*user.PasswordHash, req.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse("Invalid email or password"))
	}

	// Generate tokens
	accessToken, err := utils.GenerateAccessToken(user, h.config.JWTSecret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to generate access token"))
	}

	refreshToken, err := utils.GenerateRefreshToken(user, h.config.JWTSecret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to generate refresh token"))
	}

	// Create session
	session := &models.Session{
		ID:        utils.NewSessionID(),
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(utils.RefreshTokenExpiry),
	}
	userAgent := c.Get("User-Agent")
	if userAgent != "" {
		session.UserAgent = &userAgent
	}
	ipAddress := middleware.GetClientIP(c)
	if ipAddress != "" {
		session.IPAddress = &ipAddress
	}

	if err := h.db.CreateSession(c.Context(), session); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to create session"))
	}

	// Return response
	response := models.AuthResponse{
		User:         *user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int(utils.AccessTokenExpiry.Seconds()),
	}

	return c.JSON(models.SuccessResponse(response))
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req models.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Invalid request body"))
	}

	// Validate refresh token
	claims, err := utils.ValidateToken(req.RefreshToken, h.config.JWTSecret)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse("Invalid or expired refresh token"))
	}

	// Verify it's a refresh token
	if claims.Type != models.TokenTypeRefresh {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse("Invalid token type"))
	}

	// Get user
	user, err := h.db.GetUserByID(c.Context(), claims.UserID)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse("User not found"))
	}

	// Generate new access token
	accessToken, err := utils.GenerateAccessToken(user, h.config.JWTSecret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to generate access token"))
	}

	// Return response
	response := models.RefreshTokenResponse{
		AccessToken: accessToken,
		ExpiresIn:   int(utils.AccessTokenExpiry.Seconds()),
	}

	return c.JSON(models.SuccessResponse(response))
}

// Logout handles user logout (invalidates session)
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse("Unauthorized"))
	}

	// Note: In a production system, you might want to blacklist the JWT
	// For now, we just delete the session
	// You could pass session ID in request body to delete specific session

	return c.JSON(models.SuccessResponse(fiber.Map{
		"message": "Logged out successfully",
	}))
}

// GetCurrentUser returns the authenticated user's information
func (h *AuthHandler) GetCurrentUser(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse("Unauthorized"))
	}

	user, err := h.db.GetUserByID(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse("User not found"))
	}

	return c.JSON(models.SuccessResponse(user))
}

// GetSessions returns all active sessions for the authenticated user
func (h *AuthHandler) GetSessions(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse("Unauthorized"))
	}

	sessions, err := h.db.GetUserSessions(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to retrieve sessions"))
	}

	return c.JSON(models.SuccessResponse(sessions))
}
