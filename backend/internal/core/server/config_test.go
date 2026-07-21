package server_test

import (
	"errors"
	"strconv"
	"strings"
	"testing"

	"github.com/azaviyalov/null3/backend/internal/core/server"
)

func TestGetConfig(t *testing.T) {
	tests := []struct {
		name        string
		address     string
		enableCORS  string
		frontendURL string
		want        server.Config
	}{
		{
			name: "defaults",
			want: server.Config{
				Address:     "localhost:8080",
				FrontendURL: "http://localhost:4200",
			},
		},
		{
			name:       "custom address",
			address:    "127.0.0.1:9090",
			enableCORS: "false",
			want: server.Config{
				Address:     "127.0.0.1:9090",
				FrontendURL: "http://localhost:4200",
			},
		},
		{
			name:       "CORS enabled with default frontend URL",
			enableCORS: "true",
			want: server.Config{
				Address:     "localhost:8080",
				EnableCORS:  true,
				FrontendURL: "http://localhost:4200",
			},
		},
		{
			name:        "CORS enabled with custom frontend URL",
			enableCORS:  "true",
			frontendURL: "https://example.test",
			want: server.Config{
				Address:     "localhost:8080",
				EnableCORS:  true,
				FrontendURL: "https://example.test",
			},
		},
		{
			name:        "custom frontend URL without CORS",
			enableCORS:  "false",
			frontendURL: "https://example.test",
			want: server.Config{
				Address:     "localhost:8080",
				FrontendURL: "https://example.test",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("ADDRESS", tt.address)
			t.Setenv("ENABLE_CORS", tt.enableCORS)
			t.Setenv("FRONTEND_URL", tt.frontendURL)

			got, err := server.GetConfig()

			if err != nil {
				t.Fatalf("GetConfig() error = %v", err)
			}
			if got != tt.want {
				t.Fatalf("GetConfig() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestGetConfigRejectsInvalidEnableCORS(t *testing.T) {
	t.Setenv("ADDRESS", "127.0.0.1:9090")
	t.Setenv("ENABLE_CORS", "sometimes")
	t.Setenv("FRONTEND_URL", "https://example.test")

	_, err := server.GetConfig()

	if err == nil {
		t.Fatal("GetConfig() error = nil, want an error")
	}
	if !strings.Contains(err.Error(), "parse ENABLE_CORS") {
		t.Fatalf("GetConfig() error = %q, want variable context", err)
	}
	var parseError *strconv.NumError
	if !errors.As(err, &parseError) {
		t.Fatalf("GetConfig() error = %v, want strconv.NumError cause", err)
	}
}
