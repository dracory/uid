package uid

import (
	"crypto/rand"
	"strconv"
	"strings"
	"time"
)

// HumanUid generates a 32-character time-prefixed unique ID.
//
// Format (conceptual): YYYYMMDDHHMMSSMMMMMMM + random suffix, truncated to 32.
//
// Example: 20250831151133000012345678901234 (length: 32)
//
// Parameters:
// - None
//
// Returns:
// - A 32-character uppercase numeric string suitable for human-readable IDs
func HumanUid() string {
	time.Sleep(1 * time.Nanosecond)
	r, _ := rand.Prime(rand.Reader, 64)
	id := time.Now().UTC().Format("20060102150405.0000000")
	id = strings.ReplaceAll(id, ".", "")
	id += r.String()
	return id[0:32]
}

// NanoUid generates a 23-character time-prefixed unique ID.
//
// Format (conceptual): YYYYMMDDHHMMSSMMMMMMM + random suffix, truncated to 23.
//
// Example: 20250831151133000012345 (length: 23)
//
// Parameters:
// - None
//
// Returns:
// - A 23-character numeric string
func NanoUid() string {
	time.Sleep(time.Nanosecond) // as its a nanoseconds based ID we need at least a nanosecond between the generations to avoid collisions
	r, _ := rand.Prime(rand.Reader, 64)
	id := time.Now().UTC().Format("20060102150405.0000000")
	id = strings.ReplaceAll(id, ".", "")
	id += r.String()
	return id[0:23]
}

// MicroUid generates a 20-character time-prefixed unique ID.
//
// Format (conceptual): YYYYMMDDHHMMSSMMMMMMM + random suffix, truncated to 20.
//
// Example: 20250831151133000012 (length: 20)
//
// Parameters:
// - None
//
// Returns:
// - A 20-character numeric string
func MicroUid() string {
	time.Sleep(time.Microsecond) // as its a microseconds based ID we need at least a microsecond between the generations to avoid collisions
	r, _ := rand.Prime(rand.Reader, 64)
	id := time.Now().UTC().Format("20060102150405.0000000")
	id = strings.ReplaceAll(id, ".", "")
	id += r.String()
	return id[0:20]
}

// SecUid generates a 14-character time-based ID.
//
// Format: YYYYMMDDHHMMSS
//
// Example: 20250831151133 (length: 14)
//
// Parameters:
// - None
//
// Returns:
// - A 14-character numeric string representing UTC date/time to the second
func SecUid() string {
	time.Sleep(time.Second) // as its a seconds based ID we need at least a second between the generations to avoid collisions
	r, _ := rand.Prime(rand.Reader, 64)
	id := time.Now().UTC().Format("20060102150405.0000000")
	id = strings.ReplaceAll(id, ".", "")
	id += r.String()
	return id[0:14]
}

// Timestamp returns the current Unix timestamp in seconds as a string.
//
// Example: 1725111153 (length: 10)
//
// Parameters:
// - None
//
// Returns:
// - Unix timestamp in seconds (base-10 string)
func Timestamp() string {
	time.Sleep(time.Second) // as its a seconds based ID we need at least a second between the generations to avoid collisions
	now := time.Now().UTC().Unix()
	return strconv.FormatInt(now, 10)
}

// TimestampMicro returns the current Unix timestamp in microseconds as a string.
//
// Example: 1725111153123456 (length: 16)
//
// Parameters:
// - None
//
// Returns:
// - Unix timestamp in microseconds (base-10 string)
func TimestampMicro() string {
	time.Sleep(time.Microsecond) // as its a microseconds based ID we need at least a microsecond between the generations to avoid collisions
	now := time.Now().UTC().UnixMicro()
	return strconv.FormatInt(now, 10)
}

// TimestampNano returns the current Unix timestamp in nanoseconds as a string.
//
// Example: 1725111153123456789 (length: 19)
//
// Parameters:
// - None
//
// Returns:
// - Unix timestamp in nanoseconds (base-10 string)
func TimestampNano() string {
	time.Sleep(time.Nanosecond) // as its a nanoseconds based ID we need at least a nanosecond between the generations to avoid collisions
	now := time.Now().UTC().UnixNano()
	return strconv.FormatInt(now, 10)
}
