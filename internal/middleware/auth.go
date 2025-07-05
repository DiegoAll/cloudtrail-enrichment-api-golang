package middleware

import (
	"cloudtrail-enrichment-api-golang/internal/pkg/logger"
	"cloudtrail-enrichment-api-golang/internal/pkg/token" // Importa el paquete token
	"cloudtrail-enrichment-api-golang/internal/pkg/utils"
	"cloudtrail-enrichment-api-golang/services" // Importa el paquete services
	"context"                                   // Importa context
	"net/http"
)

// Middleware contiene las dependencias para los middlewares.
type Middleware struct {
	JWTService  *token.JWTService    // CAMBIO: Ahora es *token.JWTService
	AuthService services.AuthService // Añade la dependencia del AuthService
}

// NewMiddleware crea una nueva instancia de Middleware.
// CAMBIO: Recibe *token.JWTService en lugar de *token.JWTToken
func NewMiddleware(jwtService *token.JWTService, authService services.AuthService) *Middleware {
	return &Middleware{
		JWTService:  jwtService,  // CAMBIO: Asigna jwtService
		AuthService: authService, // Asigna el servicio de autenticación
	}
}

// AuthTokenMiddleware es un middleware que valida el token JWT de la solicitud.
func (mw *Middleware) AuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.InfoLog.Println("Ejecutando AuthTokenMiddleware...")

		// CAMBIO: Llama al método de JWTService
		tokenString, err := mw.JWTService.ExtractJWTToken(r)
		if err != nil {
			payload := utils.JSONResponse{
				Error:   true,
				Message: "Token no proporcionado o formato inválido: " + err.Error(),
			}
			logger.ErrorLog.Printf("Error al extraer token: %v", err)
			utils.WriteJSON(w, http.StatusUnauthorized, payload)
			return
		}

		// Llama al servicio de autenticación para validar el token y verificar el usuario
		// El servicio de autenticación a su vez usará JWTService.ValidJWTToken
		userClaims, err := mw.AuthService.ValidateTokenForMiddleware(r.Context(), tokenString)
		if err != nil {
			payload := utils.JSONResponse{
				Error:   true,
				Message: "Autenticación fallida: " + err.Error(),
			}
			logger.ErrorLog.Printf("Validación de token fallida: %v", err)
			utils.WriteJSON(w, http.StatusUnauthorized, payload)
			return
		}

		// Opcional: Puedes añadir los claims del usuario al contexto de la solicitud
		// para que los handlers posteriores puedan acceder a ellos.
		ctx := context.WithValue(r.Context(), "userClaims", userClaims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
