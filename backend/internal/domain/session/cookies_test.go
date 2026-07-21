package session_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/azaviyalov/null3/backend/internal/domain/session"
	"github.com/labstack/echo/v4"
)

func TestUserSessionCookies(t *testing.T) {
	config := session.Config{SecureCookies: true, RefreshTokenExpiration: 48 * time.Hour}
	tokens := &session.UserSessionTokens{
		AccessToken:  "access-value",
		RefreshToken: &session.RefreshToken{Value: "refresh-value"},
	}
	cookies := recordCookies(t, func(c echo.Context) {
		session.SetUserSessionCookies(c, config, tokens)
	})

	if len(cookies) != 2 {
		t.Fatalf("cookie count = %d, want 2", len(cookies))
	}
	byName := cookiesByName(cookies)
	assertCookie(t, byName[session.UserCookieName], http.Cookie{
		Name:     session.UserCookieName,
		Value:    "access-value",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
	assertCookie(t, byName[session.UserRefreshCookieName], http.Cookie{
		Name:     session.UserRefreshCookieName,
		Value:    "refresh-value",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int((48 * time.Hour).Seconds()),
	})

	cleared := recordCookies(t, func(c echo.Context) {
		session.ClearUserSessionCookies(c, config)
	})
	if len(cleared) != 2 {
		t.Fatalf("cleared cookie count = %d, want 2", len(cleared))
	}
	for _, cookie := range cleared {
		if cookie.Value != "" || cookie.MaxAge != -1 || cookie.Path != "/" || !cookie.HttpOnly || !cookie.Secure || cookie.SameSite != http.SameSiteLaxMode {
			t.Errorf("cleared cookie %q has unexpected attributes", cookie.Name)
		}
	}
}

func TestAdminCookie(t *testing.T) {
	config := session.Config{SecureCookies: true}
	cookies := recordCookies(t, func(c echo.Context) {
		session.SetAdminCookie(c, config, "admin-value", 30*time.Minute)
	})
	if len(cookies) != 1 {
		t.Fatalf("cookie count = %d, want 1", len(cookies))
	}
	assertCookie(t, cookies[0], http.Cookie{
		Name:     session.AdminCookieName,
		Value:    "admin-value",
		Path:     "/api/admin",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int((30 * time.Minute).Seconds()),
	})

	cleared := recordCookies(t, func(c echo.Context) {
		session.ClearAdminCookie(c, config)
	})
	if len(cleared) != 1 {
		t.Fatalf("cleared cookie count = %d, want 1", len(cleared))
	}
	if cleared[0].Name != session.AdminCookieName || cleared[0].Value != "" || cleared[0].MaxAge != -1 || cleared[0].Path != "/api/admin" || !cleared[0].HttpOnly || !cleared[0].Secure || cleared[0].SameSite != http.SameSiteLaxMode {
		t.Error("ClearAdminCookie() returned unexpected attributes")
	}
}

func recordCookies(t *testing.T, setCookies func(echo.Context)) []*http.Cookie {
	t.Helper()
	e := echo.New()
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()
	setCookies(e.NewContext(request, recorder))
	return recorder.Result().Cookies()
}

func cookiesByName(cookies []*http.Cookie) map[string]*http.Cookie {
	byName := make(map[string]*http.Cookie, len(cookies))
	for _, cookie := range cookies {
		byName[cookie.Name] = cookie
	}
	return byName
}

func assertCookie(t *testing.T, cookie *http.Cookie, want http.Cookie) {
	t.Helper()
	if cookie == nil {
		t.Fatalf("cookie %q is missing", want.Name)
	}
	if cookie.Value != want.Value {
		t.Errorf("cookie %q value does not match", want.Name)
	}
	if cookie.Name != want.Name || cookie.Path != want.Path || cookie.HttpOnly != want.HttpOnly || cookie.Secure != want.Secure || cookie.SameSite != want.SameSite || cookie.MaxAge != want.MaxAge {
		t.Errorf(
			"cookie attributes = name %q path %q HttpOnly %t Secure %t SameSite %v MaxAge %d; want name %q path %q HttpOnly %t Secure %t SameSite %v MaxAge %d",
			cookie.Name, cookie.Path, cookie.HttpOnly, cookie.Secure, cookie.SameSite, cookie.MaxAge,
			want.Name, want.Path, want.HttpOnly, want.Secure, want.SameSite, want.MaxAge,
		)
	}
}
