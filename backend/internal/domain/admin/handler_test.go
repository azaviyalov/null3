package admin_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/azaviyalov/null3/backend/internal/domain/account"
	"github.com/azaviyalov/null3/backend/internal/domain/session"
	"github.com/azaviyalov/null3/backend/internal/testutil"
	"github.com/labstack/echo/v4"
)

func TestAdminAuthenticationHTTPFlow(t *testing.T) {
	testutil.SkipIntegration(t)
	environment := newAdminTestEnvironment(t)

	unauthorizedResponse := testutil.JSONRequest(t, environment.echo, http.MethodGet, "/api/admin/auth/me", nil)
	if unauthorizedResponse.Code != http.StatusUnauthorized {
		t.Fatalf("unauthorized me status = %d, want %d", unauthorizedResponse.Code, http.StatusUnauthorized)
	}

	tests := []struct {
		name        string
		body        string
		wantStatus  int
		wantMessage string
	}{
		{name: "malformed JSON", body: `{"password":`, wantStatus: http.StatusBadRequest},
		{name: "empty password", body: `{}`, wantStatus: http.StatusUnauthorized, wantMessage: "Incorrect admin credentials."},
		{name: "wrong password", body: `{"password":"incorrect-admin-password"}`, wantStatus: http.StatusUnauthorized, wantMessage: "Incorrect admin credentials."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := testutil.JSONRequest(t, environment.echo, http.MethodPost, "/api/admin/auth/login", tt.body)
			if response.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", response.Code, tt.wantStatus)
			}
			if tt.wantMessage != "" {
				var body account.MessageResponse
				testutil.DecodeJSON(t, response, &body)
				if body.Message != tt.wantMessage {
					t.Errorf("message = %q, want %q", body.Message, tt.wantMessage)
				}
			}
			if len(response.Result().Cookies()) != 0 {
				t.Fatal("rejected login set a cookie")
			}
		})
	}

	adminCookie := loginAdmin(t, environment)
	if !adminCookie.HttpOnly || !adminCookie.Secure || adminCookie.Path != "/api/admin" {
		t.Errorf("admin cookie attributes = HttpOnly %t Secure %t Path %q", adminCookie.HttpOnly, adminCookie.Secure, adminCookie.Path)
	}
	if adminCookie.SameSite != http.SameSiteLaxMode || adminCookie.MaxAge != 30*60 {
		t.Errorf("admin cookie = SameSite %v MaxAge %d", adminCookie.SameSite, adminCookie.MaxAge)
	}

	meResponse := testutil.JSONRequest(t, environment.echo, http.MethodGet, "/api/admin/auth/me", nil, adminCookie)
	if meResponse.Code != http.StatusOK {
		t.Fatalf("authenticated me status = %d, want %d", meResponse.Code, http.StatusOK)
	}

	userToken, err := environment.sessionService.GenerateUserAccessToken(42)
	if err != nil {
		t.Fatalf("generate user access token: %v", err)
	}
	userScopeCookie := &http.Cookie{Name: session.AdminCookieName, Value: userToken}
	wrongScopeResponse := testutil.JSONRequest(t, environment.echo, http.MethodGet, "/api/admin/auth/me", nil, userScopeCookie)
	if wrongScopeResponse.Code != http.StatusUnauthorized {
		t.Fatalf("user-scope token status = %d, want %d", wrongScopeResponse.Code, http.StatusUnauthorized)
	}

	assertAdminIsStateless(t, environment)

	logoutResponse := testutil.JSONRequest(t, environment.echo, http.MethodPost, "/api/admin/auth/logout", nil, adminCookie)
	if logoutResponse.Code != http.StatusOK {
		t.Fatalf("logout status = %d, want %d", logoutResponse.Code, http.StatusOK)
	}
	clearedCookie := testutil.ResponseCookie(t, logoutResponse, session.AdminCookieName)
	if clearedCookie.Value != "" || clearedCookie.MaxAge != -1 {
		t.Fatal("logout did not expire the admin cookie")
	}
}

