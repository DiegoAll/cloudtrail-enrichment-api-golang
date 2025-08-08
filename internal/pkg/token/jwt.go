package token

import (
	"cloudtrail-enrichment-api-golang/internal/config"
	"cloudtrail-enrichment-api-golang/internal/pkg/logger"
	"cloudtrail-enrichment-api-golang/models"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// User represents the user information contained in the JWT token.
type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// JWTToken represents the structure of the JWT token to be serialized in the payload.
// This struct must not contain dependencies like Config or TokenRepository. (package Config Issue)
type JWTToken struct {
	UserID    int       `json:"user_id,omitempty"`
	Email     string    `json:"email,omitempty"`
	Token     string    `json:"token"`
	TokenHash string    `json:"token_hash"`
	Expiry    time.Time `json:"expiry"`
	Role      string    `json:"role,omitempty"`
	jwt.RegisteredClaims
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TokenDBRepository is an interface for token persistence operations.
// Defined here so that the token package doesn’t directly depend on the repository package.
type TokenDBRepository interface {
	InsertToken(ctx context.Context, token *models.Token) error
	GetTokenByTokenHash(ctx context.Context, tokenHash string) (*models.Token, error)
	DeleteTokensByUserID(ctx context.Context, userID int) error
	GetTokenByToken(ctx context.Context, tokenString string) (*models.Token, error)
	GetUserForToken(ctx context.Context, userID int) (*models.User, error)
}

// JWTService is the struct that holds the dependencies and methods related to JWT tokens.
type JWTService struct {
	Config          *config.Config
	TokenRepository TokenDBRepository
}

// NewJWTService creates a new instance of the JWT token service.
func NewJWTService(cfg *config.Config, tokenRepo TokenDBRepository) *JWTService {
	return &JWTService{
		Config:          cfg,
		TokenRepository: tokenRepo,
	}
}

// GetByToken takes a plain text token and looks it up in the database. Returns a pointer to the Token model.
func (j *JWTService) GetByToken(ctx context.Context, plainText string) (*models.Token, error) {
	dbTimeout := j.Config.DatabaseConfig.DBTimeout
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	token, err := j.TokenRepository.GetTokenByToken(ctx, plainText)
	if err != nil {
		logger.ErrorLog.Printf("Error getting token by value in DB: %v", err)
		return nil, err
	}
	return token, nil
}

// GetUserForToken retrieves a user from the database given a persisted token.
func (j *JWTService) GetUserForToken(ctx context.Context, token models.Token) (*models.User, error) {
	dbTimeout := j.Config.DatabaseConfig.DBTimeout
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	user, err := j.TokenRepository.GetUserForToken(ctx, token.UserID)
	if err != nil {
		logger.ErrorLog.Printf("Error retrieving user for token in DB: %v", err)
		return nil, err
	}
	return user, nil
}

// InsertJWT inserts a new JWT token into the database, revoking any previous ones for the same user.
func (j *JWTService) InsertJWT(ctx context.Context, tokenData models.Token, user models.User) error {
	dbTimeout := j.Config.DatabaseConfig.DBTimeout
	fmt.Println("DB TIMEOUT", dbTimeout)
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	// Remove any existing token for this user
	err := j.TokenRepository.DeleteTokensByUserID(ctx, tokenData.UserID)
	if err != nil {
		logger.ErrorLog.Printf("Error deleting existing tokens for user %d: %v", tokenData.UserID, err)
		return err
	}

	tokenData.Email = user.Email
	tokenData.Role = user.Role

	// Insert the new token
	err = j.TokenRepository.InsertToken(ctx, &tokenData)
	if err != nil {
		logger.ErrorLog.Printf("Error inserting new token for user %d: %v", tokenData.UserID, err)
		return err
	}

	return nil
}

// DeleteByJWTToken deletes a token from the database given its plain text value.
func (j *JWTService) DeleteByJWTToken(ctx context.Context, plainText string) error {
	dbTimeout := j.Config.DatabaseConfig.DBTimeout
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	tokenToDelete, err := j.TokenRepository.GetTokenByToken(ctx, plainText)
	if err != nil {
		logger.ErrorLog.Printf("Token to delete not found: %v", err)
		return errors.New("token to delete not found")
	}

	err = j.TokenRepository.DeleteTokensByUserID(ctx, tokenToDelete.UserID)
	if err != nil {
		logger.ErrorLog.Printf("Error deleting token by plain text for user %d: %v", tokenToDelete.UserID, err)
		return err
	}

	return nil
}

// GenerateJWTToken generates a new JWT token for a given user and persists it.
func (j *JWTService) GenerateJWTToken(ctx context.Context, userID int, email string, role string) (string, time.Time, error) {
	expiry := time.Now().Add(j.Config.AuthConfig.TokenDuration)

	// Here we use JWTToken only for the claims to be serialized.
	claims := &JWTToken{
		UserID: userID,
		Email:  email,
		Role:   role,
		Expiry: expiry,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   email,
			ExpiresAt: jwt.NewNumericDate(expiry),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "The Organization",
			Audience:  []string{"The Users"},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretkey := j.Config.AuthConfig.JWTSecret
	if secretkey == "" {
		return "", time.Time{}, errors.New("JWT secret key not configured")
	}

	signedToken, err := token.SignedString([]byte(secretkey))
	if err != nil {
		logger.ErrorLog.Printf("Error signing JWT token: %v", err)
		return "", time.Time{}, err
	}

	tokenHash := sha256.Sum256([]byte(signedToken))
	tokenHashString := fmt.Sprintf("%x", tokenHash)

	dbToken := models.Token{
		UserID:    userID,
		Email:     email,
		Token:     signedToken,
		TokenHash: tokenHashString,
		Expiry:    expiry,
		Role:      role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = j.InsertJWT(ctx, dbToken, models.User{ID: userID, Email: email, Role: role})
	if err != nil {
		logger.ErrorLog.Printf("Error persisting token in DB for user %s: %v", email, err)
		return "", time.Time{}, fmt.Errorf("error persisting token: %w", err)
	}

	return signedToken, expiry, nil
}

// ExtractJWTToken extracts the JWT token from the Authorization header.
func (j *JWTService) ExtractJWTToken(r *http.Request) (string, error) {
	authorizationHeader := r.Header.Get("Authorization")
	if authorizationHeader == "" {
		return "", errors.New("authorization header not provided")
	}

	headerParts := strings.Split(authorizationHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("invalid token format")
	}

	tokenString := headerParts[1]
	return tokenString, nil
}

// ValidJWTToken validates the JWT token and returns user information, also verifying its existence in the database.
func (j *JWTService) ValidJWTToken(ctx context.Context, tokenString string) (*User, error) {
	secretkey := j.Config.AuthConfig.JWTSecret
	if secretkey == "" {
		return nil, errors.New("JWT secret key not configured")
	}

	token, err := jwt.ParseWithClaims(tokenString, &JWTToken{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretkey), nil
	})

	if err != nil {
		logger.ErrorLog.Printf("Error parsing or validating JWT token: %v", err)
		return nil, fmt.Errorf("invalid or expired token: %w", err)
	}

	claims, ok := token.Claims.(*JWTToken)
	if !ok || !token.Valid {
		logger.ErrorLog.Println("Invalid JWT token claims or token not valid.")
		return nil, errors.New("invalid token")
	}

	if time.Now().After(claims.Expiry) {
		logger.InfoLog.Printf("Token expired for user: %s", claims.Email)
		return nil, errors.New("token expired")
	}

	tokenHash := sha256.Sum256([]byte(tokenString))
	tokenHashString := fmt.Sprintf("%x", tokenHash)

	persistedToken, err := j.TokenRepository.GetTokenByTokenHash(ctx, tokenHashString)
	if err != nil {
		logger.ErrorLog.Printf("Error looking up token in database by hash: %v", err)
		return nil, errors.New("token not found or revoked")
	}

	if persistedToken.UserID != claims.UserID || persistedToken.Email != claims.Email || persistedToken.Role != claims.Role || time.Now().After(persistedToken.Expiry) {
		logger.ErrorLog.Printf("Persisted token inconsistency or expired token in DB for user: %s", claims.Email)
		return nil, errors.New("invalid or revoked token")
	}

	return &User{
		ID:    claims.UserID,
		Email: claims.Email,
		Role:  claims.Role,
	}, nil
}

// AuthenticateJWTToken extracts and validates the JWT token from the request.
func (j *JWTService) AuthenticateJWTToken(r *http.Request) (*User, error) {
	tokenString, err := j.ExtractJWTToken(r)
	if err != nil {
		return nil, err
	}

	user, err := j.ValidJWTToken(r.Context(), tokenString)
	if err != nil {
		return nil, err
	}

	return user, nil
}
