package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	repo           *Repository
	config         Config
	stubUserConfig StubUserConfig
}

func NewService(repo *Repository, config Config, stubUserConfig StubUserConfig) *Service {
	return &Service{
		repo:           repo,
		config:         config,
		stubUserConfig: stubUserConfig,
	}
}

func (s *Service) Authenticate(req LoginRequest) (*UserResponse, *UserTokenData, error) {
	// For simplicity, we use a stub function to check credentials
	// Don't use this in production!

	login := s.stubUserConfig.Login
	password := s.stubUserConfig.Password

	// Allow case-insensitive login, trim spaces
	if strings.TrimSpace(strings.ToLower(req.Login)) != strings.TrimSpace(strings.ToLower(login)) ||
		req.Password != password {
		return nil, nil, ErrInvalidCredentials
	}

	jwtToken, err := s.GenerateJWT(s.stubUserConfig.UserID)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %v", ErrJWTGenerationFailed, err)
	}

	refreshToken, err := s.CreateRefreshToken(s.stubUserConfig.UserID)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %v", ErrRefreshTokenCreationFailed, err)
	}

	res := &UserResponse{
		ID: s.stubUserConfig.UserID,
	}

	tokenData := NewUserTokenData(jwtToken, refreshToken)

	return res, tokenData, nil
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
	if claims.ExpiresAt.Time.Before(now) || claims.ExpiresAt.Time.Equal(now) {
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

func (s *Service) CreateRefreshToken(userID uint) (*RefreshToken, error) {
	now := time.Now()

	tokenString, err := generateRandomRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRefreshTokenCreationFailed, err)
	}
	token := &RefreshToken{
		UserID:    userID,
		Value:     tokenString,
		CreatedAt: now,
		ExpiresAt: now.Add(s.config.RefreshTokenExpiration),
	}
	createdToken, err := s.repo.SaveRefreshToken(token)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRefreshTokenCreationFailed, err)
	}
	return createdToken, nil
}

func (s *Service) GetRefreshToken(tokenString string) (*RefreshToken, error) {
	return s.repo.GetRefreshToken(tokenString)
}

func (s *Service) InvalidateRefreshToken(tokenString string) error {
	token, err := s.repo.GetRefreshToken(tokenString)
	if err != nil {
		return err
	}

	return s.repo.DeleteRefreshToken(token)
}

func generateRandomRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate random refresh token: %v", err)
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
