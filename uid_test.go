package uid

import (
	"testing"
)

// helper to assert expected length and hyphen positions
func assertHyphenPositions(t *testing.T, s string, wantLen int, positions []int) {
	t.Helper()
	if len(s) != wantLen {
		t.Fatalf("length = %d, want %d; value=%s", len(s), wantLen, s)
	}
	for _, i := range positions {
		if i < 0 || i >= len(s) {
			t.Fatalf("hyphen index %d out of range for %q", i, s)
		}
		if s[i] != '-' {
			t.Fatalf("expected hyphen at index %d, got %q; value=%s", i, s[i], s)
		}
	}
}

func TestHumanUid(t *testing.T) {
	humanUid := HumanUid()
	humanUid2 := HumanUid()

	if humanUid == "" {
		t.Fatal("Human UID must not be null")
	}

	if len(humanUid) != 32 {
		t.Fatal("Human UID length must be 32 characters")
	}

	if humanUid == humanUid2 {
		t.Fatal("Human UID and Human UID 2 must not be the same")
	}

	if humanUid > humanUid2 {
		t.Fatal("Human UID 1 must be smaller than Human UID 2")
	}
}

func TestHumanUidFormatted(t *testing.T) {
	hf := HumanUid(true)
	if hf == "" {
		t.Fatal("Human UID (formatted) must not be null")
	}
	// formatted variant: groups 8-4-4-16 => hyphens at 8,13,18; total length 35
	assertHyphenPositions(t, hf, 35, []int{8, 13, 18})
}

func TestMicroUid(t *testing.T) {
	microUid := MicroUid()
	microUid2 := MicroUid()

	if microUid == "" {
		t.Fatal("Micro UID must not be null")
	}

	if len(microUid) != 20 {
		t.Fatal("Micro UID length must be 20 charcters")
	}

	if microUid == microUid2 {
		t.Fatal("Micro UID and Micro UID 2 must not be the same")
	}

	if microUid > microUid2 {
		t.Fatal("Micro UID 1 must be smaller than Micro UID 2")
	}
}

func TestMicroUidFormatted(t *testing.T) {
	mf := MicroUid(true)
	if mf == "" {
		t.Fatal("Micro UID (formatted) must not be null")
	}
	// formatted variant: groups 8-6-6 => hyphens at 8,15; total length 22
	assertHyphenPositions(t, mf, 22, []int{8, 15})
}

func TestNanoUid(t *testing.T) {
	nanoUid := NanoUid()
	nanoUid2 := NanoUid()

	if nanoUid == "" {
		t.Fatal("Nano UID must not be null")
	}

	if len(nanoUid) != 23 {
		t.Fatal("Nano UID length must be 23 charcters")
	}

	if nanoUid == nanoUid2 {
		t.Fatal("Nano UID and Nano UID 2 must not be the same")
	}

	if nanoUid > nanoUid2 {
		t.Fatal("Nano UID 1 must be smaller than Nano UID 2")
	}
}

func TestNanoUidFormatted(t *testing.T) {
	nf := NanoUid(true)
	if nf == "" {
		t.Fatal("Nano UID (formatted) must not be null")
	}
	// formatted variant: groups 8-6-6-3 => hyphens at 8,15,22; total length 26
	assertHyphenPositions(t, nf, 26, []int{8, 15, 22})
}

func TestSecUid(t *testing.T) {
	secUid := SecUid()
	// time.Sleep(time.Second) // as its a seconds based ID we need at least a second between the generation
	secUid2 := SecUid()

	if secUid == "" {
		t.Fatal("Sec UID must not be null")
	}

	if len(secUid) != 14 {
		t.Fatal("Sec UID length must be 14 characters")
	}

	if secUid == secUid2 {
		t.Fatal("Sec UID and sec UID 2 must not be the same")
	}

	if secUid > secUid2 {
		t.Fatal("Sec UID 1 must be smaller than sec UID 2")
	}
}

func TestSecUidFormatted(t *testing.T) {
	sf := SecUid(true)
	if sf == "" {
		t.Fatal("Sec UID (formatted) must not be null")
	}
	// formatted variant: groups 8-6 => hyphen at 8; total length 15
	assertHyphenPositions(t, sf, 15, []int{8})
}

func TestMinUid(t *testing.T) {
	u := MinUid()
	if u == "" {
		t.Fatal("Min UID must not be empty")
	}
	if len(u) != 12 {
		t.Fatalf("Min UID length must be 12 characters, got %d", len(u))
	}
}

