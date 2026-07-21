package admin_test

import (
	"strings"
	"testing"

	"github.com/azaviyalov/null3/backend/internal/domain/admin"
)

func TestGetConfig(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{name: "password", password: "configured-admin-password"},
		{name: "missing password", wantErr: true},
		{name: "blank password", password: "  ", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("ADMIN_PASSWORD", tt.password)

			config, err := admin.GetConfig()

			if tt.wantErr {
				if err == nil || !strings.Contains(err.Error(), "ADMIN_PASSWORD must be set and non-empty") {
					t.Fatalf("GetConfig() error = %v, want required-password error", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("GetConfig() error = %v", err)
			}
			if config.Password != tt.password {
				t.Error("GetConfig() did not preserve the configured password")
			}
		})
	}
}
