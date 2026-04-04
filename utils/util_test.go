package utils

import "testing"

func TestIsValidURL(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected bool
	}{
		{"valid http", "http://example.com", true},
		{"valid https", "https://example.com", true},
		{"missing scheme", "example.com", false},
		{"empty string", "", false},
		{"non http/https scheme", "ftp://example.com", false},
		{"no host", "http://", false},
		{"malformed url 1", "http:/example.com", false},
		{"malformed url 2", "http:/examplecom", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidURL(tt.value)

			if result != tt.expected {
				t.Errorf("IsValidURL(%q) = %v, want %v", tt.value, result, tt.expected)
			}
		})
	}
}
