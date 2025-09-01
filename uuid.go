package uid

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"net"
	"strings"
	"sync"
	"time"
)

// Uuid returns a random UUID (version 4) without hyphens.
//
// Example: 550e8400e29b41d4a716446655440000 (length: 32)
//
// https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)
//
// Parameters:
// - None
//
// Returns:
// - A random UUID (version 4) without hyphens
func Uuid() string {
	return strings.ReplaceAll(UuidFormatted(), "-", "")
}

// UuidFormatted returns a random UUID (version 4) string with hyphens.
//
// Example: 550e8400-e29b-41d4-a716-446655440000 (length: 36)
//
// https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)
//
// Parameters:
// - None
//
// Returns:
// - A random UUID (version 4) with hyphens
func UuidFormatted() string {
	return bytesToUUIDString(newV4(), true)
}

// UuidV1 returns a version 1 (time-based) UUID without hyphens.
//
// Example: 6ba7b8109dad11d180b400c04fd430c8 (length: 32)
//
// https://en.wikipedia.org/wiki/Universally_unique_identifier#Versions_1_and_6_(date-time_and_MAC_address)
//
// Parameters:
// - None
//
// Returns:
// - A UUID v1 (time-based) without hyphens
func UuidV1() string {
	return strings.ReplaceAll(UuidV1Formatted(), "-", "")
}

// UuidV1Formatted returns a version 1 (time-based) UUID with hyphens.
//
// Example: 6ba7b810-9dad-11d1-80b4-00c04fd430c8 (length: 36)
//
// https://en.wikipedia.org/wiki/Universally_unique_identifier#Versions_1_and_6_(date-time_and_MAC_address)
//
// Parameters:
// - None
//
// Returns:
// - A UUID v1 (time-based) with hyphens
func UuidV1Formatted() string {
	return bytesToUUIDString(newV1(), true)
}

// UuidV3 returns a version 3 (MD5 name-based) UUID without hyphens.
// Provide a 16-byte namespace UUID and arbitrary data.
//
// Example (no hyphens): 3d813cbb47fb32ba91df831e1593ac29 (length: 32)
//
// https://en.wikipedia.org/wiki/Universally_unique_identifier#Versions_3_and_5_(namespace_name-based)
//
// Parameters:
// - namespace: a 16-byte UUID (as bytes) used as the namespace
// - data: the name bytes to hash
//
// Returns:
// - The UUID v3 as a 32-character string without hyphens, or an error
func UuidV3(namespace string, data []byte) (string, error) {
	uid, err := UuidV3Formatted(namespace, data)

	if err != nil {
		return "", err
	}

	return strings.ReplaceAll(uid, "-", ""), nil
}

// UuidV3Formatted returns a version 3 (MD5 name-based) UUID with hyphens.
// Provide a 16-byte namespace UUID and arbitrary data.
//
// Example: 3d813cbb-47fb-32ba-91df-831e1593ac29 (length: 36)
//
// https://en.wikipedia.org/wiki/Universally_unique_identifier#Versions_3_and_5_(namespace_name-based)
//
// Parameters:
// - namespace: a 16-byte UUID (as bytes) used as the namespace
// - data: the name bytes to hash
//
// Returns:
// - The UUID v3 as a 36-character string with hyphens, or an error
func UuidV3Formatted(namespace string, data []byte) (string, error) {
	ns := []byte(namespace)
	if len(ns) != 16 {
		return "", errors.New("namespace must be 16 bytes")
	}
	h := md5.New()
	h.Write(ns)
	h.Write(data)
	sum := h.Sum(nil)[:16]
	setVersion(sum, 3)
	setVariantRFC4122(sum)
	return bytesToUUIDString(sum, true), nil
}

// UuidV4 returns a random UUID (version 4) without hyphens.
//
// Example: 550e8400e29b41d4a716446655440000 (length: 32)
//
// https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)
//
// Parameters:
// - None
//
// Returns:
// - A random UUID (version 4) without hyphens
func UuidV4() string {
	return strings.ReplaceAll(UuidV4Formatted(), "-", "")
}

// UuidV4Formatted returns a random UUID (version 4) with hyphens.
//
// Example: 550e8400-e29b-41d4-a716-446655440000 (length: 36)
//
// https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_4_(random)
//
// Parameters:
// - None
//
// Returns:
// - A random UUID (version 4) with hyphens
func UuidV4Formatted() string {
	return bytesToUUIDString(newV4(), true)
}

