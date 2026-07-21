package frontend_test

import (
	"errors"
	"strconv"
	"strings"
	"testing"

	"github.com/azaviyalov/null3/backend/internal/core/frontend"
)

func TestGetConfig(t *testing.T) {
	tests := []struct {
		name               string
		enableFrontendDist string
		apiURL             string
		want               frontend.Config
	}{
		{
			name: "defaults",
			want: frontend.Config{APIURL: "http://localhost:8080/api"},
		},
		{
			name:               "enabled with default API URL",
			enableFrontendDist: "true",
			want: frontend.Config{
				EnableFrontendDist: true,
				APIURL:             "http://localhost:8080/api",
			},
		},
		{
			name:               "enabled with custom API URL",
			enableFrontendDist: "true",
			apiURL:             "https://example.test/api",
			want: frontend.Config{
				EnableFrontendDist: true,
				APIURL:             "https://example.test/api",
			},
		},
		{
			name:               "disabled ignores custom API URL",
			enableFrontendDist: "false",
			apiURL:             "https://example.test/api",
			want:               frontend.Config{APIURL: "http://localhost:8080/api"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("ENABLE_FRONTEND_DIST", tt.enableFrontendDist)
			t.Setenv("API_URL", tt.apiURL)

			got, err := frontend.GetConfig()

			if err != nil {
				t.Fatalf("GetConfig() error = %v", err)
			}
			if got != tt.want {
				t.Fatalf("GetConfig() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestGetConfigRejectsInvalidEnableFrontendDist(t *testing.T) {
	t.Setenv("ENABLE_FRONTEND_DIST", "sometimes")
	t.Setenv("API_URL", "https://example.test/api")

	got, err := frontend.GetConfig()

	if err == nil {
		t.Fatal("GetConfig() error = nil, want an error")
	}
	if !strings.Contains(err.Error(), "parse ENABLE_FRONTEND_DIST") {
		t.Fatalf("GetConfig() error = %q, want variable context", err)
	}
	var parseError *strconv.NumError
	if !errors.As(err, &parseError) {
		t.Fatalf("GetConfig() error = %v, want strconv.NumError cause", err)
	}
	want := frontend.Config{APIURL: "http://localhost:8080/api"}
	if got != want {
		t.Fatalf("GetConfig() = %+v, want defaults %+v", got, want)
	}
}
