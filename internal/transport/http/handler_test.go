package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"backend-challenge/internal/application"
	"backend-challenge/internal/domain"
	jwtinfra "backend-challenge/internal/infrastructure/jwt"
	"backend-challenge/internal/transport/authctx"
	transport "backend-challenge/internal/transport/http"

	"backend-challenge/internal/infrastructure/memory"

	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
)

type fakeRepo struct {
	createFn   func(context.Context, domain.User) (domain.User, error)
	getByEmail func(context.Context, string) (domain.User, error)
	getByID    func(context.Context, string) (domain.User, error)
	listFn     func(context.Context) ([]domain.User, error)
	updateFn   func(context.Context, string, domain.UpdateUser) (domain.User, error)
	deleteFn   func(context.Context, string) error
	countFn    func(context.Context) (int64, error)
}

func (f *fakeRepo) Create(ctx context.Context, user domain.User) (domain.User, error) {
	if f.createFn != nil {
		return f.createFn(ctx, user)
	}
	return user, nil
}

func (f *fakeRepo) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	if f.getByEmail != nil {
		return f.getByEmail(ctx, email)
	}
	return domain.User{}, application.ErrNotFound
}

func (f *fakeRepo) GetByID(ctx context.Context, id string) (domain.User, error) {
	if f.getByID != nil {
		return f.getByID(ctx, id)
	}
	return domain.User{}, application.ErrNotFound
}

func (f *fakeRepo) List(ctx context.Context) ([]domain.User, error) {
	if f.listFn != nil {
		return f.listFn(ctx)
	}
	return nil, nil
}

func (f *fakeRepo) Update(ctx context.Context, id string, update domain.UpdateUser) (domain.User, error) {
	if f.updateFn != nil {
		return f.updateFn(ctx, id, update)
	}
	return domain.User{}, nil
}

func (f *fakeRepo) Delete(ctx context.Context, id string) error {
	if f.deleteFn != nil {
		return f.deleteFn(ctx, id)
	}
	return nil
}

func (f *fakeRepo) Count(ctx context.Context) (int64, error) {
	if f.countFn != nil {
		return f.countFn(ctx)
	}
	return 0, nil
}

