package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/azaviyalov/null3/backend/internal/core"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	clockSkew              = time.Minute
	adminUserID       uint = 1
	adminLogin             = "admin"
	adminEmail             = "admin@example.com"
	adminPassword          = "password"
	inviteExpiration       = 24 * time.Hour
	minPasswordLength      = 8
	maxPasswordLength      = 72
)

var loginPattern = regexp.MustCompile(`^[a-z0-9_-]{3,32}$`)

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

func (s *Service) SeedAdminUser(ctx context.Context) error {
	_, err := s.repo.GetUserByID(ctx, adminUserID)
	if err == nil {
		return nil
	}
	if !errors.Is(err, core.ErrItemNotFound) {
		return err
	}

	passwordHash, err := hashPassword(adminPassword)
	if err != nil {
		return fmt.Errorf("failed to hash admin password: %w", err)
	}

	_, err = s.repo.CreateUser(ctx, &User{
		ID:           adminUserID,
		Login:        adminLogin,
		Email:        adminEmail,
		PasswordHash: passwordHash,
	})
	if err != nil {
		if isUniqueConstraintError(err) {
			return nil
		}
		return fmt.Errorf("failed to seed admin user: %w", err)
	}

	return nil
}

func (s *Service) AuthenticateUser(ctx context.Context, req LoginRequest) (*UserResponse, *UserTokenData, error) {
	user, err := s.authenticateByLogin(ctx, req)
	if err != nil {
		return nil, nil, err
	}
	if s.IsAdmin(user) {
		return nil, nil, ErrInvalidCredentials
	}

	tokenData, err := s.createUserSession(ctx, user)
	if err != nil {
		return nil, nil, err
	}
	return NewUserResponse(user), tokenData, nil
}

func (s *Service) AuthenticateAdmin(ctx context.Context, req LoginRequest) (*UserResponse, *UserTokenData, error) {
	user, err := s.authenticateByLogin(ctx, req)
	if err != nil {
		return nil, nil, err
	}
	if !s.IsAdmin(user) {
		return nil, nil, ErrInvalidCredentials
	}

	tokenData, err := s.createAdminSession(ctx, user)
	if err != nil {
		return nil, nil, err
	}
	return NewUserResponse(user), tokenData, nil
}

func (s *Service) authenticateByLogin(ctx context.Context, req LoginRequest) (*User, error) {
	login := normalizeLogin(req.Login)
	if login == "" || strings.TrimSpace(req.Password) == "" {
		return nil, ErrInvalidCredentials
	}

	user, err := s.repo.GetUserByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, core.ErrItemNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}
	return user, nil
}

func (s *Service) GetUserByID(ctx context.Context, id uint) (*User, error) {
	return s.repo.GetUserByID(ctx, id)
}

