package services

import (
	"cloudtrail-enrichment-api-golang/internal/pkg/logger"
	"cloudtrail-enrichment-api-golang/internal/pkg/token"
	"cloudtrail-enrichment-api-golang/internal/repository"
	"cloudtrail-enrichment-api-golang/models"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// AuthService defines the interface for authentication service operations.
type AuthService interface {
	RegisterUser(ctx context.Context, payload *models.RegisterPayload) (*models.User, error)
	AuthenticateUser(ctx context.Context, email, password string) (*models.User, *token.JWTToken, error)
	ValidateTokenForMiddleware(ctx context.Context, tokenString string) (*token.User, error)
}

// DefaultAuthService is the default implementation of AuthService.
type DefaultAuthService struct {
	repo       repository.AuthRepository
	jwtService *token.JWTService
}

// NewAuthService creates a new instance of DefaultAuthService.
func NewAuthService(repo repository.AuthRepository, jwtService *token.JWTService) *DefaultAuthService {
	return &DefaultAuthService{
		repo:       repo,
		jwtService: jwtService,
	}
}

// RegisterUser registers a new user in the system.
func (s *DefaultAuthService) RegisterUser(ctx context.Context, payload *models.RegisterPayload) (*models.User, error) {
	// Check if user already exists by email
	_, err := s.repo.GetUserByEmail(ctx, payload.Email)
	if err == nil {
		logger.InfoLog.Printf("Attempt to register with duplicate email: %s", payload.Email)
		return nil, errors.New("email is already registered")
	}
	if err != nil && err != sql.ErrNoRows {
		logger.ErrorLog.Printf("Error checking existing email: %v", err)
		return nil, fmt.Errorf("error checking email: %w", err)
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.ErrorLog.Printf("Error hashing password: %v", err)
		return nil, fmt.Errorf("error processing password: %w", err)
	}

	newUser := &models.User{
		UUID:         uuid.NewString(),
		Email:        payload.Email,
		PasswordHash: string(hashedPassword),
		Role:         payload.Role, // Validate or assign a default role (Pending Check)
	}

	err = s.repo.InsertUser(ctx, newUser)
	if err != nil {
		logger.ErrorLog.Printf("Service error inserting new user: %v", err)
		return nil, fmt.Errorf("error registering user: %w", err)
	}

	logger.InfoLog.Printf("User %s successfully registered with UUID: %s", newUser.Email, newUser.UUID)
	return newUser, nil
}

// AuthenticateUser authenticates a user and generates a JWT token.
func (s *DefaultAuthService) AuthenticateUser(ctx context.Context, email, password string) (*models.User, *token.JWTToken, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.InfoLog.Printf("Failed authentication attempt: user not found with email %s", email)
			return nil, nil, errors.New("invalid credentials")
		}
		logger.ErrorLog.Printf("Error getting user by email %s to authenticate: %v", email, err)
		return nil, nil, fmt.Errorf("authentication error: %w", err)
	}

	// Compare hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		logger.InfoLog.Printf("Failed authentication attempt for %s: incorrect password", email)
		return nil, nil, errors.New("invalid credentials")
	}

	// Generate the JWT token. The logic for persistence and revocation of previous tokens is now inside GenerateJWTToken in the token package.
	signedToken, expiry, err := s.jwtService.GenerateJWTToken(ctx, user.ID, user.Email, user.Role)
	if err != nil {
		logger.ErrorLog.Printf("Error generating JWT for user %s: %v", user.Email, err)
		return nil, nil, fmt.Errorf("error generating token: %w", err)
	}

	// Create a simple JWTToken with data for the response, without dependencies.
	jwtData := &token.JWTToken{
		UserID: user.ID,
		Email:  user.Email,
		Token:  signedToken,
		Expiry: expiry,
		Role:   user.Role,
	}

	logger.InfoLog.Printf("User %s authenticated and token generated.", user.Email)
	return user, jwtData, nil
}

// ValidateTokenForMiddleware validates a JWT token and verifies the existence of the user and token validity in the DB. This method is invoked by the middleware.
func (s *DefaultAuthService) ValidateTokenForMiddleware(ctx context.Context, tokenString string) (*token.User, error) {
	// Validate the token using logic from the token package
	jwtClaims, err := s.jwtService.ValidJWTToken(ctx, tokenString) // Pass the context
	if err != nil {
		logger.ErrorLog.Printf("Invalid, expired, or not found JWT token in DB: %v", err)
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// Verify that the token's user really exists in our DB. This is useful if users can be deleted or disabled after a token is generated.
	_, err = s.repo.GetUserByEmail(ctx, jwtClaims.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.ErrorLog.Printf("User associated with token (%s) not found in DB.", jwtClaims.Email)
			return nil, errors.New("user associated with token not found")
		}
		logger.ErrorLog.Printf("Error verifying token user in DB: %v", err)
		return nil, fmt.Errorf("user verification error: %w", err)
	}

	logger.InfoLog.Printf("JWT token successfully validated for user: %s", jwtClaims.Email)
	return jwtClaims, nil
}
