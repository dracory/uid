package uid

import (
	"strings"
	"sync/atomic"
	"time"
)

// Package-level state for the short ID generators.
var (
	// shortIDAtomicState packs (microsecond timestamp << 4) | 4-bit counter.
	shortIDAtomicState atomic.Uint64

	// nowMicro is the clock used by the generators. It is a variable so tests
	// can inject controlled timestamps (e.g. backwards clock steps).
	// Tests that override it MUST NOT use t.Parallel() and must restore it
	// (and shortIDAtomicState) afterwards.
	nowMicro = func() int64 { return time.Now().UnixMicro() }
)

const (
	safeShortAlphabet       = "0123456789bcdfrt"
	hexShortAlphabet        = "0123456789abcdef"
	crockfordLowerAlphabet  = "0123456789abcdefghjkmnpqrstvwxyz"
	shortIDLength           = 11 // base32 digits needed for the 55-bit composite
	shortIDHexEncodedLength = 14 // base16 digits needed for the 55-bit composite
)

// ShortID creates a new 11-character lowercase short ID.
// It encodes the current microsecond timestamp in Crockford Base32.
// Thread-safe via a lock-free CAS loop; survives backwards clock steps.
//
// Business Logic:
//   - Uses time.Now().UnixMicro() as the uniqueness base.
//   - Maintains a package-level atomic (timestamp << 4) | counter state.
//   - Counter increments while the timestamp does not advance; when it would
//     exceed 15, the timestamp is bumped by one microsecond instead of
//     sleeping, so callers never stall and IDs never regress.
//   - Packs timestamp and counter into a single int64: (ts << 4) | counter.
//   - Encodes the composite value with a lowercase Crockford Base32 alphabet.
//
// Notes:
//   - Uniqueness is per-process only: two processes generating IDs in the same
//     microsecond can collide. The IDs are also predictable and reveal their
//     creation time, so they are suitable for sortable database keys but NOT
//     for security tokens or unguessable public identifiers.
//
// Returns:
//   - An 11-character lowercase string (e.g. "sa4rc789wxg").
//
// Example:
//   - id := uid.ShortID() // "sa4rc789wxg"
func ShortID() string {
	return encodeInt64(nextShortIDValue(), crockfordLowerAlphabet)
}

// encodeInt64 encodes a non-negative int64 in the given alphabet without
// allocations beyond the final string. The alphabet is expected to be
// ascending so that lexicographic order matches numeric order.
func encodeInt64(n int64, alphabet string) string {
	if n == 0 {
		return string(alphabet[0])
	}
	base := int64(len(alphabet))
	var buf [64]byte // worst case: 55 bits at base 2
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = alphabet[n%base]
		n /= base
	}
	return string(buf[i:])
}

// encodeHexFixed encodes val into exactly width lowercase hex-style digits,
// most-significant digit first. It only works for 16-character alphabets
// (it shifts by 4 bits per digit) and width must be <= 14.
func encodeHexFixed(val int64, alphabet string, width int) string {
	var buf [shortIDHexEncodedLength]byte
	for i := width - 1; i >= 0; i-- {
		buf[i] = alphabet[val&0x0f]
		val >>= 4
	}
	return string(buf[:width])
}

// NormalizeID normalizes an ID to lowercase and trims whitespace for consistent lookups.
//
// Parameters:
//   - id: the ID string to normalize.
//
// Returns:
//   - The trimmed, lowercased ID.
//
// Example:
//   - norm := uid.NormalizeID("  SA4RC789WXG  ") // "sa4rc789wxg"
func NormalizeID(id string) string {
	return strings.ToLower(strings.TrimSpace(id))
}

// IsShortID checks whether the given string matches the format produced by
// ShortID: the correct length and only lowercase Crockford Base32 characters.
// It does NOT cover the SafeShortID/HexShortID/ShortIDCustom formats.
//
// Parameters:
//   - id: the ID string to check.
//
// Returns:
//   - true if id looks like a ShortID() output; false otherwise.
//
// Example:
//   - ok := uid.IsShortID("sa4rc789wxg") // true
func IsShortID(id string) bool {
	// 11 chars today; 12 once the composite exceeds 2^55 (from ~2041).
	if len(id) != shortIDLength && len(id) != shortIDLength+1 {
		return false
	}
	for i := 0; i < len(id); i++ {
		if strings.IndexByte(crockfordLowerAlphabet, id[i]) < 0 {
			return false
		}
	}
	return true
}

