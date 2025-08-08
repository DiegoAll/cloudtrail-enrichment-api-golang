package postgresql

import (
	"cloudtrail-enrichment-api-golang/internal/pkg/logger"
	"cloudtrail-enrichment-api-golang/models"
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// AuthPostgresRepository implements the AuthRepository interface for PostgreSQL.
type AuthPostgresRepository struct {
	db *sql.DB
}

// NewAuthPostgresRepository creates a new instance of AuthPostgresRepository.
func NewAuthPostgresRepository(db *sql.DB) *AuthPostgresRepository {
	return &AuthPostgresRepository{db: db}
}

// InsertUser inserts a new user into the database.
func (r *AuthPostgresRepository) InsertUser(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (uuid, email, password_hash, role, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	err := r.db.QueryRowContext(ctx, query, user.UUID, user.Email, user.PasswordHash, user.Role, user.CreatedAt, user.UpdatedAt).Scan(&user.ID)
	if err != nil {
		logger.ErrorLog.Printf("Error inserting user into database: %v", err)
		return fmt.Errorf("error inserting user: %w", err)
	}
	logger.InfoLog.Printf("User inserted into DB with ID: %d and UUID: %s", user.ID, user.UUID)
	return nil
}

// GetUserByEmail retrieves a user by their email address.
func (r *AuthPostgresRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT id, uuid, email, password_hash, role, created_at, updated_at FROM users WHERE email = $1`
	row := r.db.QueryRowContext(ctx, query, email)
	user := &models.User{}
	err := row.Scan(&user.ID, &user.UUID, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.InfoLog.Printf("User with email %s not found in DB.", email)
			return nil, sql.ErrNoRows
		}
		logger.ErrorLog.Printf("Error getting user by email %s from database: %v", email, err)
		return nil, fmt.Errorf("error getting user by email: %w", err)
	}
	logger.InfoLog.Printf("User with email %s retrieved from DB.", email)
	return user, nil
}

// GetUserByUUID retrieves a user by their UUID.
func (r *AuthPostgresRepository) GetUserByUUID(ctx context.Context, uuid string) (*models.User, error) {
	query := `SELECT id, uuid, email, password_hash, role, created_at, updated_at FROM users WHERE uuid = $1`
	row := r.db.QueryRowContext(ctx, query, uuid)
	user := &models.User{}
	err := row.Scan(&user.ID, &user.UUID, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.InfoLog.Printf("User with UUID %s not found in DB.", uuid)
			return nil, sql.ErrNoRows
		}
		logger.ErrorLog.Printf("Error getting user by UUID %s from database: %v", uuid, err)
		return nil, fmt.Errorf("error getting user by UUID: %w", err)
	}
	logger.InfoLog.Printf("User with UUID %s retrieved from DB.", uuid)
	return user, nil
}

// InsertToken inserts a new token into the database.
func (r *AuthPostgresRepository) InsertToken(ctx context.Context, token *models.Token) error {
	// Ensure 'role' is inserted
	query := `INSERT INTO tokens (user_id, email, token, token_hash, expiry, created_at, updated_at, role) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	token.CreatedAt = time.Now()
	token.UpdatedAt = time.Now()
	_, err := r.db.ExecContext(ctx, query, token.UserID, token.Email, token.Token, token.TokenHash, token.Expiry, token.CreatedAt, token.UpdatedAt, token.Role)
	if err != nil {
		logger.ErrorLog.Printf("Error inserting token into database: %v", err)
		return fmt.Errorf("error inserting token: %w", err)
	}
	logger.InfoLog.Printf("Token inserted into DB for UserID: %d", token.UserID)
	return nil
}

// GetTokenByTokenHash retrieves a token by its hash from the database.
func (r *AuthPostgresRepository) GetTokenByTokenHash(ctx context.Context, tokenHash string) (*models.Token, error) {
	query := `SELECT id, user_id, email, token, token_hash, expiry, created_at, updated_at, role FROM tokens WHERE token_hash = $1`
	row := r.db.QueryRowContext(ctx, query, tokenHash)
	token := &models.Token{}
	err := row.Scan(&token.ID, &token.UserID, &token.Email, &token.Token, &token.TokenHash, &token.Expiry, &token.CreatedAt, &token.UpdatedAt, &token.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.InfoLog.Printf("Token with hash %s not found in DB.", tokenHash)
			return nil, sql.ErrNoRows
		}
		logger.ErrorLog.Printf("Error getting token by hash %s from database: %v", tokenHash, err)
		return nil, fmt.Errorf("error getting token by hash: %w", err)
	}
	logger.InfoLog.Printf("Token with hash %s retrieved from DB for UserID: %d", tokenHash, token.UserID)
	return token, nil
}

// DeleteTokensByUserID deletes all tokens associated with a specific UserID.
func (r *AuthPostgresRepository) DeleteTokensByUserID(ctx context.Context, userID int) error {
	query := `DELETE FROM tokens WHERE user_id = $1`
	fmt.Println("Executing query IN REPO:", query, "with UserID:", userID)

	_, err := r.db.ExecContext(ctx, query, userID)
	fmt.Println("ERROR EXECUTING DELETE:", err)
	if err != nil {
		logger.ErrorLog.Printf("Error deleting tokens for user %d from database: %v", userID, err)
		return fmt.Errorf("error deleting tokens by UserID: %w", err)
	}
	logger.InfoLog.Printf("Tokens deleted for UserID: %d", userID)
	return nil
}

// GetTokenByToken retrieves a token by its plain text value from the database.
func (r *AuthPostgresRepository) GetTokenByToken(ctx context.Context, tokenString string) (*models.Token, error) {
	query := `SELECT id, user_id, email, token, token_hash, expiry, created_at, updated_at, role FROM tokens WHERE token = $1`
	row := r.db.QueryRowContext(ctx, query, tokenString)
	token := &models.Token{}
	err := row.Scan(&token.ID, &token.UserID, &token.Email, &token.Token, &token.TokenHash, &token.Expiry, &token.CreatedAt, &token.UpdatedAt, &token.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.InfoLog.Printf("Token with value %s not found in DB.", tokenString)
			return nil, sql.ErrNoRows
		}
		logger.ErrorLog.Printf("Error getting token by value %s from database: %v", tokenString, err)
		return nil, fmt.Errorf("error getting token by value: %w", err)
	}
	logger.InfoLog.Printf("Token with value %s retrieved from DB for UserID: %d", tokenString, token.UserID)
	return token, nil
}

// GetUserForToken retrieves a user by their ID.
func (r *AuthPostgresRepository) GetUserForToken(ctx context.Context, userID int) (*models.User, error) {
	query := `SELECT id, uuid, email, password_hash, role, created_at, updated_at FROM users WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, userID)
	user := &models.User{}
	err := row.Scan(&user.ID, &user.UUID, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.InfoLog.Printf("User with ID %d not found in DB.", userID)
			return nil, sql.ErrNoRows
		}
		logger.ErrorLog.Printf("Error getting user by ID %d from database: %v", userID, err)
		return nil, fmt.Errorf("error getting user by ID: %w", err)
	}
	logger.InfoLog.Printf("User with ID %d retrieved from DB.", userID)
	return user, nil
}
