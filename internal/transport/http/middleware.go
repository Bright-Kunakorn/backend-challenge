package http

import (
	"log"
	stdhttp "net/http"
	"strings"
	"time"

	jwtinfra "backend-challenge/internal/infrastructure/jwt"
	"backend-challenge/internal/transport/authctx"
)

// LoggingMiddleware logs method, path, status, and duration.
func LoggingMiddleware(next stdhttp.Handler) stdhttp.Handler {
	return stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		start := time.Now()
		ww := &statusWriter{ResponseWriter: w, status: stdhttp.StatusOK}
		next.ServeHTTP(ww, r)
		duration := time.Since(start)
		log.Printf("%s %s %d %s", r.Method, r.URL.Path, ww.status, duration)
	})
}

// AuthMiddleware ensures requests have a valid JWT.
func AuthMiddleware(manager *jwtinfra.Manager) func(stdhttp.Handler) stdhttp.Handler {
	return func(next stdhttp.Handler) stdhttp.Handler {
		return stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				stdhttp.Error(w, "missing authorization header", stdhttp.StatusUnauthorized)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				stdhttp.Error(w, "invalid authorization header", stdhttp.StatusUnauthorized)
				return
			}

			userID, err := manager.ValidateToken(parts[1])
			if err != nil {
				stdhttp.Error(w, "invalid token", stdhttp.StatusUnauthorized)
				return
			}

			ctx := authctx.WithUserID(r.Context(), userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

type statusWriter struct {
	stdhttp.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}
