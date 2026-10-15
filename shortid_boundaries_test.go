package uid

import (
	"strings"
	"sync"
	"testing"
	"time"
)

// microsAt returns the Unix microsecond timestamp for midnight UTC on the given date.
func microsAt(y int, m time.Month, d int) int64 {
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC).UnixMicro()
}

// pinClock replaces the generator clock and resets the shared generator state
// for the duration of the test. Both are restored on cleanup.
// Tests using it must not call t.Parallel().
func pinClock(t *testing.T, clock func() int64) {
	t.Helper()
	oldClock := nowMicro
	oldState := shortIDAtomicState.Load()
	nowMicro = clock
	shortIDAtomicState.Store(0)
	t.Cleanup(func() {
		nowMicro = oldClock
		shortIDAtomicState.Store(oldState)
	})
}

// TestShortIDGenerators_ConcurrentUniqueness hammers every generator that
// shares the atomic state from many goroutines and verifies there are no
// duplicates and that each format has its documented length.
// HexShortID and ShortIDCustom("") produce the same alphabet and width, so
// they are checked against one shared set.
func TestShortIDGenerators_ConcurrentUniqueness(t *testing.T) {
	const goroutines, perGoroutine = 16, 5000

	type entry struct{ family, id string }
	wantLen := map[string]int{"short": 11, "safe": 14, "hex": 14}

	results := make([][]entry, goroutines)
	var wg sync.WaitGroup
	for g := 0; g < goroutines; g++ {
		wg.Add(1)
		go func(g int) {
			defer wg.Done()
			out := make([]entry, 0, perGoroutine*4)
			for i := 0; i < perGoroutine; i++ {
				out = append(out,
					entry{"short", ShortID()},
					entry{"safe", SafeShortID()},
					entry{"hex", HexShortID()},
					entry{"hex", ShortIDCustom("")},
				)
			}
			results[g] = out
		}(g)
	}
	wg.Wait()

	seen := make(map[entry]struct{}, goroutines*perGoroutine*4)
	for _, out := range results {
		for _, e := range out {
			if len(e.id) != wantLen[e.family] {
				t.Fatalf("%s id %q has length %d, want %d", e.family, e.id, len(e.id), wantLen[e.family])
			}
			if _, dup := seen[e]; dup {
				t.Fatalf("duplicate %s id %q", e.family, e.id)
			}
			seen[e] = struct{}{}
		}
	}
}

// TestNextShortIDValue_StrictlyIncreasingOnChaoticClock feeds the generator a
// clock that repeatedly jumps forwards and backwards by up to 1ms and checks
// the composite value never repeats or decreases.
func TestNextShortIDValue_StrictlyIncreasingOnChaoticClock(t *testing.T) {
	base := microsAt(2026, time.October, 16)
	var calls int64
	pinClock(t, func() int64 {
		calls++
		return base + (calls*7919)%2000 - 1000
	})

	prev := int64(-1)
	for i := 0; i < 20000; i++ {
		v := nextShortIDValue()
		if v <= prev {
			t.Fatalf("call %d: value %d does not exceed previous value %d", i, v, prev)
		}
		prev = v
	}
}

// TestNextShortIDValue_CounterOverflowAdvancesTimestamp pins the clock and
// checks the exact layout: 16 values per timestamp, then the timestamp moves
// forward by one microsecond (no sleeping, no wraparound to a used value).
func TestNextShortIDValue_CounterOverflowAdvancesTimestamp(t *testing.T) {
	ts := microsAt(2026, time.October, 16)
	pinClock(t, func() int64 { return ts })

	for k := 0; k < 1000; k++ {
		v := nextShortIDValue()
		wantTS, wantCounter := ts+int64(k/16), int64(k%16)
		if v>>4 != wantTS || v&0x0f != wantCounter {
			t.Fatalf("call %d: got ts=%d counter=%d, want ts=%d counter=%d",
				k, v>>4, v&0x0f, wantTS, wantCounter)
		}
	}
}

