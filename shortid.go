package uid

import (
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// Package-level state for ShortID functions.
var (
	shortIDMutex     sync.Mutex
	lastShortIDStamp int64
	shortIDCounter   int64

	atomicCounter uint32
)

const (
	safeShortAlphabet = "0123456789bcdfrt"
	hexShortAlphabet  = "0123456789abcdef"
)

// ShortID creates a new 11-character lowercase short ID.
// It encodes the current microsecond timestamp in Crockford Base32.
// Thread-safe via mutex to prevent duplicate IDs under concurrency.
//
// Business Logic:
//   - Uses time.Now().UnixMicro() as the uniqueness base.
//   - Maintains package-level lastShortIDStamp and shortIDCounter across calls.
//   - Counter increments when the timestamp equals the previous call;
//     resets to 0 when the timestamp changes.
//   - The counter is 4 bits (0-15), giving 16 unique IDs per timestamp tick.
//   - On platforms with coarse timer resolution (e.g. ~1ms on Windows),
//     UnixMicro can return the same value across several calls in a row.
//   - Once the counter exceeds 15 it would overflow and wrap back to 0,
//     producing the same composite value and therefore a duplicate ID.
//   - Sleeping 1ms guarantees the timestamp advances to the next tick.
//   - Packs timestamp and counter into a single int64: (ts << 4) | counter.
//     The timestamp occupies the high bits; the counter occupies the low 4 bits.
//   - Encodes the composite value using the Crockford Base32 alphabet.
//   - Converts the result to lowercase.
//
// Parameters:
//   - None
//
// Returns:
//   - An 11-character lowercase string (e.g. "sa4rc789wxg").
//
// Example:
//   - id := uid.ShortID() // "sa4rc789wxg"
func ShortID() string {
	shortIDMutex.Lock()
	defer shortIDMutex.Unlock()

	ts := time.Now().UnixMicro()

	if ts == lastShortIDStamp {
		shortIDCounter++
	} else {
		lastShortIDStamp = ts
		shortIDCounter = 0
	}

	if shortIDCounter > 15 {
		time.Sleep(1 * time.Millisecond)
		ts = time.Now().UnixMicro()
		lastShortIDStamp = ts
		shortIDCounter = 0
	}

	composite := (ts << 4) | shortIDCounter
	return strings.ToLower(encodeCrockfordInt64(composite))
}

// encodeCrockfordInt64 is a private, allocation-light fast path for encoding
// a non-negative int64 using the Crockford Base32 alphabet.
//
// Design Decision (Option a):
// This intentionally does NOT go through the shared encodeBase/math-big path
// in shorten.go: ShortID is meant to be called frequently and its
// input is always a fixed-size int64 (timestamp+counter composite), never an
// arbitrary-precision number. Using big.Int here would add unnecessary
// allocation and conversion overhead for a value that always fits in a
// machine word. It reuses the alphabetCrockford constant already defined in
// shorten.go rather than redefining a second copy.
func encodeCrockfordInt64(n int64) string {
	if n == 0 {
		return string(alphabetCrockford[0])
	}

	var buf [13]byte // max length for a base32-encoded int64
	i := len(buf) - 1
	for n > 0 {
		buf[i] = alphabetCrockford[n&0x1f]
		n >>= 5
		i--
	}
	return string(buf[i+1:])
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

// IsShortID checks whether the given string matches a known short ID length.
//
// Parameters:
//   - id: the ID string to check.
//
// Returns:
//   - true if the length matches a short ID produced by ShortID (11 characters);
//     false otherwise.
//
// Example:
//   - ok := uid.IsShortID("sa4rc789wxg") // true
func IsShortID(id string) bool {
	length := len(id)
	return length == 11
}

// SafeShortID generates a fixed 13-character lowercase short ID using a vowel-free Base 16 alphabet.
// Core logic: microsecond timestamp (time.Now().UnixMicro()) shifted left by 4 bits,
// combined with a 4-bit atomic rolling counter (atomic.AddUint32(&counter, 1) & 0x0F).
// Alphabet (Base 16, vowel-free): "0123456789bcdfrt".
// Replaces vowels ('a' and 'e' from standard hex) with hard consonants ('r' and 't')
// to ensure zero accidental English words/profanities are generated, while avoiding
// phonetic confusion and visual ambiguity.
// Performance: Uses bitwise extraction (>> 4 and & 0x0F) for speed and zero allocations.
//
// Returns:
//   - A 13-character lowercase string (e.g. "5fr88tc139r8t").
func SafeShortID() string {
	ts := time.Now().UnixMicro()
	cnt := int64(atomic.AddUint32(&atomicCounter, 1) & 0x0F)
	val := (ts << 4) | cnt

	var buf [13]byte
	for i := 12; i >= 0; i-- {
		buf[i] = safeShortAlphabet[val&0x0f]
		val >>= 4
	}
	return string(buf[:])
}

// HexShortID generates a standard 13-character lowercase hexadecimal short ID.
// Uses the same bitwise timestamp and 4-bit atomic rolling counter logic as SafeShortID,
// but formatted using the standard lowercase hex alphabet: "0123456789abcdef".
//
// Returns:
//   - A 13-character lowercase hexadecimal string.
func HexShortID() string {
	ts := time.Now().UnixMicro()
	cnt := int64(atomic.AddUint32(&atomicCounter, 1) & 0x0F)
	val := (ts << 4) | cnt

	var buf [13]byte
	for i := 12; i >= 0; i-- {
		buf[i] = hexShortAlphabet[val&0x0f]
		val >>= 4
	}
	return string(buf[:])
}

// ShortIDCustom generates a time-ordered short ID using a user-provided custom alphabet string.
// Encodes the combined (timestamp << 4) | counter value using modulo arithmetic against
// the length of the provided alphabet string.
//
// Parameters:
//   - alphabet: a string containing characters to use for encoding. If empty, defaults to "0123456789abcdef".
//
// Returns:
//   - A time-ordered short ID encoded with the custom alphabet.
func ShortIDCustom(alphabet string) string {
	if alphabet == "" {
		alphabet = hexShortAlphabet
	}

	ts := time.Now().UnixMicro()
	cnt := int64(atomic.AddUint32(&atomicCounter, 1) & 0x0F)
	val := (ts << 4) | cnt

	if val == 0 {
		return string(alphabet[0])
	}

	base := int64(len(alphabet))
	var buf [32]byte
	i := len(buf)
	for val > 0 {
		i--
		buf[i] = alphabet[val%base]
		val /= base
	}
	return string(buf[i:])
}
