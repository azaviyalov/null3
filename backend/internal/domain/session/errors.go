package session

import "errors"

var (
	ErrJWTGenerationFailed        = errors.New("failed to generate JWT")
	ErrJWTInvalid                 = errors.New("invalid JWT")
	ErrJWTExpired                 = errors.New("JWT expired")
	ErrJWTInvalidClaims           = errors.New("invalid JWT claims")
	ErrRefreshTokenInvalid        = errors.New("invalid refresh token")
	ErrRefreshTokenCreationFailed = errors.New("failed to create refresh token")
)
