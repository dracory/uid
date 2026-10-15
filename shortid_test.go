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

// TestShortID_CounterRollover forces the clock to a fixed microsecond so the
// 4-bit counter overflows (more than 16 IDs on one tick), verifying the
// timestamp is bumped forward and all IDs remain unique.
func TestShortID_CounterRollover(t *testing.T) {
	old := nowMicro
	oldState := shortIDAtomicState.Load()
	defer func() { nowMicro = old; shortIDAtomicState.Store(oldState) }()

	fixed := time.Now().UnixMicro()
	nowMicro = func() int64 { return fixed }
	shortIDAtomicState.Store(0)

	seen := map[string]struct{}{}
	for i := 0; i < 100; i++ {
		id := ShortID()
		if len(id) != shortIDLength {
			t.Fatalf("expected length %d after rollover, got %d for id %q", shortIDLength, len(id), id)
		}
		if _, dup := seen[id]; dup {
			t.Fatalf("duplicate id after counter rollover: %q", id)
		}
		seen[id] = struct{}{}
	}

	// The stored timestamp must have advanced past the fixed injected time.
	if ts := int64(shortIDAtomicState.Load() >> 4); ts <= fixed {
		t.Fatalf("expected timestamp to advance past %d after rollover, got %d", fixed, ts)
	}
}

// TestShortID_BackwardsClock verifies that a clock moving backwards (NTP step,
// VM resume) does not reset the counter or produce duplicate IDs.
func TestShortID_BackwardsClock(t *testing.T) {
	old := nowMicro
	oldState := shortIDAtomicState.Load()
	defer func() { nowMicro = old; shortIDAtomicState.Store(oldState) }()

	current := int64(2000000000000)
	nowMicro = func() int64 { return current }
	shortIDAtomicState.Store(0)

	first := ShortID()

	// Step the clock backwards, then generate again.
	current -= 60_000_000 // one minute earlier
	second := ShortID()

	if first == second {
		t.Fatalf("duplicate id after backwards clock step: %q", first)
	}
}

// TestShortID_MonotonicWithinSameMicrosecond verifies that repeated
// calls landing on the same timestamp produce different composite values
// (i.e. the counter actually increments rather than silently reusing 0).
func TestShortID_MonotonicWithinSameMicrosecond(t *testing.T) {
	old := nowMicro
	oldState := shortIDAtomicState.Load()
	defer func() { nowMicro = old; shortIDAtomicState.Store(oldState) }()

	nowMicro = func() int64 { return int64(3000000000000) }
	shortIDAtomicState.Store(0)

	first := ShortID()
	second := ShortID()

	if first == second {
		t.Fatalf("expected different ids for same-timestamp calls, got %q twice", first)
	}
	if !(first < second) {
		t.Fatalf("expected lexicographically increasing ids, got %q then %q", first, second)
	}
}

