package session

import "errors"

var (
	ErrActorNotAuthenticated      = errors.New("actor not authenticated")
	ErrActorInvalidType           = errors.New("invalid actor type")
	ErrAdminAccessRequired        = errors.New("admin access required")
	ErrUserScopeRequired          = errors.New("user scope required")
	ErrJWTGenerationFailed        = errors.New("failed to generate JWT")
	ErrJWTInvalid                 = errors.New("invalid JWT")
	ErrJWTExpired                 = errors.New("JWT expired")
	ErrJWTInvalidClaims           = errors.New("invalid JWT claims")
	ErrRefreshTokenInvalid        = errors.New("invalid refresh token")
	ErrRefreshTokenExpired        = errors.New("refresh token expired")
	ErrRefreshTokenCreationFailed = errors.New("failed to create refresh token")
)
