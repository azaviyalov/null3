package auth

import "errors"

var (
	ErrInvalidCredentials         = errors.New("invalid credentials")
	ErrAdminAccessRequired        = errors.New("admin access required")
	ErrUserNotAuthenticated       = errors.New("user not authenticated")
	ErrUserInvalidType            = errors.New("invalid user type")
	ErrJWTGenerationFailed        = errors.New("failed to generate JWT")
	ErrJWTInvalid                 = errors.New("invalid JWT")
	ErrJWTExpired                 = errors.New("JWT expired")
	ErrJWTInvalidClaims           = errors.New("invalid JWT claims")
	ErrRefreshTokenInvalid        = errors.New("invalid refresh token")
	ErrRefreshTokenExpired        = errors.New("refresh token expired")
	ErrRefreshTokenCreationFailed = errors.New("failed to create refresh token")
	ErrLoginAlreadyTaken          = errors.New("login already taken")
	ErrEmailAlreadyTaken          = errors.New("email already taken")
	ErrInviteInvalid              = errors.New("invalid invite")
	ErrInviteExpired              = errors.New("invite expired")
	ErrInviteAlreadyUsed          = errors.New("invite already used")
	ErrPasswordResetTokenInvalid  = errors.New("invalid password reset token")
	ErrPasswordResetTokenExpired  = errors.New("password reset token expired")
)
