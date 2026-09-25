package uid

import (
	"strings"
	"sync"
	"time"
)

// Package-level state for ShortID. Guarded by shortIDMutex to keep
// ID generation collision-free under concurrent use.
var (
	shortIDMutex     sync.Mutex
	lastShortIDStamp int64
	shortIDCounter   int64
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
