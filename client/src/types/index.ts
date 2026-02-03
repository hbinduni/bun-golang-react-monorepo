// Client-side TypeScript types matching Go server models

// ============================================================================
// User Roles & Authentication
// ============================================================================

export type UserRole = 'admin' | 'user' | 'moderator'

export type OAuthProvider = 'google' | 'facebook' | 'twitter'

export type AuthMethod = 'email' | 'oauth'

// ============================================================================
// Core Entity Types
// ============================================================================

export interface User {
  id: string // TypeID: user_xxx
  email: string
  name: string
  role: UserRole
  emailVerified: boolean
  avatarUrl?: string
  createdAt: Date
  updatedAt: Date
}

export interface PublicUser {
  id: string
  name: string
  avatarUrl?: string
  createdAt: Date
}

export type ItemStatus = 'active' | 'completed' | 'archived'

export interface Item {
  id: string // TypeID: item_xxx
  userId: string
  title: string
  description: string
  status: ItemStatus
  createdAt: Date
  updatedAt: Date
}

export interface OAuthAccount {
  id: string // TypeID: oauth_xxx
  userId: string
  provider: OAuthProvider
  providerAccountId: string
  expiresAt?: Date
  createdAt: Date
  updatedAt: Date
}

export interface Session {
  id: string // TypeID: sess_xxx
  userId: string
  userAgent?: string
  ipAddress?: string
  expiresAt: Date
  createdAt: Date
}

// ============================================================================
// Authentication Request/Response Types
// ============================================================================

export interface LoginRequest {
  email: string
  password: string
}

export interface RegisterRequest {
  email: string
  password: string
  name: string
}

export interface AuthResponse {
  user: User
  accessToken: string
  refreshToken: string
  expiresIn: number // seconds
}

export interface RefreshTokenRequest {
  refreshToken: string
}

export interface RefreshTokenResponse {
  accessToken: string
  expiresIn: number // seconds
}

export interface OAuthUrlResponse {
  url: string
  state: string
}

export interface OAuthCallbackParams {
  code: string
  state: string
}

// ============================================================================
// JWT Types
// ============================================================================

export type TokenType = 'access' | 'refresh'

export interface JwtPayload {
  sub: string // User ID
  email: string
  role: UserRole
  type: TokenType
  iat: number
  exp: number
}

// ============================================================================
// API Response Types
// ============================================================================

export interface ApiResponse<T> {
  success: boolean
  data?: T
  error?: string
  message?: string
}

export interface PaginatedResponse<T> {
  success: boolean
  data: T[]
  pagination: {
    page: number
    limit: number
    total: number
    totalPages: number
  }
}

export interface ApiError {
  success: false
  error: string
  code?: string
  details?: Record<string, string[]>
}