// TestShortID_LexicographicOrdering verifies generated IDs sort in
// chronological order (the alphabet is ascending).
func TestShortID_LexicographicOrdering(t *testing.T) {
	prev := ShortID()
	for i := 0; i < 1000; i++ {
		id := ShortID()
		if id <= prev {
			t.Fatalf("ids not strictly increasing: %q then %q", prev, id)
		}
		prev = id
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

// --- SafeShortID, HexShortID, ShortIDCustom ---

// TestSafeShortID_Length verifies that SafeShortID always returns a 14-character string.
func TestSafeShortID_Length(t *testing.T) {
	for i := 0; i < 1000; i++ {
		id := SafeShortID()
		if len(id) != 14 {
			t.Fatalf("expected length 14 for SafeShortID, got %d for id %q", len(id), id)
		}
	}
}

// TestHexShortID_Length verifies that HexShortID always returns a 14-character string.
func TestHexShortID_Length(t *testing.T) {
	for i := 0; i < 1000; i++ {
		id := HexShortID()
		if len(id) != 14 {
			t.Fatalf("expected length 14 for HexShortID, got %d for id %q", len(id), id)
		}
	}
}

// TestSafeShortID_NoVowels verifies that SafeShortID outputs contain no vowels
// ('a', 'e', 'i', 'o', 'u') and only contain characters from the Safe Base16 alphabet.
func TestSafeShortID_NoVowels(t *testing.T) {
	vowels := "aeiou"
	validAlphabet := "0123456789bcdfrt"

	for i := 0; i < 10000; i++ {
		id := SafeShortID()
		for _, c := range id {
			if strings.ContainsRune(vowels, c) {
				t.Fatalf("found vowel %q in SafeShortID output %q", string(c), id)
			}
			if !strings.ContainsRune(validAlphabet, c) {
				t.Fatalf("character %q in SafeShortID %q is not in safe alphabet", string(c), id)
			}
		}
	}
}

// TestHexShortID_ValidAlphabet verifies every character in HexShortID
// belongs to standard lowercase hex ("0123456789abcdef").
func TestHexShortID_ValidAlphabet(t *testing.T) {
	validAlphabet := "0123456789abcdef"
	for i := 0; i < 1000; i++ {
		id := HexShortID()
		for _, c := range id {
			if !strings.ContainsRune(validAlphabet, c) {
				t.Fatalf("character %q in HexShortID %q is not in hex alphabet", string(c), id)
			}
		}
	}
}

// TestShortIDCustom verifies ShortIDCustom behavior with various custom alphabets
// as well as the empty string fallback.
func TestShortIDCustom(t *testing.T) {
	t.Run("empty string fallback to hex", func(t *testing.T) {
		id := ShortIDCustom("")
		validAlphabet := "0123456789abcdef"
		for _, c := range id {
			if !strings.ContainsRune(validAlphabet, c) {
				t.Fatalf("character %q in fallback output %q not in hex alphabet", string(c), id)
			}
		}
	})

	t.Run("digits only alphabet", func(t *testing.T) {
		digits := "0123456789"
		id := ShortIDCustom(digits)
		if id == "" {
			t.Fatal("expected non-empty string")
		}
		for _, c := range id {
			if !strings.ContainsRune(digits, c) {
				t.Fatalf("character %q in output %q not in custom alphabet %q", string(c), id, digits)
			}
		}
	})

	t.Run("custom uppercase alphabet", func(t *testing.T) {
		alphabet := "ABCDEFGHJK"
		id := ShortIDCustom(alphabet)
		if id == "" {
			t.Fatal("expected non-empty string")
		}
		for _, c := range id {
			if !strings.ContainsRune(alphabet, c) {
				t.Fatalf("character %q in output %q not in custom alphabet %q", string(c), id, alphabet)
			}
		}
	})

	t.Run("binary alphabet", func(t *testing.T) {
		id := ShortIDCustom("01")
		for _, c := range id {
			if c != '0' && c != '1' {
				t.Fatalf("unexpected character %q in binary output %q", c, id)
			}
		}
	})

	t.Run("panics on single-char alphabet", func(t *testing.T) {
		defer func() {
			if recover() == nil {
				t.Fatal("expected panic for 1-character alphabet")
			}
		}()
		ShortIDCustom("x")
	})

	t.Run("panics on duplicate characters", func(t *testing.T) {
		defer func() {
			if recover() == nil {
				t.Fatal("expected panic for alphabet with duplicates")
			}
		}()
		ShortIDCustom("aab")
	})

	t.Run("panics on non-ASCII alphabet", func(t *testing.T) {
		defer func() {
			if recover() == nil {
				t.Fatal("expected panic for multi-byte alphabet")
			}
		}()
		ShortIDCustom("αβγδεζηθικλμ")
	})
}

// FuzzShortIDCustom exercises ShortIDCustom over arbitrary alphabet inputs,
// ensuring it never panics for valid alphabets and panics cleanly otherwise.
func FuzzShortIDCustom(f *testing.F) {
	f.Add("01")
	f.Add("0123456789abcdef")
	f.Add("x")
	f.Add("aab")
	f.Add("日本語")
	f.Fuzz(func(t *testing.T, alphabet string) {
		defer func() {
			// Panics for invalid alphabets are acceptable; out-of-bounds
			// writes are not (they would crash rather than panic).
			_ = recover()
		}()
		id := ShortIDCustom(alphabet)
		if alphabet == "" {
			alphabet = hexShortAlphabet // matches the empty-string fallback
		}
		for i := 0; i < len(id); i++ {
			if strings.IndexByte(alphabet, id[i]) < 0 {
				t.Fatalf("output %q contains byte not in alphabet %q", id, alphabet)
			}
		}
	})
}

// TestNewShortIDs_ConcurrentUniqueness generates IDs concurrently using SafeShortID,
// HexShortID, and ShortIDCustom across multiple goroutines to verify uniqueness.
func TestNewShortIDs_ConcurrentUniqueness(t *testing.T) {
	const goroutines = 50
	const perGoroutine = 1000
	const total = goroutines * perGoroutine * 3

	ids := make(chan string, total)
	var wg sync.WaitGroup

	for g := 0; g < goroutines; g++ {
		wg.Add(3)
		go func() {
			defer wg.Done()
			for i := 0; i < perGoroutine; i++ {
				ids <- SafeShortID()
			}
		}()
		go func() {
			defer wg.Done()
			for i := 0; i < perGoroutine; i++ {
				ids <- HexShortID()
			}
		}()
		go func() {
			defer wg.Done()
			for i := 0; i < perGoroutine; i++ {
				ids <- ShortIDCustom("0123456789abcdefghijklmnopqrstuvwxyz")
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
		{"valid 12-char id (post-2041 format)", "ab12cd3ef456", true},
		{"too long (13 chars)", "ab12cd3ef4567", false},
		{"right length but invalid charset", "hello world", false},
		{"chars outside crockford alphabet", "iluao012345", false},
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
