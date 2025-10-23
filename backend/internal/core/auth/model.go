package auth

import (
	"time"

	"github.com/azaviyalov/null3/backend/internal/core/logging"
)

type JWT struct {
	Value  string
	UserID uint
}

func (j JWT) ToFieldValue() logging.FieldValue {
	return logging.CombineFields(
		logging.NewField("value", logging.NewStringValue("[REDACTED]")),
		logging.NewField("user_id", logging.NewUint64Value(uint64(j.UserID))),
	)
}

type RefreshToken struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null;index"`
	Value     string    `gorm:"not null;uniqueIndex"`
	CreatedAt time.Time `gorm:"not null"`
	ExpiresAt time.Time `gorm:"not null;index"`
}

func (r RefreshToken) ToFieldValue() logging.FieldValue {
	return logging.CombineFields(
		logging.NewField("value", logging.NewStringValue("[REDACTED]")),
		logging.NewField("user_id", logging.NewUint64Value(uint64(r.UserID))),
		logging.NewField("created_at", logging.NewTimeValue(r.CreatedAt)),
		logging.NewField("expires_at", logging.NewTimeValue(r.ExpiresAt)),
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

func (r LoginRequest) ToFieldValue() logging.FieldValue {
	return logging.CombineFields(
		logging.NewField("login", logging.NewStringValue(r.Login)),
		logging.NewField("password", logging.NewStringValue("[REDACTED]")),
	)
}

type UserResponse struct {
	ID uint `json:"id"`
}
