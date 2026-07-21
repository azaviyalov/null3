package main

import (
	"bytes"
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestHasVariable(t *testing.T) {
	tests := []struct {
		name string
		data string
		key  string
		want bool
	}{
		{name: "plain declaration", data: "JWT_SECRET=value", key: jwtSecretKey, want: true},
		{name: "export declaration", data: "  export ADMIN_PASSWORD = value  ", key: adminPasswordKey, want: true},
		{name: "declaration among other lines", data: "FIRST=value\nJWT_SECRET=value\nLAST=value", key: jwtSecretKey, want: true},
		{name: "commented declaration", data: "# JWT_SECRET=value", key: jwtSecretKey, want: false},
		{name: "similar key", data: "JWT_SECRET_OLD=value", key: jwtSecretKey, want: false},
		{name: "line without assignment", data: "JWT_SECRET", key: jwtSecretKey, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasVariable([]byte(tt.data), tt.key); got != tt.want {
				t.Fatalf("hasVariable() = %t, want %t", got, tt.want)
			}
		})
	}
}

func TestRunCreatesSecretFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".env")

	if err := run(path); err != nil {
		t.Fatalf("run() error = %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat generated file: %v", err)
	}
	if got := info.Mode().Perm(); got != 0o600 {
		t.Errorf("generated file mode = %o, want 600", got)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read generated file: %v", err)
	}
	lines := strings.Split(strings.TrimSuffix(string(data), "\n"), "\n")
	if len(lines) != 2 {
		t.Fatalf("generated declaration count = %d, want 2", len(lines))
	}
	for i, wantKey := range []string{jwtSecretKey, adminPasswordKey} {
		key, value, found := strings.Cut(lines[i], "=")
		if !found {
			t.Fatalf("generated declaration %d has no assignment", i)
		}
		if key != wantKey {
			t.Errorf("generated key %d = %q, want %q", i, key, wantKey)
		}
		decoded, err := base64.RawURLEncoding.DecodeString(value)
		if err != nil {
			t.Fatalf("generated value for %s is not raw URL-safe base64: %v", wantKey, err)
		}
		if len(decoded) != secretSize {
			t.Errorf("decoded value size for %s = %d, want %d", wantKey, len(decoded), secretSize)
		}
	}
}

func TestRunAddsNewlineBeforeDeclarations(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(path, []byte("EXISTING=value"), 0o600); err != nil {
		t.Fatalf("write existing file: %v", err)
	}

	if err := run(path); err != nil {
		t.Fatalf("run() error = %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read generated file: %v", err)
	}
	if !bytes.HasPrefix(data, []byte("EXISTING=value\nJWT_SECRET=")) {
		t.Fatal("generated declarations do not start on a new line")
	}
}

func TestRunRejectsExistingSecretsWithoutChangingFile(t *testing.T) {
	tests := []struct {
		name string
		key  string
		data string
	}{
		{name: "JWT secret", key: jwtSecretKey, data: "JWT_SECRET=existing\n"},
		{name: "admin password", key: adminPasswordKey, data: "export ADMIN_PASSWORD=existing\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), ".env")
			before := []byte(tt.data)
			if err := os.WriteFile(path, before, 0o600); err != nil {
				t.Fatalf("write existing file: %v", err)
			}

			err := run(path)

			if err == nil {
				t.Fatal("run() error = nil, want an error")
			}
			if !strings.Contains(err.Error(), tt.key+" already exists") {
				t.Fatalf("run() error = %q, want existing-key context", err)
			}
			after, readErr := os.ReadFile(path)
			if readErr != nil {
				t.Fatalf("read unchanged file: %v", readErr)
			}
			if !bytes.Equal(after, before) {
				t.Fatal("run() changed a file containing an existing secret")
			}
		})
	}
}
