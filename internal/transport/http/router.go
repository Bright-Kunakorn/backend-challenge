package http

import (
	stdhttp "net/http"

	jwtinfra "backend-challenge/internal/infrastructure/jwt"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.elastic.co/apm/module/apmchi/v2"
)

// NewRouter wires routes and middleware.
func NewRouter(handler *Handler, jwtManager *jwtinfra.Manager) stdhttp.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(apmchi.Middleware())
	r.Use(LoggingMiddleware)

	r.Post("/auth/register", handler.Register)
	r.Post("/auth/login", handler.Login)

	r.Group(func(group chi.Router) {
		group.Use(AuthMiddleware(jwtManager))
		group.Get("/users", handler.ListUsers)
		group.Get("/users/{id}", handler.GetUser)
		group.Patch("/users/{id}", handler.UpdateUser)
		group.Delete("/users/{id}", handler.DeleteUser)
	})

	return r
}