func (s *Service) IsAdmin(user *User) bool {
	return user != nil && user.ID == adminUserID
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

func (s *Service) RefreshUserSession(ctx context.Context, tokenString string) (*UserResponse, *UserTokenData, error) {
	return s.refreshSession(ctx, tokenString, false)
}

func (s *Service) RefreshAdminSession(ctx context.Context, tokenString string) (*UserResponse, *UserTokenData, error) {
	return s.refreshSession(ctx, tokenString, true)
}

func (s *Service) refreshSession(ctx context.Context, tokenString string, requireAdmin bool) (*UserResponse, *UserTokenData, error) {
	token, err := s.repo.GetRefreshToken(ctx, tokenString)
	if err != nil {
		if errors.Is(err, core.ErrItemNotFound) {
			return nil, nil, ErrRefreshTokenInvalid
		}
		return nil, nil, err
	}

	if token.ExpiresAt.Before(time.Now()) {
		_ = s.repo.DeleteRefreshToken(ctx, token)
		return nil, nil, ErrRefreshTokenInvalid
	}

	user, err := s.repo.GetUserByID(ctx, token.UserID)
	if err != nil {
		if errors.Is(err, core.ErrItemNotFound) {
			return nil, nil, ErrRefreshTokenInvalid
		}
		return nil, nil, err
	}

	if requireAdmin != s.IsAdmin(user) {
		return nil, nil, ErrRefreshTokenInvalid
	}

	if err := s.repo.DeleteRefreshToken(ctx, token); err != nil {
		return nil, nil, err
	}

	var tokenData *UserTokenData
	if requireAdmin {
		tokenData, err = s.createAdminSession(ctx, user)
	} else {
		tokenData, err = s.createUserSession(ctx, user)
	}
	if err != nil {
		return nil, nil, err
	}

	return NewUserResponse(user), tokenData, nil
}

func (s *Service) CreateRefreshToken(ctx context.Context, userID uint) (*RefreshToken, error) {
	return s.createRefreshTokenWithRepo(ctx, s.repo, userID)
}

func (s *Service) createRefreshTokenWithRepo(ctx context.Context, repo *Repository, userID uint) (*RefreshToken, error) {
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

func (s *Service) CreateInvite(ctx context.Context, adminUser *User) (string, *Invite, error) {
	if !s.IsAdmin(adminUser) {
		return "", nil, ErrAdminAccessRequired
	}

	now := time.Now()
	rawToken, err := generateRandomToken()
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate invite token: %w", err)
	}

	invite := &Invite{
		TokenHash:       hashToken(rawToken),
		CreatedByUserID: adminUser.ID,
		CreatedAt:       now,
		ExpiresAt:       now.Add(inviteExpiration),
	}
	createdInvite, err := s.repo.CreateInvite(ctx, invite)
	if err != nil {
		return "", nil, err
	}

	return rawToken, createdInvite, nil
}

func (s *Service) ValidateInvite(ctx context.Context, rawToken string) (*Invite, error) {
	invite, err := s.repo.GetInviteByHash(ctx, hashToken(rawToken))
	if err != nil {
		if errors.Is(err, core.ErrItemNotFound) {
			return nil, ErrInviteInvalid
		}
		return nil, err
	}

	if invite.UsedAt != nil {
		return nil, ErrInviteAlreadyUsed
	}
	if invite.ExpiresAt.Before(time.Now()) {
		return nil, ErrInviteExpired
	}

	return invite, nil
}

func (s *Service) RegisterWithInvite(ctx context.Context, rawToken string, req InviteRegistrationRequest) (*UserResponse, *UserTokenData, error) {
	login := normalizeLogin(req.Login)
	email := normalizeEmail(req.Email)

	if err := validateLogin(login); err != nil {
		return nil, nil, err
	}
	if err := validatePassword(req.Password); err != nil {
		return nil, nil, err
	}

	passwordHash, err := hashPassword(req.Password)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to hash password: %w", err)
	}

	tokenHash := hashToken(rawToken)
	var createdUser *User
	var refreshToken *RefreshToken

	err = s.repo.WithTx(ctx, func(repo *Repository) error {
		invite, err := repo.GetInviteByHash(ctx, tokenHash)
		if err != nil {
			if errors.Is(err, core.ErrItemNotFound) {
				return ErrInviteInvalid
			}
			return err
		}
		if invite.UsedAt != nil {
			return ErrInviteAlreadyUsed
		}
		if invite.ExpiresAt.Before(time.Now()) {
			return ErrInviteExpired
		}

		if err := s.ensureUserIdentityAvailable(ctx, repo, login, email); err != nil {
			return err
		}

		createdUser, err = repo.CreateUser(ctx, &User{
			Login:        login,
			Email:        email,
			PasswordHash: passwordHash,
		})
		if err != nil {
			if isUniqueConstraintError(err) {
				if _, loginErr := repo.GetUserByLogin(ctx, login); loginErr == nil {
					return ErrLoginAlreadyTaken
				}
				if _, emailErr := repo.GetUserByEmail(ctx, email); emailErr == nil {
					return ErrEmailAlreadyTaken
				}
			}
			return err
		}

		now := time.Now()
		invite.UsedAt = &now
		invite.RegisteredUserID = &createdUser.ID
		if err := repo.SaveInvite(ctx, invite); err != nil {
			return err
		}

		refreshToken, err = s.createRefreshTokenWithRepo(ctx, repo, createdUser.ID)
		return err
	})
	if err != nil {
		return nil, nil, err
	}

	jwtToken, err := s.GenerateJWT(createdUser.ID)
	if err != nil {
		return nil, nil, err
	}

	return NewUserResponse(createdUser), NewUserTokenData(jwtToken, refreshToken), nil
}

