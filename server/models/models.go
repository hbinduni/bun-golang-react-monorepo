package models

import (
	"time"
)

// ============================================================================
// User Models
// ============================================================================

type UserRole string

const (
	RoleAdmin     UserRole = "admin"
	RoleUser      UserRole = "user"
	RoleModerator UserRole = "moderator"
)

type User struct {
	ID            string    `json:"id"` // TypeID: user_xxx
	Email         string    `json:"email"`
	PasswordHash  *string   `json:"-"` // Never send to client
	Name          string    `json:"name"`
	AvatarURL     *string   `json:"avatarUrl,omitempty"`
	Role          UserRole  `json:"role"`
	EmailVerified bool      `json:"emailVerified"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// PublicUser is a safe user profile for public display
type PublicUser struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	AvatarURL *string   `json:"avatarUrl,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}

// ============================================================================
// Item Models
// ============================================================================

type ItemStatus string

const (
	ItemStatusActive    ItemStatus = "active"
	ItemStatusCompleted ItemStatus = "completed"
	ItemStatusArchived  ItemStatus = "archived"
)

type Item struct {
	ID          string     `json:"id"` // TypeID: item_xxx
	UserID      string     `json:"userId"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      ItemStatus `json:"status"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

// ============================================================================
// OAuth Models
// ============================================================================

type OAuthProvider string

const (
	OAuthProviderGoogle   OAuthProvider = "google"
	OAuthProviderFacebook OAuthProvider = "facebook"
	OAuthProviderTwitter  OAuthProvider = "twitter"
)

type OAuthAccount struct {
	ID                string        `json:"id"` // TypeID: oauth_xxx
	UserID            string        `json:"userId"`
	Provider          OAuthProvider `json:"provider"`
	ProviderAccountID string        `json:"providerAccountId"`
	AccessToken       *string       `json:"-"` // Never send to client
	RefreshToken      *string       `json:"-"` // Never send to client
	ExpiresAt         *time.Time    `json:"expiresAt,omitempty"`
	CreatedAt         time.Time     `json:"createdAt"`
	UpdatedAt         time.Time     `json:"updatedAt"`
}

// ============================================================================
// Session Models
// ============================================================================

type Session struct {
	ID        string    `json:"id"` // TypeID: sess_xxx
	UserID    string    `json:"userId"`
	UserAgent *string   `json:"userAgent,omitempty"`
	IPAddress *string   `json:"ipAddress,omitempty"`
	ExpiresAt time.Time `json:"expiresAt"`
	CreatedAt time.Time `json:"createdAt"`
}

// ============================================================================
// Authentication Request/Response Models
// ============================================================================

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	User         User   `json:"user"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int    `json:"expiresIn"` // seconds
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type RefreshTokenResponse struct {
	AccessToken string `json:"accessToken"`
	ExpiresIn   int    `json:"expiresIn"` // seconds
}

// ============================================================================
// JWT Models
// ============================================================================

type TokenType string

const (
	TokenTypeAccess  TokenType = "access"
	TokenTypeRefresh TokenType = "refresh"
)

type JWTPayload struct {
	Subject   string    `json:"sub"` // User ID
	Email     string    `json:"email"`
	Role      UserRole  `json:"role"`
	Type      TokenType `json:"type"`
	IssuedAt  int64     `json:"iat"`
	ExpiresAt int64     `json:"exp"`
}
