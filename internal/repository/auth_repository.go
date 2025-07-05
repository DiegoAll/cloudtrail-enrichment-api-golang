package repository

import (
	"cloudtrail-enrichment-api-golang/models"
	"context"
)

// AuthRepository define la interfaz para las operaciones relacionadas con la autenticación
// y ahora también para la gestión de tokens.
type AuthRepository interface {
	InsertUser(ctx context.Context, user *models.User) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByUUID(ctx context.Context, uuid string) (*models.User, error)

	// Métodos para la gestión de tokens
	InsertToken(ctx context.Context, token *models.Token) error
	GetTokenByTokenHash(ctx context.Context, tokenHash string) (*models.Token, error)
	DeleteTokensByUserID(ctx context.Context, userID int) error
	GetTokenByToken(ctx context.Context, tokenString string) (*models.Token, error)
	GetUserForToken(ctx context.Context, userID int) (*models.User, error) // Método para obtener usuario por ID
}

var AuthRepo AuthRepository

func SetAuthRepository(repo AuthRepository) {
	AuthRepo = repo
}

func InsertUser(ctx context.Context, user *models.User) error {
	return AuthRepo.InsertUser(ctx, user)
}

func GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return AuthRepo.GetUserByEmail(ctx, email)
}

func GetUserByUUID(ctx context.Context, uuid string) (*models.User, error) {
	return AuthRepo.GetUserByUUID(ctx, uuid)
}

// Nuevas funciones para tokens
func InsertToken(ctx context.Context, token *models.Token) error {
	return AuthRepo.InsertToken(ctx, token)
}

func GetTokenByTokenHash(ctx context.Context, tokenHash string) (*models.Token, error) {
	return AuthRepo.GetTokenByTokenHash(ctx, tokenHash)
}

func DeleteTokensByUserID(ctx context.Context, userID int) error {
	return AuthRepo.DeleteTokensByUserID(ctx, userID)
}

func GetTokenByToken(ctx context.Context, tokenString string) (*models.Token, error) {
	return AuthRepo.GetTokenByToken(ctx, tokenString)
}

func GetUserForToken(ctx context.Context, userID int) (*models.User, error) {
	return AuthRepo.GetUserForToken(ctx, userID)
}