func TestRegisterHandlerSuccess(t *testing.T) {
	repo := &fakeRepo{
		getByEmail: func(context.Context, string) (domain.User, error) {
			return domain.User{}, application.ErrNotFound
		},
		createFn: func(context.Context, domain.User) (domain.User, error) {
			return domain.User{ID: "1", Name: "Test", Email: "test@example.com", CreatedAt: time.Now()}, nil
		},
	}
	service := application.NewUserService(repo)
	jwtManager := jwtinfra.NewManager("secret", time.Hour, "issuer")
	handler := transport.NewHandler(service, jwtManager)

	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(`{"name":"Test","email":"test@example.com","password":"strongpass"}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.Register(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status 201 got %d", rr.Code)
	}

	var body map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("parse response: %v", err)
	}
	if body["token"] == "" {
		t.Fatalf("expected token in response")
	}
}

func TestRegisterHandlerError(t *testing.T) {
	repo := &fakeRepo{
		getByEmail: func(context.Context, string) (domain.User, error) {
			return domain.User{}, application.ErrNotFound
		},
		createFn: func(context.Context, domain.User) (domain.User, error) {
			return domain.User{}, application.ErrDuplicateEmail
		},
	}
	service := application.NewUserService(repo)
	jwtManager := jwtinfra.NewManager("secret", time.Hour, "issuer")
	handler := transport.NewHandler(service, jwtManager)

	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(`{"name":"Test","email":"test@example.com","password":"strongpass"}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.Register(rr, req)

	if rr.Code != http.StatusConflict {
		t.Fatalf("expected 409 got %d", rr.Code)
	}
}

func TestRegisterHandlerInvalidPayload(t *testing.T) {
	repo := &fakeRepo{}
	service := application.NewUserService(repo)
	jwtManager := jwtinfra.NewManager("secret", time.Hour, "issuer")
	handler := transport.NewHandler(service, jwtManager)

	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(`{`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.Register(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 got %d", rr.Code)
	}
}

func TestLoginHandler(t *testing.T) {
	repo := &fakeRepo{
		getByEmail: func(context.Context, string) (domain.User, error) {
			return domain.User{ID: "1", Email: "test@example.com", Password: hashPassword(t, "pass")}, nil
		},
	}
	service := application.NewUserService(repo)
	jwtManager := jwtinfra.NewManager("secret", time.Hour, "issuer")
	handler := transport.NewHandler(service, jwtManager)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString(`{"email":"test@example.com","password":"pass"}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.Login(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rr.Code)
	}
}

func TestLoginHandlerError(t *testing.T) {
	repo := &fakeRepo{}
	service := application.NewUserService(repo)
	jwtManager := jwtinfra.NewManager("secret", time.Hour, "issuer")
	handler := transport.NewHandler(service, jwtManager)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBufferString(`{"email":"test@example.com","password":"wrong"}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.Login(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d", rr.Code)
	}
}

func TestListAndGetHandlers(t *testing.T) {
	repo := &fakeRepo{
		listFn: func(context.Context) ([]domain.User, error) {
			return []domain.User{{ID: "1", Name: "A", Email: "a@example.com"}}, nil
		},
		getByID: func(context.Context, string) (domain.User, error) {
			return domain.User{ID: "1", Name: "A", Email: "a@example.com"}, nil
		},
	}
	service := application.NewUserService(repo)
	jwtManager := jwtinfra.NewManager("secret", time.Hour, "issuer")
	handler := transport.NewHandler(service, jwtManager)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	handler.ListUsers(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rr.Code)
	}

	rr = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/users/1", nil)
	req = withRouteParam(req, "id", "1")
	req = req.WithContext(authctx.WithUserID(req.Context(), "1"))
	handler.GetUser(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rr.Code)
	}
}

func TestListHandlerError(t *testing.T) {
	repo := &fakeRepo{
		listFn: func(context.Context) ([]domain.User, error) {
			return nil, errors.New("list failed")
		},
	}
	service := application.NewUserService(repo)
	jwtManager := jwtinfra.NewManager("secret", time.Hour, "issuer")
	handler := transport.NewHandler(service, jwtManager)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	handler.ListUsers(rr, req)
	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 got %d", rr.Code)
	}
}

func TestGetHandlerNotFound(t *testing.T) {
	repo := &fakeRepo{
		getByID: func(context.Context, string) (domain.User, error) {
			return domain.User{}, application.ErrNotFound
		},
	}
	service := application.NewUserService(repo)
	jwtManager := jwtinfra.NewManager("secret", time.Hour, "issuer")
	handler := transport.NewHandler(service, jwtManager)

	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	req = withRouteParam(req, "id", "1")
	rr := httptest.NewRecorder()
	handler.GetUser(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404 got %d", rr.Code)
	}
}

func TestUpdateAndDeleteHandlers(t *testing.T) {
	repo := &fakeRepo{
		updateFn: func(context.Context, string, domain.UpdateUser) (domain.User, error) {
			return domain.User{ID: "1", Name: "Updated", Email: "a@example.com"}, nil
		},
		deleteFn: func(context.Context, string) error {
			return nil
		},
	}
	service := application.NewUserService(repo)
	jwtManager := jwtinfra.NewManager("secret", time.Hour, "issuer")
	handler := transport.NewHandler(service, jwtManager)

	req := httptest.NewRequest(http.MethodPatch, "/users/1", bytes.NewBufferString(`{"name":"Updated"}`))
	req = withRouteParam(req, "id", "1")
	req = req.WithContext(authctx.WithUserID(req.Context(), "1"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.UpdateUser(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rr.Code)
	}

	req = httptest.NewRequest(http.MethodDelete, "/users/1", nil)
	req = withRouteParam(req, "id", "1")
	req = req.WithContext(authctx.WithUserID(req.Context(), "1"))
	rr = httptest.NewRecorder()
	handler.DeleteUser(rr, req)
	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected 204 got %d", rr.Code)
	}
}

