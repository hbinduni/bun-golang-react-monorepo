package database

import (
	"context"
	"fmt"

	"github.com/binduni/bun-golang-react-monorepo/server/models"
)

// ============================================================================
// User Queries
// ============================================================================

func (db *DB) CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (id, email, password_hash, name, avatar_url, role, email_verified)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING created_at, updated_at
	`
	return db.Pool.QueryRow(ctx, query,
		user.ID, user.Email, user.PasswordHash, user.Name, user.AvatarURL, user.Role, user.EmailVerified,
	).Scan(&user.CreatedAt, &user.UpdatedAt)
}

func (db *DB) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	query := `
		SELECT id, email, password_hash, name, avatar_url, role, email_verified, created_at, updated_at
		FROM users WHERE id = $1
	`
	err := db.Pool.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.AvatarURL,
		&user.Role, &user.EmailVerified, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (db *DB) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	query := `
		SELECT id, email, password_hash, name, avatar_url, role, email_verified, created_at, updated_at
		FROM users WHERE email = $1
	`
	err := db.Pool.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.AvatarURL,
		&user.Role, &user.EmailVerified, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// ============================================================================
// Session Queries
// ============================================================================

func (db *DB) CreateSession(ctx context.Context, session *models.Session) error {
	query := `
		INSERT INTO sessions (id, user_id, user_agent, ip_address, expires_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING created_at
	`
	return db.Pool.QueryRow(ctx, query,
		session.ID, session.UserID, session.UserAgent, session.IPAddress, session.ExpiresAt,
	).Scan(&session.CreatedAt)
}

func (db *DB) GetSessionByID(ctx context.Context, id string) (*models.Session, error) {
	var session models.Session
	query := `
		SELECT id, user_id, user_agent, ip_address, expires_at, created_at
		FROM sessions WHERE id = $1
	`
	err := db.Pool.QueryRow(ctx, query, id).Scan(
		&session.ID, &session.UserID, &session.UserAgent, &session.IPAddress,
		&session.ExpiresAt, &session.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (db *DB) GetUserSessions(ctx context.Context, userID string) ([]*models.Session, error) {
	query := `
		SELECT id, user_id, user_agent, ip_address, expires_at, created_at
		FROM sessions WHERE user_id = $1 AND expires_at > NOW()
		ORDER BY created_at DESC
	`
	rows, err := db.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []*models.Session
	for rows.Next() {
		var session models.Session
		if err := rows.Scan(
			&session.ID, &session.UserID, &session.UserAgent, &session.IPAddress,
			&session.ExpiresAt, &session.CreatedAt,
		); err != nil {
			return nil, err
		}
		sessions = append(sessions, &session)
	}
	return sessions, rows.Err()
}

func (db *DB) DeleteSession(ctx context.Context, id string) error {
	query := `DELETE FROM sessions WHERE id = $1`
	result, err := db.Pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("session not found")
	}
	return nil
}

// ============================================================================
// Item Queries
// ============================================================================

func (db *DB) CreateItem(ctx context.Context, item *models.Item) error {
	query := `
		INSERT INTO items (id, user_id, title, description, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING created_at, updated_at
	`
	return db.Pool.QueryRow(ctx, query,
		item.ID, item.UserID, item.Title, item.Description, item.Status,
	).Scan(&item.CreatedAt, &item.UpdatedAt)
}

func (db *DB) GetItemByID(ctx context.Context, id string) (*models.Item, error) {
	var item models.Item
	query := `
		SELECT id, user_id, title, description, status, created_at, updated_at
		FROM items WHERE id = $1
	`
	err := db.Pool.QueryRow(ctx, query, id).Scan(
		&item.ID, &item.UserID, &item.Title, &item.Description,
		&item.Status, &item.CreatedAt, &item.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (db *DB) GetUserItems(ctx context.Context, userID string) ([]*models.Item, error) {
	query := `
		SELECT id, user_id, title, description, status, created_at, updated_at
		FROM items WHERE user_id = $1
		ORDER BY created_at DESC
	`
	rows, err := db.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*models.Item
	for rows.Next() {
		var item models.Item
		if err := rows.Scan(
			&item.ID, &item.UserID, &item.Title, &item.Description,
			&item.Status, &item.CreatedAt, &item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	return items, rows.Err()
}

func (db *DB) UpdateItem(ctx context.Context, item *models.Item) error {
	query := `
		UPDATE items
		SET title = $2, description = $3, status = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING updated_at
	`
	return db.Pool.QueryRow(ctx, query,
		item.ID, item.Title, item.Description, item.Status,
	).Scan(&item.UpdatedAt)
}

func (db *DB) DeleteItem(ctx context.Context, id string) error {
	query := `DELETE FROM items WHERE id = $1`
	result, err := db.Pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("item not found")
	}
	return nil
}

// ============================================================================
// OAuth Queries
// ============================================================================

func (db *DB) CreateOAuthAccount(ctx context.Context, account *models.OAuthAccount) error {
	query := `
		INSERT INTO oauth_accounts (id, user_id, provider, provider_account_id, access_token, refresh_token, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING created_at, updated_at
	`
	return db.Pool.QueryRow(ctx, query,
		account.ID, account.UserID, account.Provider, account.ProviderAccountID,
		account.AccessToken, account.RefreshToken, account.ExpiresAt,
	).Scan(&account.CreatedAt, &account.UpdatedAt)
}

func (db *DB) GetOAuthAccount(ctx context.Context, provider models.OAuthProvider, providerAccountID string) (*models.OAuthAccount, error) {
	var account models.OAuthAccount
	query := `
		SELECT id, user_id, provider, provider_account_id, access_token, refresh_token, expires_at, created_at, updated_at
		FROM oauth_accounts WHERE provider = $1 AND provider_account_id = $2
	`
	err := db.Pool.QueryRow(ctx, query, provider, providerAccountID).Scan(
		&account.ID, &account.UserID, &account.Provider, &account.ProviderAccountID,
		&account.AccessToken, &account.RefreshToken, &account.ExpiresAt,
		&account.CreatedAt, &account.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &account, nil
}
