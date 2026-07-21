package frontend_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/azaviyalov/null3/backend/internal/core/frontend"
	"github.com/labstack/echo/v4"
)

func TestRegisterRoutesLeavesApplicationRoutesAvailable(t *testing.T) {
	tests := []struct {
		name    string
		enabled bool
	}{
		{name: "frontend disabled", enabled: false},
		{name: "frontend enabled", enabled: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			e.GET("/health", func(c echo.Context) error {
				return c.NoContent(http.StatusNoContent)
			})
			frontend.RegisterRoutes(e, frontend.Config{EnableFrontendDist: tt.enabled})

			request := httptest.NewRequest(http.MethodGet, "/health", nil)
			response := httptest.NewRecorder()
			e.ServeHTTP(response, request)

			if response.Code != http.StatusNoContent {
				t.Fatalf("status = %d, want %d", response.Code, http.StatusNoContent)
			}
		})
	}
}
