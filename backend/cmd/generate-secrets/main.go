package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
)

const (
	secretSize       = 32
	jwtSecretKey     = "JWT_SECRET"
	adminPasswordKey = "ADMIN_PASSWORD"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: generate-secrets <env-file>")
		os.Exit(2)
	}

	if err := run(os.Args[1]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(path string) error {
	data, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("read %s: %w", path, err)
	}

	if hasVariable(data, jwtSecretKey) {
		return fmt.Errorf("%s already exists in %s", jwtSecretKey, path)
	}
	if hasVariable(data, adminPasswordKey) {
		return fmt.Errorf("%s already exists in %s", adminPasswordKey, path)
	}

	jwtSecret, err := generateSecret()
	if err != nil {
		return fmt.Errorf("generate %s: %w", jwtSecretKey, err)
	}

	adminPassword, err := generateSecret()
	if err != nil {
		return fmt.Errorf("generate %s: %w", adminPasswordKey, err)
	}

	separator := ""
	if len(data) > 0 && data[len(data)-1] != '\n' {
		separator = "\n"
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o600)
	if err != nil {
		return fmt.Errorf("open %s: %w", path, err)
	}

	_, writeErr := fmt.Fprintf(file, "%s%s=%s\n%s=%s\n",
		separator,
		jwtSecretKey,
		jwtSecret,
		adminPasswordKey,
		adminPassword,
	)
	closeErr := file.Close()

	if writeErr != nil {
		return fmt.Errorf("write %s: %w", path, writeErr)
	}

	if closeErr != nil {
		return fmt.Errorf("close %s: %w", path, closeErr)
	}

	return nil
}

func hasVariable(data []byte, name string) bool {
	for line := range strings.SplitSeq(string(data), "\n") {
		line = strings.TrimSpace(line)
		if after, ok := strings.CutPrefix(line, "export "); ok {
			line = strings.TrimSpace(after)
		}
		key, _, found := strings.Cut(line, "=")
		if found && strings.TrimSpace(key) == name {
			return true
		}
	}
	return false
}

func generateSecret() (string, error) {
	value := make([]byte, secretSize)
	if _, err := rand.Read(value); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(value), nil
}
