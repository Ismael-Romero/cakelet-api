package factory

import (
	"errors"
	"fmt"
	"time"

	"github.com/Ismael-Romero/cakelet-suite/config"
	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken            = errors.New("security/factory: invalid or expired token")
	ErrEmptyUserID             = errors.New("security/factory: user ID cannot be empty")
	ErrTokenGeneration         = errors.New("security/factory: failed to generate token")
	ErrUnexpectedSigningMethod = errors.New("security/factory: unexpected signing method")
	ErrTokenParsing            = errors.New("security/factory: failed to parse token")
)

// TokenClaims defines the custom and registered JWT claims.
type TokenClaims struct {
	UserID  string `json:"user_id,omitempty"`
	Purpose string `json:"purpose,omitempty"`
	MFAHash string `json:"mfa_hash,omitempty"`

	jwt.RegisteredClaims
}

// TokenFactory defines the contract for JWT creation and validation.
type TokenFactory interface {
	GenerateToken(claims TokenClaims, duration time.Duration) (string, error)
	ValidateToken(tokenString string) (*TokenClaims, error)
}

// tokenFactory is the concrete JWT implementation.
type tokenFactory struct {
	secret []byte
}

// NewTokenFactory creates a new TokenFactory instance.
func NewTokenFactory(cfg *config.Config) TokenFactory {
	return &tokenFactory{
		secret: []byte(cfg.Security.JWT.SecretKey),
	}
}

// GenerateToken creates a signed JWT using HS256.
func (f *tokenFactory) GenerateToken(claims TokenClaims, duration time.Duration) (string, error) {

	if claims.UserID == "" {
		return "", ErrEmptyUserID
	}

	now := time.Now()

	claims.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	signedToken, err := token.SignedString(f.secret)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrTokenGeneration, err)
	}

	return signedToken, nil
}

// ValidateToken validates JWT signature, algorithm and claims.
func (f *tokenFactory) ValidateToken(tokenString string) (*TokenClaims, error) {

	token, err := jwt.ParseWithClaims(
		tokenString,
		&TokenClaims{},
		func(t *jwt.Token) (interface{}, error) {

			if t.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("%w: %v", ErrUnexpectedSigningMethod, t.Header["alg"])
			}

			return f.secret, nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrTokenParsing, err)
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
