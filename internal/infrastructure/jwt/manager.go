package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Manager handles JWT generation and validation.
type Manager struct {
	secret     []byte
	expiration time.Duration
	issuer     string
}

// NewManager creates a JWT manager.
func NewManager(secret string, expiration time.Duration, issuer string) *Manager {
	return &Manager{
		secret:     []byte(secret),
		expiration: expiration,
		issuer:     issuer,
	}
}

// GenerateToken issues a signed JWT with the user ID as subject.
func (m *Manager) GenerateToken(userID string) (string, error) {
	now := time.Now().UTC()
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		Issuer:    m.issuer,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(m.expiration)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

// ParseToken validates and returns JWT claims.
func (m *Manager) ParseToken(tokenString string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return m.secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

// ValidateToken returns the subject (user id) if token is valid.
func (m *Manager) ValidateToken(token string) (string, error) {
	claims, err := m.ParseToken(token)
	if err != nil {
		return "", err
	}
	if err := claims.Valid(); err != nil {
		return "", err
	}
	return claims.Subject, nil
}
