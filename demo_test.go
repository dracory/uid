package uid_test

import (
	"bytes"
	"fmt"
	"testing"
	"text/tabwriter"

	"github.com/dracory/uid"
)

func printTable(name string, original string) string {
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', tabwriter.Debug)
	fmt.Fprintln(w, "OPERATION\tRESULT\tLENGTH")
	fmt.Fprintf(w, "Original %s\t%s\t%d\n", name, original, len(original))
	s16, _ := uid.ShortenBase16(original)
	fmt.Fprintf(w, "ShortenBase16\t%s\t%d\n", s16, len(s16))
	s32, _ := uid.ShortenBase32(original)
	fmt.Fprintf(w, "ShortenBase32\t%s\t%d\n", s32, len(s32))
	scrock, _ := uid.ShortenCrockford(original)
	fmt.Fprintf(w, "ShortenCrockford\t%s\t%d\n", scrock, len(scrock))
	sz32, _ := uid.ShortenZBase32(original)
	fmt.Fprintf(w, "ShortenZBase32\t%s\t%d\n", sz32, len(sz32))
	s36, _ := uid.ShortenBase36(original)
	fmt.Fprintf(w, "ShortenBase36\t%s\t%d\n", s36, len(s36))
	s58, _ := uid.ShortenBase58(original)
	fmt.Fprintf(w, "ShortenBase58\t%s\t%d\n", s58, len(s58))
	s62, _ := uid.ShortenBase62(original)
	fmt.Fprintf(w, "ShortenBase62\t%s\t%d\n", s62, len(s62))
	s64, _ := uid.ShortenBase64(original)
	fmt.Fprintf(w, "ShortenBase64\t%s\t%d\n", s64, len(s64))
	w.Flush()
	return buf.String()
}

func TestSecUid(t *testing.T) {
	secUid := uid.SecUid()
	t.Log("\n" + printTable("SecUid", secUid))
}

func TestMicroUid(t *testing.T) {
	microUid := uid.MicroUid()
	t.Log("\n" + printTable("MicroUid", microUid))

	// Roundtrip checks
	s16, _ := uid.ShortenBase16(microUid)
	if res, _ := uid.UnshortenBase16(s16); res != microUid {
		t.Error("UnshortenBase16 failed")
	}
	s32, _ := uid.ShortenBase32(microUid)
	if res, _ := uid.UnshortenBase32(s32); res != microUid {
		t.Error("UnshortenBase32 failed")
	}
	scrock, _ := uid.ShortenCrockford(microUid)
	if res, _ := uid.UnshortenCrockford(scrock); res != microUid {
		t.Error("UnshortenCrockford failed")
	}
	s64, _ := uid.ShortenBase64(microUid)
	if res, _ := uid.UnshortenBase64(s64); res != microUid {
		t.Error("UnshortenBase64 failed")
	}
}

func TestNanoUid(t *testing.T) {
	nanoUid := uid.NanoUid()
	t.Log("\n" + printTable("NanoUid", nanoUid))
}

func TestUidShorten(t *testing.T) {
	humanUid := uid.HumanUid()
	t.Log("\n" + printTable("HumanUid", humanUid))

	// Roundtrip checks
	s16, _ := uid.ShortenBase16(humanUid)
	if res, _ := uid.UnshortenBase16(s16); res != humanUid {
		t.Error("UnshortenBase16 failed")
	}
	s32, _ := uid.ShortenBase32(humanUid)
	if res, _ := uid.UnshortenBase32(s32); res != humanUid {
		t.Error("UnshortenBase32 failed")
	}
	scrock, _ := uid.ShortenCrockford(humanUid)
	if res, _ := uid.UnshortenCrockford(scrock); res != humanUid {
		t.Error("UnshortenCrockford failed")
	}
	s64, _ := uid.ShortenBase64(humanUid)
	if res, _ := uid.UnshortenBase64(s64); res != humanUid {
		t.Error("UnshortenBase64 failed")
	}
}

func TestTimestamp(t *testing.T) {
	ts := uid.Timestamp()
	t.Log("\n" + printTable("Timestamp", ts))
}

func TestTimestampMicro(t *testing.T) {
	tsu := uid.TimestampMicro()
	t.Log("\n" + printTable("TimestampMicro", tsu))
}

func TestTimestampNano(t *testing.T) {
	tsn := uid.TimestampNano()
	t.Log("\n" + printTable("TimestampNano", tsn))
}
