package auth

import (
	"log/slog"
	"time"
)

type JWT struct {
	Value  string
	UserID uint
}

func (j JWT) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("value", "[REDACTED]"),
		slog.Uint64("userID", uint64(j.UserID)),
	)
}

type RefreshToken struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null;index"`
	Value     string    `gorm:"not null;uniqueIndex"`
	CreatedAt time.Time `gorm:"not null"`
	ExpiresAt time.Time `gorm:"not null;index"`
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
	JWT          *JWT
	RefreshToken *RefreshToken
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
