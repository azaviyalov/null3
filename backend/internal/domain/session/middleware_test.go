package session_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/azaviyalov/null3/backend/internal/domain/session"
	"github.com/labstack/echo/v4"
)

func TestUserJWTMiddleware(t *testing.T) {
	service := session.NewService(nil, session.Config{JWTSecret: testJWTSecret, JWTExpiration: time.Hour})
	userToken, err := service.GenerateUserAccessToken(42)
	if err != nil {
		t.Fatalf("GenerateUserAccessToken() error = %v", err)
	}
	adminToken, err := service.GenerateAdminAccessToken(time.Hour)
	if err != nil {
		t.Fatalf("GenerateAdminAccessToken() error = %v", err)
	}

	tests := []struct {
		name         string
		token        string
		validateUser func(context.Context, uint) error
		wantStatus   int
		wantRun      bool
	}{
		{name: "missing cookie", validateUser: acceptUser, wantStatus: http.StatusUnauthorized},
		{name: "malformed token", token: "malformed", validateUser: acceptUser, wantStatus: http.StatusUnauthorized},
		{name: "admin token", token: adminToken, validateUser: acceptUser, wantStatus: http.StatusUnauthorized},
		{name: "unknown user", token: userToken, validateUser: func(context.Context, uint) error { return errors.New("not found") }, wantStatus: http.StatusUnauthorized},
		{name: "valid user", token: userToken, validateUser: acceptUser, wantStatus: http.StatusNoContent, wantRun: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			ran := false
			e.GET("/private", func(c echo.Context) error {
				ran = true
				if userID := session.GetUserID(c); userID != 42 {
					t.Errorf("GetUserID() = %d, want 42", userID)
				}
				return c.NoContent(http.StatusNoContent)
			}, session.UserJWTMiddleware(service, tt.validateUser))

			request := httptest.NewRequest(http.MethodGet, "/private", nil)
			if tt.token != "" {
				request.AddCookie(&http.Cookie{Name: session.UserCookieName, Value: tt.token})
			}
			recorder := httptest.NewRecorder()
			e.ServeHTTP(recorder, request)

			if recorder.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", recorder.Code, tt.wantStatus)
			}
			if ran != tt.wantRun {
				t.Errorf("downstream ran = %t, want %t", ran, tt.wantRun)
			}
		})
	}
}

func TestAdminJWTMiddleware(t *testing.T) {
	service := session.NewService(nil, session.Config{JWTSecret: testJWTSecret, JWTExpiration: time.Hour})
	adminToken, err := service.GenerateAdminAccessToken(time.Hour)
	if err != nil {
		t.Fatalf("GenerateAdminAccessToken() error = %v", err)
	}
	userToken, err := service.GenerateUserAccessToken(42)
	if err != nil {
		t.Fatalf("GenerateUserAccessToken() error = %v", err)
	}

	tests := []struct {
		name       string
		token      string
		wantStatus int
		wantRun    bool
	}{
		{name: "missing cookie", wantStatus: http.StatusUnauthorized},
		{name: "malformed token", token: "malformed", wantStatus: http.StatusUnauthorized},
		{name: "user token", token: userToken, wantStatus: http.StatusUnauthorized},
		{name: "valid admin", token: adminToken, wantStatus: http.StatusNoContent, wantRun: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			ran := false
			e.GET("/admin", func(c echo.Context) error {
				ran = true
				return c.NoContent(http.StatusNoContent)
			}, session.AdminJWTMiddleware(service))

			request := httptest.NewRequest(http.MethodGet, "/admin", nil)
			if tt.token != "" {
				request.AddCookie(&http.Cookie{Name: session.AdminCookieName, Value: tt.token})
			}
			recorder := httptest.NewRecorder()
			e.ServeHTTP(recorder, request)

			if recorder.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", recorder.Code, tt.wantStatus)
			}
			if ran != tt.wantRun {
				t.Errorf("downstream ran = %t, want %t", ran, tt.wantRun)
			}
		})
	}
}

func acceptUser(context.Context, uint) error {
	return nil
}
