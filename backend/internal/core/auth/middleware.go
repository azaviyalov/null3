package auth

import (
	"fmt"
	"strconv"

	"github.com/labstack/echo/v4"
)

type User struct {
	ID    uint   `json:"id"`
	Token string `json:"token"`
}

func JWTMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie("jwt_token")
			if err != nil {
				return echo.ErrUnauthorized.WithInternal(err)
			}
			tokenStr := cookie.Value
			claims, err := ParseJWT(tokenStr)
			if err != nil {
				return echo.ErrUnauthorized.WithInternal(err)
			}
			userID, err := strconv.ParseUint(claims.Subject, 10, 64)
			if err != nil {
				return echo.ErrUnauthorized.WithInternal(err)
			}
			c.Set("user", &User{
				ID:    uint(userID),
				Token: tokenStr,
			})
			return next(c)
		}
	}
}

func GetUser(c echo.Context) (*User, error) {
	user := c.Get("user")
	if user == nil {
		return nil, fmt.Errorf("user not nil: %w", ErrUserNotAuthenticated)
	}
	u, ok := user.(*User)
	if !ok {
		return nil, fmt.Errorf("invalid user type: expected *User, got %T: %w", user, ErrUserInvalidType)
	}
	return u, nil
}
