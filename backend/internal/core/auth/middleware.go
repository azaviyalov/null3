package auth

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

const (
	cookieName = "jwt"
)

func JWTMiddleware(config Config, service *Service) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie(cookieName)
			if err != nil {
				return echo.ErrUnauthorized.WithInternal(err)
			}
			tokenStr := cookie.Value
			jwt, err := service.ParseJWT(tokenStr)
			if err != nil {
				return echo.ErrUnauthorized.WithInternal(err)
			}
			setUser(c, &User{ID: jwt.UserID})
			return next(c)
		}
	}
}

func GetUser(c echo.Context) (*User, error) {
	user := c.Get(echoUserKey)
	if user == nil {
		return nil, fmt.Errorf("user is nil: %w", ErrUserNotAuthenticated)
	}
	u, ok := user.(*User)
	if !ok {
		return nil, fmt.Errorf("invalid user type: expected *User, got %T: %w", user, ErrUserInvalidType)
	}
	return u, nil
}

func setUser(c echo.Context, user *User) {
	c.Set(echoUserKey, user)
}
