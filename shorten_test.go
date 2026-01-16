package uid

import (
	"testing"
)

func TestShortenUnshorten(t *testing.T) {
	testCases := []string{
		"0",
		"1",
		"12345",
		"12345678901234567890",             // 20 digits
		"20250831151133000012345678901234", // 32 digits
	}

	for _, tc := range testCases {
		t.Run("Base16_"+tc, func(t *testing.T) {
			short, err := ShortenBase16(tc)
			if err != nil {
				t.Fatalf("ShortenBase16 failed: %v", err)
			}
			unshort, _ := UnshortenBase16(short)
			if unshort != tc {
				t.Errorf("Base16 failed for %s: got %s -> %s", tc, short, unshort)
			}
		})

		t.Run("Base32_"+tc, func(t *testing.T) {
			short, err := ShortenBase32(tc)
			if err != nil {
				t.Fatalf("ShortenBase32 failed: %v", err)
			}
			unshort, _ := UnshortenBase32(short)
			if unshort != tc {
				t.Errorf("Base32 failed for %s: got %s -> %s", tc, short, unshort)
			}
		})

		t.Run("Crockford_"+tc, func(t *testing.T) {
			short, err := ShortenCrockford(tc)
			if err != nil {
				t.Fatalf("ShortenCrockford failed: %v", err)
			}
			unshort, _ := UnshortenCrockford(short)
			if unshort != tc {
				t.Errorf("Crockford failed for %s: got %s -> %s", tc, short, unshort)
			}
		})

		t.Run("Base64_"+tc, func(t *testing.T) {
			short, err := ShortenBase64(tc)
			if err != nil {
				t.Fatalf("ShortenBase64 failed: %v", err)
			}
			unshort, _ := UnshortenBase64(short)
			if unshort != tc {
				t.Errorf("Base64 failed for %s: got %s -> %s", tc, short, unshort)
			}
		})

		t.Run("Base36_"+tc, func(t *testing.T) {
			short, err := ShortenBase36(tc)
			if err != nil {
				t.Fatalf("ShortenBase36 failed: %v", err)
			}
			unshort, _ := UnshortenBase36(short)
			if unshort != tc {
				t.Errorf("Base36 failed for %s: got %s -> %s", tc, short, unshort)
			}
		})

		t.Run("Base58_"+tc, func(t *testing.T) {
			short, err := ShortenBase58(tc)
			if err != nil {
				t.Fatalf("ShortenBase58 failed: %v", err)
			}
			unshort, _ := UnshortenBase58(short)
			if unshort != tc {
				t.Errorf("Base58 failed for %s: got %s -> %s", tc, short, unshort)
			}
		})

		t.Run("Base62_"+tc, func(t *testing.T) {
			short, err := ShortenBase62(tc)
			if err != nil {
				t.Fatalf("ShortenBase62 failed: %v", err)
			}
			unshort, _ := UnshortenBase62(short)
			if unshort != tc {
				t.Errorf("Base62 failed for %s: got %s -> %s", tc, short, unshort)
			}
		})

		t.Run("ZBase32_"+tc, func(t *testing.T) {
			short, err := ShortenZBase32(tc)
			if err != nil {
				t.Fatalf("ShortenZBase32 failed: %v", err)
			}
			unshort, _ := UnshortenZBase32(short)
			if unshort != tc {
				t.Errorf("ZBase32 failed for %s: got %s -> %s", tc, short, unshort)
			}
		})
	}
}

func TestCrockfordNormalization(t *testing.T) {
	id := "12345678901234567890"
	short, _ := ShortenCrockford(id)

	// Test lowercase
	unshort, _ := UnshortenCrockford(short)
	if unshort != id {
		t.Error("Crockford normalization failed for standard")
	}

	// Test O -> 0, I/L -> 1
	// Alphabet: 0123456789ABCDEFGHJKMNPQRSTVWXYZ
	// 0 is index 0, 1 is index 1
	// Let's manually construct a string with O, I, L
	normalized, _ := UnshortenCrockford("1O-IL")
	// 1 (1) O (0) - (ignored) I (1) L (1) => base32: [1, 0, 1, 1]
	// 1*32^3 + 0*32^2 + 1*32^1 + 1*32^0 = 32768 + 0 + 32 + 1 = 32801
	if normalized != "32801" {
		t.Errorf("Crockford normalization failed: got %s, want 32801", normalized)
	}
}

func TestUIDShortening(t *testing.T) {
	// Generate a real HumanUid and shorten it
	id := HumanUid()
	t.Logf("Original ID: %s (%d digits)", id, len(id))

	b16, _ := ShortenBase16(id)
	t.Logf("Base16: %s (%d chars)", b16, len(b16))

	b32, _ := ShortenBase32(id)
	t.Logf("Base32: %s (%d chars)", b32, len(b32))

	crock, _ := ShortenCrockford(id)
	t.Logf("Crockford: %s (%d chars)", crock, len(crock))

	b64, _ := ShortenBase64(id)
	t.Logf("Base64: %s (%d chars)", b64, len(b64))

	b36, _ := ShortenBase36(id)
	t.Logf("Base36: %s (%d chars)", b36, len(b36))

	b58, _ := ShortenBase58(id)
	t.Logf("Base58: %s (%d chars)", b58, len(b58))

	b62, _ := ShortenBase62(id)
	t.Logf("Base62: %s (%d chars)", b62, len(b62))

	z32, _ := ShortenZBase32(id)
	t.Logf("ZBase32: %s (%d chars)", z32, len(z32))

	uz32, _ := UnshortenZBase32(z32)
	if uz32 != id {
		t.Errorf("ZBase32 roundtrip failed for HumanUid")
	}

	if len(z32) >= len(id) {
		t.Errorf("Shortening didn't reduce length: original %d, z32 %d", len(id), len(z32))
	}
}

func TestShortenInvalid(t *testing.T) {
	_, err := ShortenBase58("abc")
	if err == nil {
		t.Error("ShortenBase58 should fail for non-numeric input")
	}
	if err != ErrInvalidInput {
		t.Errorf("Expected ErrInvalidInput, got %v", err)
	}
}
