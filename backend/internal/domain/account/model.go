package account

import "time"

type User struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Login        string    `json:"login" gorm:"not null;uniqueIndex"`
	Email        string    `json:"email" gorm:"not null;uniqueIndex"`
	PasswordHash string    `json:"-" gorm:"not null"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type PasswordResetToken struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null;index"`
	TokenHash string    `gorm:"not null;uniqueIndex"`
	CreatedAt time.Time `gorm:"not null"`
	ExpiresAt time.Time `gorm:"not null;index"`
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

type LoginRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
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
