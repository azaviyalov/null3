package account

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/azaviyalov/null3/backend/internal/core"
	"github.com/azaviyalov/null3/backend/internal/domain/session"
	"golang.org/x/crypto/bcrypt"
)

const (
	inviteExpiration  = 24 * time.Hour
	minPasswordLength = 8
	maxPasswordLength = 72
)

var loginPattern = regexp.MustCompile(`^[a-z0-9_-]{3,32}$`)

type Service struct {
	repo           *Repository
	sessionService *session.Service
	config         Config
}

func NewService(repo *Repository, sessionService *session.Service, config Config) *Service {
	return &Service{
		repo:           repo,
		sessionService: sessionService,
		config:         config,
	}
}

func (s *Service) AuthenticateUser(ctx context.Context, req LoginRequest) (*UserResponse, *session.UserSessionTokens, error) {
	user, err := s.authenticateByLogin(ctx, req)
	if err != nil {
		return nil, nil, err
	}
	tokenData, err := s.createUserSession(ctx, user)
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

func (s *Service) RefreshUserSession(ctx context.Context, tokenString string) (*UserResponse, *session.UserSessionTokens, error) {
	sessionRepo := s.repo.SessionRepository()

	token, err := sessionRepo.GetRefreshToken(ctx, tokenString)
	if err != nil {
		if errors.Is(err, core.ErrItemNotFound) {
			return nil, nil, session.ErrRefreshTokenInvalid
		}
		return nil, nil, err
	}

	if token.ExpiresAt.Before(time.Now()) {
		_ = sessionRepo.DeleteRefreshToken(ctx, token)
		return nil, nil, session.ErrRefreshTokenInvalid
	}

	user, err := s.repo.GetUserByID(ctx, token.UserID)
	if err != nil {
		if errors.Is(err, core.ErrItemNotFound) {
			return nil, nil, session.ErrRefreshTokenInvalid
		}
		return nil, nil, err
	}

	if err := sessionRepo.DeleteRefreshToken(ctx, token); err != nil {
		return nil, nil, err
	}

	tokenData, err := s.createUserSession(ctx, user)
	if err != nil {
		return nil, nil, err
	}

	return NewUserResponse(user), tokenData, nil
}

func (s *Service) CreateInvite(ctx context.Context) (string, *Invite, error) {
	now := time.Now()
	rawToken, err := generateRandomToken()
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate invite token: %w", err)
	}

	invite := &Invite{
		TokenHash: hashToken(rawToken),
		CreatedAt: now,
		ExpiresAt: now.Add(inviteExpiration),
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

func (s *Service) RegisterWithInvite(ctx context.Context, rawToken string, req InviteRegistrationRequest) (*UserResponse, *session.UserSessionTokens, error) {
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
	var refreshToken *session.RefreshToken

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

		refreshToken, err = s.sessionService.CreateRefreshTokenWithRepo(ctx, repo.SessionRepository(), createdUser.ID)
		return err
	})
	if err != nil {
		return nil, nil, err
	}

	accessToken, err := s.sessionService.GenerateUserAccessToken(createdUser.ID)
	if err != nil {
		return nil, nil, err
	}

	return NewUserResponse(createdUser), &session.UserSessionTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
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
		if err := repo.UpdateUserPassword(ctx, user.ID, passwordHash); err != nil {
			return err
		}
		if err := repo.DeletePasswordResetTokensByUser(ctx, user.ID); err != nil {
			return err
		}
		if err := repo.SessionRepository().DeleteRefreshTokensByUser(ctx, user.ID); err != nil {
			return err
		}
		return nil
	})
}

func (s *Service) createUserSession(ctx context.Context, user *User) (*session.UserSessionTokens, error) {
	accessToken, err := s.sessionService.GenerateUserAccessToken(user.ID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.sessionService.CreateRefreshToken(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return &session.UserSessionTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) ensureUserIdentityAvailable(ctx context.Context, repo *Repository, login, email string) error {
	_, err := repo.GetUserByLogin(ctx, login)
	if err == nil {
		return ErrLoginAlreadyTaken
	}
	if !errors.Is(err, core.ErrItemNotFound) {
		return err
	}

	_, err = repo.GetUserByEmail(ctx, email)
	if err == nil {
		return ErrEmailAlreadyTaken
	}
	if !errors.Is(err, core.ErrItemNotFound) {
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
		return "", fmt.Errorf("generate random token: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func hashToken(rawToken string) string {
	sum := sha256.Sum256([]byte(rawToken))
	return hex.EncodeToString(sum[:])
}
