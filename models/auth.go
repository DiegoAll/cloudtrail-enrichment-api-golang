package models

import "time"

// User represents the user information.
type User struct {
	// ID is the unique auto-incrementing identifier in the database.
	ID int `json:"id"`
	// UUID is the universally unique identifier generated for the user.
	UUID         string    `json:"uuid"`
	Email        string    `json:"email"`
	Password     string    `json:"password,omitempty"` // To receive the password in the request, not the hash
	PasswordHash string    `json:"-"`                  // The hashed password, not exposed in JSON
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Token represents a JWT token, to be stored if necessary or used as a structure when generating/validating.
type Token struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
	TokenHash string    `json:"token_hash"` // Hash of the token for secure storage or revocation
	Expiry    time.Time `json:"expiry"`
	Role      string    `json:"role"` // Role field to match JWTToken and the DB
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// LoginPayload is the structure for login requests.
type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterPayload is the structure for user registration requests.
type RegisterPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}
