package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type sessionCookies struct {
	jwtCookieName     string
	refreshCookieName string
}

var (
	userSessionCookies = sessionCookies{
		jwtCookieName:     userCookieName,
		refreshCookieName: userRefreshCookieName,
	}
	adminSessionCookies = sessionCookies{
		jwtCookieName:     adminCookieName,
		refreshCookieName: adminRefreshCookieName,
	}
)

func setSessionCookies(c echo.Context, config Config, scope sessionCookies, tokenData *UserTokenData) {
	c.SetCookie(&http.Cookie{
		Name:     scope.jwtCookieName,
		Value:    tokenData.JWT.Value,
		HttpOnly: true,
		Secure:   config.SecureCookies,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})
	c.SetCookie(&http.Cookie{
		Name:     scope.refreshCookieName,
		Value:    tokenData.RefreshToken.Value,
		HttpOnly: true,
		Secure:   config.SecureCookies,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(config.RefreshTokenExpiration.Seconds()),
	})
}

func clearSessionCookies(c echo.Context, config Config, scope sessionCookies) {
	c.SetCookie(&http.Cookie{
		Name:     scope.jwtCookieName,
		Value:    "",
		HttpOnly: true,
		Secure:   config.SecureCookies,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
	c.SetCookie(&http.Cookie{
		Name:     scope.refreshCookieName,
		Value:    "",
		HttpOnly: true,
		Secure:   config.SecureCookies,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
}