func TestUpdateForbidden(t *testing.T) {
	repo := &fakeRepo{}
	service := application.NewUserService(repo)
	handler := transport.NewHandler(service, jwtinfra.NewManager("secret", time.Hour, "issuer"))

	req := httptest.NewRequest(http.MethodPatch, "/users/1", bytes.NewBufferString(`{"name":"Updated"}`))
	req = withRouteParam(req, "id", "1")
	req = req.WithContext(authctx.WithUserID(req.Context(), "other"))
	rr := httptest.NewRecorder()
	handler.UpdateUser(rr, req)
	if rr.Code != http.StatusForbidden {
		t.Fatalf("expected 403 got %d", rr.Code)
	}
}

func TestUpdateHandlerError(t *testing.T) {
	repo := &fakeRepo{
		updateFn: func(context.Context, string, domain.UpdateUser) (domain.User, error) {
			return domain.User{}, application.ErrDuplicateEmail
		},
	}
	service := application.NewUserService(repo)
	handler := transport.NewHandler(service, jwtinfra.NewManager("secret", time.Hour, "issuer"))

	req := httptest.NewRequest(http.MethodPatch, "/users/1", bytes.NewBufferString(`{"email":"dup@example.com"}`))
	req = withRouteParam(req, "id", "1")
	req = req.WithContext(authctx.WithUserID(req.Context(), "1"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler.UpdateUser(rr, req)
	if rr.Code != http.StatusConflict {
		t.Fatalf("expected 409 got %d", rr.Code)
	}
}

func TestDeleteHandlerError(t *testing.T) {
	repo := &fakeRepo{
		deleteFn: func(context.Context, string) error {
			return application.ErrNotFound
		},
	}
	service := application.NewUserService(repo)
	handler := transport.NewHandler(service, jwtinfra.NewManager("secret", time.Hour, "issuer"))

	req := httptest.NewRequest(http.MethodDelete, "/users/1", nil)
	req = withRouteParam(req, "id", "1")
	req = req.WithContext(authctx.WithUserID(req.Context(), "1"))
	rr := httptest.NewRecorder()
	handler.DeleteUser(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404 got %d", rr.Code)
	}
}

func TestAuthMiddleware(t *testing.T) {
	manager := jwtinfra.NewManager("secret", time.Hour, "issuer")
	token, err := manager.GenerateToken("abc")
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	middleware := transport.AuthMiddleware(manager)
	called := false
	h := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		if id, ok := authctx.UserIDFromContext(r.Context()); !ok || id != "abc" {
			t.Fatalf("expected user id in context")
		}
	}))
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	h.ServeHTTP(httptest.NewRecorder(), req)
	if !called {
		t.Fatalf("expected handler to be called")
	}

	rr := httptest.NewRecorder()
	middleware(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})).ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/protected", nil))
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d", rr.Code)
	}
}

func TestAuthMiddlewareInvalidHeader(t *testing.T) {
	manager := jwtinfra.NewManager("secret", time.Hour, "issuer")
	middleware := transport.AuthMiddleware(manager)
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Token invalid")
	rr := httptest.NewRecorder()
	middleware(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})).ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d", rr.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid.token")
	rr = httptest.NewRecorder()
	middleware(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})).ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d", rr.Code)
	}
}

func TestLoggingMiddleware(t *testing.T) {
	handler := transport.LoggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	}))

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/", nil))
	if rr.Code != http.StatusTeapot {
		t.Fatalf("expected 418 got %d", rr.Code)
	}
}

func TestNewRouterIntegration(t *testing.T) {
	repo := memory.NewUserRepository()
	service := application.NewUserService(repo)
	manager := jwtinfra.NewManager("secret", time.Hour, "issuer")
	handler := transport.NewHandler(service, manager)
	router := transport.NewRouter(handler, manager)

	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(`{"name":"Test","email":"route@example.com","password":"pass12345"}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201 got %d", rr.Code)
	}
}

func hashPassword(t *testing.T, password string) string {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	return string(hash)
}

func withRouteParam(r *http.Request, key, value string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(key, value)
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}
