package main

import (
	"cloudtrail-enrichment-api-golang/internal/pkg/utils"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Grupo de rutas para la versión 1 de la API
	mux.Route("/v1", func(r chi.Router) {
		// Rutas públicas de la V1

		r.Get("/health", app.systemController.HealthCheck) // Ruta de salud pública
		// r.Get("/health", app.systemController.HealthCheck)
		// r.Get("/public", app.authController.PublicRouteHandler) // Ejemplo de ruta pública
		// r.Post("/register", app.authController.RegisterUser)
		// r.Post("/authenticate", app.authController.AuthenticateUser)

		// Rutas protegidas por el middleware de autenticación de la V1
		// r.Route("/products", func(r chi.Router) {
		// r.Use(app.middleware.AuthTokenMiddleware)
		// r.Use(app.middleware.AuthTokenMiddleware)
		// r.Post("/", app.productsController.CreateProduct)
		// r.Get("/", app.productsController.ListProducts)
		// r.Get("/{uuid}", app.productsController.GetProductByUUID)
		// r.Put("/{uuid}", app.productsController.UpdateProductByUUID)
		// r.Delete("/{uuid}", app.productsController.DeleteProductByUUID)
		// })

		r.Route("/admin", func(r chi.Router) {
			r.Use(app.middleware.AuthTokenMiddleware)
			// Puedes añadir middleware de autorización aquí si tienes roles
			r.Get("/dashboard", app.AdminDashboard)
		})
	})

	return mux
}

func (app *application) Health(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status": "OK",
		"uptime": "server is healthy",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// AdminDashboard es un handler de ejemplo para una ruta protegida.
func (app *application) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"message": "Bienvenido al panel de administración!",
	}
	utils.WriteJSON(w, http.StatusOK, utils.JSONResponse{
		Error:   false,
		Message: "Dashboard de administración",
		Data:    response,
	})
}
