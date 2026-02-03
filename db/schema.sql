-- Database schema for the application with TypeID and Authentication
-- TypeIDs are K-sortable, type-safe identifiers generated application-side
-- Format: {prefix}_{26-char-base32-uuidv7} (e.g., user_01h2xcejqtf2nbrexx3vqjhp41)

-- ============================================================================
-- Users Table
-- ============================================================================

CREATE TABLE IF NOT EXISTS users (
  -- TypeID format: user_xxx... (max ~30 chars)
  id VARCHAR(30) PRIMARY KEY,

  -- Authentication
  email VARCHAR(255) UNIQUE NOT NULL,
  password_hash VARCHAR(60), -- bcrypt hash (null for OAuth-only users)

  -- Profile
  name VARCHAR(255) NOT NULL,
  avatar_url TEXT,

  -- Authorization
  role VARCHAR(20) NOT NULL DEFAULT 'user' CHECK (role IN ('admin', 'user', 'moderator')),

  -- Email verification
  email_verified BOOLEAN NOT NULL DEFAULT FALSE,

  -- Timestamps
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ============================================================================
-- OAuth Accounts Table - Links external OAuth providers to users
-- ============================================================================

CREATE TABLE IF NOT EXISTS oauth_accounts (
  -- TypeID format: oauth_xxx...
  id VARCHAR(30) PRIMARY KEY,
  user_id VARCHAR(30) NOT NULL REFERENCES users(id) ON DELETE CASCADE,

  -- Provider info
  provider VARCHAR(20) NOT NULL CHECK (provider IN ('google', 'facebook', 'twitter')),
  provider_account_id VARCHAR(255) NOT NULL,

  -- Tokens (stored encrypted in production)
  access_token TEXT,
  refresh_token TEXT,
  expires_at TIMESTAMP WITH TIME ZONE,

  -- Timestamps
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

  -- Prevent duplicate provider accounts per user
  UNIQUE(provider, provider_account_id),
  UNIQUE(user_id, provider)
);

-- ============================================================================
-- Sessions Table - Tracks active user sessions for refresh tokens
-- ============================================================================

CREATE TABLE IF NOT EXISTS sessions (
  -- TypeID format: sess_xxx...
  id VARCHAR(30) PRIMARY KEY,
  user_id VARCHAR(30) NOT NULL REFERENCES users(id) ON DELETE CASCADE,

  -- Session metadata
  user_agent TEXT,
  ip_address VARCHAR(45), -- IPv6 max length

  -- Expiration
  expires_at TIMESTAMP WITH TIME ZONE NOT NULL,

  -- Timestamps
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ============================================================================
-- Items Table - Example resource owned by users
-- ============================================================================

CREATE TABLE IF NOT EXISTS items (
  -- TypeID format: item_xxx...
  id VARCHAR(30) PRIMARY KEY,
  user_id VARCHAR(30) NOT NULL REFERENCES users(id) ON DELETE CASCADE,

  -- Content
  title VARCHAR(255) NOT NULL,
  description TEXT,
  status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'completed', 'archived')),

  -- Timestamps
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ============================================================================
-- Indexes for Performance
-- ============================================================================

-- Users
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);

-- OAuth Accounts
CREATE INDEX IF NOT EXISTS idx_oauth_accounts_user_id ON oauth_accounts(user_id);
CREATE INDEX IF NOT EXISTS idx_oauth_accounts_provider ON oauth_accounts(provider, provider_account_id);

-- Sessions
CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON sessions(expires_at);

-- Items
CREATE INDEX IF NOT EXISTS idx_items_user_id ON items(user_id);
CREATE INDEX IF NOT EXISTS idx_items_status ON items(status);
CREATE INDEX IF NOT EXISTS idx_items_created_at ON items(created_at DESC);

-- ============================================================================
-- Trigger: Auto-update updated_at timestamp
-- ============================================================================

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = CURRENT_TIMESTAMP;
  RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply trigger to tables with updated_at
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
CREATE TRIGGER update_users_updated_at
  BEFORE UPDATE ON users
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_oauth_accounts_updated_at ON oauth_accounts;
CREATE TRIGGER update_oauth_accounts_updated_at
  BEFORE UPDATE ON oauth_accounts
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_items_updated_at ON items;
CREATE TRIGGER update_items_updated_at
  BEFORE UPDATE ON items
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at_column();

-- ============================================================================
-- Helper: Clean up expired sessions (run periodically via cron/scheduler)
-- ============================================================================

CREATE OR REPLACE FUNCTION cleanup_expired_sessions()
RETURNS INTEGER AS $$
DECLARE
  deleted_count INTEGER;
BEGIN
  DELETE FROM sessions WHERE expires_at < CURRENT_TIMESTAMP;
  GET DIAGNOSTICS deleted_count = ROW_COUNT;
  RETURN deleted_count;
END;
$$ LANGUAGE plpgsql;
