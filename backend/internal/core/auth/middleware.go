package auth

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

const (
	userCookieName         = "jwt"
	userRefreshCookieName  = "refresh_token"
	adminCookieName        = "admin_jwt"
	adminRefreshCookieName = "admin_refresh_token"
)

func UserJWTMiddleware(service *Service) echo.MiddlewareFunc {
	return jwtMiddleware(service, userCookieName, false)
}

func AdminJWTMiddleware(service *Service) echo.MiddlewareFunc {
	return jwtMiddleware(service, adminCookieName, true)
}

func jwtMiddleware(service *Service, cookieName string, requireAdmin bool) echo.MiddlewareFunc {
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
			user, err := service.GetUserByID(c.Request().Context(), jwt.UserID)
			if err != nil {
				return echo.ErrUnauthorized.WithInternal(err)
			}
			if requireAdmin && !service.IsAdmin(user) {
				return echo.ErrForbidden.WithInternal(ErrAdminAccessRequired)
			}
			if !requireAdmin && service.IsAdmin(user) {
				return echo.ErrUnauthorized.WithInternal(ErrInvalidCredentials)
			}
			setUser(c, user)
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
