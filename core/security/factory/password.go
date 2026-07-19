package factory

import (
	"errors"
	"fmt"

	"github.com/Ismael-Romero/cakelet-suite/config"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidHashFormat = errors.New("security/factory: invalid or corrupted password hash format")
	ErrInternalCrypto    = errors.New("security/factory: internal cryptographic operation failed")
)

type PasswordFactory struct {
	workFactor int
}

// NewPasswordFactory creates a new instance of PasswordFactory.
// It allows configuring the optimal work factor (cost) for the runtime environment
func NewPasswordFactory(cfg *config.Config) *PasswordFactory {
	workFactor := cfg.Security.Bcrypt.WorkFactor

	if workFactor < bcrypt.MinCost || workFactor > bcrypt.MaxCost {
		workFactor = bcrypt.DefaultCost
	}

	return &PasswordFactory{
		workFactor: workFactor,
	}
}

// HashPassword takes a plaintext password and generates a BCrypt-encoded hash.
// It automatically applies a unique random salt and the configured work factor.
func (f *PasswordFactory) HashPassword(plaintext string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plaintext), f.workFactor)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrInternalCrypto, err)
	}

	return string(bytes), nil
}

// VerifyPassword validates whether a plaintext password matches a previously stored BCrypt hash.
// Returns true if it matches, false if it does not (without an error).
// Returns an error only if the hash format is invalid or corrupt.
func (f *PasswordFactory) VerifyPassword(plaintext, hashedValue string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedValue), []byte(plaintext))
	if err == nil {
		return true, nil
	}

	// In accordance with security requirements, we distinguish authentication failures
	// from internal failures related to format or cryptographic configuration.
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return false, nil
	}

	// Any other error indicates that the stored hash is corrupt or invalid.
	return false, fmt.Errorf("%w: %v", ErrInvalidHashFormat, err)
}
