package blogstore

import (
	"strings"

	"github.com/dracory/uid"
)

// GenerateShortID generates a new 9-character shortened ID using TimestampMicro
func GenerateShortID() string {
	timestampMicro := uid.TimestampMicro()
	shortened, _ := uid.ShortenCrockford(timestampMicro)
	return strings.ToLower(shortened)
}

// ShortenID shortens any numeric ID string using Crockford Base32
func ShortenID(id string) string {
	if id == "" {
		return ""
	}

	// If already short (9-11 chars), return as-is
	if len(id) <= 11 {
		return strings.ToLower(id)
	}

	// Shorten long IDs
	shortened, err := uid.ShortenCrockford(id)
	if err != nil {
		return id
	}

	return strings.ToLower(shortened)
}

// UnshortenID attempts to unshorten a Crockford Base32 ID
func UnshortenID(shortID string) (string, error) {
	return uid.UnshortenCrockford(strings.ToUpper(shortID))
}

// IsShortID checks if an ID appears to be shortened (9-21 chars, alphanumeric)
func IsShortID(id string) bool {
	if len(id) < 9 || len(id) > 21 {
		return false
	}

	for _, c := range id {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')) {
			return false
		}
	}

	return true
}

// NormalizeID normalizes an ID for lookup (lowercase)
func NormalizeID(id string) string {
	return strings.ToLower(strings.TrimSpace(id))
}
