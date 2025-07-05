package postgresql

import (
	"cloudtrail-enrichment-api-golang/internal/pkg/logger"
	"cloudtrail-enrichment-api-golang/models"
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq" // Importa el driver de PostgreSQL
)

// AuthPostgresRepository implementa la interfaz AuthRepository para PostgreSQL.
type AuthPostgresRepository struct {
	db *sql.DB
}

// NewAuthPostgresRepository crea una nueva instancia de AuthPostgresRepository.
// Recibe un *sql.DB ya inicializado para compartir la conexión.
func NewAuthPostgresRepository(db *sql.DB) *AuthPostgresRepository {
	return &AuthPostgresRepository{db: db}
}

// InsertUser inserta un nuevo usuario en la base de datos.
func (r *AuthPostgresRepository) InsertUser(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (uuid, email, password_hash, role, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	err := r.db.QueryRowContext(ctx, query, user.UUID, user.Email, user.PasswordHash, user.Role, user.CreatedAt, user.UpdatedAt).Scan(&user.ID)
	if err != nil {
		logger.ErrorLog.Printf("Error al insertar usuario en la base de datos: %v", err)
		return fmt.Errorf("error al insertar usuario: %w", err)
	}
	logger.InfoLog.Printf("Usuario insertado en DB con ID: %d y UUID: %s", user.ID, user.UUID)
	return nil
}

// GetUserByEmail recupera un usuario por su dirección de correo electrónico.
func (r *AuthPostgresRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT id, uuid, email, password_hash, role, created_at, updated_at FROM users WHERE email = $1`
	row := r.db.QueryRowContext(ctx, query, email)
	user := &models.User{}
	err := row.Scan(&user.ID, &user.UUID, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.InfoLog.Printf("Usuario con email %s no encontrado en la DB.", email)
			return nil, sql.ErrNoRows
		}
		logger.ErrorLog.Printf("Error al obtener usuario por email %s desde la base de datos: %v", email, err)
		return nil, fmt.Errorf("error al obtener usuario por email: %w", err)
	}
	logger.InfoLog.Printf("Usuario con email %s obtenido de DB.", email)
	return user, nil
}

// GetUserByUUID recupera un usuario por su UUID.
func (r *AuthPostgresRepository) GetUserByUUID(ctx context.Context, uuid string) (*models.User, error) {
	query := `SELECT id, uuid, email, password_hash, role, created_at, updated_at FROM users WHERE uuid = $1`
	row := r.db.QueryRowContext(ctx, query, uuid)
	user := &models.User{}
	err := row.Scan(&user.ID, &user.UUID, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.InfoLog.Printf("Usuario con UUID %s no encontrado en la DB.", uuid)
			return nil, sql.ErrNoRows
		}
		logger.ErrorLog.Printf("Error al obtener usuario por UUID %s desde la base de datos: %v", uuid, err)
		return nil, fmt.Errorf("error al obtener usuario por UUID: %w", err)
	}
	logger.InfoLog.Printf("Usuario con UUID %s obtenido de DB.", uuid)
	return user, nil
}

// InsertToken inserta un nuevo token en la base de datos.
func (r *AuthPostgresRepository) InsertToken(ctx context.Context, token *models.Token) error {
	// Asegurarse de que se inserta el 'role'
	query := `INSERT INTO tokens (user_id, email, token, token_hash, expiry, created_at, updated_at, role) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	token.CreatedAt = time.Now()
	token.UpdatedAt = time.Now()
	_, err := r.db.ExecContext(ctx, query, token.UserID, token.Email, token.Token, token.TokenHash, token.Expiry, token.CreatedAt, token.UpdatedAt, token.Role)
	if err != nil {
		logger.ErrorLog.Printf("Error al insertar token en la base de datos: %v", err)
		return fmt.Errorf("error al insertar token: %w", err)
	}
	logger.InfoLog.Printf("Token insertado en DB para UserID: %d", token.UserID)
	return nil
}

// GetTokenByTokenHash recupera un token por su hash desde la base de datos.
func (r *AuthPostgresRepository) GetTokenByTokenHash(ctx context.Context, tokenHash string) (*models.Token, error) {
	query := `SELECT id, user_id, email, token, token_hash, expiry, created_at, updated_at, role FROM tokens WHERE token_hash = $1`
	row := r.db.QueryRowContext(ctx, query, tokenHash)
	token := &models.Token{}
	err := row.Scan(&token.ID, &token.UserID, &token.Email, &token.Token, &token.TokenHash, &token.Expiry, &token.CreatedAt, &token.UpdatedAt, &token.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.InfoLog.Printf("Token con hash %s no encontrado en la DB.", tokenHash)
			return nil, sql.ErrNoRows
		}
		logger.ErrorLog.Printf("Error al obtener token por hash %s desde la base de datos: %v", tokenHash, err)
		return nil, fmt.Errorf("error al obtener token por hash: %w", err)
	}
	logger.InfoLog.Printf("Token con hash %s obtenido de DB para UserID: %d", tokenHash, token.UserID)
	return token, nil
}

