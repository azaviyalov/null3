package session

import (
	"context"

	"github.com/labstack/echo/v4"
)

func UserJWTMiddleware(service *Service, validateUser func(context.Context, uint) error) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie(UserCookieName)
			if err != nil {
				return echo.ErrUnauthorized.WithInternal(err)
			}
			userID, err := service.ParseUserAccessToken(cookie.Value)
			if err != nil {
				return echo.ErrUnauthorized.WithInternal(err)
			}
			if err := validateUser(c.Request().Context(), userID); err != nil {
				return echo.ErrUnauthorized.WithInternal(err)
			}
			setUserID(c, userID)
			return next(c)
		}
	}
}

func AdminJWTMiddleware(service *Service) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie(AdminCookieName)
			if err != nil {
				return echo.ErrUnauthorized.WithInternal(err)
			}

			if err := service.ValidateAdminAccessToken(cookie.Value); err != nil {
				return echo.ErrUnauthorized.WithInternal(err)
			}

			return next(c)
		}
	}
}
