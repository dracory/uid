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
    human := uid.HumanUid()      // length: 32

    // NanoUid generates a UID (23 digits)
    // Format: YYYYMMDD-HHMMSS-MMMMMM-NNN
    nano := uid.NanoUid()        // length: 23

    // MicroUid generates a UID (20 digits)
    // Format: YYYYMMDD-HHMMSS-MMMMMM
    micro := uid.MicroUid()      // length: 20

    // SecUid generates a UID (14 digits)
    // Format: YYYYMMDD-HHMMSS
    sec := uid.SecUid()          // length: 14

    // Unix timestamps as strings
    ts := uid.Timestamp()        // seconds, length: 10
    tsu := uid.TimestampMicro()  // microseconds, length: 16
    tsn := uid.TimestampNano()   // nanoseconds, length: 19

    // UUIDs (via github.com/google/uuid)
    u := uid.Uuid()              // v4 without hyphens, length: 32
    uf := uid.UuidFormatted()    // v4 with hyphens, length: 36

    fmt.Println(human, nano, micro, sec, ts, tsu, tsn, u, uf)
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

All UUID functions are thin wrappers around github.com/google/uuid.

- Uuid() → v4 without hyphens
  Example: 550e8400e29b41d4a716446655440000 (length: 32)

- UuidFormatted() → v4 with hyphens
  Example: 550e8400-e29b-41d4-a716-446655440000 (length: 36)

- UuidV1() / UuidV1Formatted() → version 1 (time-based)
  Examples: 6ba7b8109dad11d180b400c04fd430c8 (32) • 6ba7b810-9dad-11d1-80b4-00c04fd430c8 (36)

- UuidV3(namespace, data) / UuidV3Formatted(namespace, data) → version 3 (MD5 name-based)
  Examples: 3d813cbb47fb32ba91df831e1593ac29 (32) • 3d813cbb-47fb-32ba-91df-831e1593ac29 (36)

- UuidV4() / UuidV4Formatted() → version 4 (random)
  Examples: 550e8400e29b41d4a716446655440000 (32) • 550e8400-e29b-41d4-a716-446655440000 (36)

- UuidV5(namespace, data) / UuidV5Formatted(namespace, data) → version 5 (SHA-1 name-based)
  Examples: 21f7f8de80515b8986800195ef798b6a (32) • 21f7f8de-8051-5b89-8680-0195ef798b6a (36)

- UuidV6() / UuidV6Formatted() → version 6 (time-ordered)
  Examples: 1ed0c9e48f7b6b2c9c3b6a6c7a9d5e12 (32) • 1ed0c9e4-8f7b-6b2c-9c3b-6a6c7a9d5e12 (36)

- UuidV7() / UuidV7Formatted() → version 7 (Unix time-based)
  Examples: 01890f5f3d9c7a0e8a7b6c5d4e3f2a10 (32) • 01890f5f-3d9c-7a0e-8a7b-6c5d4e3f2a10 (36)

## Change Log
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
