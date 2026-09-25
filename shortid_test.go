package uid

import (
	"strings"
	"sync"
	"testing"
	"time"
)

// TestShortID_Length verifies every generated ID is exactly 11
// characters long, regardless of how many are generated in sequence.
func TestShortID_Length(t *testing.T) {
	for i := 0; i < 1000; i++ {
		id := ShortID()
		if len(id) != 11 {
			t.Fatalf("expected length 11, got %d for id %q", len(id), id)
		}
	}
}

// TestShortID_Lowercase verifies the output is always lowercase.
func TestShortID_Lowercase(t *testing.T) {
	id := ShortID()
	if id != strings.ToLower(id) {
		t.Fatalf("expected lowercase id, got %q", id)
	}
}

// TestShortID_ValidAlphabet verifies every character in the ID
// belongs to the Crockford Base32 alphabet (case-insensitively).
func TestShortID_ValidAlphabet(t *testing.T) {
	lowerAlphabet := strings.ToLower(alphabetCrockford)
	id := ShortID()
	for _, c := range id {
		if !strings.ContainsRune(lowerAlphabet, c) {
			t.Fatalf("character %q in id %q is not in the Crockford alphabet", c, id)
		}
	}
}

// TestShortID_SequentialUniqueness generates a large number of IDs
// back-to-back on a single goroutine (the worst case for hitting the same
// microsecond repeatedly) and asserts none collide.
func TestShortID_SequentialUniqueness(t *testing.T) {
	const n = 50000
	seen := make(map[string]struct{}, n)

	for i := 0; i < n; i++ {
		id := ShortID()
		if _, exists := seen[id]; exists {
			t.Fatalf("duplicate id generated: %q (iteration %d)", id, i)
		}
		seen[id] = struct{}{}
	}
}

// TestShortID_ConcurrentUniqueness generates IDs from many
// goroutines simultaneously and asserts the mutex-guarded counter logic
// prevents any duplicates.
func TestShortID_ConcurrentUniqueness(t *testing.T) {
	const goroutines = 50
	const perGoroutine = 1000
	const total = goroutines * perGoroutine

	ids := make(chan string, total)
	var wg sync.WaitGroup

	for g := 0; g < goroutines; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < perGoroutine; i++ {
				ids <- ShortID()
			}
		}()
	}

	wg.Wait()
	close(ids)

	seen := make(map[string]struct{}, total)
	for id := range ids {
		if _, exists := seen[id]; exists {
			t.Fatalf("duplicate id generated under concurrency: %q", id)
		}
		seen[id] = struct{}{}
	}

	if len(seen) != total {
		t.Fatalf("expected %d unique ids, got %d", total, len(seen))
	}
}

// TestShortID_CounterRollover forces more than 16 calls within
// (as close as possible to) the same microsecond by directly manipulating
// the package-level state, verifying the overflow fallback path does not
// panic and still yields a valid, correctly-shaped ID.
func TestShortID_CounterRollover(t *testing.T) {
	shortIDMutex.Lock()
	lastShortIDStamp = time.Now().UnixMicro()
	shortIDCounter = 16 // one past the max (0-15)
	shortIDMutex.Unlock()

	id := ShortID()

	if len(id) != 11 {
		t.Fatalf("expected length 11 after rollover, got %d for id %q", len(id), id)
	}

	shortIDMutex.Lock()
	counterAfter := shortIDCounter
	shortIDMutex.Unlock()

	if counterAfter > 15 {
		t.Fatalf("expected counter to reset to <=15 after rollover, got %d", counterAfter)
	}
}

// TestShortID_MonotonicWithinSameMicrosecond verifies that repeated
// calls landing on the same timestamp produce different composite values
// (i.e. the counter actually increments rather than silently reusing 0).
func TestShortID_MonotonicWithinSameMicrosecond(t *testing.T) {
	shortIDMutex.Lock()
	fixedTS := time.Now().UnixMicro()
	lastShortIDStamp = fixedTS
	shortIDCounter = -1 // so the first call below increments to 0
	shortIDMutex.Unlock()

	first := ShortID()

	shortIDMutex.Lock()
	// Force the same timestamp again to simulate a second call landing on
	// the identical microsecond tick.
	lastShortIDStamp = fixedTS
	shortIDMutex.Unlock()

	second := ShortID()

	if first == second {
		t.Fatalf("expected different ids for same-timestamp calls, got %q twice", first)
	}
}

// --- NormalizeID ---

func TestNormalizeID(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"already lowercase", "ab12cd3ef45", "ab12cd3ef45"},
		{"uppercase", "AB12CD3EF45", "ab12cd3ef45"},
		{"mixed case", "aB12cD3eF45", "ab12cd3ef45"},
		{"leading/trailing whitespace", "  ab12cd3ef45  ", "ab12cd3ef45"},
		{"tabs and newlines", "\tab12cd3ef45\n", "ab12cd3ef45"},
		{"empty string", "", ""},
		{"whitespace only", "   ", ""},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := NormalizeID(c.in)
			if got != c.want {
				t.Errorf("NormalizeID(%q) = %q, want %q", c.in, got, c.want)
			}
		})
	}
}

// --- IsShortID ---

func TestIsShortID(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want bool
	}{
		{"valid 11-char id", "ab12cd3ef45", true},
		{"generated id", ShortID(), true},
		{"empty string", "", false},
		{"too short", "ab12cd3", false},
		{"too long", "ab12cd3ef45xyz", false},
		{"one char short", "ab12cd3ef4", false},
		{"one char long", "ab12cd3ef456", false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := IsShortID(c.in)
			if got != c.want {
				t.Errorf("IsShortID(%q) = %v, want %v", c.in, got, c.want)
			}
		})
	}
}