// UuidV5 returns a version 5 (SHA-1 name-based) UUID without hyphens.
// Provide a 16-byte namespace UUID and arbitrary data.
//
// Example (no hyphens): 21f7f8de80515b8986800195ef798b6a (length: 32)
//
// https://en.wikipedia.org/wiki/Universally_unique_identifier#Versions_3_and_5_(namespace_name-based)
//
// Parameters:
// - namespace: a 16-byte UUID (as bytes) used as the namespace
// - data: the name bytes to hash
//
// Returns:
// - The UUID v5 as a 32-character string without hyphens, or an error
func UuidV5(namespace string, data []byte) (string, error) {
	uid, err := UuidV5Formatted(namespace, data)

	if err != nil {
		return "", err
	}

	return strings.ReplaceAll(uid, "-", ""), nil
}

// UuidV5Formatted returns a version 5 (SHA-1 name-based) UUID with hyphens.
// Provide a 16-byte namespace UUID and arbitrary data.
//
// Example: 21f7f8de-8051-5b89-8680-0195ef798b6a (length: 36)
//
// https://en.wikipedia.org/wiki/Universally_unique_identifier#Versions_3_and_5_(namespace_name-based)
//
// Parameters:
// - namespace: a 16-byte UUID (as bytes) used as the namespace
// - data: the name bytes to hash
//
// Returns:
// - The UUID v5 as a 36-character string with hyphens, or an error
func UuidV5Formatted(namespace string, data []byte) (string, error) {
	ns := []byte(namespace)
	if len(ns) != 16 {
		return "", errors.New("namespace must be 16 bytes")
	}
	h := sha1.New()
	h.Write(ns)
	h.Write(data)
	sum := h.Sum(nil)[:16]
	setVersion(sum, 5)
	setVariantRFC4122(sum)
	return bytesToUUIDString(sum, true), nil
}

// UuidV6 returns a version 6 (time-ordered) UUID without hyphens.
//
// Example: 1ed0c9e48f7b6b2c9c3b6a6c7a9d5e12 (length: 32)
//
// Draft: https://en.wikipedia.org/wiki/Universally_unique_identifier#Versions_1_and_6_(date-time_and_MAC_address)
//
// Parameters:
// - None
//
// Returns:
// - A UUID v6 (time-ordered) without hyphens
func UuidV6() string {
	return strings.ReplaceAll(UuidV6Formatted(), "-", "")
}

// UuidV6Formatted returns a version 6 (time-ordered) UUID with hyphens.
//
// Example: 1ed0c9e4-8f7b-6b2c-9c3b-6a6c7a9d5e12 (length: 36)
//
// Draft:https://en.wikipedia.org/wiki/Universally_unique_identifier#Versions_1_and_6_(date-time_and_MAC_address)
//
// Parameters:
// - None
//
// Returns:
// - A UUID v6 (time-ordered) with hyphens
func UuidV6Formatted() string {
	return bytesToUUIDString(newV6(), true)
}

// UuidV7 returns a version 7 (Unix time-based) UUID without hyphens.
//
// Example: 01890f5f3d9c7a0e8a7b6c5d4e3f2a10 (length: 32)
//
// Draft: https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_7_(timestamp_and_random)
//
// Parameters:
// - None
//
// Returns:
// - A UUID v7 (Unix time-based) without hyphens
func UuidV7() string {
	return strings.ReplaceAll(UuidV7Formatted(), "-", "")
}

// UuidV7Formatted returns a version 7 (Unix time-based) UUID with hyphens.
//
// Example: 01890f5f-3d9c-7a0e-8a7b-6c5d4e3f2a10 (length: 36)
//
// Draft: https://en.wikipedia.org/wiki/Universally_unique_identifier#Version_7_(timestamp_and_random)
//
// Parameters:
// - None
//
// Returns:
// - A UUID v7 (Unix time-based) with hyphens
func UuidV7Formatted() string {
	return bytesToUUIDString(newV7(), true)
}

// ---- Internal implementation ----

var (
	onceInit   sync.Once
	nodeIDData [6]byte
	clockSeq   uint16 // 14-bit
	mu         sync.Mutex
	lastTime   uint64 // 100-ns intervals since 1582
)

const gregorianToUnix100ns = uint64(122192928000000000)

func initState() {
	// Initialize node ID
	if nid, ok := systemNodeID(); ok {
		copy(nodeIDData[:], nid)
	} else {
		// Random multicast node per RFC 4122
		if _, err := rand.Read(nodeIDData[:]); err == nil {
			nodeIDData[0] |= 0x01 // multicast bit
		}
	}
	// Initialize clock sequence randomly (14-bit)
	var b [2]byte
	if _, err := rand.Read(b[:]); err == nil {
		clockSeq = binary.BigEndian.Uint16(b[:]) & 0x3FFF
	} else {
		clockSeq = uint16(time.Now().UnixNano()) & 0x3FFF
	}
}

