package testutil

import (
	"golang.org/x/crypto/bcrypt"
	"testing"
)

// GeneratePasswordHash 生成密码的hash值并验证
func GeneratePasswordHash(t *testing.T, password string) string {
	// 生成hash
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Failed to generate hash: %v", err)
	}

	// 验证生成的hash
	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err != nil {
		t.Fatalf("Failed to verify generated hash: %v", err)
	}

	return string(hash)
}