// TestShortIDs_WidthAndOrderingAcrossEras checks widths and chronological
// (lexicographic) ordering at dates around the known format boundaries:
//   - 2032-06-09: where the old 13-digit Hex/Safe IDs would have wrapped.
//   - 2041-05-10: the composite passes 2^55, so ShortID grows from 11 to 12
//     digits. Ordering is only asserted within each width, because a 12-digit
//     ID sorts lexicographically before an 11-digit one.
//   - ~2112: the composite passes 2^56, the limit of 14 hex digits.
func TestShortIDs_WidthAndOrderingAcrossEras(t *testing.T) {
	eras := []struct {
		name     string
		ts       int64
		shortLen int
	}{
		{"2026", microsAt(2026, time.October, 16), 11},
		{"2032-06-09 (past old 13-digit wrap)", microsAt(2032, time.June, 9), 11},
		{"2041-05-09 (last 11-digit ShortID days)", microsAt(2041, time.May, 9), 11},
		{"2041-05-11 (first 12-digit ShortID days)", microsAt(2041, time.May, 11), 12},
		{"2100", microsAt(2100, time.January, 1), 12},
		{"2111", microsAt(2111, time.January, 1), 12},
	}

	var current int64
	pinClock(t, func() int64 { return current })

	var prevShort, prevHex, prevSafe string
	for i, era := range eras {
		current = era.ts

		shortIDAtomicState.Store(0)
		short := ShortID()
		shortIDAtomicState.Store(0)
		hex := HexShortID()
		shortIDAtomicState.Store(0)
		safe := SafeShortID()

		if len(short) != era.shortLen {
			t.Errorf("%s: ShortID %q has length %d, want %d", era.name, short, len(short), era.shortLen)
		}
		if !IsShortID(short) {
			t.Errorf("%s: IsShortID(%q) = false", era.name, short)
		}
		if len(hex) != 14 || len(safe) != 14 {
			t.Errorf("%s: HexShortID %q / SafeShortID %q should both be 14 characters", era.name, hex, safe)
		}

		if i > 0 {
			if era.shortLen == eras[i-1].shortLen && short <= prevShort {
				t.Errorf("%s: ShortID %q does not sort after %q", era.name, short, prevShort)
			}
			if hex <= prevHex {
				t.Errorf("%s: HexShortID %q does not sort after %q", era.name, hex, prevHex)
			}
			if safe <= prevSafe {
				t.Errorf("%s: SafeShortID %q does not sort after %q", era.name, safe, prevSafe)
			}
		}
		prevShort, prevHex, prevSafe = short, hex, safe
	}
}

// TestShortIDCustom_RoundTripAllAlphabetSizes decodes the generated ID back
// into a number for every alphabet size from 2 to 94 (all printable ASCII
// except space) and compares it with the expected composite value. This proves
// there is no truncation and no buffer overflow for small bases.
func TestShortIDCustom_RoundTripAllAlphabetSizes(t *testing.T) {
	ts := microsAt(2026, time.October, 16)
	pinClock(t, func() int64 { return ts })

	var sb strings.Builder
	for c := byte('!'); c <= '~'; c++ {
		sb.WriteByte(c)
	}
	all := sb.String()

	for n := 2; n <= len(all); n++ {
		shortIDAtomicState.Store(0)
		alphabet := all[:n]
		id := ShortIDCustom(alphabet)

		var decoded uint64
		for i := 0; i < len(id); i++ {
			idx := strings.IndexByte(alphabet, id[i])
			if idx < 0 {
				t.Fatalf("alphabet size %d: character %q is not in the alphabet", n, id[i])
			}
			decoded = decoded*uint64(n) + uint64(idx)
		}
		if want := uint64(ts) << 4; decoded != want {
			t.Fatalf("alphabet size %d: id %q decodes to %d, want %d", n, id, decoded, want)
		}
	}
}

// TestShortIDCustom_RejectsInvalidAlphabets verifies the documented panics.
func TestShortIDCustom_RejectsInvalidAlphabets(t *testing.T) {
	invalid := map[string]string{
		"single character":     "a",
		"duplicate characters": "aab",
		"non-ASCII character":  "abé",
	}
	for name, alphabet := range invalid {
		t.Run(name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatalf("ShortIDCustom(%q) did not panic", alphabet)
				}
			}()
			_ = ShortIDCustom(alphabet)
		})
	}

	// The smallest valid alphabet must work.
	if id := ShortIDCustom("01"); id == "" {
		t.Fatal("ShortIDCustom(\"01\") returned an empty id")
	}
}

// TestSafeShortID_IsIndexMappingOfHexShortID verifies SafeShortID is exactly
// HexShortID with each hex digit mapped to the same position in the safe
// alphabet (so both encode the same underlying value and sort identically).
func TestSafeShortID_IsIndexMappingOfHexShortID(t *testing.T) {
	ts := microsAt(2026, time.October, 16)
	pinClock(t, func() int64 { return ts })

	hex := HexShortID()
	shortIDAtomicState.Store(0)
	safe := SafeShortID()

	if len(hex) != len(safe) {
		t.Fatalf("length mismatch: hex %q, safe %q", hex, safe)
	}
	for i := range hex {
		idx := strings.IndexByte(hexShortAlphabet, hex[i])
		if idx < 0 || safeShortAlphabet[idx] != safe[i] {
			t.Fatalf("position %d: hex %q and safe %q do not map by index", i, hex, safe)
		}
	}
}
