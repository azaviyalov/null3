package account_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/azaviyalov/null3/backend/internal/core/server"
	"github.com/azaviyalov/null3/backend/internal/domain/account"
	"github.com/azaviyalov/null3/backend/internal/domain/session"
	"github.com/azaviyalov/null3/backend/internal/testutil"
	"github.com/labstack/echo/v4"
)

func TestAccountSessionHTTPFlow(t *testing.T) {
	testutil.SkipIntegration(t)
	environment := newAccountTestEnvironment(t)
	user := createTestUser(t, environment, "journal_user", "person@example.test")
	e := newAccountTestServer(t, environment)

	loginResponse := testutil.JSONRequest(t, e, http.MethodPost, "/api/auth/login", `{
		"login":"JOURNAL_USER",
		"password":"correct-password"
	}`)
	if loginResponse.Code != http.StatusOK {
		t.Fatalf("login status = %d, want %d", loginResponse.Code, http.StatusOK)
	}
	if contentType := loginResponse.Header().Get(echo.HeaderContentType); !strings.HasPrefix(contentType, echo.MIMEApplicationJSON) {
		t.Errorf("login Content-Type = %q, want JSON", contentType)
	}
	var loginUser account.UserResponse
	testutil.DecodeJSON(t, loginResponse, &loginUser)
	assertUserResponse(t, &loginUser, user)
	accessCookie := testutil.ResponseCookie(t, loginResponse, session.UserCookieName)
	refreshCookie := testutil.ResponseCookie(t, loginResponse, session.UserRefreshCookieName)
	if !accessCookie.HttpOnly || !refreshCookie.HttpOnly {
		t.Fatal("login cookies are not HttpOnly")
	}

	meResponse := testutil.JSONRequest(t, e, http.MethodGet, "/api/auth/me", nil, accessCookie)
	if meResponse.Code != http.StatusOK {
		t.Fatalf("me status = %d, want %d", meResponse.Code, http.StatusOK)
	}

	refreshResponse := testutil.JSONRequest(t, e, http.MethodPost, "/api/auth/refresh", nil, refreshCookie)
	if refreshResponse.Code != http.StatusOK {
		t.Fatalf("refresh status = %d, want %d", refreshResponse.Code, http.StatusOK)
	}
	newAccessCookie := testutil.ResponseCookie(t, refreshResponse, session.UserCookieName)
	newRefreshCookie := testutil.ResponseCookie(t, refreshResponse, session.UserRefreshCookieName)
	if newRefreshCookie.Value == refreshCookie.Value {
		t.Fatal("refresh endpoint did not rotate the refresh token")
	}

	reusedResponse := testutil.JSONRequest(t, e, http.MethodPost, "/api/auth/refresh", nil, refreshCookie)
	if reusedResponse.Code != http.StatusUnauthorized {
		t.Fatalf("reused refresh status = %d, want %d", reusedResponse.Code, http.StatusUnauthorized)
	}

	logoutResponse := testutil.JSONRequest(t, e, http.MethodPost, "/api/auth/logout", nil, newAccessCookie, newRefreshCookie)
	if logoutResponse.Code != http.StatusOK {
		t.Fatalf("logout status = %d, want %d", logoutResponse.Code, http.StatusOK)
	}
	clearedAccessCookie := testutil.ResponseCookie(t, logoutResponse, session.UserCookieName)
	clearedRefreshCookie := testutil.ResponseCookie(t, logoutResponse, session.UserRefreshCookieName)
	if clearedAccessCookie.MaxAge != -1 || clearedRefreshCookie.MaxAge != -1 {
		t.Fatal("logout endpoint did not expire both session cookies")
	}

	var refreshTokenCount int64
	if err := environment.database.Model(&session.RefreshToken{}).Where("user_id = ?", user.ID).Count(&refreshTokenCount).Error; err != nil {
		t.Fatalf("count refresh tokens: %v", err)
	}
	if refreshTokenCount != 0 {
		t.Fatalf("refresh token count = %d, want 0", refreshTokenCount)
	}
}

