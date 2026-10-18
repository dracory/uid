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
	safeShortAlphabet       = "0123456789bcdfgh"
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
// Alphabet (Base 16, vowel-free): "0123456789bcdfgh".
// Every letter of standard hex is shifted to the next sequential consonant
// ('a'->'b', 'b'->'c', 'c'->'d', 'd'->'f', 'e'->'g', 'f'->'h') to avoid
// generating accidental English words; this greatly reduces but does not
// eliminate them (leet-speak like "b00b" is still reachable).
// The alphabet is in ascending byte order, so generated IDs sort chronologically.
//
// Returns:
//   - A 14-character lowercase string (e.g. "05fh88gd139f8h").
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

// encodeFixed encodes val into exactly width digits using alphabet,
// zero-padded on the left if necessary.
func encodeFixed(val int64, alphabet string, width int) string {
	base := int64(len(alphabet))
	var buf [16]byte
	for i := width - 1; i >= 0; i-- {
		buf[i] = alphabet[val%base]
		val /= base
	}
	return string(buf[:width])
}

func timeValSec(t time.Time) int64 {
	y := int64(t.Year() % 100)
	m := int64(t.Month())
	d := int64(t.Day())
	hh := int64(t.Hour())
	mm := int64(t.Minute())
	ss := int64(t.Second())
	return y*10000000000 + m*100000000 + d*1000000 + hh*10000 + mm*100 + ss
}

func timeValMin(t time.Time) int64 {
	y := int64(t.Year() % 100)
	m := int64(t.Month())
	d := int64(t.Day())
	hh := int64(t.Hour())
	mm := int64(t.Minute())
	return y*100000000 + m*1000000 + d*10000 + hh*100 + mm
}

func timeValHour(t time.Time) int64 {
	y := int64(t.Year() % 100)
	m := int64(t.Month())
	d := int64(t.Day())
	hh := int64(t.Hour())
	return y*1000000 + m*10000 + d*100 + hh
}

func timeValDay(t time.Time) int64 {
	y := int64(t.Year() % 100)
	m := int64(t.Month())
	d := int64(t.Day())
	return y*10000 + m*100 + d
}

func timeValWeek(t time.Time) int64 {
	isoYear, isoWeek := t.ISOWeek()
	y := int64(isoYear % 100)
	w := int64(isoWeek)
	return y*100 + w
}

func timeValMonth(t time.Time) int64 {
	y := int64(t.Year() % 100)
	m := int64(t.Month())
	return y*100 + m
}

func timeValYear(t time.Time) int64 {
	return int64(t.Year() % 100)
}

// Crockford Base32 precision short IDs (fixed widths)

// ShortIDSec creates an 8-character Crockford Base32 short ID to second precision (YYMMDDHHMMSS).
func ShortIDSec() string {
	return encodeFixed(timeValSec(time.Now().UTC()), crockfordLowerAlphabet, 8)
}

// ShortIDMin creates a 7-character Crockford Base32 short ID to minute precision (YYMMDDHHMM).
func ShortIDMin() string {
	return encodeFixed(timeValMin(time.Now().UTC()), crockfordLowerAlphabet, 7)
}

// ShortIDHour creates a 6-character Crockford Base32 short ID to hour precision (YYMMDDHH).
func ShortIDHour() string {
	return encodeFixed(timeValHour(time.Now().UTC()), crockfordLowerAlphabet, 6)
}

// ShortIDDay creates a 4-character Crockford Base32 short ID to day precision (YYMMDD).
func ShortIDDay() string {
	return encodeFixed(timeValDay(time.Now().UTC()), crockfordLowerAlphabet, 4)
}

// ShortIDWeek creates a 3-character Crockford Base32 short ID to ISO week precision (YYWW).
func ShortIDWeek() string {
	return encodeFixed(timeValWeek(time.Now().UTC()), crockfordLowerAlphabet, 3)
}

// ShortIDMonth creates a 3-character Crockford Base32 short ID to month precision (YYMM).
func ShortIDMonth() string {
	return encodeFixed(timeValMonth(time.Now().UTC()), crockfordLowerAlphabet, 3)
}

// ShortIDYear creates a 2-character Crockford Base32 short ID to year precision (YY).
func ShortIDYear() string {
	return encodeFixed(timeValYear(time.Now().UTC()), crockfordLowerAlphabet, 2)
}

// Safe Base16 precision short IDs (fixed widths, vowel-free)

// SafeShortIDSec creates a 10-character Safe Base16 short ID to second precision.
func SafeShortIDSec() string {
	return encodeFixed(timeValSec(time.Now().UTC()), safeShortAlphabet, 10)
}

// SafeShortIDMin creates a 9-character Safe Base16 short ID to minute precision.
func SafeShortIDMin() string {
	return encodeFixed(timeValMin(time.Now().UTC()), safeShortAlphabet, 9)
}

// SafeShortIDHour creates a 7-character Safe Base16 short ID to hour precision.
func SafeShortIDHour() string {
	return encodeFixed(timeValHour(time.Now().UTC()), safeShortAlphabet, 7)
}

// SafeShortIDDay creates a 5-character Safe Base16 short ID to day precision.
func SafeShortIDDay() string {
	return encodeFixed(timeValDay(time.Now().UTC()), safeShortAlphabet, 5)
}

// SafeShortIDWeek creates a 4-character Safe Base16 short ID to ISO week precision.
func SafeShortIDWeek() string {
	return encodeFixed(timeValWeek(time.Now().UTC()), safeShortAlphabet, 4)
}

// SafeShortIDMonth creates a 4-character Safe Base16 short ID to month precision.
func SafeShortIDMonth() string {
	return encodeFixed(timeValMonth(time.Now().UTC()), safeShortAlphabet, 4)
}

// SafeShortIDYear creates a 2-character Safe Base16 short ID to year precision.
func SafeShortIDYear() string {
	return encodeFixed(timeValYear(time.Now().UTC()), safeShortAlphabet, 2)
}

// Hex Base16 precision short IDs (fixed widths, lowercase hex)

// HexShortIDSec creates a 10-character Hex Base16 short ID to second precision.
func HexShortIDSec() string {
	return encodeFixed(timeValSec(time.Now().UTC()), hexShortAlphabet, 10)
}

// HexShortIDMin creates a 9-character Hex Base16 short ID to minute precision.
func HexShortIDMin() string {
	return encodeFixed(timeValMin(time.Now().UTC()), hexShortAlphabet, 9)
}

// HexShortIDHour creates a 7-character Hex Base16 short ID to hour precision.
func HexShortIDHour() string {
	return encodeFixed(timeValHour(time.Now().UTC()), hexShortAlphabet, 7)
}

// HexShortIDDay creates a 5-character Hex Base16 short ID to day precision.
func HexShortIDDay() string {
	return encodeFixed(timeValDay(time.Now().UTC()), hexShortAlphabet, 5)
}

// HexShortIDWeek creates a 4-character Hex Base16 short ID to ISO week precision.
func HexShortIDWeek() string {
	return encodeFixed(timeValWeek(time.Now().UTC()), hexShortAlphabet, 4)
}

// HexShortIDMonth creates a 4-character Hex Base16 short ID to month precision.
func HexShortIDMonth() string {
	return encodeFixed(timeValMonth(time.Now().UTC()), hexShortAlphabet, 4)
}

// HexShortIDYear creates a 2-character Hex Base16 short ID to year precision.
func HexShortIDYear() string {
	return encodeFixed(timeValYear(time.Now().UTC()), hexShortAlphabet, 2)
}
