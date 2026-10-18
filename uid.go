package uid

import (
	"crypto/rand"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// HumanUid generates a 32-character time-prefixed unique ID.
//
// Format (conceptual): YYYYMMDDHHMMSSMMMMMMM + random suffix, truncated to 32.
//
// Example (unformatted): 20250831151133000012345678901234 (length: 32)
// Example (formatted): 20171119-0849-2665-991498485465 (length: 35)
//
// Parameters:
// - formatted: when true, include hyphens in groups 8-4-4-16 (length becomes 35)
//
// Returns:
// - A 32-character uppercase numeric string suitable for human-readable IDs
func HumanUid(formatted ...bool) string {
	time.Sleep(1 * time.Nanosecond)

	r, _ := rand.Prime(rand.Reader, 64)

	id := time.Now().UTC().Format("20060102150405.0000000")
	id = strings.ReplaceAll(id, ".", "")
	id += r.String()

	s := id[0:32]
	withHyphens := len(formatted) > 0 && formatted[0]
	if withHyphens {
		return formatWithHyphens(s, []int{8, 4, 4, 16})
	}
	return s
}

// MinUid generates a 12-character time-based ID to minute precision.
//
// Format: YYYYMMDDHHMM
//
// Example (unformatted): 202508311511 (length: 12)
// Example (formatted): 20250831-1511 (length: 13)
//
// Parameters:
// - formatted: when true, include hyphens in groups 8-4 (length becomes 13)
//
// Returns:
// - A 12-character numeric string representing UTC date/time to the minute
func MinUid(formatted ...bool) string {
	s := time.Now().UTC().Format("200601021504")
	withHyphens := len(formatted) > 0 && formatted[0]
	if withHyphens {
		return formatWithHyphens(s, []int{8, 4})
	}
	return s
}

// HourUid generates a 10-character time-based ID to hour precision.
//
// Format: YYYYMMDDHH
//
// Example (unformatted): 2025083115 (length: 10)
// Example (formatted): 20250831-15 (length: 11)
//
// Parameters:
// - formatted: when true, include hyphens in groups 8-2 (length becomes 11)
//
// Returns:
// - A 10-character numeric string representing UTC date/time to the hour
func HourUid(formatted ...bool) string {
	s := time.Now().UTC().Format("2006010215")
	withHyphens := len(formatted) > 0 && formatted[0]
	if withHyphens {
		return formatWithHyphens(s, []int{8, 2})
	}
	return s
}

// DayUid generates an 8-character time-based ID to day precision.
//
// Format: YYYYMMDD
//
// Example (unformatted): 20250831 (length: 8)
// Example (formatted): 2025-08-31 (length: 10)
//
// Parameters:
// - formatted: when true, include hyphens in groups 4-2-2 (length becomes 10)
//
// Returns:
// - An 8-character numeric string representing UTC date to the day
func DayUid(formatted ...bool) string {
	s := time.Now().UTC().Format("20060102")
	withHyphens := len(formatted) > 0 && formatted[0]
	if withHyphens {
		return formatWithHyphens(s, []int{4, 2, 2})
	}
	return s
}

// WeekUid generates a 6-character time-based ID to ISO week precision.
//
// Format: YYYYWW
//
// Example (unformatted): 202535 (length: 6)
// Example (formatted): 2025-35 (length: 7)
//
// Parameters:
// - formatted: when true, include hyphens in groups 4-2 (length becomes 7)
//
// Returns:
// - A 6-character numeric string representing UTC ISO year and week
func WeekUid(formatted ...bool) string {
	isoYear, isoWeek := time.Now().UTC().ISOWeek()
	s := fmt.Sprintf("%04d%02d", isoYear, isoWeek)
	withHyphens := len(formatted) > 0 && formatted[0]
	if withHyphens {
		return formatWithHyphens(s, []int{4, 2})
	}
	return s
}

// MonthUid generates a 6-character time-based ID to month precision.
//
// Format: YYYYMM
//
// Example (unformatted): 202508 (length: 6)
// Example (formatted): 2025-08 (length: 7)
//
// Parameters:
// - formatted: when true, include hyphens in groups 4-2 (length becomes 7)
//
// Returns:
// - A 6-character numeric string representing UTC date to the month
func MonthUid(formatted ...bool) string {
	s := time.Now().UTC().Format("200601")
	withHyphens := len(formatted) > 0 && formatted[0]
	if withHyphens {
		return formatWithHyphens(s, []int{4, 2})
	}
	return s
}

// NanoUid generates a 23-character time-prefixed unique ID.
//
// Format (conceptual): YYYYMMDDHHMMSSMMMMMMM + random suffix, truncated to 23.
//
// Example (unformatted): 20250831151133000012345 (length: 23)
// Example (formatted): 20171119-084926-659914-984 (length: 26)
//
// Parameters:
// - formatted: when true, include hyphens in groups 8-6-6-3 (length becomes 26)
//
// Returns:
// - A 23-character numeric string
func NanoUid(formatted ...bool) string {
	time.Sleep(time.Nanosecond) // as its a nanoseconds based ID we need at least a nanosecond between the generations to avoid collisions

	r, _ := rand.Prime(rand.Reader, 64)

	id := time.Now().UTC().Format("20060102150405.0000000")
	id = strings.ReplaceAll(id, ".", "")
	id += r.String()

	s := id[0:23]
	withHyphens := len(formatted) > 0 && formatted[0]
	if withHyphens {
		return formatWithHyphens(s, []int{8, 6, 6, 3})
	}
	return s
}

// MicroUid generates a 20-character time-prefixed unique ID.
//
// Format (conceptual): YYYYMMDDHHMMSSMMMMMMM + random suffix, truncated to 20.
//
// Example (unformatted): 20250831151133000012 (length: 20)
// Example (formatted): 20171119-084926-659914 (length: 22)
//
// Parameters:
// - formatted: when true, include hyphens in groups 8-6-6 (length becomes 22)
//
// Returns:
// - A 20-character numeric string
func MicroUid(formatted ...bool) string {
	time.Sleep(time.Microsecond) // as its a microseconds based ID we need at least a microsecond between the generations to avoid collisions

	r, _ := rand.Prime(rand.Reader, 64)

	id := time.Now().UTC().Format("20060102150405.0000000")
	id = strings.ReplaceAll(id, ".", "")
	id += r.String()

	s := id[0:20]
	withHyphens := len(formatted) > 0 && formatted[0]
	if withHyphens {
		return formatWithHyphens(s, []int{8, 6, 6})
	}
	return s
}

// SecUid generates a 14-character time-based ID.
//
// Format: YYYYMMDDHHMMSS
//
// Example (unformatted): 20250831151133 (length: 14)
// Example (formatted): 20171119-084926 (length: 15)
//
// Parameters:
// - formatted: when true, include hyphens in groups 8-6 (length becomes 15)
//
// Returns:
// - A 14-character numeric string representing UTC date/time to the second
func SecUid(formatted ...bool) string {
	time.Sleep(time.Second) // as its a seconds based ID we need at least a second between the generations to avoid collisions

	r, _ := rand.Prime(rand.Reader, 64)

	id := time.Now().UTC().Format("20060102150405.0000000")
	id = strings.ReplaceAll(id, ".", "")
	id += r.String()

	s := id[0:14]
	withHyphens := len(formatted) > 0 && formatted[0]
	if withHyphens {
		return formatWithHyphens(s, []int{8, 6})
	}
	return s
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

// formatWithHyphens inserts hyphens into s grouped by the provided sizes.
// Example: formatWithHyphens("20171119084926659914", []int{8,6,6}) => "20171119-084926-659914".
func formatWithHyphens(s string, groups []int) string {
	var b strings.Builder
	b.Grow(len(s) + len(groups) - 1)
	pos := 0
	for i, g := range groups {
		if i > 0 {
			b.WriteByte('-')
		}
		end := pos + g
		if end > len(s) {
			end = len(s)
		}
		b.WriteString(s[pos:end])
		pos = end
		if pos >= len(s) {
			break
		}
	}
	// Append any remaining characters (should not happen if groups sum to len(s))
	if pos < len(s) {
		b.WriteString(s[pos:])
	}
	return b.String()
}
