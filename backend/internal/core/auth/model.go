package auth

import (
	"time"

	"github.com/azaviyalov/null3/backend/internal/core/logging"
)

type User struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Login        string    `json:"login" gorm:"not null;uniqueIndex"`
	Email        string    `json:"email" gorm:"not null;uniqueIndex"`
	PasswordHash string    `json:"-" gorm:"not null"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (u User) ToFieldValue() logging.FieldValue {
	return logging.CombineFields(
		logging.NewField("id", logging.NewUint64Value(uint64(u.ID))),
		logging.NewField("login", logging.NewStringValue(u.Login)),
		logging.NewField("email", logging.NewStringValue(u.Email)),
	)
}

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

type PasswordResetToken struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null;index"`
	TokenHash string    `gorm:"not null;uniqueIndex"`
	CreatedAt time.Time `gorm:"not null"`
	ExpiresAt time.Time `gorm:"not null;index"`
}

func (t PasswordResetToken) ToFieldValue() logging.FieldValue {
	return logging.CombineFields(
		logging.NewField("user_id", logging.NewUint64Value(uint64(t.UserID))),
		logging.NewField("created_at", logging.NewTimeValue(t.CreatedAt)),
		logging.NewField("expires_at", logging.NewTimeValue(t.ExpiresAt)),
	)
}

type Invite struct {
	ID               uint       `gorm:"primaryKey"`
	TokenHash        string     `gorm:"not null;uniqueIndex"`
	CreatedByUserID  uint       `gorm:"not null;index"`
	CreatedAt        time.Time  `gorm:"not null"`
	ExpiresAt        time.Time  `gorm:"not null;index"`
	UsedAt           *time.Time `gorm:"index"`
	RegisteredUserID *uint      `gorm:"index"`
}

func (i Invite) ToFieldValue() logging.FieldValue {
	fields := []logging.Field{
		logging.NewField("id", logging.NewUint64Value(uint64(i.ID))),
		logging.NewField("created_by_user_id", logging.NewUint64Value(uint64(i.CreatedByUserID))),
		logging.NewField("created_at", logging.NewTimeValue(i.CreatedAt)),
		logging.NewField("expires_at", logging.NewTimeValue(i.ExpiresAt)),
	}
	if i.RegisteredUserID != nil {
		fields = append(fields, logging.NewField("registered_user_id", logging.NewUint64Value(uint64(*i.RegisteredUserID))))
	}
	if i.UsedAt != nil {
		fields = append(fields, logging.NewField("used_at", logging.NewTimeValue(*i.UsedAt)))
	}
	return logging.CombineFields(fields...)
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

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordRequest struct {
	Token    string `json:"token" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type InviteRegistrationRequest struct {
	Login    string `json:"login" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserResponse struct {
	ID    uint   `json:"id"`
	Login string `json:"login"`
	Email string `json:"email"`
}

func NewUserResponse(user *User) *UserResponse {
	return &UserResponse{
		ID:    user.ID,
		Login: user.Login,
		Email: user.Email,
	}
}

type MessageResponse struct {
	Message string `json:"message"`
}

type ForgotPasswordResponse struct {
	Message  string `json:"message"`
	ResetURL string `json:"reset_url,omitempty"`
}

type InviteValidationResponse struct {
	ExpiresAt time.Time `json:"expires_at"`
}

type InviteResponse struct {
	InviteURL string    `json:"invite_url"`
	ExpiresAt time.Time `json:"expires_at"`
}
