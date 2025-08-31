package uid

import "testing"

// helper to assert UUID length and version nibble
func assertLenAndVersion(t *testing.T, s string, wantLen int, wantVersion byte, withHyphens bool) {
    t.Helper()
    if len(s) != wantLen {
        t.Fatalf("length = %d, want %d; value=%s", len(s), wantLen, s)
    }
    idx := 12 // position of version nibble in hex when no hyphens
    if withHyphens {
        idx = 14 // first char of 3rd group in 8-4-4-4-12
    }
    if s[idx] != wantVersion {
        t.Fatalf("version nibble at index %d = %c, want %c; value=%s", idx, s[idx], wantVersion, s)
    }
}

func TestUuid(t *testing.T) {
    uuid1 := Uuid()
    uuid2 := Uuid()

    if uuid1 == "" {
        t.Fatal("Uuid must not be null")
    }

    if len(uuid1) != 32 {
        t.Fatal("Uuid length must be 32 characters, found: ", len(uuid1))
    }

    if uuid1 == uuid2 {
        t.Fatal("Uuid 1 and Timestamp 2 must not be the same")
    }
}

func TestUuidV1(t *testing.T) {
    a := UuidV1()
    b := UuidV1()
    if a == "" || b == "" {
        t.Fatal("UuidV1 must not be empty")
    }
    if a == b {
        t.Fatal("UuidV1 values must differ")
    }
    assertLenAndVersion(t, a, 32, '1', false)
}

func TestUuidV1Formatted(t *testing.T) {
    a := UuidV1Formatted()
    b := UuidV1Formatted()
    if a == "" || b == "" {
        t.Fatal("UuidV1Formatted must not be empty")
    }
    if a == b {
        t.Fatal("UuidV1Formatted values must differ")
    }
    assertLenAndVersion(t, a, 36, '1', true)
}

func TestUuidV3(t *testing.T) {
    ns := string([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16})
    got, err := UuidV3(ns, []byte("name"))
    if err != nil {
        t.Fatalf("UuidV3 error: %v", err)
    }
    assertLenAndVersion(t, got, 32, '3', false)
}

func TestUuidV3Formatted(t *testing.T) {
    ns := string([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16})
    got, err := UuidV3Formatted(ns, []byte("name"))
    if err != nil {
        t.Fatalf("UuidV3Formatted error: %v", err)
    }
    assertLenAndVersion(t, got, 36, '3', true)
}

func TestUuidV3_InvalidNamespace(t *testing.T) {
    // invalid length namespace -> expect error
    _, err := UuidV3("short", []byte("name"))
    if err == nil {
        t.Fatal("UuidV3 expected error for invalid namespace length")
    }
}

func TestUuidV4_Explicit(t *testing.T) {
    a := UuidV4()
    b := UuidV4()
    if a == "" || b == "" {
        t.Fatal("UuidV4 must not be empty")
    }
    if a == b {
        t.Fatal("UuidV4 values must differ")
    }
    assertLenAndVersion(t, a, 32, '4', false)
}

func TestUuidV4Formatted_Explicit(t *testing.T) {
    a := UuidV4Formatted()
    b := UuidV4Formatted()
    if a == "" || b == "" {
        t.Fatal("UuidV4Formatted must not be empty")
    }
    if a == b {
        t.Fatal("UuidV4Formatted values must differ")
    }
    assertLenAndVersion(t, a, 36, '4', true)
}

func TestUuidV5(t *testing.T) {
    ns := string([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16})
    got, err := UuidV5(ns, []byte("name"))
    if err != nil {
        t.Fatalf("UuidV5 error: %v", err)
    }
    assertLenAndVersion(t, got, 32, '5', false)
}

func TestUuidV5Formatted(t *testing.T) {
    ns := string([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16})
    got, err := UuidV5Formatted(ns, []byte("name"))
    if err != nil {
        t.Fatalf("UuidV5Formatted error: %v", err)
    }
    assertLenAndVersion(t, got, 36, '5', true)
}

func TestUuidV5_InvalidNamespace(t *testing.T) {
    _, err := UuidV5("short", []byte("name"))
    if err == nil {
        t.Fatal("UuidV5 expected error for invalid namespace length")
    }
}

func TestUuidV6(t *testing.T) {
    a := UuidV6()
    b := UuidV6()
    if a == "" || b == "" {
        t.Fatal("UuidV6 must not be empty")
    }
    if a == b {
        t.Fatal("UuidV6 values must differ")
    }
    assertLenAndVersion(t, a, 32, '6', false)
}

func TestUuidV6Formatted(t *testing.T) {
    a := UuidV6Formatted()
    b := UuidV6Formatted()
    if a == "" || b == "" {
        t.Fatal("UuidV6Formatted must not be empty")
    }
    if a == b {
        t.Fatal("UuidV6Formatted values must differ")
    }
    assertLenAndVersion(t, a, 36, '6', true)
}

func TestUuidV7(t *testing.T) {
    a := UuidV7()
    b := UuidV7()
    if a == "" || b == "" {
        t.Fatal("UuidV7 must not be empty")
    }
    if a == b {
        t.Fatal("UuidV7 values must differ")
    }
    assertLenAndVersion(t, a, 32, '7', false)
}

func TestUuidV7Formatted(t *testing.T) {
    a := UuidV7Formatted()
    b := UuidV7Formatted()
    if a == "" || b == "" {
        t.Fatal("UuidV7Formatted must not be empty")
    }
    if a == b {
        t.Fatal("UuidV7Formatted values must differ")
    }
    assertLenAndVersion(t, a, 36, '7', true)
}

func TestUuidFormatted(t *testing.T) {
	uuid1 := UuidFormatted()
	uuid2 := UuidFormatted()

	if uuid1 == "" {
		t.Fatal("Uuid must not be null")
	}

	if len(uuid1) != 36 {
		t.Fatal("Uuid length must be 36 characters, found: ", len(uuid1))
	}

	if uuid1 == uuid2 {
		t.Fatal("Uuid 1 and Timestamp 2 must not be the same")
	}
}
