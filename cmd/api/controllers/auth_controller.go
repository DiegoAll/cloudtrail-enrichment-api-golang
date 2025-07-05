package controllers

import (
	"cloudtrail-enrichment-api-golang/internal/pkg/logger"
	"cloudtrail-enrichment-api-golang/internal/pkg/utils"
	"cloudtrail-enrichment-api-golang/models"
	"cloudtrail-enrichment-api-golang/services"
	"errors"
	"net/http"
	// Asegúrate de importar uuid si lo usas para generar UUIDs de usuario
)

// AuthController maneja las solicitudes HTTP relacionadas con la autenticación.
type AuthController struct {
	service services.AuthService
}

// NewAuthController crea una nueva instancia de AuthController.
func NewAuthController(s services.AuthService) *AuthController {
	return &AuthController{
		service: s,
	}
}

// RegisterUser maneja el registro de nuevos usuarios.
func (ac *AuthController) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var payload models.RegisterPayload
	err := utils.ReadJSON(w, r, &payload)
	if err != nil {
		logger.ErrorLog.Printf("Error al decodificar payload de registro: %v", err)
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Validaciones básicas de entrada
	if payload.Email == "" || payload.Password == "" {
		logger.ErrorLog.Println("Email o contraseña vacíos en el registro.")
		utils.ErrorJSON(w, errors.New("email y contraseña son requeridos"), http.StatusBadRequest)
		return
	}
	// Asignar un rol por defecto si no se especifica o validar el rol
	if payload.Role == "" {
		payload.Role = "user" // Rol por defecto
	}

	// Llama al servicio para registrar el usuario
	user, err := ac.service.RegisterUser(r.Context(), &payload)
	if err != nil {
		logger.ErrorLog.Printf("Error al registrar usuario en el servicio: %v", err)
		utils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	response := utils.JSONResponse{
		Error:   false,
		Message: "Usuario registrado exitosamente",
		Data: map[string]string{
			"uuid":  user.UUID,
			"email": user.Email,
			"role":  user.Role,
		},
	}
	utils.WriteJSON(w, http.StatusCreated, response)
	logger.InfoLog.Printf("Usuario %s registrado con éxito.", user.Email)
}

// AuthenticateUser maneja el inicio de sesión de usuarios y la generación de tokens.
func (ac *AuthController) AuthenticateUser(w http.ResponseWriter, r *http.Request) {
	var payload models.LoginPayload
	err := utils.ReadJSON(w, r, &payload)
	if err != nil {
		logger.ErrorLog.Printf("Error al decodificar payload de autenticación: %v", err)
		utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	if payload.Email == "" || payload.Password == "" {
		logger.ErrorLog.Println("Email o contraseña vacíos en la autenticación.")
		utils.ErrorJSON(w, errors.New("email y contraseña son requeridos"), http.StatusBadRequest)
		return
	}

	// Llama al servicio para autenticar al usuario y obtener el token
	user, jwtToken, err := ac.service.AuthenticateUser(r.Context(), payload.Email, payload.Password)
	if err != nil {
		logger.ErrorLog.Printf("Error de autenticación en el servicio para %s: %v", payload.Email, err)
		utils.ErrorJSON(w, err, http.StatusUnauthorized) // 401 Unauthorized
		return
	}

	response := utils.JSONResponse{
		Error:   false,
		Message: "Autenticación exitosa",
		Data: map[string]interface{}{
			"user_uuid": user.UUID,
			"email":     user.Email,
			"role":      user.Role,
			"token":     jwtToken.Token, // Incluye el token JWT en la respuesta
			"expiry":    jwtToken.Expiry,
		},
	}
	utils.WriteJSON(w, http.StatusOK, response)
	logger.InfoLog.Printf("Usuario %s autenticado exitosamente y token generado.", user.Email)
}

// PublicRouteHandler es un ejemplo de un handler para una ruta pública.
func (ac *AuthController) PublicRouteHandler(w http.ResponseWriter, r *http.Request) {
	response := utils.JSONResponse{
		Error:   false,
		Message: "Esta es una ruta pública, accesible sin autenticación.",
		Data:    nil,
	}
	utils.WriteJSON(w, http.StatusOK, response)
	logger.InfoLog.Println("Acceso a ruta pública.")
}
