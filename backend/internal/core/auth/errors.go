package auth

import "errors"

var (
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrTokenGenerationFailed = errors.New("could not generate token")
	ErrTokenInvalid          = errors.New("invalid token")
	ErrTokenInvalidClaims    = errors.New("invalid token claims")
	ErrUserNotAuthenticated  = errors.New("user not authenticated")
	ErrUserInvalidType       = errors.New("invalid user type")
)
