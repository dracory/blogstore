package blogstore

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateShortID(t *testing.T) {
	id := GenerateShortID()

	assert.NotEmpty(t, id)
	assert.LessOrEqual(t, len(id), 11)
	assert.GreaterOrEqual(t, len(id), 9)
	assert.Equal(t, id, strings.ToLower(id), "ID should be lowercase")
}

func TestGenerateShortID_Uniqueness(t *testing.T) {
	ids := make(map[string]bool)

	for i := 0; i < 100; i++ {
		id := GenerateShortID()
		assert.False(t, ids[id], "Generated duplicate ID: %s", id)
		ids[id] = true
	}
}

func TestShortenID(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		maxLen   int
		minLen   int
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
				assert.Empty(t, result)
			} else {
				assert.LessOrEqual(t, len(result), tt.maxLen)
				if tt.minLen > 0 {
					assert.GreaterOrEqual(t, len(result), tt.minLen)
				}
				assert.Equal(t, result, strings.ToLower(result), "Result should be lowercase")
			}
		})
	}
}

func TestShortenID_Idempotent(t *testing.T) {
	longID := "20260116055547619570214289007495"
	shortened1 := ShortenID(longID)
	shortened2 := ShortenID(shortened1)

	assert.Equal(t, shortened1, shortened2, "Shortening an already short ID should return the same ID")
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
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, result)
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
			assert.Equal(t, tt.expected, result)
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
			assert.Equal(t, tt.expected, result)
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
			assert.NotEqual(t, original, shortened, "ID should be shortened")
			assert.Less(t, len(shortened), len(original), "Shortened ID should be shorter")

			unshortened, err := UnshortenID(shortened)
			assert.NoError(t, err)
			assert.Equal(t, original, unshortened, "Unshortened ID should match original")
		})
	}
}
