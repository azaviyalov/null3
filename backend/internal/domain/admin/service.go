package admin

import (
	"crypto/sha256"
	"crypto/subtle"
	"errors"
	"time"

	"github.com/azaviyalov/null3/backend/internal/domain/session"
)

const adminAccessTokenTTL = 30 * time.Minute

var ErrInvalidCredentials = errors.New("invalid credentials")

type Service struct {
	passwordHash [sha256.Size]byte
	tokens       *session.Service
}

func NewService(password string, tokens *session.Service) *Service {
	return &Service{passwordHash: sha256.Sum256([]byte(password)), tokens: tokens}
}

func (s *Service) Authenticate(password string) (string, error) {
	candidate := sha256.Sum256([]byte(password))
	if subtle.ConstantTimeCompare(candidate[:], s.passwordHash[:]) != 1 || password == "" {
		return "", ErrInvalidCredentials
	}
	return s.tokens.GenerateAdminAccessToken(adminAccessTokenTTL)
}
