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

const clockSkew = time.Minute

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

func (s *Service) GenerateJWT(userID uint) (*JWT, error) {
	userIDStr := strconv.FormatUint(uint64(userID), 10)
	exp := time.Now().Add(s.config.JWTExpiration)
	claims := jwt.RegisteredClaims{
		Issuer:    "null3",
		Subject:   userIDStr,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(exp),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrJWTGenerationFailed, err)
	}

	return &JWT{
		Value:  tokenStr,
		UserID: userID,
	}, nil
}

func (s *Service) ParseJWT(tokenStr string) (*JWT, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(s.config.JWTSecret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrJWTInvalid, err)
	}
	if !token.Valid {
		return nil, ErrJWTInvalidClaims
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || claims.Subject == "" {
		return nil, fmt.Errorf("%w: invalid user ID in JWT", ErrJWTInvalidClaims)
	}

	userID, err := strconv.ParseUint(claims.Subject, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid user ID in JWT: %v", ErrJWTInvalidClaims, err)
	}
	if userID == 0 {
		return nil, fmt.Errorf("%w: user ID cannot be zero", ErrJWTInvalidClaims)
	}

	now := time.Now()
	if claims.ExpiresAt.Time.Before(now.Add(clockSkew)) {
		return nil, fmt.Errorf("%w: JWT has expired", ErrJWTExpired)
	}
	if claims.Issuer != "null3" {
		return nil, fmt.Errorf("%w: invalid JWT issuer", ErrJWTInvalidClaims)
	}

	return &JWT{
		Value:  tokenStr,
		UserID: uint(userID),
	}, nil
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
