package auth

import "errors"

var (
	ErrInvalidCredentials         = errors.New("invalid credentials")
	ErrUserNotAuthenticated       = errors.New("user not authenticated")
	ErrUserInvalidType            = errors.New("invalid user type")
	ErrJWTGenerationFailed        = errors.New("failed to generate JWT")
	ErrJWTInvalid                 = errors.New("invalid JWT")
	ErrJWTExpired                 = errors.New("JWT expired")
	ErrJWTInvalidClaims           = errors.New("invalid JWT claims")
	ErrRefreshTokenInvalid        = errors.New("invalid refresh token")
	ErrRefreshTokenExpired        = errors.New("refresh token expired")
	ErrRefreshTokenCreationFailed = errors.New("failed to create refresh token")
)
