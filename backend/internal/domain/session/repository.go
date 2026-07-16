package session

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/azaviyalov/null3/backend/internal/core"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetRefreshToken(ctx context.Context, tokenString string) (*RefreshToken, error) {
	var token RefreshToken
	err := r.db.WithContext(ctx).Where("value = ?", tokenString).First(&token).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, core.ErrItemNotFound
		}
		return nil, fmt.Errorf("get refresh token: %w", err)
	}
	return &token, nil
}

func (r *Repository) SaveRefreshToken(ctx context.Context, token *RefreshToken) (*RefreshToken, error) {
	if err := r.db.WithContext(ctx).Save(token).Error; err != nil {
		return nil, fmt.Errorf("save refresh token: %w", err)
	}
	return token, nil
}

func (r *Repository) DeleteRefreshToken(ctx context.Context, token *RefreshToken) error {
	if err := r.db.WithContext(ctx).Delete(token).Error; err != nil {
		return fmt.Errorf("delete refresh token: %w", err)
	}
	return nil
}

func (r *Repository) DeleteRefreshTokensByUser(ctx context.Context, userID uint) error {
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&RefreshToken{}).Error; err != nil {
		return fmt.Errorf("delete refresh tokens for user %d: %w", userID, err)
	}
	return nil
}

func (r *Repository) DeleteExpiredRefreshTokens(ctx context.Context) error {
	now := time.Now()
	if err := r.db.WithContext(ctx).Where("expires_at < ?", now).Delete(&RefreshToken{}).Error; err != nil {
		return fmt.Errorf("delete expired refresh tokens: %w", err)
	}
	return nil
}
