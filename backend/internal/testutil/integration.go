package testutil

import "testing"

func SkipIntegration(t testing.TB) {
	t.Helper()
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
}
