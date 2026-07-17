package session

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/azaviyalov/null3/backend/internal/core"
	"github.com/golang-jwt/jwt/v5"
)

const (
	clockSkew  = time.Minute
	userScope  = "user"
	adminScope = "admin"
)

type accessTokenClaims struct {
	Scope string `json:"scope"`
	jwt.RegisteredClaims
}

type Service struct {
	repo   *Repository
	config Config
}

func NewService(repo *Repository, config Config) *Service {
	return &Service{
		repo:   repo,
		config: config,
	}
}

func (s *Service) GenerateUserAccessToken(userID uint) (string, error) {
	return s.generateAccessToken(strconv.FormatUint(uint64(userID), 10), userScope, s.config.JWTExpiration)
}

func (s *Service) GenerateAdminAccessToken(expiration time.Duration) (string, error) {
	return s.generateAccessToken("admin", adminScope, expiration)
}

func (s *Service) generateAccessToken(subject, scope string, expiration time.Duration) (string, error) {
	now := time.Now()
	tokenClaims := accessTokenClaims{
		Scope: scope,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "null3",
			Subject:   subject,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(expiration)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
	tokenStr, err := token.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrJWTGenerationFailed, err)
	}
	return tokenStr, nil
}

func (s *Service) parseAccessTokenClaims(tokenStr string) (*accessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &accessTokenClaims{}, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, ErrJWTInvalidClaims
		}
		return []byte(s.config.JWTSecret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrJWTInvalid, err)
	}
	if !token.Valid {
		return nil, ErrJWTInvalidClaims
	}

	tokenClaims, ok := token.Claims.(*accessTokenClaims)
	if !ok || tokenClaims.Subject == "" || tokenClaims.ExpiresAt == nil {
		return nil, ErrJWTInvalidClaims
	}
	if tokenClaims.ExpiresAt.Time.Before(time.Now().Add(clockSkew)) {
		return nil, fmt.Errorf("%w: JWT has expired", ErrJWTExpired)
	}
	if tokenClaims.Issuer != "null3" {
		return nil, fmt.Errorf("%w: invalid JWT issuer", ErrJWTInvalidClaims)
	}
	return tokenClaims, nil
}

func (s *Service) ParseUserAccessToken(tokenStr string) (uint, error) {
	tokenClaims, err := s.parseAccessTokenClaims(tokenStr)
	if err != nil {
		return 0, err
	}
	if tokenClaims.Scope != userScope {
		return 0, fmt.Errorf("%w: user scope required", ErrJWTInvalidClaims)
	}
	userID, err := strconv.ParseUint(tokenClaims.Subject, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%w: invalid user ID in JWT: %v", ErrJWTInvalidClaims, err)
	}
	if userID == 0 {
		return 0, fmt.Errorf("%w: user ID cannot be zero", ErrJWTInvalidClaims)
	}

	return uint(userID), nil
}

func (s *Service) ValidateAdminAccessToken(tokenStr string) error {
	tokenClaims, err := s.parseAccessTokenClaims(tokenStr)
	if err != nil {
		return err
	}
	if tokenClaims.Scope != adminScope || tokenClaims.Subject != "admin" {
		return fmt.Errorf("%w: admin subject and scope required", ErrJWTInvalidClaims)
	}
	return nil
}

func (s *Service) CreateRefreshToken(ctx context.Context, userID uint) (*RefreshToken, error) {
	return s.CreateRefreshTokenWithRepo(ctx, s.repo, userID)
}

func (s *Service) CreateRefreshTokenWithRepo(ctx context.Context, repo *Repository, userID uint) (*RefreshToken, error) {
	now := time.Now()

	tokenString, err := generateRandomToken()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRefreshTokenCreationFailed, err)
	}

	token := &RefreshToken{
		UserID:    userID,
		Value:     tokenString,
		CreatedAt: now,
		ExpiresAt: now.Add(s.config.RefreshTokenExpiration),
	}

	createdToken, err := repo.SaveRefreshToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRefreshTokenCreationFailed, err)
	}

	return createdToken, nil
}

func (s *Service) InvalidateRefreshToken(ctx context.Context, tokenString string) error {
	token, err := s.repo.GetRefreshToken(ctx, tokenString)
	if err != nil {
		if errors.Is(err, core.ErrItemNotFound) {
			return nil
		}
		return err
	}

	return s.repo.DeleteRefreshToken(ctx, token)
}

func (s *Service) DeleteExpiredRefreshTokens(ctx context.Context) error {
	return s.repo.DeleteExpiredRefreshTokens(ctx)
}

func generateRandomToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate random token: %v", err)
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}
