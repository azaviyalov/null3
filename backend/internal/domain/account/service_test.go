package account_test

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/azaviyalov/null3/backend/internal/core"
	"github.com/azaviyalov/null3/backend/internal/domain/account"
	"github.com/azaviyalov/null3/backend/internal/domain/session"
	"github.com/azaviyalov/null3/backend/internal/testutil"
	"golang.org/x/crypto/bcrypt"
)

func TestServiceRegistrationValidation(t *testing.T) {
	testutil.SkipIntegration(t)
	tests := []struct {
		name     string
		login    string
		password string
		wantErr  bool
	}{
		{name: "minimum lengths", login: "abc", password: strings.Repeat("p", 8)},
		{name: "maximum lengths", login: strings.Repeat("a", 32), password: strings.Repeat("p", 72)},
		{name: "login too short", login: "ab", password: testPassword, wantErr: true},
		{name: "login too long", login: strings.Repeat("a", 33), password: testPassword, wantErr: true},
		{name: "unsupported login character", login: "user.name", password: testPassword, wantErr: true},
		{name: "password too short", login: "journal_user", password: strings.Repeat("p", 7), wantErr: true},
		{name: "password too long", login: "journal_user", password: strings.Repeat("p", 73), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			environment := newAccountTestEnvironment(t)
			rawToken, _, err := environment.service.CreateInvite(t.Context())
			if err != nil {
				t.Fatalf("CreateInvite() error = %v", err)
			}

			_, _, err = environment.service.RegisterWithInvite(t.Context(), rawToken, account.InviteRegistrationRequest{
				Login:    tt.login,
				Email:    "person@example.test",
				Password: tt.password,
			})

			if tt.wantErr && !errors.Is(err, core.ErrInvalidItem) {
				t.Fatalf("RegisterWithInvite() error = %v, want ErrInvalidItem", err)
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("RegisterWithInvite() error = %v", err)
			}
		})
	}
}

func TestServiceAuthenticateUser(t *testing.T) {
	testutil.SkipIntegration(t)
	environment := newAccountTestEnvironment(t)
	user := createTestUser(t, environment, "journal_user", "person@example.test")

	t.Run("valid credentials", func(t *testing.T) {
		response, tokens, err := environment.service.AuthenticateUser(t.Context(), account.LoginRequest{
			Login:    "  JOURNAL_USER  ",
			Password: testPassword,
		})

		if err != nil {
			t.Fatalf("AuthenticateUser() error = %v", err)
		}
		assertUserResponse(t, response, user)
		if tokens == nil || tokens.AccessToken == "" || tokens.RefreshToken == nil || tokens.RefreshToken.Value == "" {
			t.Fatal("AuthenticateUser() returned incomplete session tokens")
		}
		if _, err := environment.repository.SessionRepository().GetRefreshToken(t.Context(), tokens.RefreshToken.Value); err != nil {
			t.Fatalf("get persisted refresh token: %v", err)
		}
	})

	t.Run("invalid credentials have one public error", func(t *testing.T) {
		tests := []struct {
			name    string
			request account.LoginRequest
		}{
			{name: "unknown login", request: account.LoginRequest{Login: "unknown", Password: testPassword}},
			{name: "wrong password", request: account.LoginRequest{Login: user.Login, Password: "incorrect-password"}},
			{name: "blank login", request: account.LoginRequest{Login: " ", Password: testPassword}},
			{name: "blank password", request: account.LoginRequest{Login: user.Login, Password: " "}},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				response, tokens, err := environment.service.AuthenticateUser(t.Context(), tt.request)
				if !errors.Is(err, account.ErrInvalidCredentials) {
					t.Fatalf("AuthenticateUser() error = %v, want account.ErrInvalidCredentials", err)
				}
				if response != nil || tokens != nil {
					t.Fatal("AuthenticateUser() returned data for invalid credentials")
				}
			})
		}
	})
}

