package services

import (
	"cloudtrail-enrichment-api-golang/internal/pkg/logger"
	"cloudtrail-enrichment-api-golang/internal/pkg/token" // Importa el paquete token
	"cloudtrail-enrichment-api-golang/internal/repository"
	"cloudtrail-enrichment-api-golang/models"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt" // Para hashing de contraseñas
)

// AuthService define la interfaz para las operaciones de servicio de autenticación.
type AuthService interface {
	RegisterUser(ctx context.Context, payload *models.RegisterPayload) (*models.User, error)
	// CAMBIO: El segundo retorno de AuthenticateUser ahora es *token.JWTToken (solo claims)
	AuthenticateUser(ctx context.Context, email, password string) (*models.User, *token.JWTToken, error)
	ValidateTokenForMiddleware(ctx context.Context, tokenString string) (*token.User, error)
}

// DefaultAuthService es la implementación predeterminada de AuthService.
type DefaultAuthService struct {
	repo       repository.AuthRepository
	jwtService *token.JWTService // CAMBIO: Ahora inyectamos *token.JWTService
}

// NewAuthService crea una nueva instancia de DefaultAuthService.
// CAMBIO: Recibe *token.JWTService en lugar de *token.JWTToken
func NewAuthService(repo repository.AuthRepository, jwtService *token.JWTService) *DefaultAuthService {
	return &DefaultAuthService{
		repo:       repo,
		jwtService: jwtService, // CAMBIO: Asigna jwtService
	}
}

// RegisterUser registra un nuevo usuario en el sistema.
func (s *DefaultAuthService) RegisterUser(ctx context.Context, payload *models.RegisterPayload) (*models.User, error) {
	// Verificar si el usuario ya existe por email
	_, err := s.repo.GetUserByEmail(ctx, payload.Email)
	if err == nil {
		logger.InfoLog.Printf("Intento de registro con email duplicado: %s", payload.Email)
		return nil, errors.New("el email ya está registrado")
	}
	if err != nil && err != sql.ErrNoRows {
		logger.ErrorLog.Printf("Error al verificar email existente: %v", err)
		return nil, fmt.Errorf("error al verificar email: %w", err)
	}

	// Hashear la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.ErrorLog.Printf("Error al hashear la contraseña: %v", err)
		return nil, fmt.Errorf("error al procesar contraseña: %w", err)
	}

	newUser := &models.User{
		UUID:         uuid.NewString(),
		Email:        payload.Email,
		PasswordHash: string(hashedPassword),
		Role:         payload.Role, // Asegúrate de validar o asignar un rol por defecto
	}

	err = s.repo.InsertUser(ctx, newUser)
	if err != nil {
		logger.ErrorLog.Printf("Error en el servicio al insertar nuevo usuario: %v", err)
		return nil, fmt.Errorf("error al registrar usuario: %w", err)
	}

	logger.InfoLog.Printf("Usuario %s registrado exitosamente con UUID: %s", newUser.Email, newUser.UUID)
	return newUser, nil
}

// AuthenticateUser autentica a un usuario y genera un token JWT.
// CAMBIO: El segundo retorno ahora es *token.JWTToken (solo claims)
func (s *DefaultAuthService) AuthenticateUser(ctx context.Context, email, password string) (*models.User, *token.JWTToken, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.InfoLog.Printf("Intento de autenticación fallido: usuario no encontrado con email %s", email)
			return nil, nil, errors.New("credenciales inválidas")
		}
		logger.ErrorLog.Printf("Error al obtener usuario por email %s para autenticar: %v", email, err)
		return nil, nil, fmt.Errorf("error al autenticar: %w", err)
	}

	// Comparar la contraseña hasheada
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		logger.InfoLog.Printf("Intento de autenticación fallido para %s: contraseña incorrecta", email)
		return nil, nil, errors.New("credenciales inválidas")
	}

	// Generar el token JWT. La lógica de persistencia y revocación de tokens anteriores
	// ahora está dentro de GenerateJWTToken en el paquete token.
	// CAMBIO: Llamamos a GenerateJWTToken del jwtService
	signedToken, expiry, err := s.jwtService.GenerateJWTToken(ctx, user.ID, user.Email, user.Role)
	if err != nil {
		logger.ErrorLog.Printf("Error al generar JWT para usuario %s: %v", user.Email, err)
		return nil, nil, fmt.Errorf("error al generar token: %w", err)
	}

	// Creamos un JWTToken simple con los datos para la respuesta, sin las dependencias.
	jwtData := &token.JWTToken{
		UserID: user.ID,
		Email:  user.Email,
		Token:  signedToken,
		Expiry: expiry,
		Role:   user.Role,
	}

	logger.InfoLog.Printf("Usuario %s autenticado y token generado.", user.Email)
	return user, jwtData, nil
}

// ValidateTokenForMiddleware valida un token JWT y verifica la existencia del usuario y la validez del token en la DB.
// Este método es invocado por el middleware.
func (s *DefaultAuthService) ValidateTokenForMiddleware(ctx context.Context, tokenString string) (*token.User, error) {
	// Validar el token usando la lógica del paquete token
	// CAMBIO: Llamamos a ValidJWTToken del jwtService
	jwtClaims, err := s.jwtService.ValidJWTToken(ctx, tokenString) // Pasa el contexto
	if err != nil {
		logger.ErrorLog.Printf("Token JWT inválido, expirado o no encontrado en DB: %v", err)
		return nil, fmt.Errorf("token inválido: %w", err)
	}

	// Opcional pero recomendado: Verificar que el usuario del token realmente existe en nuestra DB.
	// Esto es útil si los usuarios pueden ser eliminados o deshabilitados después de generar un token.
	_, err = s.repo.GetUserByEmail(ctx, jwtClaims.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.ErrorLog.Printf("Usuario asociado al token (%s) no encontrado en la DB.", jwtClaims.Email)
			return nil, errors.New("usuario asociado al token no encontrado")
		}
		logger.ErrorLog.Printf("Error al verificar usuario de token en DB: %v", err)
		return nil, fmt.Errorf("error de verificación de usuario: %w", err)
	}

	logger.InfoLog.Printf("Token JWT validado exitosamente para usuario: %s", jwtClaims.Email)
	return jwtClaims, nil
}