func TestMinUidFormatted(t *testing.T) {
	uf := MinUid(true)
	if uf == "" {
		t.Fatal("Min UID (formatted) must not be empty")
	}
	// formatted variant: groups 8-4 => hyphen at 8; total length 13
	assertHyphenPositions(t, uf, 13, []int{8})
}

func TestHourUid(t *testing.T) {
	u := HourUid()
	if u == "" {
		t.Fatal("Hour UID must not be empty")
	}
	if len(u) != 10 {
		t.Fatalf("Hour UID length must be 10 characters, got %d", len(u))
	}
}

func TestHourUidFormatted(t *testing.T) {
	uf := HourUid(true)
	if uf == "" {
		t.Fatal("Hour UID (formatted) must not be empty")
	}
	// formatted variant: groups 8-2 => hyphen at 8; total length 11
	assertHyphenPositions(t, uf, 11, []int{8})
}

func TestDayUid(t *testing.T) {
	u := DayUid()
	if u == "" {
		t.Fatal("Day UID must not be empty")
	}
	if len(u) != 8 {
		t.Fatalf("Day UID length must be 8 characters, got %d", len(u))
	}
}

func TestDayUidFormatted(t *testing.T) {
	uf := DayUid(true)
	if uf == "" {
		t.Fatal("Day UID (formatted) must not be empty")
	}
	// formatted variant: groups 4-2-2 => hyphens at 4, 7; total length 10
	assertHyphenPositions(t, uf, 10, []int{4, 7})
}

func TestWeekUid(t *testing.T) {
	u := WeekUid()
	if u == "" {
		t.Fatal("Week UID must not be empty")
	}
	if len(u) != 6 {
		t.Fatalf("Week UID length must be 6 characters, got %d", len(u))
	}
}

func TestWeekUidFormatted(t *testing.T) {
	uf := WeekUid(true)
	if uf == "" {
		t.Fatal("Week UID (formatted) must not be empty")
	}
	// formatted variant: groups 4-2 => hyphen at 4; total length 7
	assertHyphenPositions(t, uf, 7, []int{4})
}

func TestMonthUid(t *testing.T) {
	u := MonthUid()
	if u == "" {
		t.Fatal("Month UID must not be empty")
	}
	if len(u) != 6 {
		t.Fatalf("Month UID length must be 6 characters, got %d", len(u))
	}
}

func TestMonthUidFormatted(t *testing.T) {
	uf := MonthUid(true)
	if uf == "" {
		t.Fatal("Month UID (formatted) must not be empty")
	}
	// formatted variant: groups 4-2 => hyphen at 4; total length 7
	assertHyphenPositions(t, uf, 7, []int{4})
}

func TestTimestamp(t *testing.T) {
	ts1 := Timestamp()
	ts2 := Timestamp()

	if ts1 == "" {
		t.Fatal("Timestamp must not be null")
	}

	if len(ts1) != 10 {
		t.Fatal("Timestamp length must be 10 characters, found: ", len(ts1))
	}

	if ts1 == ts2 {
		t.Fatal("Timestamp 1 and Timestamp 2 must not be the same")
	}

	if ts1 > ts2 {
		t.Fatal("Timestamp 1 must be smaller than Timestamp 2")
	}
}

func TestTimestampMicro(t *testing.T) {
	ts1 := TimestampMicro()
	ts2 := TimestampMicro()

	if ts1 == "" {
		t.Fatal("Timestamp must not be null")
	}

	if len(ts1) != 16 {
		t.Fatal("Timestamp length must be 16 characters, found: ", len(ts1))
	}

	if ts1 == ts2 {
		t.Fatal("Timestamp 1 and Timestamp 2 must not be the same")
	}

	if ts1 > ts2 {
		t.Fatal("Timestamp 1 must be smaller than Timestamp 2")
	}
}

func TestTimestampNano(t *testing.T) {
	ts1 := TimestampNano()
	ts2 := TimestampNano()

	if ts1 == "" {
		t.Fatal("Timestamp must not be null")
	}

	if len(ts1) != 19 {
		t.Fatal("Timestamp length must be 19 characters, found: ", len(ts1))
	}

	if ts1 == ts2 {
		t.Fatal("Timestamp 1 and Timestamp 2 must not be the same")
	}

	if ts1 > ts2 {
		t.Fatal("Timestamp 1 must be smaller than Timestamp 2")
	}
}
