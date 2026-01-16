package uid

import (
	"strings"
)

// Short converts a numeric ID string to a short, human-friendly representation using Base58.
//
// Parameters:
// - id: A numeric string ID (may contain hyphens which are ignored).
//
// Returns:
// - A shortened string representation, or an error if input is invalid.
func Short(id string) (string, error) {
	clean := strings.ReplaceAll(id, "-", "")
	return ShortenBase58(clean)
}

// Unshort converts a shortened string (Base58) back to its original numeric ID.
//
// Parameters:
// - s: A shortened string ID (Base58).
//
// Returns:
// - The original numeric string ID, or an error if decoding fails.
func Unshort(s string) (string, error) {
	// Remove hyphens which might have been added for readability
	clean := strings.ReplaceAll(s, "-", "")
	if isAllDigits(clean) {
		return clean, nil
	}

	// Primary: Base58 (matching Short)
	return UnshortenBase58(s)
}

func isAllDigits(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}
