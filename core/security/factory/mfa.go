package factory

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/Ismael-Romero/cakelet-suite/config"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidChallenge    = errors.New("security/factory: invalid MFA challenge")
	ErrChallengeGeneration = errors.New("security/factory: failed to generate secure challenge")
	ErrHashGeneration      = errors.New("security/factory: failed to generate hash for OTP")
)

// MFAChallenge contains the information required to validate a second factor.
type MFAChallenge struct {
	CodeHash  string
	ExpiresAt time.Time
}

// MFAFactory defines the contract for creating and verifying MFA challenges.
type MFAFactory interface {
	// GenerateChallenge generates a cryptographically secure one-time verification challenge
	// returning both the plaintext code (for delivery) and the hashed challenge representation.
	GenerateChallenge() (string, *MFAChallenge, error)

	// ValidateChallenge validates whether a user-provided MFA code matches the previously generated challenge hash.
	ValidateChallenge(code string, challengeHash string) (bool, error)
}

// mfaFactory is the concrete implementation of MFAFactory.
type mfaFactory struct {
	otpLength  int
	expiration time.Duration
	workFactor int
}

// NewMFAFactory creates a new instance of MFAFactory.
func NewMFAFactory(cfg *config.Config) MFAFactory {
	expiration := time.Duration(cfg.Security.MFA.ExpirationMinutes) * time.Minute

	return &mfaFactory{
		otpLength:  cfg.Security.MFA.OTPLength,
		expiration: expiration,
		workFactor: cfg.Security.Bcrypt.WorkFactor,
	}
}

// GenerateChallenge creates a cryptographically secure OTP and its hashed representation.
func (f *mfaFactory) GenerateChallenge() (string, *MFAChallenge, error) {
	code, err := generateSecureOTP(f.otpLength)
	if err != nil {
		return "", nil, fmt.Errorf("%w: %v", ErrChallengeGeneration, err)
	}

	hashBytes, err := bcrypt.GenerateFromPassword([]byte(code), f.workFactor)
	if err != nil {
		return "", nil, fmt.Errorf("%w: %v", ErrHashGeneration, err)
	}

	challenge := &MFAChallenge{
		CodeHash:  string(hashBytes),
		ExpiresAt: time.Now().Add(f.expiration),
	}

	return code, challenge, nil
}

// ValidateChallenge validates whether a user-provided MFA code matches the generated challenge.
func (f *mfaFactory) ValidateChallenge(code string, challengeHash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(challengeHash), []byte(code))
	if err == nil {
		return true, nil
	}

	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return false, nil
	}

	return false, fmt.Errorf("%w: %v", ErrInvalidChallenge, err)
}

// generateSecureOTP is a helper function that generates a cryptographically secure random numeric string.
func generateSecureOTP(length int) (string, error) {
	if length <= 0 {
		return "", errors.New("OTP length must be greater than 0")
	}

	max := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(length)), nil)

	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}

	format := fmt.Sprintf("%%0%dd", length)
	return fmt.Sprintf(format, n), nil
}
