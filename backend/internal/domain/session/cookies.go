package session

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

const (
	UserCookieName        = "jwt"
	UserRefreshCookieName = "refresh_token"
	AdminCookieName       = "admin_jwt"
)

func SetUserSessionCookies(c echo.Context, config Config, tokens *UserSessionTokens) {
	c.SetCookie(&http.Cookie{
		Name:     UserCookieName,
		Value:    tokens.AccessToken,
		HttpOnly: true,
		Secure:   config.SecureCookies,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})
	c.SetCookie(&http.Cookie{
		Name:     UserRefreshCookieName,
		Value:    tokens.RefreshToken.Value,
		HttpOnly: true,
		Secure:   config.SecureCookies,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(config.RefreshTokenExpiration.Seconds()),
	})
}

func SetAdminCookie(c echo.Context, config Config, token string, expiration time.Duration) {
	c.SetCookie(&http.Cookie{Name: AdminCookieName, Value: token, HttpOnly: true, Secure: config.SecureCookies, Path: "/api/admin", SameSite: http.SameSiteLaxMode, MaxAge: int(expiration.Seconds())})
}

func ClearUserSessionCookies(c echo.Context, config Config) {
	c.SetCookie(&http.Cookie{Name: UserCookieName, Value: "", HttpOnly: true, Secure: config.SecureCookies, Path: "/", SameSite: http.SameSiteLaxMode, MaxAge: -1})
	c.SetCookie(&http.Cookie{Name: UserRefreshCookieName, Value: "", HttpOnly: true, Secure: config.SecureCookies, Path: "/", SameSite: http.SameSiteLaxMode, MaxAge: -1})
}

func ClearAdminCookie(c echo.Context, config Config) {
	c.SetCookie(&http.Cookie{Name: AdminCookieName, Value: "", HttpOnly: true, Secure: config.SecureCookies, Path: "/api/admin", SameSite: http.SameSiteLaxMode, MaxAge: -1})
}