func TestCreateInviteHTTP(t *testing.T) {
	testutil.SkipIntegration(t)
	environment := newAdminTestEnvironment(t)

	unauthorizedResponse := testutil.JSONRequest(t, environment.echo, http.MethodPost, "/api/admin/invites", nil)
	if unauthorizedResponse.Code != http.StatusUnauthorized {
		t.Fatalf("unauthorized create invite status = %d, want %d", unauthorizedResponse.Code, http.StatusUnauthorized)
	}
	var inviteCount int64
	if err := environment.database.Model(&account.Invite{}).Count(&inviteCount).Error; err != nil {
		t.Fatalf("count invites after unauthorized request: %v", err)
	}
	if inviteCount != 0 {
		t.Fatal("unauthorized create invite reached the handler")
	}

	adminCookie := loginAdmin(t, environment)
	request := httptest.NewRequest(http.MethodPost, "/api/admin/invites", nil)
	request.Header.Set(echo.HeaderOrigin, "https://attacker.example")
	request.Host = "attacker.example"
	request.AddCookie(adminCookie)
	response := httptest.NewRecorder()
	environment.echo.ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("create invite status = %d, want %d", response.Code, http.StatusCreated)
	}
	var body account.InviteResponse
	testutil.DecodeJSON(t, response, &body)
	if !strings.HasPrefix(body.InviteURL, "https://journal.example/invite/") {
		t.Fatal("invite response does not use the configured frontend URL")
	}
	if strings.Contains(body.InviteURL, "attacker.example") {
		t.Fatal("invite response trusted request host headers")
	}
	inviteURL, err := url.Parse(body.InviteURL)
	if err != nil {
		t.Fatalf("parse invite URL: %v", err)
	}
	rawToken := strings.TrimPrefix(inviteURL.Path, "/invite/")
	if rawToken == "" || rawToken == inviteURL.Path {
		t.Fatal("invite URL does not contain a token")
	}

	var storedInvite account.Invite
	if err := environment.database.First(&storedInvite).Error; err != nil {
		t.Fatalf("get stored invite: %v", err)
	}
	if storedInvite.TokenHash == rawToken {
		t.Fatal("create invite stored the raw token")
	}
	if !storedInvite.ExpiresAt.Equal(body.ExpiresAt) {
		t.Errorf("stored expiration = %v, response expiration = %v", storedInvite.ExpiresAt, body.ExpiresAt)
	}

	assertAdminIsStateless(t, environment)
}

func loginAdmin(t *testing.T, environment *adminTestEnvironment) *http.Cookie {
	t.Helper()

	response := testutil.JSONRequest(t, environment.echo, http.MethodPost, "/api/admin/auth/login", `{"password":"configured-admin-password"}`)
	if response.Code != http.StatusOK {
		t.Fatalf("admin login status = %d, want %d", response.Code, http.StatusOK)
	}
	if len(response.Result().Cookies()) != 1 {
		t.Fatalf("admin login cookie count = %d, want 1", len(response.Result().Cookies()))
	}
	return testutil.ResponseCookie(t, response, session.AdminCookieName)
}

func assertAdminIsStateless(t *testing.T, environment *adminTestEnvironment) {
	t.Helper()

	var userCount int64
	if err := environment.database.Model(&account.User{}).Count(&userCount).Error; err != nil {
		t.Fatalf("count users: %v", err)
	}
	var refreshTokenCount int64
	if err := environment.database.Model(&session.RefreshToken{}).Count(&refreshTokenCount).Error; err != nil {
		t.Fatalf("count refresh tokens: %v", err)
	}
	if userCount != 0 || refreshTokenCount != 0 {
		t.Errorf("admin authentication created %d users and %d refresh tokens", userCount, refreshTokenCount)
	}
}
