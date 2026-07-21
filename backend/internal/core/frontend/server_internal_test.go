package frontend

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"

	"github.com/labstack/echo/v4"
)

func TestRegisterStaticRoutesServesFrontendFiles(t *testing.T) {
	e := echo.New()
	filesystem := fstest.MapFS{
		"index.html": &fstest.MapFile{Data: []byte("api=%%API_URL%%")},
		"asset.js":   &fstest.MapFile{Data: []byte("unchanged")},
	}
	registerStaticRoutes(e, filesystem, "https://example.test/api")

	tests := []struct {
		name     string
		path     string
		wantBody string
	}{
		{name: "patched SPA fallback", path: "/diary-entries/42", wantBody: "api=https://example.test/api"},
		{name: "unchanged asset", path: "/asset.js", wantBody: "unchanged"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, tt.path, nil)
			response := httptest.NewRecorder()
			e.ServeHTTP(response, request)

			if response.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
			}
			if response.Body.String() != tt.wantBody {
				t.Fatalf("body = %q, want %q", response.Body.String(), tt.wantBody)
			}
		})
	}
}
