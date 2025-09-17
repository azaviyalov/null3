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

// User operations
func (r *Repository) CreateUser(user *User) (*User, error) {
	slog.Debug("CreateUser called", "user", user)
	if err := r.db.Create(user).Error; err != nil {
		slog.Error("db error in CreateUser", "error", err)
		return nil, err
	}
	return user, nil
}

func (r *Repository) GetUserByEmail(email string) (*User, error) {
	var user User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetUserByID(id uint) (*User, error) {
	var user User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) UpdateUser(user *User) (*User, error) {
	slog.Debug("UpdateUser called", "user", user)
	if err := r.db.Save(user).Error; err != nil {
		slog.Error("db error in UpdateUser", "error", err)
		return nil, err
	}
	return user, nil
}

func (r *Repository) DeleteUser(id uint) error {
	slog.Debug("DeleteUser called", "id", id)
	if err := r.db.Delete(&User{}, id).Error; err != nil {
		slog.Error("db error in DeleteUser", "error", err)
		return err
	}
	return nil
}

func (r *Repository) GetAllUsers() ([]*User, error) {
	var users []*User
	err := r.db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *Repository) GetRefreshTokensByUserID(userID uint) ([]*RefreshToken, error) {
	var tokens []*RefreshToken
	err := r.db.Where("user_id = ?", userID).Find(&tokens).Error
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

func (r *Repository) GetAllRefreshTokens() ([]*RefreshToken, error) {
	var tokens []*RefreshToken
	err := r.db.Find(&tokens).Error
	if err != nil {
		return nil, err
	}
	return tokens, nil
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