func TestServiceInviteLifecycle(t *testing.T) {
	testutil.SkipIntegration(t)
	environment := newAccountTestEnvironment(t)

	rawToken, created, err := environment.service.CreateInvite(t.Context())
	if err != nil {
		t.Fatalf("CreateInvite() error = %v", err)
	}
	if rawToken == "" {
		t.Fatal("CreateInvite() returned an empty token")
	}
	if created.TokenHash == rawToken {
		t.Fatal("CreateInvite() stored the raw token")
	}
	validated, err := environment.service.ValidateInvite(t.Context(), rawToken)
	if err != nil {
		t.Fatalf("ValidateInvite() error = %v", err)
	}
	if validated.ID != created.ID {
		t.Errorf("ValidateInvite() ID = %d, want %d", validated.ID, created.ID)
	}

	usedAt := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	created.UsedAt = &usedAt
	if err := environment.repository.SaveInvite(t.Context(), created); err != nil {
		t.Fatalf("mark invite used: %v", err)
	}
	if _, err := environment.service.ValidateInvite(t.Context(), rawToken); !errors.Is(err, account.ErrInviteAlreadyUsed) {
		t.Fatalf("ValidateInvite() used error = %v, want account.ErrInviteAlreadyUsed", err)
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
	if _, err := environment.service.ValidateInvite(t.Context(), expiredRawToken); !errors.Is(err, account.ErrInviteExpired) {
		t.Fatalf("ValidateInvite() expired error = %v, want account.ErrInviteExpired", err)
	}
	if _, err := environment.service.ValidateInvite(t.Context(), "unknown-invite-token"); !errors.Is(err, account.ErrInviteInvalid) {
		t.Fatalf("ValidateInvite() unknown error = %v, want account.ErrInviteInvalid", err)
	}
}

func TestServiceRegisterWithInvite(t *testing.T) {
	testutil.SkipIntegration(t)
	environment := newAccountTestEnvironment(t)
	rawToken, invite, err := environment.service.CreateInvite(t.Context())
	if err != nil {
		t.Fatalf("CreateInvite() error = %v", err)
	}

	response, tokens, err := environment.service.RegisterWithInvite(t.Context(), rawToken, account.InviteRegistrationRequest{
		Login:    "  Journal_User  ",
		Email:    "  Person@Example.TEST  ",
		Password: testPassword,
	})

	if err != nil {
		t.Fatalf("RegisterWithInvite() error = %v", err)
	}
	if response.Login != "journal_user" || response.Email != "person@example.test" {
		t.Errorf("RegisterWithInvite() identity = login %q email %q, want normalized values", response.Login, response.Email)
	}
	user, err := environment.repository.GetUserByID(t.Context(), response.ID)
	if err != nil {
		t.Fatalf("get registered user: %v", err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(testPassword)); err != nil {
		t.Fatal("registered password hash does not match")
	}
	storedInvite, err := environment.repository.GetInviteByHash(t.Context(), invite.TokenHash)
	if err != nil {
		t.Fatalf("get registered invite: %v", err)
	}
	if storedInvite.UsedAt == nil || storedInvite.RegisteredUserID == nil || *storedInvite.RegisteredUserID != user.ID {
		t.Errorf("registered invite used = %t registered user matches = %t, want both true", storedInvite.UsedAt != nil, storedInvite.RegisteredUserID != nil && *storedInvite.RegisteredUserID == user.ID)
	}
	if tokens == nil || tokens.RefreshToken == nil || tokens.RefreshToken.Value == "" {
		t.Fatal("RegisterWithInvite() returned incomplete session tokens")
	}
	if gotUserID, err := environment.sessionService.ParseUserAccessToken(tokens.AccessToken); err != nil || gotUserID != user.ID {
		t.Fatalf("parse registered access token: user ID = %d, error = %v", gotUserID, err)
	}

	if _, _, err := environment.service.RegisterWithInvite(t.Context(), rawToken, account.InviteRegistrationRequest{
		Login:    "another_user",
		Email:    "another@example.test",
		Password: testPassword,
	}); !errors.Is(err, account.ErrInviteAlreadyUsed) {
		t.Fatalf("reuse invite error = %v, want account.ErrInviteAlreadyUsed", err)
	}
}

func TestServiceRegistrationConflictDoesNotConsumeInvite(t *testing.T) {
	testutil.SkipIntegration(t)
	tests := []struct {
		name    string
		request account.InviteRegistrationRequest
		wantErr error
	}{
		{
			name: "login",
			request: account.InviteRegistrationRequest{
				Login:    "Existing_User",
				Email:    "different@example.test",
				Password: testPassword,
			},
			wantErr: account.ErrLoginAlreadyTaken,
		},
		{
			name: "email",
			request: account.InviteRegistrationRequest{
				Login:    "different_user",
				Email:    "EXISTING@EXAMPLE.TEST",
				Password: testPassword,
			},
			wantErr: account.ErrEmailAlreadyTaken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			environment := newAccountTestEnvironment(t)
			createTestUser(t, environment, "existing_user", "existing@example.test")
			rawToken, invite, err := environment.service.CreateInvite(t.Context())
			if err != nil {
				t.Fatalf("CreateInvite() error = %v", err)
			}

			_, _, err = environment.service.RegisterWithInvite(t.Context(), rawToken, tt.request)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("RegisterWithInvite() error = %v, want %v", err, tt.wantErr)
			}
			storedInvite, err := environment.repository.GetInviteByHash(t.Context(), invite.TokenHash)
			if err != nil {
				t.Fatalf("get invite after conflict: %v", err)
			}
			if storedInvite.UsedAt != nil || storedInvite.RegisteredUserID != nil {
				t.Fatal("registration conflict consumed the invite")
			}
			var userCount int64
			if err := environment.database.Model(&account.User{}).Count(&userCount).Error; err != nil {
				t.Fatalf("count users: %v", err)
			}
			if userCount != 1 {
				t.Fatalf("user count = %d, want 1", userCount)
			}
		})
	}
}