func TestAccountLoginHTTPRejections(t *testing.T) {
	testutil.SkipIntegration(t)
	environment := newAccountTestEnvironment(t)
	createTestUser(t, environment, "journal_user", "person@example.test")
	e := newAccountTestServer(t, environment)

	tests := []struct {
		name        string
		body        string
		wantStatus  int
		wantMessage string
	}{
		{
			name:       "malformed JSON",
			body:       `{"login":`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing required field",
			body:       `{"login":"journal_user"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:        "invalid credentials",
			body:        `{"login":"journal_user","password":"incorrect-password"}`,
			wantStatus:  http.StatusUnauthorized,
			wantMessage: "Incorrect login credentials.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := testutil.JSONRequest(t, e, http.MethodPost, "/api/auth/login", tt.body)
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
		})
	}
}

func TestAccountInviteHTTPFlow(t *testing.T) {
	testutil.SkipIntegration(t)
	environment := newAccountTestEnvironment(t)
	e := newAccountTestServer(t, environment)
	rawToken, _, err := environment.service.CreateInvite(t.Context())
	if err != nil {
		t.Fatalf("CreateInvite() error = %v", err)
	}
	invitePath := "/api/auth/invites/" + rawToken

	validationResponse := testutil.JSONRequest(t, e, http.MethodGet, invitePath, nil)
	if validationResponse.Code != http.StatusOK {
		t.Fatalf("invite validation status = %d, want %d", validationResponse.Code, http.StatusOK)
	}

	registrationResponse := testutil.JSONRequest(t, e, http.MethodPost, invitePath+"/register", `{
		"login":"journal_user",
		"email":"person@example.test",
		"password":"correct-password"
	}`)
	if registrationResponse.Code != http.StatusCreated {
		t.Fatalf("invite registration status = %d, want %d", registrationResponse.Code, http.StatusCreated)
	}
	testutil.ResponseCookie(t, registrationResponse, session.UserCookieName)
	testutil.ResponseCookie(t, registrationResponse, session.UserRefreshCookieName)

	reusedResponse := testutil.JSONRequest(t, e, http.MethodGet, invitePath, nil)
	if reusedResponse.Code != http.StatusBadRequest {
		t.Fatalf("used invite status = %d, want %d", reusedResponse.Code, http.StatusBadRequest)
	}
	var errorResponse account.MessageResponse
	testutil.DecodeJSON(t, reusedResponse, &errorResponse)
	if errorResponse.Message != "This invite link has already been used." {
		t.Errorf("used invite message = %q", errorResponse.Message)
	}

	invalidResponse := testutil.JSONRequest(t, e, http.MethodGet, "/api/auth/invites/unknown-token", nil)
	if invalidResponse.Code != http.StatusBadRequest {
		t.Fatalf("invalid invite status = %d, want %d", invalidResponse.Code, http.StatusBadRequest)
	}
	testutil.DecodeJSON(t, invalidResponse, &errorResponse)
	if errorResponse.Message != "This invite link is invalid." {
		t.Errorf("invalid invite message = %q", errorResponse.Message)
	}

	expiredRawToken, expiredInvite, err := environment.service.CreateInvite(t.Context())
	if err != nil {
		t.Fatalf("create expired invite: %v", err)
	}
	if err := environment.database.Model(&account.Invite{}).
		Where("id = ?", expiredInvite.ID).
		Update("expires_at", time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)).Error; err != nil {
		t.Fatalf("expire invite: %v", err)
	}
	expiredResponse := testutil.JSONRequest(t, e, http.MethodGet, "/api/auth/invites/"+expiredRawToken, nil)
	if expiredResponse.Code != http.StatusBadRequest {
		t.Fatalf("expired invite status = %d, want %d", expiredResponse.Code, http.StatusBadRequest)
	}
	testutil.DecodeJSON(t, expiredResponse, &errorResponse)
	if errorResponse.Message != "This invite link has expired." {
		t.Errorf("expired invite message = %q", errorResponse.Message)
	}

	conflictToken, _, err := environment.service.CreateInvite(t.Context())
	if err != nil {
		t.Fatalf("create conflict invite: %v", err)
	}
	conflictResponse := testutil.JSONRequest(t, e, http.MethodPost, "/api/auth/invites/"+conflictToken+"/register", `{
		"login":"journal_user",
		"email":"different@example.test",
		"password":"correct-password"
	}`)
	if conflictResponse.Code != http.StatusConflict {
		t.Fatalf("registration conflict status = %d, want %d", conflictResponse.Code, http.StatusConflict)
	}
	testutil.DecodeJSON(t, conflictResponse, &errorResponse)
	if errorResponse.Message != "That login is already in use." {
		t.Errorf("registration conflict message = %q", errorResponse.Message)
	}
}

func TestAccountPasswordRecoveryHTTPFlow(t *testing.T) {
	testutil.SkipIntegration(t)
	environment := newAccountTestEnvironment(t)
	user := createTestUser(t, environment, "journal_user", "person@example.test")
	e := newAccountTestServer(t, environment)

	unknownResponse := testutil.JSONRequest(t, e, http.MethodPost, "/api/auth/forgot-password", `{"email":"unknown@example.test"}`)
	if unknownResponse.Code != http.StatusOK {
		t.Fatalf("unknown email status = %d, want %d", unknownResponse.Code, http.StatusOK)
	}
	var unknownBody account.ForgotPasswordResponse
	testutil.DecodeJSON(t, unknownResponse, &unknownBody)

	knownRequest := httptest.NewRequest(http.MethodPost, "/api/auth/forgot-password", strings.NewReader(`{"email":"person@example.test"}`))
	knownRequest.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	knownRequest.Header.Set(echo.HeaderOrigin, "https://attacker.example")
	knownRequest.Host = "attacker.example"
	knownResponse := httptest.NewRecorder()
	e.ServeHTTP(knownResponse, knownRequest)
	if knownResponse.Code != http.StatusOK {
		t.Fatalf("known email status = %d, want %d", knownResponse.Code, http.StatusOK)
	}
	var knownBody account.ForgotPasswordResponse
	testutil.DecodeJSON(t, knownResponse, &knownBody)
	if knownBody.Message != unknownBody.Message {
		t.Fatal("password recovery messages differ for known and unknown emails")
	}
	if unknownBody.ResetURL != "" {
		t.Fatal("unknown email response contains a reset URL")
	}
	if !strings.HasPrefix(knownBody.ResetURL, "https://journal.example/reset-password?") {
		t.Fatal("known email response does not use the configured frontend URL")
	}
	if strings.Contains(knownBody.ResetURL, "attacker.example") {
		t.Fatal("known email response trusted request host headers")
	}

	resetURL, err := url.Parse(knownBody.ResetURL)
	if err != nil {
		t.Fatalf("parse reset URL: %v", err)
	}
	resetToken := resetURL.Query().Get("token")
	if resetToken == "" {
		t.Fatal("reset URL does not contain a token")
	}

	weakPasswordResponse := testutil.JSONRequest(t, e, http.MethodPost, "/api/auth/reset-password", `{
		"token":"`+resetToken+`",
		"password":"short"
	}`)
	if weakPasswordResponse.Code != http.StatusBadRequest {
		t.Fatalf("weak password status = %d, want %d", weakPasswordResponse.Code, http.StatusBadRequest)
	}
	var errorResponse account.MessageResponse
	testutil.DecodeJSON(t, weakPasswordResponse, &errorResponse)
	if errorResponse.Message != "password must be between 8 and 72 characters" {
		t.Errorf("weak password message = %q", errorResponse.Message)
	}

	resetResponse := testutil.JSONRequest(t, e, http.MethodPost, "/api/auth/reset-password", `{
		"token":"`+resetToken+`",
		"password":"new-correct-password"
	}`)
	if resetResponse.Code != http.StatusOK {
		t.Fatalf("password reset status = %d, want %d", resetResponse.Code, http.StatusOK)
	}

	reusedResponse := testutil.JSONRequest(t, e, http.MethodPost, "/api/auth/reset-password", `{
		"token":"`+resetToken+`",
		"password":"another-correct-password"
	}`)
	if reusedResponse.Code != http.StatusBadRequest {
		t.Fatalf("reused reset token status = %d, want %d", reusedResponse.Code, http.StatusBadRequest)
	}
	testutil.DecodeJSON(t, reusedResponse, &errorResponse)
	if errorResponse.Message != "This password reset link is invalid." {
		t.Errorf("reused reset token message = %q", errorResponse.Message)
	}

	expiredRawToken, err := environment.service.RequestPasswordReset(t.Context(), account.ForgotPasswordRequest{Email: "person@example.test"})
	if err != nil {
		t.Fatalf("create expired reset token: %v", err)
	}
	if err := environment.database.Model(&account.PasswordResetToken{}).
		Where("user_id = ?", user.ID).
		Update("expires_at", time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)).Error; err != nil {
		t.Fatalf("expire reset token: %v", err)
	}
	expiredResponse := testutil.JSONRequest(t, e, http.MethodPost, "/api/auth/reset-password", `{
		"token":"`+expiredRawToken+`",
		"password":"another-correct-password"
	}`)
	if expiredResponse.Code != http.StatusBadRequest {
		t.Fatalf("expired reset token status = %d, want %d", expiredResponse.Code, http.StatusBadRequest)
	}
	testutil.DecodeJSON(t, expiredResponse, &errorResponse)
	if errorResponse.Message != "This password reset link has expired." {
		t.Errorf("expired reset token message = %q", errorResponse.Message)
	}
}

func newAccountTestServer(t *testing.T, environment *accountTestEnvironment) *echo.Echo {
	t.Helper()

	testutil.DiscardLogs(t)

	e := server.NewEchoServer(server.Config{})
	handler := account.NewHandler(environment.service, environment.sessionService, environment.accountConfig, environment.sessionConfig)
	validateUser := func(ctx context.Context, userID uint) error {
		_, err := environment.service.GetUserByID(ctx, userID)
		return err
	}
	account.RegisterRoutes(e, handler, session.UserJWTMiddleware(environment.sessionService, validateUser))
	return e
}
