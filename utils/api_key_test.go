package utils

import "testing"

func TestGenerateAPIKey(t *testing.T) {
	key, err := GenerateAPIKey()
	if err != nil {
		t.Fatalf("GenerateAPIKey() error: %v", err)
	}

	if len(key) < 20 {
		t.Errorf("Expected API key to be at least 20 characters, got %d", len(key))
	}

	if !startsWith(key, "sk_live_") {
		t.Errorf("Expected API key to start with 'sk_live_', got %s", key)
	}
}

func startsWith(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

func TestGenerateMagicToken(t *testing.T) {
	token, err := GenerateMagicToken()
	if err != nil {
		t.Fatalf("GenerateMagicToken() error: %v", err)
	}

	if len(token) < 20 {
		t.Errorf("Expected magic token to be at least 20 characters, got %d", len(token))
	}
}

func TestHash(t *testing.T) {
	value := "test_value"
	hashed := Hash(value)

	if hashed == "" {
		t.Error("Expected hash to be non-empty")
	}

	if hashed == value {
		t.Error("Expected hash to be different from original value")
	}
}
