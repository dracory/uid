package uid

import "testing"

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
