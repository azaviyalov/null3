package auth

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

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

func (r *Repository) SaveRefreshToken(token *RefreshToken) (*RefreshToken, error) {
	slog.Debug("SaveRefreshToken called", "token", token)
	if err := r.db.Save(token).Error; err != nil {
		slog.Error("db error in SaveRefreshToken", "error", err)
		return nil, err
	}
	return token, nil
}

func (r *Repository) DeleteRefreshToken(token *RefreshToken) error {
	slog.Debug("DeleteRefreshToken called", "token", token)

	if err := r.db.Delete(token).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			slog.Warn("DeleteRefreshToken: token not found", "token", token)
			return nil // Token not found, nothing to delete
		}
		slog.Error("db error in DeleteRefreshToken", "error", err)
		return fmt.Errorf("db error: %w", err)
	}

	slog.Info("token deleted in DeleteRefreshToken", "token", token)
	return nil
}

func (r *Repository) DeleteExpiredRefreshTokens() error {
	slog.Debug("DeleteExpiredRefreshTokens called")

	now := time.Now()
	if err := r.db.Where("expires_at < ?", now).Delete(&RefreshToken{}).Error; err != nil {
		slog.Error("db error in DeleteExpiredRefreshTokens", "error", err)
		return fmt.Errorf("db error: %w", err)
	}
	slog.Info("expired refresh tokens deleted in DeleteExpiredRefreshTokens")
	return nil
}
