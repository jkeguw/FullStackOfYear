package testutil

import (
	"testing"
)

func TestGeneratePassword(t *testing.T) {
	password := "password123"
	hash := GeneratePasswordHash(t, password)

	t.Logf("Password: %s", password)
	t.Logf("Generated Hash: %s", hash)
}
