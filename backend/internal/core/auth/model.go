package auth

import (
	"log/slog"
	"time"
)

type JWT struct {
	Value  string `json:"-"`
	UserID uint   `json:"-"`
}

type RefreshToken struct {
	ID        uint      `gorm:"primaryKey" json:"-"`
	UserID    uint      `gorm:"not null;index" json:"-"`
	Value     string    `gorm:"not null;uniqueIndex" json:"-"`
	CreatedAt time.Time `gorm:"not null" json:"-"`
	ExpiresAt time.Time `gorm:"not null;index" json:"-"`
}

func (r RefreshToken) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("value", "[REDACTED]"),
		slog.Uint64("userID", uint64(r.UserID)),
		slog.Time("createdAt", r.CreatedAt),
		slog.Time("expiresAt", r.ExpiresAt),
	)
}

type UserTokenData struct {
	JWT          *JWT          `json:"-"`
	RefreshToken *RefreshToken `json:"-"`
}

func NewUserTokenData(jwt *JWT, refreshToken *RefreshToken) *UserTokenData {
	return &UserTokenData{
		JWT:          jwt,
		RefreshToken: refreshToken,
	}
}

type User struct {
	ID uint `json:"id"`
}

type LoginRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (r LoginRequest) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("login", r.Login),
		slog.String("password", "[REDACTED]"),
	)
}

type UserResponse struct {
	ID uint `json:"id"`
}
