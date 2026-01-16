package uid

import (
	"errors"
	"math/big"
	"strings"
)

var (
	ErrInvalidByte  = errors.New("character not in alphabet")
	ErrInvalidInput = errors.New("uid: invalid numeric input")
)

const (
	alphabet16        = "0123456789abcdef"
	alphabet32        = "abcdefghijklmnopqrstuvwxyz234567"
	alphabetCrockford = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"
	alphabet64        = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"
	alphabet36        = "0123456789abcdefghijklmnopqrstuvwxyz"
	alphabet58        = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	alphabet62        = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	alphabetZ32       = "ybndrfg8ejkmcpqxot1uwisza345h769"
)

// ShortenBase16 converts a numeric string ID to a base16 (hexadecimal) string.
//
// Benefits:
// - Standard hexadecimal representation.
// - Compatible with most systems and easily parsed.
//
// Drawbacks:
// - Least compact representation; results in the longest strings.
//
// Example: ShortenBase16("1234567890") -> "499602d2"
//
// Parameters:
// - id: A numeric string represention of a positive integer.
//
// Returns:
// - A lowercase hexadecimal string.
func ShortenBase16(id string) (string, error) {
	return encodeBase(id, alphabet16)
}

// UnshortenBase16 converts a base16 (hexadecimal) string back to its original numeric string ID.
//
// Example: UnshortenBase16("499602d2") -> "1234567890"
//
// Parameters:
// - s: A hexadecimal string.
//
// Returns:
// - The original numeric string ID, or an error if invalid characters are encountered.
func UnshortenBase16(s string) (string, error) {
	return decodeBase(s, alphabet16)
}

// ShortenBase32 converts a numeric string ID to a base32 string (RFC 4648).
//
// Benefits:
// - Case-insensitive and safe for filesystems.
// - Compact representation using only 32 characters.
//
// Drawbacks:
// - Less compact than Base58/Base62.
// - Includes 'l' and 'i' which can be visually similar in some fonts.
//
// Example: ShortenBase32("1234567890") -> "gnm6fsq"
//
// Parameters:
// - id: A numeric string represention of a positive integer.
//
// Returns:
// - A lowercase base32 string.
func ShortenBase32(id string) (string, error) {
	return encodeBase(id, alphabet32)
}

// UnshortenBase32 converts a base32 string back to its original numeric string ID.
//
// Example: UnshortenBase32("gnm6fsq") -> "1234567890"
//
// Parameters:
// - s: A base32 string (case-insensitive).
//
// Returns:
// - The original numeric string ID, or an error if invalid characters are encountered.
func UnshortenBase32(s string) (string, error) {
	return decodeBase(strings.ToLower(s), alphabet32)
}

// ShortenCrockford converts a numeric string ID to a Crockford Base32 string.
//
// Benefits:
// - Highly human-readable and designed to avoid common transcription errors.
// - Excludes ambiguous characters like I, L, O, and U.
//
// Drawbacks:
// - Less compact than Base62/Base64.
//
// Example: ShortenCrockford("1234567890") -> "14MM0M"
//
// Parameters:
// - id: A numeric string represention of a positive integer.
//
// Returns:
// - An uppercase Crockford Base32 string.
func ShortenCrockford(id string) (string, error) {
	return encodeBase(id, alphabetCrockford)
}

// UnshortenCrockford converts a Crockford Base32 string back to its original numeric string ID.
//
// Features:
// - Supports normalization: converts to uppercase, removes hyphens.
// - Maps visually similar characters (O -> 0, I/L -> 1) to reduce errors.
//
// Example: UnshortenCrockford("14MM0M") -> "1234567890"
//
// Parameters:
// - s: A Crockford Base32 string.
//
// Returns:
// - The original numeric string ID, or an error if invalid characters are encountered.
func UnshortenCrockford(s string) (string, error) {
	s = strings.ToUpper(s)
	s = strings.ReplaceAll(s, "-", "")
	s = strings.ReplaceAll(s, "O", "0")
	s = strings.ReplaceAll(s, "I", "1")
	s = strings.ReplaceAll(s, "L", "1")
	return decodeBase(s, alphabetCrockford)
}

// ShortenBase64 converts a numeric string ID to a URL-safe base64 string.
//
// Benefits:
// - Extremely compact representation.
// - Safe for use in URLs (uses '-' and '_' instead of '+' and '/').
//
// Drawbacks:
// - Case-sensitive, which can lead to errors if manually typed.
// - Contains visually similar characters (e.g., 'O' and '0', 'I' and 'l').
// - Includes '-' and '_', which can interfere with double-click selection in some UIs.
// - Not recommended for human-facing IDs; consider Base58 or Crockford instead.
//
// Example: ShortenBase64("1234567890") -> "BJlg8u"
//
// Parameters:
// - id: A numeric string represention of a positive integer.
//
// Returns:
// - A URL-safe base64 string.
func ShortenBase64(id string) (string, error) {
	return encodeBase(id, alphabet64)
}

// UnshortenBase64 converts a URL-safe base64 string back to its original numeric string ID.
//
// Example: UnshortenBase64("BJlg8u") -> "1234567890"
//
// Parameters:
// - s: A URL-safe base64 string.
//
// Returns:
// - The original numeric string ID, or an error if invalid characters are encountered.
func UnshortenBase64(s string) (string, error) {
	return decodeBase(s, alphabet64)
}

