package uid

import (
	"testing"
)

func TestShortUnshortConvenience(t *testing.T) {
	testCases := []struct {
		name     string
		original string
	}{
		{"HumanUid", HumanUid()},
		{"NanoUid", NanoUid()},
		{"MicroUid", MicroUid()},
		{"SecUid", SecUid()},
		{"TimestampNano", TimestampNano()},
		{"TimestampMicro", TimestampMicro()},
		{"Timestamp", Timestamp()},
		{"Custom20", "12345678901234567890"},
		{"Custom14", "12345678901234"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			short, err := Short(tc.original)
			if err != nil {
				t.Fatalf("%s: Short failed: %v", tc.name, err)
			}
			unshort, err := Unshort(short)
			if err != nil {
				t.Fatalf("%s: Unshort failed: %v", tc.name, err)
			}
			if unshort != tc.original {
				t.Errorf("%s failed: original %s, short %s, unshort %s", tc.name, tc.original, short, unshort)
			}
			t.Logf("%s: %s -> %s (len %d -> %d)", tc.name, tc.original, short, len(tc.original), len(short))
		})
	}
}

func TestShortUnshortFormatted(t *testing.T) {
	// Test that Short handles formatted input (with hyphens)
	original := "20250831-151133-000012-345678-901234" // conceptual HumanUid with hyphens
	// Clean version of above is 32 digits
	clean := "20250831151133000012345678901234"

	short, err := Short(original)
	if err != nil {
		t.Fatalf("Short failed: %v", err)
	}
	unshort, err := Unshort(short)
	if err != nil {
		t.Fatalf("Unshort failed: %v", err)
	}

	if unshort != clean {
		t.Errorf("Formatted HumanUid failed: got %s, want %s", unshort, clean)
	}
}

func TestShortUnshortAlreadyDigits(t *testing.T) {
	// If it's already digits, Unshort should just return it (heuristic)
	id := "99999999999999999999" // 20 digits
	res, err := Unshort(id)
	if err != nil {
		t.Fatalf("Unshort failed for already numeric string: %v", err)
	}
	if res != id {
		t.Errorf("Unshort failed for already numeric string: got %s, want %s", res, id)
	}
}

func TestShortNonNumeric(t *testing.T) {
	// Should return error
	s := "not-an-id"
	_, err := Short(s)
	if err == nil {
		t.Errorf("Short should have failed for non-numeric string")
	}
}
