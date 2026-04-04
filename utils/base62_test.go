package utils

import "testing"

func TestEncode(t *testing.T) {
	tests := []struct {
		name string
		num  int64
		want string
	}{
		{"Encode 0", 0, "0"},
		{"Encode 1", 1, "1"},
		{"Encode 61", 61, "Z"},
		{"Encode 62", 62, "10"},
		{"Encode 12345", 12345, "3d7"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Encode(tt.num); got != tt.want {
				t.Errorf("Encode(%d) = %s, want %s", tt.num, got, tt.want)
			}
		})
	}
}

func TestEncodeWithPadding(t *testing.T) {
	id := int64(12345)
	pad := 10
	result := EncodeWithPadding(id, pad)

	if len(result) < pad {
		t.Errorf("Expected result to be at least %d characters, got %d", pad, len(result))
	}

	if !startsWith(result, Encode(id)) {
		t.Errorf("Expected result to start with Base62 encoding of id, got %s", result)
	}
}