func (s *Service) RequestPasswordReset(ctx context.Context, req ForgotPasswordRequest) (string, error) {
	email := normalizeEmail(req.Email)
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, core.ErrItemNotFound) {
			return "", nil
		}
		return "", err
	}
	if s.IsAdmin(user) {
		return "", nil
	}

	rawToken, err := generateRandomToken()
	if err != nil {
		return "", fmt.Errorf("failed to generate password reset token: %w", err)
	}

	resetToken := &PasswordResetToken{
		UserID:    user.ID,
		TokenHash: hashToken(rawToken),
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(s.config.PasswordResetTokenExpiration),
	}

	err = s.repo.WithTx(ctx, func(repo *Repository) error {
		if err := repo.DeletePasswordResetTokensByUser(ctx, user.ID); err != nil {
			return err
		}
		_, err := repo.CreatePasswordResetToken(ctx, resetToken)
		return err
	})
	if err != nil {
		return "", err
	}

	return rawToken, nil
}

func (s *Service) ResetPassword(ctx context.Context, req ResetPasswordRequest) error {
	if err := validatePassword(req.Password); err != nil {
		return err
	}

	passwordHash, err := hashPassword(req.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	tokenHash := hashToken(req.Token)

	return s.repo.WithTx(ctx, func(repo *Repository) error {
		resetToken, err := repo.GetPasswordResetTokenByHash(ctx, tokenHash)
		if err != nil {
			if errors.Is(err, core.ErrItemNotFound) {
				return ErrPasswordResetTokenInvalid
			}
			return err
		}

		if resetToken.ExpiresAt.Before(time.Now()) {
			return ErrPasswordResetTokenExpired
		}

		user, err := repo.GetUserByID(ctx, resetToken.UserID)
		if err != nil {
			if errors.Is(err, core.ErrItemNotFound) {
				return ErrPasswordResetTokenInvalid
			}
			return err
		}
		if s.IsAdmin(user) {
			return ErrPasswordResetTokenInvalid
		}

		if err := repo.UpdateUserPassword(ctx, user.ID, passwordHash); err != nil {
			return err
		}
		if err := repo.DeletePasswordResetTokensByUser(ctx, user.ID); err != nil {
			return err
		}
		if err := repo.DeleteRefreshTokensByUser(ctx, user.ID); err != nil {
			return err
		}
		return nil
	})
}

func (s *Service) createUserSession(ctx context.Context, user *User) (*UserTokenData, error) {
	if s.IsAdmin(user) {
		return nil, ErrInvalidCredentials
	}
	return s.createSession(ctx, user)
}

func (s *Service) createAdminSession(ctx context.Context, user *User) (*UserTokenData, error) {
	if !s.IsAdmin(user) {
		return nil, ErrAdminAccessRequired
	}
	return s.createSession(ctx, user)
}

func (s *Service) createSession(ctx context.Context, user *User) (*UserTokenData, error) {
	jwtToken, err := s.GenerateJWT(user.ID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrJWTGenerationFailed, err)
	}

	refreshToken, err := s.CreateRefreshToken(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRefreshTokenCreationFailed, err)
	}

	return NewUserTokenData(jwtToken, refreshToken), nil
}

func (s *Service) ensureUserIdentityAvailable(ctx context.Context, repo *Repository, login, email string) error {
	_, err := repo.GetUserByLogin(ctx, login)
	if err == nil {
		return ErrLoginAlreadyTaken
	}
	if err != nil && !errors.Is(err, core.ErrItemNotFound) {
		return err
	}

	_, err = repo.GetUserByEmail(ctx, email)
	if err == nil {
		return ErrEmailAlreadyTaken
	}
	if err != nil && !errors.Is(err, core.ErrItemNotFound) {
		return err
	}

	return nil
}

func normalizeLogin(login string) string {
	return strings.TrimSpace(strings.ToLower(login))
}

func normalizeEmail(email string) string {
	return strings.TrimSpace(strings.ToLower(email))
}

func validateLogin(login string) error {
	if !loginPattern.MatchString(login) {
		return fmt.Errorf("%w: login must be 3-32 characters and contain only lowercase letters, digits, underscores, or hyphens", core.ErrInvalidItem)
	}
	return nil
}

func validatePassword(password string) error {
	length := len(password)
	if length < minPasswordLength || length > maxPasswordLength {
		return fmt.Errorf("%w: password must be between %d and %d characters", core.ErrInvalidItem, minPasswordLength, maxPasswordLength)
	}
	return nil
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func generateRandomToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate random token: %v", err)
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func hashToken(rawToken string) string {
	sum := sha256.Sum256([]byte(rawToken))
	return hex.EncodeToString(sum[:])
}
