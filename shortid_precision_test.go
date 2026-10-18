package uid

import (
	"strings"
	"testing"
	"time"
)

func TestPrecisionShortIDs_LengthsAndCharsets(t *testing.T) {
	crockfordAlph := crockfordLowerAlphabet
	safeAlph := safeShortAlphabet
	hexAlph := hexShortAlphabet

	tests := []struct {
		name     string
		fn       func() string
		wantLen  int
		alphabet string
		noVowels bool
	}{
		// Crockford Base32
		{"ShortIDSec", ShortIDSec, 8, crockfordAlph, false},
		{"ShortIDMin", ShortIDMin, 7, crockfordAlph, false},
		{"ShortIDHour", ShortIDHour, 6, crockfordAlph, false},
		{"ShortIDDay", ShortIDDay, 4, crockfordAlph, false},
		{"ShortIDWeek", ShortIDWeek, 3, crockfordAlph, false},
		{"ShortIDMonth", ShortIDMonth, 3, crockfordAlph, false},
		{"ShortIDYear", ShortIDYear, 2, crockfordAlph, false},

		// Safe Base16
		{"SafeShortIDSec", SafeShortIDSec, 10, safeAlph, true},
		{"SafeShortIDMin", SafeShortIDMin, 9, safeAlph, true},
		{"SafeShortIDHour", SafeShortIDHour, 7, safeAlph, true},
		{"SafeShortIDDay", SafeShortIDDay, 5, safeAlph, true},
		{"SafeShortIDWeek", SafeShortIDWeek, 4, safeAlph, true},
		{"SafeShortIDMonth", SafeShortIDMonth, 4, safeAlph, true},
		{"SafeShortIDYear", SafeShortIDYear, 2, safeAlph, true},

		// Hex Base16
		{"HexShortIDSec", HexShortIDSec, 10, hexAlph, false},
		{"HexShortIDMin", HexShortIDMin, 9, hexAlph, false},
		{"HexShortIDHour", HexShortIDHour, 7, hexAlph, false},
		{"HexShortIDDay", HexShortIDDay, 5, hexAlph, false},
		{"HexShortIDWeek", HexShortIDWeek, 4, hexAlph, false},
		{"HexShortIDMonth", HexShortIDMonth, 4, hexAlph, false},
		{"HexShortIDYear", HexShortIDYear, 2, hexAlph, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := 0; i < 100; i++ {
				id := tt.fn()
				if len(id) != tt.wantLen {
					t.Fatalf("%s returned length %d, want %d (id: %q)", tt.name, len(id), tt.wantLen, id)
				}
				for _, c := range id {
					if !strings.ContainsRune(tt.alphabet, c) {
						t.Fatalf("%s returned invalid char %q in id %q for alphabet %q", tt.name, string(c), id, tt.alphabet)
					}
					if tt.noVowels && strings.ContainsRune("aeiou", c) {
						t.Fatalf("%s returned vowel %q in safe id %q", tt.name, string(c), id)
					}
				}
			}
		})
	}
}

func TestPrecisionShortIDs_ValuesAtKnownDates(t *testing.T) {
	// 2027-01-01 00:00:00 UTC -> belongs to ISO year 2026 week 53
	dt := time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)

	if got := timeValSec(dt); got != 270101000000 {
		t.Errorf("timeValSec = %d, want 270101000000", got)
	}
	if got := timeValMin(dt); got != 2701010000 {
		t.Errorf("timeValMin = %d, want 2701010000", got)
	}
	if got := timeValHour(dt); got != 27010100 {
		t.Errorf("timeValHour = %d, want 27010100", got)
	}
	if got := timeValDay(dt); got != 270101 {
		t.Errorf("timeValDay = %d, want 270101", got)
	}
	if got := timeValWeek(dt); got != 2653 {
		t.Errorf("timeValWeek = %d, want 2653", got)
	}
	if got := timeValMonth(dt); got != 2701 {
		t.Errorf("timeValMonth = %d, want 2701", got)
	}
	if got := timeValYear(dt); got != 27 {
		t.Errorf("timeValYear = %d, want 27", got)
	}

	// Verify encoding for week 2653
	// 2653 in Crockford Base32 width 3 -> "2jx"
	if got := encodeFixed(2653, crockfordLowerAlphabet, 3); got != "2jx" {
		t.Errorf("encodeFixed(2653, crockford, 3) = %q, want %q", got, "2jx")
	}
	// 2653 in Safe Base16 width 4 -> "0b5f"
	if got := encodeFixed(2653, safeShortAlphabet, 4); got != "0b5f" {
		t.Errorf("encodeFixed(2653, safe, 4) = %q, want %q", got, "0b5f")
	}
	// 2653 in Hex Base16 width 4 -> "0a5d"
	if got := encodeFixed(2653, hexShortAlphabet, 4); got != "0a5d" {
		t.Errorf("encodeFixed(2653, hex, 4) = %q, want %q", got, "0a5d")
	}
}
