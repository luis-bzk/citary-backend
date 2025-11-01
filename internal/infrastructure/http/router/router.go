package router

import (
	"citary-backend/internal/infrastructure/http/handlers/auth"
	"citary-backend/internal/infrastructure/http/middleware"
	"net/http"
)

// Router manages HTTP route configuration
type Router struct {
	authHandler *auth.AuthHandler
}

// NewRouter creates a new Router instance
func NewRouter(authHandler *auth.AuthHandler) *Router {
	return &Router{
		authHandler: authHandler,
	}
}

// SetupRoutes configures all HTTP routes and returns the configured handler
func (rt *Router) SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	// Auth routes
	mux.HandleFunc("/auth/signup", rt.authHandler.SignupUser)

	// Health check route
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Apply middleware chain (order matters: Recovery -> CORS -> Logging -> routes)
	handler := middleware.Recovery(mux)
	handler = middleware.CORS(handler)
	handler = middleware.Logging(handler)

	return handler
}
