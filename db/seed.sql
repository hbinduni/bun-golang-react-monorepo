-- Seed data for the application with TypeID format
-- TypeIDs are K-sortable identifiers: {prefix}_{26-char-base32-suffix}
--
-- NOTE: These are sample TypeIDs for testing. In production,
-- TypeIDs are generated application-side using typeid-js.
--
-- Password for test users: "Password123" (bcrypt hash below)

-- ============================================================================
-- Sample Users
-- ============================================================================

INSERT INTO users (id, email, password_hash, name, role, email_verified, avatar_url) VALUES
  -- Admin user (Password: Password123)
  (
    'user_01h2xcejqtf2nbrexx3vqjhp41',
    'admin@example.com',
    '$2a$12$LQv3c1yqBwEC3QhHjbZGP.4Z0pqTHdT6MjF7NtPQw6/YvYQR5YQWG',
    'Admin User',
    'admin',
    TRUE,
    'https://api.dicebear.com/7.x/avataaars/svg?seed=admin'
  ),
  -- Regular users
  (
    'user_01h2xcejqtf2nbrexx3vqjhp42',
    'john.doe@example.com',
    '$2a$12$LQv3c1yqBwEC3QhHjbZGP.4Z0pqTHdT6MjF7NtPQw6/YvYQR5YQWG',
    'John Doe',
    'user',
    TRUE,
    'https://api.dicebear.com/7.x/avataaars/svg?seed=john'
  ),
  (
    'user_01h2xcejqtf2nbrexx3vqjhp43',
    'jane.smith@example.com',
    '$2a$12$LQv3c1yqBwEC3QhHjbZGP.4Z0pqTHdT6MjF7NtPQw6/YvYQR5YQWG',
    'Jane Smith',
    'user',
    TRUE,
    'https://api.dicebear.com/7.x/avataaars/svg?seed=jane'
  ),
  -- Moderator
  (
    'user_01h2xcejqtf2nbrexx3vqjhp44',
    'mod@example.com',
    '$2a$12$LQv3c1yqBwEC3QhHjbZGP.4Z0pqTHdT6MjF7NtPQw6/YvYQR5YQWG',
    'Moderator User',
    'moderator',
    TRUE,
    'https://api.dicebear.com/7.x/avataaars/svg?seed=mod'
  ),
  -- OAuth-only user (no password)
  (
    'user_01h2xcejqtf2nbrexx3vqjhp45',
    'oauth.user@gmail.com',
    NULL,
    'OAuth User',
    'user',
    TRUE,
    'https://lh3.googleusercontent.com/a/default-user'
  )
ON CONFLICT (email) DO UPDATE SET
  name = EXCLUDED.name,
  role = EXCLUDED.role;

-- ============================================================================
-- Sample OAuth Accounts
-- ============================================================================

INSERT INTO oauth_accounts (id, user_id, provider, provider_account_id) VALUES
  (
    'oauth_01h2xcejqtf2nbrexx3vqjhp01',
    'user_01h2xcejqtf2nbrexx3vqjhp45',
    'google',
    '117098765432109876543'
  )
ON CONFLICT (provider, provider_account_id) DO NOTHING;

-- ============================================================================
-- Sample Items
-- ============================================================================

INSERT INTO items (id, user_id, title, description, status) VALUES
  -- John's items
  (
    'item_01h2xcejqtf2nbrexx3vqjhp01',
    'user_01h2xcejqtf2nbrexx3vqjhp42',
    'Complete project setup',
    'Set up the monorepo with Bun, Hono, and React',
    'completed'
  ),
  (
    'item_01h2xcejqtf2nbrexx3vqjhp02',
    'user_01h2xcejqtf2nbrexx3vqjhp42',
    'Implement authentication',
    'Add JWT auth with social login support',
    'active'
  ),
  -- Jane's items
  (
    'item_01h2xcejqtf2nbrexx3vqjhp03',
    'user_01h2xcejqtf2nbrexx3vqjhp43',
    'Design system components',
    'Create reusable UI components with Tailwind',
    'active'
  ),
  (
    'item_01h2xcejqtf2nbrexx3vqjhp04',
    'user_01h2xcejqtf2nbrexx3vqjhp43',
    'Write API documentation',
    'Document all REST endpoints with examples',
    'active'
  ),
  -- Admin's items
  (
    'item_01h2xcejqtf2nbrexx3vqjhp05',
    'user_01h2xcejqtf2nbrexx3vqjhp41',
    'Set up CI/CD pipeline',
    'Configure GitHub Actions for automated deployments',
    'completed'
  ),
  (
    'item_01h2xcejqtf2nbrexx3vqjhp06',
    'user_01h2xcejqtf2nbrexx3vqjhp41',
    'Security audit',
    'Review authentication flow and fix vulnerabilities',
    'active'
  )
ON CONFLICT DO NOTHING;

-- ============================================================================
-- Verify data
-- ============================================================================

-- You can run these queries to verify the seed data:
-- SELECT id, email, name, role FROM users;
-- SELECT id, user_id, provider FROM oauth_accounts;
-- SELECT id, user_id, title, status FROM items;
