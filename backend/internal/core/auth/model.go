package auth

import (
	"log/slog"
	"time"

	"golang.org/x/crypto/bcrypt"
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
	ID        uint      `json:"id" gorm:"primaryKey"`
	Email     string    `json:"email" gorm:"uniqueIndex;not null" validate:"required,email"`
	Name      string    `json:"name" gorm:"not null" validate:"required"`
	Password  string    `json:"-" gorm:"not null" validate:"required,min=6"`
	CreatedAt time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null"`
}

func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func (u User) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Uint64("id", uint64(u.ID)),
		slog.String("email", u.Email),
		slog.String("name", u.Name),
		slog.String("password", "[REDACTED]"),
		slog.Time("createdAt", u.CreatedAt),
		slog.Time("updatedAt", u.UpdatedAt),
	)
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
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
}

type UpdateUserRequest struct {
	Email    *string `json:"email,omitempty" validate:"omitempty,email"`
	Name     *string `json:"name,omitempty" validate:"omitempty"`
	Password *string `json:"password,omitempty" validate:"omitempty,min=6"`
}

func (r CreateUserRequest) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("email", r.Email),
		slog.String("name", r.Name),
		slog.String("password", "[REDACTED]"),
	)
}

func (r UpdateUserRequest) LogValue() slog.Value {
	email := ""
	if r.Email != nil {
		email = *r.Email
	}
	name := ""
	if r.Name != nil {
		name = *r.Name
	}
	return slog.GroupValue(
		slog.String("email", email),
		slog.String("name", name),
		slog.String("password", "[REDACTED]"),
	)
}

type AdminLoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (r AdminLoginRequest) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("username", r.Username),
		slog.String("password", "[REDACTED]"),
	)
}