// ShortenBase36 converts a numeric string ID to a base36 string.
//
// Benefits:
// - Uses all digits and lowercase letters (0-9, a-z).
// - Case-insensitive and widely used for compact IDs.
//
// Drawbacks:
// - Includes 'l', '1', 'o', and '0', which can be visually confusing.
//
// Example: ShortenBase36("1234567890") -> "kf12oi"
//
// Parameters:
// - id: A numeric string represention of a positive integer.
//
// Returns:
// - A lowercase base36 string.
func ShortenBase36(id string) (string, error) {
	return encodeBase(id, alphabet36)
}

// UnshortenBase36 converts a base36 string back to its original numeric string ID.
//
// Example: UnshortenBase36("kf12oi") -> "1234567890"
//
// Parameters:
// - s: A base36 string (case-insensitive).
//
// Returns:
// - The original numeric string ID, or an error if invalid characters are encountered.
func UnshortenBase36(s string) (string, error) {
	return decodeBase(strings.ToLower(s), alphabet36)
}

// ShortenBase58 converts a numeric string ID to a base58 string.
//
// Benefits:
// - Human-friendly: avoids ambiguous characters like 0, O, I, and l.
// - No non-alphanumeric characters, making it safe for double-clicking to select.
//
// Drawbacks:
// - Case-sensitive; requires exact casing for decoding.
//
// Example: ShortenBase58("1234567890") -> "2V6G2p"
//
// Parameters:
// - id: A numeric string represention of a positive integer.
//
// Returns:
// - A base58 string (mixed case).
func ShortenBase58(id string) (string, error) {
	return encodeBase(id, alphabet58)
}

// UnshortenBase58 converts a base58 string back to its original numeric string ID.
//
// Example: UnshortenBase58("2V6G2p") -> "1234567890"
//
// Parameters:
// - s: A base58 string.
//
// Returns:
// - The original numeric string ID, or an error if invalid characters are encountered.
func UnshortenBase58(s string) (string, error) {
	return decodeBase(s, alphabet58)
}

// ShortenBase62 converts a numeric string ID to a base62 string.
//
// Benefits:
// - Maximum density using all standard alphanumeric characters (0-9, A-Z, a-z).
// - Very compact and widely supported by URL shorteners.
//
// Drawbacks:
// - Case-sensitive, prone to manual entry errors.
// - Includes visually similar characters (e.g., 'I', 'l', '1', 'O', '0').
//
// Example: ShortenBase62("1234567890") -> "1LY7vk"
//
// Parameters:
// - id: A numeric string represention of a positive integer.
//
// Returns:
// - A base62 string (mixed case).
func ShortenBase62(id string) (string, error) {
	return encodeBase(id, alphabet62)
}

// UnshortenBase62 converts a base62 string back to its original numeric string ID.
//
// Example: UnshortenBase62("1LY7vk") -> "1234567890"
//
// Parameters:
// - s: A base62 string.
//
// Returns:
// - The original numeric string ID, or an error if invalid characters are encountered.
func UnshortenBase62(s string) (string, error) {
	return decodeBase(s, alphabet62)
}

// ShortenZBase32 converts a numeric string ID to a z-base-32 string.
//
// Benefits:
// - Optimized for human readability and memorability.
// - Alphabet is designed to put easiest-to-recognize characters in most-used positions.
//
// Drawbacks:
// - Non-standard alphabet; not as widely supported as RFC 4648 Base32.
//
// Example: ShortenZBase32("1234567890") -> "4pk9xyo"
//
// Parameters:
// - id: A numeric string represention of a positive integer.
//
// Returns:
// - A lowercase z-base-32 string.
func ShortenZBase32(id string) (string, error) {
	return encodeBase(id, alphabetZ32)
}

// UnshortenZBase32 converts a z-base-32 string back to its original numeric string ID.
//
// Example: UnshortenZBase32("4pk9xyo") -> "1234567890"
//
// Parameters:
// - s: A z-base-32 string (case-insensitive).
//
// Returns:
// - The original numeric string ID, or an error if invalid characters are encountered.
func UnshortenZBase32(s string) (string, error) {
	return decodeBase(strings.ToLower(s), alphabetZ32)
}

// encodeBase is an internal helper that converts a numeric string ID to a periodic base string
// using the provided alphabet. It handles arbitrarily large numbers.
func encodeBase(id string, alphabet string) (string, error) {
	n := new(big.Int)
	n, ok := n.SetString(id, 10)
	if !ok {
		return "", ErrInvalidInput
	}
	if n.Sign() == 0 {
		return string(alphabet[0]), nil
	}

	base := big.NewInt(int64(len(alphabet)))
	var res []byte
	mod := new(big.Int)

	// Work on a copy as DivMod modifies the receiver
	temp := new(big.Int).Set(n)

	for temp.Sign() > 0 {
		temp.DivMod(temp, base, mod)
		res = append(res, alphabet[mod.Int64()])
	}

	// Reverse res
	for i, j := 0, len(res)-1; i < j; i, j = i+1, j-1 {
		res[i], res[j] = res[j], res[i]
	}

	return string(res), nil
}

// decodeBase is an internal helper that converts a periodic base string back to its
// original numeric string ID using the provided alphabet.
func decodeBase(s string, alphabet string) (string, error) {
	if s == "" {
		return "", nil
	}
	base := big.NewInt(int64(len(alphabet)))
	res := new(big.Int)

	alphaMap := make(map[byte]int64)
	for i := 0; i < len(alphabet); i++ {
		alphaMap[alphabet[i]] = int64(i)
	}

	for i := 0; i < len(s); i++ {
		val, ok := alphaMap[s[i]]
		if !ok {
			return "", ErrInvalidByte
		}
		res.Mul(res, base)
		res.Add(res, big.NewInt(val))
	}
	return res.String(), nil
}