// nextShortIDValue atomically generates a unique 55-bit composite value
// combining a microsecond timestamp (shifted left 4 bits) and a 4-bit
// rolling counter (0-15). Uses a CAS loop for thread safety. If the clock
// moves backwards the timestamp is held at its last value and the counter
// keeps incrementing; on counter overflow the timestamp is bumped forward
// by one microsecond rather than sleeping.
func nextShortIDValue() int64 {
	for {
		now := nowMicro()
		oldState := shortIDAtomicState.Load()
		oldTS := int64(oldState >> 4)
		oldCnt := int64(oldState & 0x0F)

		var newTS, newCnt int64
		if now > oldTS {
			newTS = now
			newCnt = 0
		} else if oldCnt < 0x0F {
			newTS = oldTS
			newCnt = oldCnt + 1
		} else {
			newTS = oldTS + 1
			newCnt = 0
		}

		newState := uint64(newTS)<<4 | uint64(newCnt)
		if shortIDAtomicState.CompareAndSwap(oldState, newState) {
			return int64(newState)
		}
	}
}

// SafeShortID generates a fixed 14-character lowercase short ID using a vowel-free Base 16 alphabet.
// Core logic: microsecond timestamp shifted left by 4 bits, combined with a 4-bit
// rolling counter, produced by nextShortIDValue.
// Alphabet (Base 16, vowel-free): "0123456789bcdfrt".
// Every letter of standard hex is shifted to a non-vowel ('a'->'b' ... 'e'->'r',
// 'f'->'t') to avoid generating accidental English words; this greatly reduces
// but does not eliminate them (leet-speak like "b00b" is still reachable).
// The alphabet is in ascending byte order, so generated IDs sort chronologically.
//
// Returns:
//   - A 14-character lowercase string (e.g. "05fr88tc139r8t").
func SafeShortID() string {
	return encodeHexFixed(nextShortIDValue(), safeShortAlphabet, shortIDHexEncodedLength)
}

// HexShortID generates a standard 14-character lowercase hexadecimal short ID.
// Uses the same timestamp + 4-bit rolling counter logic as SafeShortID,
// formatted with the standard lowercase hex alphabet "0123456789abcdef".
//
// Returns:
//   - A 14-character lowercase hexadecimal string.
func HexShortID() string {
	return encodeHexFixed(nextShortIDValue(), hexShortAlphabet, shortIDHexEncodedLength)
}

// ShortIDCustom generates a time-ordered short ID using a user-provided custom alphabet string.
// Encodes the combined (timestamp << 4) | counter value using modulo arithmetic
// against the length of the provided alphabet.
//
// Parameters:
//   - alphabet: a string of at least 2 unique single-byte (ASCII) characters.
//     If empty, defaults to "0123456789abcdef". For IDs that sort
//     chronologically the characters must be in ascending byte order.
//     Panics on an invalid alphabet.
//
// Returns:
//   - A time-ordered short ID encoded with the custom alphabet.
func ShortIDCustom(alphabet string) string {
	if alphabet == "" {
		alphabet = hexShortAlphabet
	}

	if len(alphabet) < 2 {
		panic("uid: ShortIDCustom requires an alphabet of at least 2 characters")
	}
	seen := make(map[byte]struct{}, len(alphabet))
	for i := 0; i < len(alphabet); i++ {
		c := alphabet[i]
		if c > 0x7f {
			panic("uid: ShortIDCustom alphabet must contain single-byte ASCII characters only")
		}
		if _, ok := seen[c]; ok {
			panic("uid: ShortIDCustom alphabet contains duplicate characters")
		}
		seen[c] = struct{}{}
	}

	return encodeInt64(nextShortIDValue(), alphabet)
}