func TestServiceRefreshUserSession(t *testing.T) {
	testutil.SkipIntegration(t)
	environment := newAccountTestEnvironment(t)
	user := createTestUser(t, environment, "journal_user", "person@example.test")
	oldToken, err := environment.sessionService.CreateRefreshToken(t.Context(), user.ID)
	if err != nil {
		t.Fatalf("create refresh token: %v", err)
	}

	response, newTokens, err := environment.service.RefreshUserSession(t.Context(), oldToken.Value)

	if err != nil {
		t.Fatalf("RefreshUserSession() error = %v", err)
	}
	if response.ID != user.ID {
		t.Errorf("RefreshUserSession() user ID = %d, want %d", response.ID, user.ID)
	}
	if newTokens == nil || newTokens.RefreshToken == nil || newTokens.RefreshToken.Value == "" {
		t.Fatal("RefreshUserSession() returned incomplete session tokens")
	}
	if newTokens.RefreshToken.Value == oldToken.Value {
		t.Fatal("RefreshUserSession() did not rotate the refresh token")
	}
	if _, _, err := environment.service.RefreshUserSession(t.Context(), oldToken.Value); !errors.Is(err, session.ErrRefreshTokenInvalid) {
		t.Fatalf("old refresh token error = %v, want ErrRefreshTokenInvalid", err)
	}

	if err := environment.database.Model(&session.RefreshToken{}).
		Where("user_id = ?", user.ID).
		Update("expires_at", time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)).Error; err != nil {
		t.Fatalf("expire refresh token: %v", err)
	}
	if _, _, err := environment.service.RefreshUserSession(t.Context(), newTokens.RefreshToken.Value); !errors.Is(err, session.ErrRefreshTokenInvalid) {
		t.Fatalf("expired refresh token error = %v, want ErrRefreshTokenInvalid", err)
	}
	var refreshTokenCount int64
	if err := environment.database.Model(&session.RefreshToken{}).Where("user_id = ?", user.ID).Count(&refreshTokenCount).Error; err != nil {
		t.Fatalf("count refresh tokens: %v", err)
	}
	if refreshTokenCount != 0 {
		t.Fatalf("refresh token count = %d, want 0", refreshTokenCount)
	}
}

