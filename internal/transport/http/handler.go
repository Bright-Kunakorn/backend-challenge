package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"backend-challenge/internal/application"
	"backend-challenge/internal/domain"
	jwtinfra "backend-challenge/internal/infrastructure/jwt"

	"github.com/go-chi/chi/v5"
)

// Handler bundles HTTP handlers for user operations.
type Handler struct {
	service    *application.UserService
	jwtManager *jwtinfra.Manager
}

// NewHandler builds a handler.
func NewHandler(service *application.UserService, manager *jwtinfra.Manager) *Handler {
	return &Handler{
		service:    service,
		jwtManager: manager,
	}
}

type registerRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type updateRequest struct {
	Name  *string `json:"name,omitempty"`
	Email *string `json:"email,omitempty"`
}

type authResponse struct {
	Token string           `json:"token"`
	User  domain.UserPublic `json:"user"`
}

// Register handles account creation.
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var payload registerRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	user, err := h.service.Register(r.Context(), application.RegisterInput{
		Name:     payload.Name,
		Email:    payload.Email,
		Password: payload.Password,
	})
	if err != nil {
		handleError(w, err)
		return
	}

	token, err := h.jwtManager.GenerateToken(user.ID)
	if err != nil {
		http.Error(w, "failed to create token", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, authResponse{
		Token: token,
		User:  user.Sanitize(),
	})
}

// Login authenticates a user and returns JWT.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var payload loginRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	user, err := h.service.Authenticate(r.Context(), payload.Email, payload.Password)
	if err != nil {
		handleError(w, err)
		return
	}

	token, err := h.jwtManager.GenerateToken(user.ID)
	if err != nil {
		http.Error(w, "failed to create token", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, authResponse{
		Token: token,
		User:  user.Sanitize(),
	})
}

// ListUsers returns all users (sans passwords).
func (h *Handler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.List(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}

	public := make([]domain.UserPublic, 0, len(users))
	for _, u := range users {
		public = append(public, u.Sanitize())
	}

	writeJSON(w, http.StatusOK, public)
}

// GetUser returns single user info.
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user, err := h.service.Get(r.Context(), id)
	if err != nil {
		handleError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, user.Sanitize())
}

// UpdateUser updates allowed fields.
func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var payload updateRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	updated, err := h.service.Update(r.Context(), id, application.UpdateInput{
		Name:  payload.Name,
		Email: payload.Email,
	})
	if err != nil {
		handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, updated.Sanitize())
}

// DeleteUser removes a user.
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.service.Delete(r.Context(), id); err != nil {
		handleError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, application.ErrDuplicateEmail):
		http.Error(w, err.Error(), http.StatusConflict)
	case errors.Is(err, application.ErrInvalidCredentials):
		http.Error(w, err.Error(), http.StatusUnauthorized)
	case errors.Is(err, application.ErrNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, application.ErrNoFieldsToUpdate),
		errors.Is(err, domain.ErrInvalidEmail),
		errors.Is(err, domain.ErrInvalidName),
		errors.Is(err, domain.ErrInvalidPassword):
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
