package uid

import (
	"strings"

	"github.com/google/uuid"
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
	return uuid.New().String()
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
	return uuid.Must(uuid.NewUUID()).String()
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
	id, err := uuid.FromBytes([]byte(namespace))

	if err != nil {
		return "", err
	}

	return uuid.Must(uuid.NewMD5(id, data), nil).String(), nil
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
	return uuid.New().String()
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
	id, err := uuid.FromBytes([]byte(namespace))

	if err != nil {
		return "", err
	}

	return uuid.NewSHA1(id, data).String(), nil
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
	return uuid.Must(uuid.NewV6()).String()
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
	return uuid.Must(uuid.NewV7()).String()
}
