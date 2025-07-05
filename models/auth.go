package models

import "time"

// User representa la información del usuario en el sistema.
type User struct {
	// ID es el identificador único autoincremental de la base de datos.
	ID int `json:"id"`
	// UUID es el identificador único universal generado para el usuario.
	UUID         string    `json:"uuid"`
	Email        string    `json:"email"`
	Password     string    `json:"password,omitempty"` // Para recibir la contraseña en el request, no la hash
	PasswordHash string    `json:"-"`                  // El hash de la contraseña, no se expone en JSON
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Token representa un token JWT, para ser almacenado si fuera necesario
// o para ser usado como estructura al generar/validar.
type Token struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
	TokenHash string    `json:"token_hash"` // Hash del token para almacenamiento seguro o revocación
	Expiry    time.Time `json:"expiry"`
	Role      string    `json:"role"` // Añadir el campo Role para que coincida con JWTToken y la DB
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// LoginPayload es la estructura para las solicitudes de inicio de sesión.
type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterPayload es la estructura para las solicitudes de registro de usuario.
type RegisterPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"` // Puedes definir roles predeterminados si no se envían
}
