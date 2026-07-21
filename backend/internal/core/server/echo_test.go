package server_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/azaviyalov/null3/backend/internal/core/server"
	"github.com/azaviyalov/null3/backend/internal/testutil"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func TestNewEchoServerValidator(t *testing.T) {
	e := server.NewEchoServer(server.Config{})

	valid := validationRequest{
		Details: validationDetails{Value: "present"},
	}
	if err := e.Validator.Validate(valid); err != nil {
		t.Fatalf("validate valid request: %v", err)
	}

	err := e.Validator.Validate(validationRequest{})
	if err == nil {
		t.Fatal("validate empty request error = nil, want an error")
	}
	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		t.Fatalf("validation error = %T, want validator.ValidationErrors", err)
	}
	if len(validationErrors) != 1 {
		t.Fatalf("validation error count = %d, want 1", len(validationErrors))
	}
	if validationErrors[0].Field() != "Details" || validationErrors[0].Tag() != "required" {
		t.Fatalf(
			"validation error = field %q tag %q, want Details required",
			validationErrors[0].Field(),
			validationErrors[0].Tag(),
		)
	}
}

func TestNewEchoServerCORS(t *testing.T) {
	tests := []struct {
		name            string
		config          server.Config
		wantOrigin      string
		wantCredentials string
	}{
		{
			name:   "disabled",
			config: server.Config{},
		},
		{
			name: "enabled",
			config: server.Config{
				EnableCORS:  true,
				FrontendURL: "https://example.test",
			},
			wantOrigin:      "https://example.test",
			wantCredentials: "true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testutil.DiscardLogs(t)
			e := server.NewEchoServer(tt.config)
			e.GET("/resource", func(c echo.Context) error {
				return c.String(http.StatusOK, "ok")
			})
			request := httptest.NewRequest(http.MethodOptions, "/resource", nil)
			request.Header.Set("Origin", "https://example.test")
			request.Header.Set("Access-Control-Request-Method", http.MethodGet)
			response := httptest.NewRecorder()

			e.ServeHTTP(response, request)

			if got := response.Header().Get("Access-Control-Allow-Origin"); got != tt.wantOrigin {
				t.Errorf("Access-Control-Allow-Origin = %q, want %q", got, tt.wantOrigin)
			}
			if got := response.Header().Get("Access-Control-Allow-Credentials"); got != tt.wantCredentials {
				t.Errorf("Access-Control-Allow-Credentials = %q, want %q", got, tt.wantCredentials)
			}
		})
	}
}

func TestNewEchoServerRecoversFromPanic(t *testing.T) {
	testutil.DiscardLogs(t)
	e := server.NewEchoServer(server.Config{})
	e.Logger.SetOutput(io.Discard)
	e.GET("/panic", func(echo.Context) error {
		panic("unexpected failure")
	})
	request := httptest.NewRequest(http.MethodGet, "/panic", nil)
	response := httptest.NewRecorder()

	e.ServeHTTP(response, request)

	if response.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusInternalServerError)
	}
	if requestID := response.Header().Get("X-Request-Id"); requestID == "" {
		t.Fatal("response request ID is empty")
	}
}

func TestStartServerWrapsStartupError(t *testing.T) {
	testutil.DiscardLogs(t)
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Logger.SetOutput(io.Discard)

	err := server.StartServer(e, server.Config{Address: "127.0.0.1:not-a-port"})

	if err == nil {
		t.Fatal("StartServer() error = nil, want an error")
	}
	if !strings.Contains(err.Error(), "start HTTP server") {
		t.Fatalf("StartServer() error = %q, want startup context", err)
	}
	if errors.Unwrap(err) == nil {
		t.Fatal("StartServer() error does not retain its cause")
	}
}

type validationRequest struct {
	Details validationDetails `validate:"required"`
}

type validationDetails struct {
	Value string
}
