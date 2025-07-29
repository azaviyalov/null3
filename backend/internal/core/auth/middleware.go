package auth

import (
	"errors"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type User struct {
	ID    uint   `json:"id"`
	Token string `json:"token"`
}

func JWTMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			header := c.Request().Header.Get("Authorization")
			if header == "" || !strings.HasPrefix(header, "Bearer ") {
				err := errors.New("missing or invalid Authorization header")
				return echo.ErrUnauthorized.WithInternal(err)
			}
			tokenStr := strings.TrimPrefix(header, "Bearer ")
			claims, err := ParseJWT(tokenStr)
			if err != nil {
				return echo.ErrUnauthorized.WithInternal(err)
			}
			userID, err := strconv.ParseUint(claims.Subject, 10, 64)
			if err != nil {
				return echo.ErrUnauthorized.WithInternal(errors.New("invalid user id in token"))
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
		return nil, ErrUserNotAuthenticated
	}
	u, ok := user.(*User)
	if !ok {
		return nil, errors.New("user context is not of type *User")
	}
	return u, nil
}
