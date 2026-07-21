package account_test

import (
	"strings"
	"testing"
	"time"

	"github.com/azaviyalov/null3/backend/internal/domain/account"
)

func TestGetConfig(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    time.Duration
		wantErr string
	}{
		{name: "default", want: time.Hour},
		{name: "environment override", value: "45m", want: 45 * time.Minute},
		{name: "invalid duration", value: "later", wantErr: "parse PASSWORD_RESET_TOKEN_EXPIRATION"},
		{name: "zero duration", value: "0s", wantErr: "must be a positive duration"},
		{name: "negative duration", value: "-1s", wantErr: "must be a positive duration"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("PASSWORD_RESET_TOKEN_EXPIRATION", tt.value)

			config, err := account.GetConfig()

			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("GetConfig() error = %v, want error containing %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("GetConfig() error = %v", err)
			}
			if config.PasswordResetTokenExpiration != tt.want {
				t.Errorf("PasswordResetTokenExpiration = %v, want %v", config.PasswordResetTokenExpiration, tt.want)
			}
		})
	}
}
