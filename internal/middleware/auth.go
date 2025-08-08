package middleware

import (
	"cloudtrail-enrichment-api-golang/internal/pkg/logger"
	"cloudtrail-enrichment-api-golang/internal/pkg/token"
	"cloudtrail-enrichment-api-golang/internal/pkg/utils"
	"cloudtrail-enrichment-api-golang/services"
	"context"
	"net/http"
)

type Middleware struct {
	JWTService  *token.JWTService
	AuthService services.AuthService
}

// NewMiddleware creates a new instance of Middleware.
func NewMiddleware(jwtService *token.JWTService, authService services.AuthService) *Middleware {
	return &Middleware{
		JWTService:  jwtService,
		AuthService: authService,
	}
}

// AuthTokenMiddleware is a middleware that validates the JWT token from the request.
func (mw *Middleware) AuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.InfoLog.Println("Executing AuthTokenMiddleware...")

		tokenString, err := mw.JWTService.ExtractJWTToken(r)
		if err != nil {
			payload := utils.JSONResponse{
				Error:   true,
				Message: "Token not provided or invalid format: " + err.Error(),
			}
			logger.ErrorLog.Printf("Error extracting token: %v", err)
			utils.WriteJSON(w, http.StatusUnauthorized, payload)
			return
		}

		// Calls the authentication service to validate the token and check the user.
		userClaims, err := mw.AuthService.ValidateTokenForMiddleware(r.Context(), tokenString)
		if err != nil {
			payload := utils.JSONResponse{
				Error:   true,
				Message: "Authentication failed: " + err.Error(),
			}
			logger.ErrorLog.Printf("Token validation failed: %v", err)
			utils.WriteJSON(w, http.StatusUnauthorized, payload)
			return
		}

		// WIP: Add the user's claims to the request context so that subsequent handlers can access them.
		ctx := context.WithValue(r.Context(), "userClaims", userClaims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
