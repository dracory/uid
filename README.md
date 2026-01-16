# UID (Unique ID) <a href="https://gitpod.io/#https://github.com/dracory/uid" style="float:right:"><img src="https://gitpod.io/button/open-in-gitpod.svg" alt="Open in Gitpod" loading="lazy"></a>

[![Tests Status](https://github.com/dracory/uid/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/dracory/uid/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/dracory/uid)](https://goreportcard.com/report/github.com/dracory/uid)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/dracory/uid)](https://pkg.go.dev/github.com/dracory/uid)

This package generates unique identifying strings. Largest attention is paid on human friendly unique identifiers (dated digits).

## Installation

```bash
go get -u github.com/dracory/uid
```

## Usage

```go
package main

import (
    "fmt"
    "github.com/dracory/uid"
)

func main() {
    // HumanUid generates a UID (32 digits)
    // Format: YYYYMMDD-HHMM-SSMM-MMMMNNNRRRRRRRRR
    human := uid.HumanUid()          // unformatted, length: 32
    humanF := uid.HumanUid(true)     // formatted (8-4-4-16), length: 35

    // NanoUid generates a UID (23 digits)
    // Format: YYYYMMDD-HHMMSS-MMMMMM-NNN
    nano := uid.NanoUid()            // unformatted, length: 23
    nanoF := uid.NanoUid(true)       // formatted (8-6-6-3), length: 26

    // MicroUid generates a UID (20 digits)
    // Format: YYYYMMDD-HHMMSS-MMMMMM
    micro := uid.MicroUid()          // unformatted, length: 20
    microF := uid.MicroUid(true)     // formatted (8-6-6), length: 22

    // SecUid generates a UID (14 digits)
    // Format: YYYYMMDD-HHMMSS
    sec := uid.SecUid()              // unformatted, length: 14
    secF := uid.SecUid(true)         // formatted (8-6), length: 15

    // Unix timestamps as strings
    ts := uid.Timestamp()            // seconds, length: 10
    tsu := uid.TimestampMicro()      // microseconds, length: 16
    tsn := uid.TimestampNano()       // nanoseconds, length: 19

    // UUIDs (implemented via standard library)
    u4 := uid.Uuid()                 // v4 unformatted, length: 32
    u4f := uid.Uuid(true)            // v4 formatted, length: 36
    v1 := uid.UuidV1()               // v1 unformatted, length: 32
    v1f := uid.UuidV1(true)          // v1 formatted, length: 36
    v3, _ := uid.UuidV3("1234567890abcdef", []byte("name"))     // v3 unformatted
    v3f, _ := uid.UuidV3("1234567890abcdef", []byte("name"), true)
    v5, _ := uid.UuidV5("1234567890abcdef", []byte("name"))     // v5 unformatted
    v5f, _ := uid.UuidV5("1234567890abcdef", []byte("name"), true)
    v6 := uid.UuidV6()               // v6 unformatted, length: 32
    v6f := uid.UuidV6(true)          // v6 formatted, length: 36
    v7 := uid.UuidV7()               // v7 unformatted, length: 32
    v7f := uid.UuidV7(true)          // v7 formatted, length: 36

    fmt.Println(human, humanF, nano, nanoF, micro, microF, sec, secF,
        ts, tsu, tsn, u4, u4f, v1, v1f, v3, v3f, v5, v5f, v6, v6f, v7, v7f)

    // ID Shortening
    id := "12345678901234567890"
    short, _ := uid.ShortenBase62(id)       // "1n9XvB8P6M"
    original, _ := uid.UnshortenBase62(short) // "12345678901234567890"
}
```

## Supported UID Types

It supports several types of unique identifiers. 

The type you want to use will usually depends on two considerations:

1. How random you want it to be? The longer the identifier, the more the chances of collision reduce
2. How long you want the identifier to be? The longer the identifier, reduces the readability, as well as the storage space to store it.

For most of the user cases a Micro UID (20 chars) should be fine. A human UID (32 chars) should be avoided where a human is involved as too "mind bogging" to work with.

1. Human UID (32 digits)

    Format: YYYYMMDD-HHMM-SSMM-MMMMNNNRRRRRRRRR

    2017111908492665991498485465 (with dashes: 20171119-0849-2665-991498485465)

2. Nano UID (23 digits)

    Format: YYYYMMDD-HHMMSS-MMMMMM-NNN

    Examples:

    20171119084926659914984 (with dashes: 20171119-084926-659914-984)

3. Micro UID (20 digits)

    Format: YYYYMMDD-HHMMSS-MMMMMM

    Examples:

    20171119084926659914 (with dashes: 20171119-084926-659914)

4. Seconds UID (14 digits)

    Format: YYYYMMDD-HHMMSS

    Examples:

    20171119084926 (with dashes: 20171119-084926)

5. Timestamp (10 digits)
    Unit timestamp, seconds precision

    Format: 1234567890

    Examples:

    1704524414


6. TimestampMicro (16 digits)
    Unit timestamp, microseconds precision

    Format: 1234567890123456

    Examples:

    1704524414548721

6. TimestampNano (19 digits)
    Unit timestamp, nanoseconds precision

    Format: 1234567890123456789

    Examples:

    1704524414548721308

7. Uuid (32 characters)
    Random V4 UUID. UUID (Universally Unique IDentifier), also known as GUID (Globally Unique IDentifier)

    Format: abcdef1234567890abcdef1234567890

    Examples:

    459e2999bd071151a23d643da42c2cc2

## UUID functions

UUIDs are implemented using only the Go standard library (no external deps).

- Uuid(formatted ...bool) → version 4 (random)
  Examples: 550e8400e29b41d4a716446655440000 (32) • 550e8400-e29b-41d4-a716-446655440000 (36)

- UuidV1(formatted ...bool) → version 1 (time-based)
  Examples: 6ba7b8109dad11d180b400c04fd430c8 (32) • 6ba7b810-9dad-11d1-80b4-00c04fd430c8 (36)

- UuidV3(namespace string, data []byte, formatted ...bool) → version 3 (MD5 name-based)
  Examples: 3d813cbb47fb32ba91df831e1593ac29 (32) • 3d813cbb-47fb-32ba-91df-831e1593ac29 (36)

- UuidV4(formatted ...bool) → version 4 (random)
  Examples: 550e8400e29b41d4a716446655440000 (32) • 550e8400-e29b-41d4-a716-446655440000 (36)

- UuidV5(namespace string, data []byte, formatted ...bool) → version 5 (SHA-1 name-based)
  Examples: 21f7f8de80515b8986800195ef798b6a (32) • 21f7f8de-8051-5b89-8680-0195ef798b6a (36)

- UuidV6(formatted ...bool) → version 6 (time-ordered)
  Examples: 1ed0c9e48f7b6b2c9c3b6a6c7a9d5e12 (32) • 1ed0c9e4-8f7b-6b2c-9c3b-6a6c7a9d5e12 (36)

- UuidV7(formatted ...bool) → version 7 (Unix time-based)
  Examples: 01890f5f3d9c7a0e8a7b6c5d4e3f2a10 (32) • 01890f5f-3d9c-7a0e-8a7b-6c5d4e3f2a10 (36)

## ID Shortening

The package provides functions to shorten numeric string IDs into more compact representations using various base encodings. This is useful for creating URL-friendly or human-readable IDs from large numbers.

| Method | Encoding | Alphabet | Case-Sensitive | Human-Facing |
| :--- | :--- | :--- | :--- | :--- |
| `ShortenBase16` | Hexadecimal | `0-9, a-f` | No | No |
| `ShortenBase32` | Base32 (RFC 4648) | `a-z, 2-7` | No | Yes |
| `ShortenCrockford`| Crockford Base32 | `0-9, A-Z (exc. I,L,O,U)` | No (Normalizes) | Best |
| `ShortenBase36` | Base36 | `0-9, a-z` | No | Yes |
| `ShortenBase58` | Base58 (Bitcoin) | `1-9, A-Z, a-z (exc. 0,O,I,l)` | Yes | Better |
| `ShortenBase62` | Base62 | `0-9, A-Z, a-z` | Yes | No |
| `ShortenBase64` | Base64 (URL-safe) | `A-Z, a-z, 0-9, -, _` | Yes | No |
| `ShortenZBase32` | z-base-32 | `y,b,n,d,r,f,g,8,e,j...` | No | Good |

Example:
```go
id := "12345678901234567890"
short, _ := uid.ShortenBase62(id)       // "1n9XvB8P6M"
original, _ := uid.UnshortenBase62(short) // "12345678901234567890"
```

### Comparison Table

Using a sample `HumanUid` as input:

| Operation | Result | Length |
| :--- | :--- | :--- |
| Original `HumanUid` | `20260116055547619570214289007495` | 32 |
| `ShortenBase16` | `ffb7f7358fa8ad8f0adb27e387` | 26 |
| `ShortenBase32` | `p7n7xgwh2rlmpblnspy4h` | 21 |
| `ShortenCrockford` | `FZDZQ6P7THBCF1BDJFRW7` | 21 |
| `ShortenZBase32` | `x9p9zgs84tmcxbmp1xah8` | 21 |
| `ShortenBase36` | `1ik90wbiqnpjge668igef` | 21 |
| `ShortenBase58` | `NJMxnyUkCcK1mbKaR4` | 18 |
| `ShortenBase62` | `6qz9D7Ih28OnrKPXj5` | 18 |
| `ShortenBase64` | `D_t_c1j6itjwrbJ-OH` | 18 |

Using a sample `SecUid` as input:

| Operation | Result | Length |
| :--- | :--- | :--- |
| Original `SecUid` | `20260116060140` | 14 |
| `ShortenBase16` | `126d2d0557ec` | 12 |
| `ShortenBase32` | `snuwqkv7m` | 9 |
| `ShortenCrockford` | `JDMPGANZC` | 9 |
| `ShortenZBase32` | `1pwsoki9c` | 9 |
| `ShortenBase36` | `76jda01yk` | 9 |
| `ShortenBase58` | `ABCY5ZqV` | 8 |
| `ShortenBase62` | `5kgp4HO8` | 8 |
| `ShortenBase64` | `Em0tBVfs` | 8 |

Using a sample `MicroUid` as input:

| Operation | Result | Length |
| :--- | :--- | :--- |
| Original `MicroUid` | `20260116060209548481` | 20 |
| `ShortenBase16` | `1192a6436ccfa0cc1` | 17 |
| `ShortenBase32` | `rskteg3gpudgb` | 13 |
| `ShortenCrockford` | `HJAK46V6FM361` | 13 |
| `ShortenZBase32` | `t1kurg5gxwdgb` | 13 |
| `ShortenBase36` | `49xd7r0cxnoht` | 13 |
| `ShortenBase58` | `p2fBzoDdpyN` | 11 |
| `ShortenBase62` | `O8dXskGYQfx` | 11 |
| `ShortenBase64` | `RkqZDbM-gzB` | 11 |

Using a sample `NanoUid` as input:

| Operation | Result | Length |
| :--- | :--- | :--- |
| Original `NanoUid` | `20260116060418792413315` | 23 |
| `ShortenBase16` | `44a4d97764168aad883` | 19 |
| `ShortenBase32` | `rfe3f3wifukvwed` | 15 |
| `ShortenCrockford` | `H54V5VP85MANP43` | 15 |
| `ShortenZBase32` | `tfr5f5sefwkisrd` | 15 |
| `ShortenBase36` | `3arqv3acnfg822r` | 15 |
| `ShortenBase58` | `EyqbSbbMe7wtW` | 13 |
| `ShortenBase62` | `6HLJqV3vQ4U1j` | 13 |
| `ShortenBase64` | `ESk2XdkFoqtiD` | 13 |

Using a sample `Timestamp` (seconds) as input:

| Operation | Result | Length |
| :--- | :--- | :--- |
| Original `Timestamp` | `1768543518` | 10 |
| `ShortenBase16` | `6969d51e` | 8 |
| `ShortenBase32` | `buwtvi6` | 7 |
| `ShortenCrockford` | `1MPKN8Y` | 7 |
| `ShortenZBase32` | `bwsuie6` | 7 |
| `ShortenBase36` | `t8y0wu` | 6 |
| `ShortenBase58` | `3hHFNd` | 6 |
| `ShortenBase62` | `1vgcxS` | 6 |
| `ShortenBase64` | `BpadUe` | 6 |

Using a sample `TimestampMicro` (microseconds) as input:

| Operation | Result | Length |
| :--- | :--- | :--- |
| Original `TimestampMicro` | `1768543534819239` | 16 |
| `ShortenBase16` | `6487b2129a7a7` | 13 |
| `ShortenBase32` | `bsipmqstj5h` | 11 |
| `ShortenCrockford` | `1J8FCGJK9X7` | 11 |
| `ShortenZBase32` | `b1exco1uj78` | 11 |
| `ShortenBase36` | `hew9om6uuf` | 10 |
| `ShortenBase58` | `Eoye75bHY` | 9 |
| `ShortenBase62` | `86CCRTsQx` | 9 |
| `ShortenBase64` | `GSHshKaen` | 9 |

Using a sample `TimestampNano` (nanoseconds) as input:

| Operation | Result | Length |
| :--- | :--- | :--- |
| Original `TimestampNano` | `1768543548681502200` | 19 |
| `ShortenBase16` | `188b20fcc4f83df8` | 16 |
| `ShortenBase32` | `brcza7tcpqppy` | 13 |
| `ShortenCrockford` | `1H2S0ZK2FGFFR` | 13 |
| `ShortenZBase32` | `btn3y9unxoxxa` | 13 |
| `ShortenBase36` | `dfpsd65np7nc` | 12 |
| `ShortenBase58` | `576wZ6w7qyR` | 11 |
| `ShortenBase62` | `26dwmwO2cbw` | 11 |
| `ShortenBase64` | `BiLIPzE-D34` | 11 |






## Change Log
2026.01.16 - Add ID shortening and unshortening (Base16,32,36,58,62,64,Crockford,ZBase32)

2025.09.01 - Add optional hyphen formatting

2025.08.31 - Move UUID functions/tests into separate files

2024.01.06 - Added Timestamp and Uuid functions

2021.12.19 - Master branch changed to main

2021.12.19 - Added tests

## Similar Packages

- https://github.com/jaevor/go-nanoid (random 21 characters)
- https://github.com/zheng-ji/goSnowFlake (timestamp-workerid-sequence)
- https://github.com/damdo/randid (random IDs)
- https://github.com/matthewmueller/uid.go (shorcodes)
- https://github.com/aohorodnyk/uid (random IDs)
- https://github.com/google/uuid (UUIDs)
- https://github.com/oklog/ulid (???)
- https://github.com/chilts/sid (serial IDs)
- https://datatracker.ietf.org/doc/html/draft-ietf-uuidrev-rfc4122bis (GUID 6?)
- https://github.com/jetpack-io/typeid (typeid)
