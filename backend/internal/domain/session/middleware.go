package session

import (
	"context"

	"github.com/labstack/echo/v4"
)

type ActorResolver func(ctx context.Context, userID uint) (*Actor, error)

func UserJWTMiddleware(service *Service, resolveActor ActorResolver) echo.MiddlewareFunc {
	return jwtMiddleware(service, resolveActor, UserCookieName, false)
}

func AdminJWTMiddleware(service *Service, resolveActor ActorResolver) echo.MiddlewareFunc {
	return jwtMiddleware(service, resolveActor, AdminCookieName, true)
}

func jwtMiddleware(service *Service, resolveActor ActorResolver, cookieName string, requireAdmin bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie(cookieName)
			if err != nil {
				return echo.ErrUnauthorized.WithInternal(err)
			}

			token, err := service.ParseJWT(cookie.Value)
			if err != nil {
				return echo.ErrUnauthorized.WithInternal(err)
			}

			actor, err := resolveActor(c.Request().Context(), token.UserID)
			if err != nil {
				return echo.ErrUnauthorized.WithInternal(err)
			}

			if requireAdmin && !actor.IsAdmin {
				return echo.ErrForbidden.WithInternal(ErrAdminAccessRequired)
			}
			if !requireAdmin && actor.IsAdmin {
				return echo.ErrUnauthorized.WithInternal(ErrUserScopeRequired)
			}

			setActor(c, actor)
			return next(c)
		}
	}
}
