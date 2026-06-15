package blogstore

import (
	"errors"
	"math/big"
	"strings"
	"time"

	neatuid "github.com/dracory/neat/support/uid"
)

const crockfordAlphabet = "0123456789abcdefghjkmnpqrstvwxyz"

// GenerateShortID generates a new 11-character shortened ID using TimestampMicro
func GenerateShortID() string {
	time.Sleep(1 * time.Millisecond)
	return strings.ToLower(neatuid.GenerateShortID())
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

	// Check if purely numeric - if so, always shorten
	isNumeric := true
	for _, c := range id {
		if c < '0' || c > '9' {
			isNumeric = false
			break
		}
	}
	if isNumeric {
		shortened, err := shortenCrockford(id)
		if err != nil {
			return id
		}
		return strings.ToLower(shortened)
	}

	// If already appears to be a Crockford Base32 string (alphanumeric) AND short enough, return as-is
	// This ensures idempotency for already-shortened IDs (max 21 chars for 64-bit in base32)
	isCrockford := true
	for _, c := range id {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')) {
			isCrockford = false
			break
		}
	}
	if isCrockford && len(id) <= 21 {
		return strings.ToLower(id)
	}

	// Shorten long IDs
	shortened, err := shortenCrockford(id)
	if err != nil {
		return id
	}

	return strings.ToLower(shortened)
}

// UnshortenID attempts to unshorten a Crockford Base32 ID
func UnshortenID(shortID string) (string, error) {
	return unshortenCrockford(strings.ToUpper(shortID))
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

// shortenCrockford encodes a numeric string to Crockford Base32
func shortenCrockford(id string) (string, error) {
	num := new(big.Int)
	_, success := num.SetString(id, 10)
	if !success {
		return "", errors.New("invalid numeric string")
	}

	// Convert to base32
	var result []byte
	base := big.NewInt(32)
	zero := big.NewInt(0)

	for num.Cmp(zero) > 0 {
		rem := new(big.Int)
		num.DivMod(num, base, rem)
		result = append(result, crockfordAlphabet[rem.Int64()])
	}

	// Reverse the result
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return string(result), nil
}

// unshortenCrockford decodes a Crockford Base32 string to numeric
func unshortenCrockford(shortID string) (string, error) {
	// Create reverse lookup map
	decodeMap := make(map[byte]int64)
	for i, c := range crockfordAlphabet {
		decodeMap[byte(c)] = int64(i)
	}

	num := big.NewInt(0)
	base := big.NewInt(32)

	for _, c := range shortID {
		val, ok := decodeMap[byte(strings.ToLower(string(c))[0])]
		if !ok {
			return "", errors.New("invalid Crockford Base32 character")
		}
		num.Mul(num, base)
		num.Add(num, big.NewInt(val))
	}

	return num.String(), nil
}
