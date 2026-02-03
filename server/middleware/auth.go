package middleware

import (
	"strings"

	"github.com/binduni/bun-golang-react-monorepo/server/models"
	"github.com/binduni/bun-golang-react-monorepo/server/utils"
	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware validates JWT tokens and attaches user info to context
func AuthMiddleware(jwtSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse("Missing authorization header"))
		}

		tokenString, err := utils.ExtractBearerToken(authHeader)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse("Invalid authorization header format"))
		}

		claims, err := utils.ValidateToken(tokenString, jwtSecret)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse("Invalid or expired token"))
		}

		// Verify it's an access token
		if claims.Type != models.TokenTypeAccess {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse("Invalid token type"))
		}

		// Store user info in context
		c.Locals("userID", claims.UserID)
		c.Locals("userEmail", claims.Email)
		c.Locals("userRole", claims.Role)

		return c.Next()
	}
}

// GetUserID retrieves the authenticated user ID from context
func GetUserID(c *fiber.Ctx) string {
	if userID, ok := c.Locals("userID").(string); ok {
		return userID
	}
	return ""
}

// GetUserRole retrieves the authenticated user role from context
func GetUserRole(c *fiber.Ctx) models.UserRole {
	if role, ok := c.Locals("userRole").(models.UserRole); ok {
		return role
	}
	return models.RoleUser
}

// RequireRole middleware ensures the user has a specific role
func RequireRole(allowedRoles ...models.UserRole) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := GetUserRole(c)

		for _, role := range allowedRoles {
			if userRole == role {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(models.ErrorResponse("Insufficient permissions"))
	}
}

// OptionalAuth middleware that doesn't fail if token is missing
func OptionalAuth(jwtSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Next()
		}

		tokenString, err := utils.ExtractBearerToken(authHeader)
		if err != nil {
			return c.Next()
		}

		claims, err := utils.ValidateToken(tokenString, jwtSecret)
		if err != nil || claims.Type != models.TokenTypeAccess {
			return c.Next()
		}

		// Store user info in context
		c.Locals("userID", claims.UserID)
		c.Locals("userEmail", claims.Email)
		c.Locals("userRole", claims.Role)

		return c.Next()
	}
}

// GetClientIP extracts the client IP address
func GetClientIP(c *fiber.Ctx) string {
	// Check X-Forwarded-For header first (for proxies)
	if xff := c.Get("X-Forwarded-For"); xff != "" {
		// Take the first IP if multiple
		if idx := strings.Index(xff, ","); idx != -1 {
			return strings.TrimSpace(xff[:idx])
		}
		return strings.TrimSpace(xff)
	}

	// Check X-Real-IP header
	if xri := c.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fall back to RemoteIP
	return c.IP()
}
