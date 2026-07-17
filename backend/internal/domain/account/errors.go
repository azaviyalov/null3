package account

import "errors"

var (
	ErrInvalidCredentials        = errors.New("invalid credentials")
	ErrLoginAlreadyTaken         = errors.New("login already taken")
	ErrEmailAlreadyTaken         = errors.New("email already taken")
	ErrInviteInvalid             = errors.New("invalid invite")
	ErrInviteExpired             = errors.New("invite expired")
	ErrInviteAlreadyUsed         = errors.New("invite already used")
	ErrPasswordResetTokenInvalid = errors.New("invalid password reset token")
	ErrPasswordResetTokenExpired = errors.New("password reset token expired")
)