func systemNodeID() ([]byte, bool) {
	ifs, err := net.Interfaces()
	if err != nil {
		return nil, false
	}
	for _, iface := range ifs {
		hw := iface.HardwareAddr
		if len(hw) == 6 {
			b := make([]byte, 6)
			copy(b, hw)
			return b, true
		}
	}
	return nil, false
}

func now100ns() uint64 {
	ns := uint64(time.Now().UnixNano())
	return ns/100 + gregorianToUnix100ns
}

func setVariantRFC4122(b []byte) {
	b[8] &= 0x3F
	b[8] |= 0x80 // 10xx xxxx
}

func setVersion(b []byte, ver int) {
	b[6] &= 0x0F
	b[6] |= byte(ver<<4) & 0xF0
}

func newV4() []byte {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		// fallback: timestamp-based randomness
		binary.BigEndian.PutUint64(b[0:8], uint64(time.Now().UnixNano()))
		binary.BigEndian.PutUint64(b[8:16], uint64(time.Now().UnixNano()))
	}
	setVersion(b, 4)
	setVariantRFC4122(b)
	return b
}

func newV1() []byte {
	onceInit.Do(initState)
	b := make([]byte, 16)

	mu.Lock()
	t := now100ns()
	if t <= lastTime {
		clockSeq = (clockSeq + 1) & 0x3FFF
	}
	lastTime = t
	cs := clockSeq
	mu.Unlock()

	// time fields per RFC 4122
	tl := uint32(t & 0xFFFFFFFF)
	tm := uint16((t >> 32) & 0xFFFF)
	th := uint16((t >> 48) & 0x0FFF)
	th |= 0x1000 // version 1

	binary.BigEndian.PutUint32(b[0:4], tl)
	binary.BigEndian.PutUint16(b[4:6], tm)
	binary.BigEndian.PutUint16(b[6:8], th)

	// clock seq with variant
	b[8] = byte((cs>>8)&0x3F) | 0x80 // variant 10
	b[9] = byte(cs)

	copy(b[10:], nodeIDData[:])
	return b
}

func newV6() []byte {
	onceInit.Do(initState)
	b := make([]byte, 16)

	mu.Lock()
	t := now100ns()
	if t <= lastTime {
		clockSeq = (clockSeq + 1) & 0x3FFF
	}
	lastTime = t
	cs := clockSeq
	mu.Unlock()

	// Reorder v1 timestamp into v6 (time-ordered) layout
	th := uint32(t >> 28)                 // top 32 bits
	tm := uint16((t >> 12) & 0xFFFF)      // next 16 bits
	tl := uint16(t & 0x0FFF)              // low 12 bits
	tl |= 0x6000                          // set version 6

	binary.BigEndian.PutUint32(b[0:4], th)
	binary.BigEndian.PutUint16(b[4:6], tm)
	binary.BigEndian.PutUint16(b[6:8], tl)

	// clock seq with variant
	b[8] = byte((cs>>8)&0x3F) | 0x80 // variant 10
	b[9] = byte(cs)

	copy(b[10:], nodeIDData[:])
	return b
}

func newV7() []byte {
	b := make([]byte, 16)
	// 48-bit Unix ms timestamp
	ts := uint64(time.Now().UnixMilli())
	b[0] = byte(ts >> 40)
	b[1] = byte(ts >> 32)
	b[2] = byte(ts >> 24)
	b[3] = byte(ts >> 16)
	b[4] = byte(ts >> 8)
	b[5] = byte(ts)

	// 12 bits random (A), 62 bits random (B)
	var r [10]byte
	if _, err := rand.Read(r[:]); err != nil {
		// fallback
		binary.BigEndian.PutUint64(r[2:], uint64(time.Now().UnixNano()))
	}

	// set version 7: upper nibble of b[6]
	b[6] = 0x70 | (r[0] & 0x0F)
	b[7] = r[1]

	// variant in b[8]
	b[8] = (r[2] & 0x3F) | 0x80
	copy(b[9:], r[3:])
	return b
}

func bytesToUUIDString(b []byte, withHyphens bool) string {
	if !withHyphens {
		dst := make([]byte, hex.EncodedLen(len(b)))
		hex.Encode(dst, b)
		return string(dst)
	}
	// 8-4-4-4-12
	hexstr := make([]byte, hex.EncodedLen(len(b)))
	hex.Encode(hexstr, b)
	// insert hyphens
	out := make([]byte, 36)
	copy(out[0:8], hexstr[0:8])
	out[8] = '-'
	copy(out[9:13], hexstr[8:12])
	out[13] = '-'
	copy(out[14:18], hexstr[12:16])
	out[18] = '-'
	copy(out[19:23], hexstr[16:20])
	out[23] = '-'
	copy(out[24:36], hexstr[20:32])
	return string(out)
}
