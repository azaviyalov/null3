package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/azaviyalov/null3/backend/internal/core/logging"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetRefreshToken(tokenString string) (*RefreshToken, error) {
	var token RefreshToken
	err := r.db.Where("value = ?", tokenString).First(&token).Error
	if err != nil {
		return nil, err
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
