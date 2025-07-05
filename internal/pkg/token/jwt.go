package token

import (
	"cloudtrail-enrichment-api-golang/internal/config"
	"cloudtrail-enrichment-api-golang/internal/pkg/logger" // Importar logger
	"cloudtrail-enrichment-api-golang/models"
	"context"
	"crypto/sha256" // Importar para SHA256
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// User representa la información del usuario contenida en el token JWT.
type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// JWTToken representa la estructura del token JWT que se serializará en el payload.
// Este struct NO DEBE CONTENER DEPENDENCIAS como Config o TokenRepository.
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

// TokenDBRepository es una interfaz para las operaciones de persistencia del token.
// Se define aquí para que el paquete token no tenga una dependencia directa al paquete repository.
type TokenDBRepository interface {
	InsertToken(ctx context.Context, token *models.Token) error
	GetTokenByTokenHash(ctx context.Context, tokenHash string) (*models.Token, error)
	DeleteTokensByUserID(ctx context.Context, userID int) error
	GetTokenByToken(ctx context.Context, tokenString string) (*models.Token, error) // Método nuevo para obtener por token directamente
	GetUserForToken(ctx context.Context, userID int) (*models.User, error)          // Método para obtener usuario por ID
}

// JWTService es el struct que contendrá las dependencias
// y los métodos de negocio relacionados con los tokens JWT.
type JWTService struct {
	Config          *config.Config
	TokenRepository TokenDBRepository
}

// NewJWTService crea una nueva instancia del servicio de token JWT.
func NewJWTService(cfg *config.Config, tokenRepo TokenDBRepository) *JWTService {
	return &JWTService{
		Config:          cfg,
		TokenRepository: tokenRepo,
	}
}

// GetByToken toma un token en texto plano y busca el token completo en la base de datos.
// Devuelve un puntero al modelo Token.
func (j *JWTService) GetByToken(ctx context.Context, plainText string) (*models.Token, error) {
	dbTimeout := j.Config.DatabaseConfig.DBTimeout
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	token, err := j.TokenRepository.GetTokenByToken(ctx, plainText)
	if err != nil {
		logger.ErrorLog.Printf("Error al obtener token por valor en DB: %v", err)
		return nil, err
	}
	return token, nil
}

// GetUserForToken obtiene un usuario de la base de datos dado un token persistido.
func (j *JWTService) GetUserForToken(ctx context.Context, token models.Token) (*models.User, error) {
	dbTimeout := j.Config.DatabaseConfig.DBTimeout
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	user, err := j.TokenRepository.GetUserForToken(ctx, token.UserID)
	if err != nil {
		logger.ErrorLog.Printf("Error al obtener usuario para el token en DB: %v", err)
		return nil, err
	}
	return user, nil
}

// InsertJWT inserta un nuevo token JWT en la base de datos, revocando los anteriores para el mismo usuario.
func (j *JWTService) InsertJWT(ctx context.Context, tokenData models.Token, user models.User) error {
	dbTimeout := j.Config.DatabaseConfig.DBTimeout
	fmt.Println("DB TIMEOUT", dbTimeout)
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	// Eliminar cualquier token existente para este usuario
	err := j.TokenRepository.DeleteTokensByUserID(ctx, tokenData.UserID)
	if err != nil {
		logger.ErrorLog.Printf("Error al eliminar tokens existentes para el usuario %d: %v", tokenData.UserID, err)
		return err
	}

	// Asignar el email del usuario al token
	tokenData.Email = user.Email
	tokenData.Role = user.Role

	// Insertar el nuevo token
	err = j.TokenRepository.InsertToken(ctx, &tokenData)
	if err != nil {
		logger.ErrorLog.Printf("Error al insertar el nuevo token para el usuario %d: %v", tokenData.UserID, err)
		return err
	}

	return nil
}

// DeleteByJWTToken elimina un token de la base de datos dado su valor en texto plano.
func (j *JWTService) DeleteByJWTToken(ctx context.Context, plainText string) error {
	dbTimeout := j.Config.DatabaseConfig.DBTimeout
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	tokenToDelete, err := j.TokenRepository.GetTokenByToken(ctx, plainText)
	if err != nil {
		logger.ErrorLog.Printf("Token a eliminar no encontrado: %v", err)
		return errors.New("token a eliminar no encontrado")
	}

	err = j.TokenRepository.DeleteTokensByUserID(ctx, tokenToDelete.UserID)
	if err != nil {
		logger.ErrorLog.Printf("Error al eliminar token por texto plano para el usuario %d: %v", tokenToDelete.UserID, err)
		return err
	}

	return nil
}

// GenerateJWTToken genera un nuevo token JWT para un usuario dado y lo persiste.
func (j *JWTService) GenerateJWTToken(ctx context.Context, userID int, email string, role string) (string, time.Time, error) {
	expiry := time.Now().Add(j.Config.AuthConfig.TokenDuration)

	// Aquí usamos JWTToken solo para los claims que se serializarán.
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
			Issuer:    "g3notype",
			Audience:  []string{"mis-usuarios"},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretkey := j.Config.AuthConfig.JWTSecret
	if secretkey == "" {
		return "", time.Time{}, errors.New("clave secreta JWT no configurada")
	}

	signedToken, err := token.SignedString([]byte(secretkey))
	if err != nil {
		logger.ErrorLog.Printf("Error al firmar el token JWT: %v", err)
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
		logger.ErrorLog.Printf("Error al persistir el token en la DB para el usuario %s: %v", email, err)
		return "", time.Time{}, fmt.Errorf("error al persistir el token: %w", err)
	}

	return signedToken, expiry, nil
}

// ExtractJWTToken extrae el token JWT del encabezado Authorization.
func (j *JWTService) ExtractJWTToken(r *http.Request) (string, error) {
	authorizationHeader := r.Header.Get("Authorization")
	if authorizationHeader == "" {
		return "", errors.New("encabezado de autorización no proporcionado")
	}

	headerParts := strings.Split(authorizationHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("formato de token inválido")
	}

	tokenString := headerParts[1]
	return tokenString, nil
}

// ValidJWTToken valida el token JWT y devuelve la información del usuario,
// además de verificar su existencia en la base de datos.
func (j *JWTService) ValidJWTToken(ctx context.Context, tokenString string) (*User, error) {
	secretkey := j.Config.AuthConfig.JWTSecret
	if secretkey == "" {
		return nil, errors.New("clave secreta JWT no configurada")
	}

	token, err := jwt.ParseWithClaims(tokenString, &JWTToken{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de firma inesperado: %v", token.Header["alg"])
		}
		return []byte(secretkey), nil
	})

	if err != nil {
		logger.ErrorLog.Printf("Error al parsear o validar token JWT: %v", err)
		return nil, fmt.Errorf("token inválido o expirado: %w", err)
	}

	claims, ok := token.Claims.(*JWTToken)
	if !ok || !token.Valid {
		logger.ErrorLog.Println("Claims de token JWT inválidos o token no válido.")
		return nil, errors.New("token inválido")
	}

	if time.Now().After(claims.Expiry) {
		logger.InfoLog.Printf("Token expirado para el usuario: %s", claims.Email)
		return nil, errors.New("token expirado")
	}

	tokenHash := sha256.Sum256([]byte(tokenString))
	tokenHashString := fmt.Sprintf("%x", tokenHash)

	persistedToken, err := j.TokenRepository.GetTokenByTokenHash(ctx, tokenHashString)
	if err != nil {
		logger.ErrorLog.Printf("Error al buscar el token en la base de datos por hash: %v", err)
		return nil, errors.New("token no encontrado o revocado")
	}

	if persistedToken.UserID != claims.UserID || persistedToken.Email != claims.Email || persistedToken.Role != claims.Role || time.Now().After(persistedToken.Expiry) {
		logger.ErrorLog.Printf("Inconsistencia en el token persistido o token expirado en DB para usuario: %s", claims.Email)
		return nil, errors.New("token inválido o revocado")
	}

	return &User{
		ID:    claims.UserID,
		Email: claims.Email,
		Role:  claims.Role,
	}, nil
}

// AuthenticateJWTToken extrae y valida el token JWT de la solicitud.
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
