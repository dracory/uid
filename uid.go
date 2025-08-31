package uid

import (
	"crypto/rand"
	"strconv"
	"strings"
	"time"
)

// HumanUid generates a UUID (32 digits) Format: YYYYMMDD-HHMM-SSMM-MMMMNNNRRRRRRRRR
func HumanUid() string {
	time.Sleep(1 * time.Nanosecond)
	r, _ := rand.Prime(rand.Reader, 64)
	id := time.Now().UTC().Format("20060102150405.0000000")
	id = strings.ReplaceAll(id, ".", "")
	id += r.String()
	return id[0:32]
}

// NanoUid generates a UID (23 digits) Format: YYYYMMDD-HHMMSS-MMMMMM-NNN
func NanoUid() string {
	time.Sleep(time.Nanosecond) // as its a nanoseconds based ID we need at least a nanosecond between the generations to avoid collisions
	r, _ := rand.Prime(rand.Reader, 64)
	id := time.Now().UTC().Format("20060102150405.0000000")
	id = strings.ReplaceAll(id, ".", "")
	id += r.String()
	return id[0:23]
}

// MicroUid generates a UID (20 digits) Format: YYYYMMDD-HHMMSS-MMMMMM
func MicroUid() string {
	time.Sleep(time.Microsecond) // as its a microseconds based ID we need at least a microsecond between the generations to avoid collisions
	r, _ := rand.Prime(rand.Reader, 64)
	id := time.Now().UTC().Format("20060102150405.0000000")
	id = strings.ReplaceAll(id, ".", "")
	id += r.String()
	return id[0:20]
}

// SecUid generates UID (14 digits) Format: YYYYMMDD-HHMMSS
func SecUid() string {
	time.Sleep(time.Second) // as its a seconds based ID we need at least a second between the generations to avoid collisions
	r, _ := rand.Prime(rand.Reader, 64)
	id := time.Now().UTC().Format("20060102150405.0000000")
	id = strings.ReplaceAll(id, ".", "")
	id += r.String()
	return id[0:14]
}

func Timestamp() string {
	time.Sleep(time.Second) // as its a seconds based ID we need at least a second between the generations to avoid collisions
	now := time.Now().UTC().Unix()
	return strconv.FormatInt(now, 10)
}

func TimestampMicro() string {
	time.Sleep(time.Microsecond) // as its a microseconds based ID we need at least a microsecond between the generations to avoid collisions
	now := time.Now().UTC().UnixMicro()
	return strconv.FormatInt(now, 10)
}

func TimestampNano() string {
	time.Sleep(time.Nanosecond) // as its a nanoseconds based ID we need at least a nanosecond between the generations to avoid collisions
	now := time.Now().UTC().UnixNano()
	return strconv.FormatInt(now, 10)
}
