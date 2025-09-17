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

const clockSkew = time.Minute * 1

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
	// Try to authenticate with real user data first
	user, err := s.repo.GetUserByEmail(req.Login)
	if err == nil && user.CheckPassword(req.Password) {
		jwtToken, err := s.GenerateJWT(user.ID)
		if err != nil {
			return nil, nil, fmt.Errorf("%w: %v", ErrJWTGenerationFailed, err)
		}

		refreshToken, err := s.CreateRefreshToken(user.ID)
		if err != nil {
			return nil, nil, fmt.Errorf("%w: %v", ErrRefreshTokenCreationFailed, err)
		}

		res := &UserResponse{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
		}

		tokenData := NewUserTokenData(jwtToken, refreshToken)
		return res, tokenData, nil
	}

	// Fall back to stub user for compatibility (legacy login)
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

func (s *Service) AuthenticateAdmin(req AdminLoginRequest) (*UserTokenData, error) {
	if req.Username != s.config.AdminUsername || req.Password != s.config.AdminPassword {
		return nil, ErrInvalidCredentials
	}

	// Use a special admin user ID (0) to distinguish from regular users
	const adminUserID = 0
	jwtToken, err := s.GenerateAdminJWT(adminUserID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrJWTGenerationFailed, err)
	}

	refreshToken, err := s.CreateRefreshToken(adminUserID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRefreshTokenCreationFailed, err)
	}

	tokenData := NewUserTokenData(jwtToken, refreshToken)
	return tokenData, nil
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

func (s *Service) GenerateAdminJWT(userID uint) (*JWT, error) {
	userIDStr := strconv.FormatUint(uint64(userID), 10)
	exp := time.Now().Add(s.config.JWTExpiration)
	claims := jwt.RegisteredClaims{
		Issuer:    "null3-admin",
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
	now := time.Now()
	if claims.ExpiresAt.Time.Before(now.Add(clockSkew)) {
		return nil, fmt.Errorf("%w: JWT has expired", ErrJWTExpired)
	}
	if claims.Issuer != "null3" && claims.Issuer != "null3-admin" {
		return nil, fmt.Errorf("%w: invalid JWT issuer", ErrJWTInvalidClaims)
	}

	return &JWT{
		Value:  tokenStr,
		UserID: uint(userID),
	}, nil
}

func (s *Service) IsAdminJWT(tokenStr string) bool {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(s.config.JWTSecret), nil
	})
	if err != nil {
		return false
	}
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return false
	}
	return claims.Issuer == "null3-admin"
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

// User management methods
func (s *Service) CreateUser(req CreateUserRequest) (*UserResponse, error) {
	// Check if user already exists
	existingUser, err := s.repo.GetUserByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	user := &User{
		Email: req.Email,
		Name:  req.Name,
	}

	if err := user.SetPassword(req.Password); err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	createdUser, err := s.repo.CreateUser(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	return &UserResponse{
		ID:    createdUser.ID,
		Email: createdUser.Email,
		Name:  createdUser.Name,
	}, nil
}

func (s *Service) GetUserByID(id uint) (*UserResponse, error) {
	user, err := s.repo.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	return &UserResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	}, nil
}

func (s *Service) GetAllUsers() ([]*UserResponse, error) {
	users, err := s.repo.GetAllUsers()
	if err != nil {
		return nil, err
	}

	var userResponses []*UserResponse
	for _, user := range users {
		userResponses = append(userResponses, &UserResponse{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
		})
	}

	return userResponses, nil
}

func (s *Service) UpdateUser(id uint, req UpdateUserRequest) (*UserResponse, error) {
	user, err := s.repo.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Password != nil {
		if err := user.SetPassword(*req.Password); err != nil {
			return nil, fmt.Errorf("failed to hash password: %v", err)
		}
	}

	user.UpdatedAt = time.Now()

	updatedUser, err := s.repo.UpdateUser(user)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %v", err)
	}

	return &UserResponse{
		ID:    updatedUser.ID,
		Email: updatedUser.Email,
		Name:  updatedUser.Name,
	}, nil
}

func (s *Service) DeleteUser(id uint) error {
	return s.repo.DeleteUser(id)
}

func (s *Service) GetRefreshTokensForUser(userID uint) ([]*RefreshToken, error) {
	return s.repo.GetRefreshTokensByUserID(userID)
}

func (s *Service) GetAllRefreshTokens() ([]*RefreshToken, error) {
	return s.repo.GetAllRefreshTokens()
}
