package session

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

const (
	UserCookieName         = "jwt"
	UserRefreshCookieName  = "refresh_token"
	AdminCookieName        = "admin_jwt"
	AdminRefreshCookieName = "admin_refresh_token"
)

type cookieScope struct {
	jwtCookieName     string
	refreshCookieName string
}

var (
	userCookieScope = cookieScope{
		jwtCookieName:     UserCookieName,
		refreshCookieName: UserRefreshCookieName,
	}
	adminCookieScope = cookieScope{
		jwtCookieName:     AdminCookieName,
		refreshCookieName: AdminRefreshCookieName,
	}
)

func SetUserSessionCookies(c echo.Context, config Config, tokenData *TokenData) {
	setSessionCookies(c, config, userCookieScope, tokenData)
}

func SetAdminSessionCookies(c echo.Context, config Config, tokenData *TokenData) {
	setSessionCookies(c, config, adminCookieScope, tokenData)
}

func ClearUserSessionCookies(c echo.Context, config Config) {
	clearSessionCookies(c, config, userCookieScope)
}

func ClearAdminSessionCookies(c echo.Context, config Config) {
	clearSessionCookies(c, config, adminCookieScope)
}

func setSessionCookies(c echo.Context, config Config, scope cookieScope, tokenData *TokenData) {
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

func clearSessionCookies(c echo.Context, config Config, scope cookieScope) {
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
