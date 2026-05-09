package blogstore

import (
	"strings"
	"testing"
)

func TestGenerateShortID(t *testing.T) {
	id := GenerateShortID()

	if id == "" {
		t.Errorf("Expected non-empty ID")
	}
	if len(id) > 11 {
		t.Errorf("Expected ID length <= 11, got %d", len(id))
	}
	if len(id) < 9 {
		t.Errorf("Expected ID length >= 9, got %d", len(id))
	}
	if id != strings.ToLower(id) {
		t.Errorf("ID should be lowercase")
	}
}

func TestGenerateShortID_Uniqueness(t *testing.T) {
	ids := make(map[string]bool)

	for i := 0; i < 100; i++ {
		id := GenerateShortID()
		if ids[id] {
			t.Errorf("Generated duplicate ID: %s", id)
		}
		ids[id] = true
	}
}

func TestShortenID(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		maxLen int
		minLen int
	}{
		{"Empty", "", 0, 0},
		{"Short ID already", "abc123def", 11, 9},
		{"Long HumanUid", "20260116055547619570214289007495", 21, 21},
		{"TimestampMicro", "1768543534819239", 11, 11},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ShortenID(tt.input)

			if tt.input == "" {
				if result != "" {
					t.Errorf("Expected empty result for empty input")
				}
			} else {
				if len(result) > tt.maxLen {
					t.Errorf("Expected result length <= %d, got %d", tt.maxLen, len(result))
				}
				if tt.minLen > 0 && len(result) < tt.minLen {
					t.Errorf("Expected result length >= %d, got %d", tt.minLen, len(result))
				}
				if result != strings.ToLower(result) {
					t.Errorf("Result should be lowercase")
				}
			}
		})
	}
}

func TestShortenID_Idempotent(t *testing.T) {
	longID := "20260116055547619570214289007495"
	shortened1 := ShortenID(longID)
	shortened2 := ShortenID(shortened1)

	if shortened1 != shortened2 {
		t.Errorf("Shortening an already short ID should return the same ID")
	}
}

func TestUnshortenID(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		shouldErr bool
	}{
		{"Valid short ID", "abc123def", false},
		{"Lowercase Crockford", "fzdzq6p7thbcf1bdjfrw7", false},
		{"Uppercase Crockford", "FZDZQ6P7THBCF1BDJFRW7", false},
		{"Mixed case", "FzDzQ6p7ThBcF1bDjFrW7", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := UnshortenID(tt.input)

			if tt.shouldErr {
				if err == nil {
					t.Errorf("Expected error")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if result == "" {
					t.Errorf("Expected non-empty result")
				}
			}
		})
	}
}

func TestIsShortID(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"9-char short ID", "abc123def", true},
		{"11-char short ID", "abc123defgh", true},
		{"21-char short ID", "fzdzq6p7thbcf1bdjfrw7", true},
		{"32-char long ID", "20260116055547619570214289007495", false},
		{"Too short", "abc", false},
		{"Too long", "abcdefghijklmnopqrstuvwxyz", false},
		{"With special chars", "abc-123", false},
		{"With spaces", "abc 123", false},
		{"Empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsShortID(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestNormalizeID(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Lowercase", "abc123", "abc123"},
		{"Uppercase", "ABC123", "abc123"},
		{"Mixed case", "AbC123", "abc123"},
		{"With leading spaces", "  abc123", "abc123"},
		{"With trailing spaces", "abc123  ", "abc123"},
		{"With both spaces", "  abc123  ", "abc123"},
		{"Empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizeID(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestShortenAndUnshortenRoundTrip(t *testing.T) {
	originalIDs := []string{
		"20260116055547619570214289007495",
		"1768543534819239",
	}

	for _, original := range originalIDs {
		t.Run(original, func(t *testing.T) {
			shortened := ShortenID(original)
			if original == shortened {
				t.Errorf("ID should be shortened")
			}
			if len(shortened) >= len(original) {
				t.Errorf("Shortened ID should be shorter")
			}

			unshortened, err := UnshortenID(shortened)
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if original != unshortened {
				t.Errorf("Unshortened ID should match original")
			}
		})
	}
}