func TestServicePasswordReset(t *testing.T) {
	testutil.SkipIntegration(t)
	environment := newAccountTestEnvironment(t)
	user := createTestUser(t, environment, "journal_user", "person@example.test")
	refreshToken, err := environment.sessionService.CreateRefreshToken(t.Context(), user.ID)
	if err != nil {
		t.Fatalf("create refresh token: %v", err)
	}

	unknownToken, err := environment.service.RequestPasswordReset(t.Context(), account.ForgotPasswordRequest{Email: "unknown@example.test"})
	if err != nil {
		t.Fatalf("RequestPasswordReset() unknown email error = %v", err)
	}
	if unknownToken != "" {
		t.Fatal("RequestPasswordReset() returned a token for an unknown email")
	}

	firstToken, err := environment.service.RequestPasswordReset(t.Context(), account.ForgotPasswordRequest{Email: " PERSON@EXAMPLE.TEST "})
	if err != nil {
		t.Fatalf("RequestPasswordReset() error = %v", err)
	}
	secondToken, err := environment.service.RequestPasswordReset(t.Context(), account.ForgotPasswordRequest{Email: user.Email})
	if err != nil {
		t.Fatalf("replace password reset token: %v", err)
	}
	if firstToken == "" || secondToken == "" || firstToken == secondToken {
		t.Fatal("RequestPasswordReset() did not generate distinct tokens")
	}
	if err := environment.service.ResetPassword(t.Context(), account.ResetPasswordRequest{
		Token:    firstToken,
		Password: "new-correct-password",
	}); !errors.Is(err, account.ErrPasswordResetTokenInvalid) {
		t.Fatalf("first password reset token error = %v, want ErrPasswordResetTokenInvalid", err)
	}
	var storedToken account.PasswordResetToken
	if err := environment.database.Where("user_id = ?", user.ID).First(&storedToken).Error; err != nil {
		t.Fatalf("get replacement password reset token: %v", err)
	}
	if storedToken.TokenHash == firstToken || storedToken.TokenHash == secondToken {
		t.Fatal("RequestPasswordReset() stored the raw token")
	}

	if err := environment.service.ResetPassword(t.Context(), account.ResetPasswordRequest{
		Token:    "unknown-reset-token",
		Password: "new-correct-password",
	}); !errors.Is(err, account.ErrPasswordResetTokenInvalid) {
		t.Fatalf("ResetPassword() unknown token error = %v, want account.ErrPasswordResetTokenInvalid", err)
	}

	if err := environment.service.ResetPassword(t.Context(), account.ResetPasswordRequest{
		Token:    secondToken,
		Password: "new-correct-password",
	}); err != nil {
		t.Fatalf("ResetPassword() error = %v", err)
	}
	updatedUser, err := environment.repository.GetUserByID(t.Context(), user.ID)
	if err != nil {
		t.Fatalf("get user after password reset: %v", err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(updatedUser.PasswordHash), []byte("new-correct-password")); err != nil {
		t.Fatal("new password hash does not match")
	}
	var resetTokenCount int64
	if err := environment.database.Model(&account.PasswordResetToken{}).Where("user_id = ?", user.ID).Count(&resetTokenCount).Error; err != nil {
		t.Fatalf("count password reset tokens: %v", err)
	}
	if resetTokenCount != 0 {
		t.Fatalf("password reset token count = %d, want 0", resetTokenCount)
	}
	if _, err := environment.repository.SessionRepository().GetRefreshToken(t.Context(), refreshToken.Value); !errors.Is(err, core.ErrItemNotFound) {
		t.Fatalf("revoked refresh token error = %v, want ErrItemNotFound", err)
	}
	if _, _, err := environment.service.AuthenticateUser(t.Context(), account.LoginRequest{
		Login:    user.Login,
		Password: testPassword,
	}); !errors.Is(err, account.ErrInvalidCredentials) {
		t.Fatalf("old password authentication error = %v, want account.ErrInvalidCredentials", err)
	}
	if _, _, err := environment.service.AuthenticateUser(t.Context(), account.LoginRequest{
		Login:    user.Login,
		Password: "new-correct-password",
	}); err != nil {
		t.Fatalf("new password authentication error = %v", err)
	}
}

func TestServiceRejectsExpiredPasswordResetToken(t *testing.T) {
	testutil.SkipIntegration(t)
	environment := newAccountTestEnvironment(t)
	user := createTestUser(t, environment, "journal_user", "person@example.test")
	rawToken, err := environment.service.RequestPasswordReset(t.Context(), account.ForgotPasswordRequest{Email: user.Email})
	if err != nil {
		t.Fatalf("create expired password reset token: %v", err)
	}
	if err := environment.database.Model(&account.PasswordResetToken{}).
		Where("user_id = ?", user.ID).
		Update("expires_at", time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)).Error; err != nil {
		t.Fatalf("expire password reset token: %v", err)
	}

	err = environment.service.ResetPassword(t.Context(), account.ResetPasswordRequest{
		Token:    rawToken,
		Password: "new-correct-password",
	})

	if !errors.Is(err, account.ErrPasswordResetTokenExpired) {
		t.Fatalf("ResetPassword() error = %v, want account.ErrPasswordResetTokenExpired", err)
	}
	unchangedUser, err := environment.repository.GetUserByID(t.Context(), user.ID)
	if err != nil {
		t.Fatalf("get user after rejected password reset: %v", err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(unchangedUser.PasswordHash), []byte(testPassword)); err != nil {
		t.Fatal("expired password reset changed the password")
	}
}
