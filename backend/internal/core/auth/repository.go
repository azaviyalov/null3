package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/azaviyalov/null3/backend/internal/core"
	"github.com/azaviyalov/null3/backend/internal/core/logging"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) WithTx(ctx context.Context, fn func(repo *Repository) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(&Repository{db: tx})
	})
}

func (r *Repository) GetUserByID(ctx context.Context, id uint) (*User, error) {
	var user User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, core.ErrItemNotFound
		}
		return nil, fmt.Errorf("db error: %w", err)
	}
	return &user, nil
}

func (r *Repository) GetUserByLogin(ctx context.Context, login string) (*User, error) {
	var user User
	if err := r.db.WithContext(ctx).Where("login = ?", login).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, core.ErrItemNotFound
		}
		return nil, fmt.Errorf("db error: %w", err)
	}
	return &user, nil
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, core.ErrItemNotFound
		}
		return nil, fmt.Errorf("db error: %w", err)
	}
	return &user, nil
}

func (r *Repository) CreateUser(ctx context.Context, user *User) (*User, error) {
	logging.Debug(ctx, "CreateUser called", "login", user.Login, "email", user.Email)
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		logging.Error(ctx, "db error in CreateUser", "error", err, "login", user.Login, "email", user.Email)
		return nil, fmt.Errorf("db error: %w", err)
	}
	return user, nil
}

func (r *Repository) UpdateUserPassword(ctx context.Context, userID uint, passwordHash string) error {
	logging.Debug(ctx, "UpdateUserPassword called", "user_id", userID)
	if err := r.db.WithContext(ctx).Model(&User{}).Where("id = ?", userID).Update("password_hash", passwordHash).Error; err != nil {
		logging.Error(ctx, "db error in UpdateUserPassword", "error", err, "user_id", userID)
		return fmt.Errorf("db error: %w", err)
	}
	return nil
}

func (r *Repository) GetRefreshToken(ctx context.Context, tokenString string) (*RefreshToken, error) {
	var token RefreshToken
	err := r.db.WithContext(ctx).Where("value = ?", tokenString).First(&token).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, core.ErrItemNotFound
		}
		return nil, fmt.Errorf("db error: %w", err)
	}
	return &token, nil
}

func (r *Repository) SaveRefreshToken(ctx context.Context, token *RefreshToken) (*RefreshToken, error) {
	logging.Debug(ctx, "SaveRefreshToken called", "user_id", token.UserID, "expires_at", token.ExpiresAt)
	if err := r.db.WithContext(ctx).Save(token).Error; err != nil {
		logging.Error(ctx, "db error in SaveRefreshToken", "error", err, "user_id", token.UserID)
		return nil, fmt.Errorf("db error: %w", err)
	}
	return token, nil
}

func (r *Repository) DeleteRefreshToken(ctx context.Context, token *RefreshToken) error {
	logging.Debug(ctx, "DeleteRefreshToken called", "user_id", token.UserID)

	if err := r.db.WithContext(ctx).Delete(token).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logging.Info(ctx, "DeleteRefreshToken: token not found, nothing to delete", "user_id", token.UserID)
			return nil
		}
		logging.Error(ctx, "db error in DeleteRefreshToken", "error", err, "user_id", token.UserID)
		return fmt.Errorf("db error: %w", err)
	}

	logging.Info(ctx, "token deleted in DeleteRefreshToken", "user_id", token.UserID)
	return nil
}

func (r *Repository) DeleteRefreshTokensByUser(ctx context.Context, userID uint) error {
	logging.Debug(ctx, "DeleteRefreshTokensByUser called", "user_id", userID)
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&RefreshToken{}).Error; err != nil {
		logging.Error(ctx, "db error in DeleteRefreshTokensByUser", "error", err, "user_id", userID)
		return fmt.Errorf("db error: %w", err)
	}
	return nil
}

func (r *Repository) DeleteExpiredRefreshTokens(ctx context.Context) error {
	logging.Debug(ctx, "DeleteExpiredRefreshTokens called")

	now := time.Now()
	if err := r.db.WithContext(ctx).Where("expires_at < ?", now).Delete(&RefreshToken{}).Error; err != nil {
		logging.Error(ctx, "db error in DeleteExpiredRefreshTokens", "error", err)
		return fmt.Errorf("db error: %w", err)
	}
	logging.Info(ctx, "expired refresh tokens deleted in DeleteExpiredRefreshTokens")
	return nil
}

func (r *Repository) CreatePasswordResetToken(ctx context.Context, token *PasswordResetToken) (*PasswordResetToken, error) {
	logging.Debug(ctx, "CreatePasswordResetToken called", "user_id", token.UserID, "expires_at", token.ExpiresAt)
	if err := r.db.WithContext(ctx).Create(token).Error; err != nil {
		logging.Error(ctx, "db error in CreatePasswordResetToken", "error", err, "user_id", token.UserID)
		return nil, fmt.Errorf("db error: %w", err)
	}
	return token, nil
}

func (r *Repository) GetPasswordResetTokenByHash(ctx context.Context, tokenHash string) (*PasswordResetToken, error) {
	var token PasswordResetToken
	if err := r.db.WithContext(ctx).Where("token_hash = ?", tokenHash).First(&token).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, core.ErrItemNotFound
		}
		return nil, fmt.Errorf("db error: %w", err)
	}
	return &token, nil
}

func (r *Repository) DeletePasswordResetToken(ctx context.Context, token *PasswordResetToken) error {
	if err := r.db.WithContext(ctx).Delete(token).Error; err != nil {
		logging.Error(ctx, "db error in DeletePasswordResetToken", "error", err, "user_id", token.UserID)
		return fmt.Errorf("db error: %w", err)
	}
	return nil
}

func (r *Repository) DeletePasswordResetTokensByUser(ctx context.Context, userID uint) error {
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&PasswordResetToken{}).Error; err != nil {
		logging.Error(ctx, "db error in DeletePasswordResetTokensByUser", "error", err, "user_id", userID)
		return fmt.Errorf("db error: %w", err)
	}
	return nil
}

func (r *Repository) CreateInvite(ctx context.Context, invite *Invite) (*Invite, error) {
	logging.Debug(ctx, "CreateInvite called", "created_by_user_id", invite.CreatedByUserID, "expires_at", invite.ExpiresAt)
	if err := r.db.WithContext(ctx).Create(invite).Error; err != nil {
		logging.Error(ctx, "db error in CreateInvite", "error", err, "created_by_user_id", invite.CreatedByUserID)
		return nil, fmt.Errorf("db error: %w", err)
	}
	return invite, nil
}

func (r *Repository) GetInviteByHash(ctx context.Context, tokenHash string) (*Invite, error) {
	var invite Invite
	if err := r.db.WithContext(ctx).Where("token_hash = ?", tokenHash).First(&invite).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, core.ErrItemNotFound
		}
		return nil, fmt.Errorf("db error: %w", err)
	}
	return &invite, nil
}

func (r *Repository) SaveInvite(ctx context.Context, invite *Invite) error {
	if err := r.db.WithContext(ctx).Save(invite).Error; err != nil {
		logging.Error(ctx, "db error in SaveInvite", "error", err, "invite_id", invite.ID)
		return fmt.Errorf("db error: %w", err)
	}
	return nil
}

func isUniqueConstraintError(err error) bool {
	return strings.Contains(strings.ToLower(err.Error()), "unique constraint failed")
}