// DeleteTokensByUserID elimina todos los tokens asociados a un UserID específico.
func (r *AuthPostgresRepository) DeleteTokensByUserID(ctx context.Context, userID int) error {
	query := `DELETE FROM tokens WHERE user_id = $1`
	fmt.Println("Ejecutando consulta EN REPO:", query, "con UserID:", userID)

	_, err := r.db.ExecContext(ctx, query, userID)
	fmt.Println("ERROR EJECUTANDO DELETE:", err)
	if err != nil {
		logger.ErrorLog.Printf("Error al eliminar tokens para el usuario %d desde la base de datos: %v", userID, err)
		return fmt.Errorf("error al eliminar tokens por UserID: %w", err)
	}
	logger.InfoLog.Printf("Tokens eliminados para UserID: %d", userID)
	return nil
}

// GetTokenByToken recupera un token por su valor de texto plano desde la base de datos.
func (r *AuthPostgresRepository) GetTokenByToken(ctx context.Context, tokenString string) (*models.Token, error) {
	query := `SELECT id, user_id, email, token, token_hash, expiry, created_at, updated_at, role FROM tokens WHERE token = $1`
	row := r.db.QueryRowContext(ctx, query, tokenString)
	token := &models.Token{}
	err := row.Scan(&token.ID, &token.UserID, &token.Email, &token.Token, &token.TokenHash, &token.Expiry, &token.CreatedAt, &token.UpdatedAt, &token.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.InfoLog.Printf("Token con valor %s no encontrado en la DB.", tokenString)
			return nil, sql.ErrNoRows
		}
		logger.ErrorLog.Printf("Error al obtener token por valor %s desde la base de datos: %v", tokenString, err)
		return nil, fmt.Errorf("error al obtener token por valor: %w", err)
	}
	logger.InfoLog.Printf("Token con valor %s obtenido de DB para UserID: %d", tokenString, token.UserID)
	return token, nil
}

// GetUserForToken recupera un usuario por su ID.
func (r *AuthPostgresRepository) GetUserForToken(ctx context.Context, userID int) (*models.User, error) {
	query := `SELECT id, uuid, email, password_hash, role, created_at, updated_at FROM users WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, userID)
	user := &models.User{}
	err := row.Scan(&user.ID, &user.UUID, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.InfoLog.Printf("Usuario con ID %d no encontrado en la DB.", userID)
			return nil, sql.ErrNoRows
		}
		logger.ErrorLog.Printf("Error al obtener usuario por ID %d desde la base de datos: %v", userID, err)
		return nil, fmt.Errorf("error al obtener usuario por ID: %w", err)
	}
	logger.InfoLog.Printf("Usuario con ID %d obtenido de DB.", userID)
	return user, nil
}

// Puedes añadir una función para insertar tokens si decides persistir tokens de refresco o gestionar listas negras.
// func (r *AuthPostgresRepository) InsertToken(ctx context.Context, token *models.Token) error {
// 	query := `INSERT INTO tokens (user_id, token, token_hash, expiry, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`
// 	token.CreatedAt = time.Now()
// 	token.UpdatedAt = time.Now()
// 	_, err := r.db.ExecContext(ctx, query, token.UserID, token.Token, token.TokenHash, token.Expiry, token.CreatedAt, token.UpdatedAt)
// 	if err != nil {
// 		logger.ErrorLog.Printf("Error al insertar token en la base de datos: %v", err)
// 		return fmt.Errorf("error al insertar token: %w", err)
// 	}
// 	return nil
// }
