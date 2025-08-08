package controllers

import (
	"cloudtrail-enrichment-api-golang/internal/pkg/logger"
	"cloudtrail-enrichment-api-golang/internal/pkg/utils"
	"cloudtrail-enrichment-api-golang/models"
	"cloudtrail-enrichment-api-golang/services"
	"errors"
	"net/http"
)

// AuthController Handles HTTP requests related to authentication.
type AuthController struct {
	service services.AuthService
}

// NewAuthController Creates a new instance of AuthController.
func NewAuthController(s services.AuthService) *AuthController {
	return &AuthController{
		service: s,
	}
}

// RegisterUser Handles the registration of new users.
func (ac *AuthController) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var payload models.RegisterPayload
	err := utils.ReadJSON(w, r, &payload)
	if err != nil {
		logger.ErrorLog.Printf("Error decoding registration payload: %v", err)
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Basic input validations.
	if payload.Email == "" || payload.Password == "" {
		logger.ErrorLog.Println("Email or password empty during registration.")
		utils.ErrorJSON(w, errors.New("email and password are required"), http.StatusBadRequest)
		return
	}
	// Assign a default role if not specified or validate the role.
	if payload.Role == "" {
		payload.Role = "user" // Default role.
	}

	// Calls the service to register the user.
	user, err := ac.service.RegisterUser(r.Context(), &payload)
	if err != nil {
		logger.ErrorLog.Printf("Error registering user in service: %v", err)
		utils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	response := utils.JSONResponse{
		Error:   false,
		Message: "User registered successfully",
		Data: map[string]string{
			"uuid":  user.UUID,
			"email": user.Email,
			"role":  user.Role,
		},
	}
	utils.WriteJSON(w, http.StatusCreated, response)
	logger.InfoLog.Printf("User %s registered successfully.", user.Email)
}

// AuthenticateUser handles user login and token generation.
func (ac *AuthController) AuthenticateUser(w http.ResponseWriter, r *http.Request) {
	var payload models.LoginPayload
	err := utils.ReadJSON(w, r, &payload)
	if err != nil {
		logger.ErrorLog.Printf("Error decoding authentication payload: %v", err)
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	if payload.Email == "" || payload.Password == "" {
		logger.ErrorLog.Println("Email or password empty during authentication.")
		utils.ErrorJSON(w, errors.New("email and password are required"), http.StatusBadRequest)
		return
	}

	// Calls the service to authenticate the user and get the token.
	user, jwtToken, err := ac.service.AuthenticateUser(r.Context(), payload.Email, payload.Password)
	if err != nil {
		logger.ErrorLog.Printf("Authentication error in service for %s: %v", payload.Email, err)
		utils.ErrorJSON(w, err, http.StatusUnauthorized) // 401 Unauthorized
		return
	}

	response := utils.JSONResponse{
		Error:   false,
		Message: "Authentication successful",
		Data: map[string]interface{}{
			"user_uuid": user.UUID,
			"email":     user.Email,
			"role":      user.Role,
			"token":     jwtToken.Token, // Includes the JWT token in the response.
			"expiry":    jwtToken.Expiry,
		},
	}
	utils.WriteJSON(w, http.StatusOK, response)
	logger.InfoLog.Printf("User %s successfully authenticated and token generated.", user.Email)
}

// PublicRouteHandler is an example of a handler for a public route.
func (ac *AuthController) PublicRouteHandler(w http.ResponseWriter, r *http.Request) {
	response := utils.JSONResponse{
		Error:   false,
		Message: "This is a public route, accessible without authentication.",
		Data:    nil,
	}
	utils.WriteJSON(w, http.StatusOK, response)
	logger.InfoLog.Println("Access to public route.")
}
