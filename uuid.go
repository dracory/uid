package uid

import (
	"strings"

	"github.com/google/uuid"
)

func Uuid() string {
	return strings.ReplaceAll(UuidFormatted(), "-", "")
}

func UuidFormatted() string {
	return uuid.New().String()
}

func UuidV1() string {
	return strings.ReplaceAll(UuidV1Formatted(), "-", "")
}

func UuidV1Formatted() string {
	return uuid.Must(uuid.NewUUID()).String()
}

func UuidV3(namespace string, data []byte) (string, error) {
	uid, err := UuidV3Formatted(namespace, data)

	if err != nil {
		return "", err
	}

	return strings.ReplaceAll(uid, "-", ""), nil
}

func UuidV3Formatted(namespace string, data []byte) (string, error) {
	id, err := uuid.FromBytes([]byte(namespace))

	if err != nil {
		return "", err
	}

	return uuid.Must(uuid.NewMD5(id, data), nil).String(), nil
}

func UuidV4() string {
	return strings.ReplaceAll(UuidV4Formatted(), "-", "")
}

func UuidV4Formatted() string {
	return uuid.New().String()
}

func UuidV5(namespace string, data []byte) (string, error) {
	uid, err := UuidV5Formatted(namespace, data)

	if err != nil {
		return "", err
	}

	return strings.ReplaceAll(uid, "-", ""), nil
}

func UuidV5Formatted(namespace string, data []byte) (string, error) {
	id, err := uuid.FromBytes([]byte(namespace))

	if err != nil {
		return "", err
	}

	return uuid.NewSHA1(id, data).String(), nil
}

func UuidV6() string {
	return strings.ReplaceAll(UuidV6Formatted(), "-", "")
}

func UuidV6Formatted() string {
	return uuid.Must(uuid.NewV6()).String()
}

func UuidV7() string {
	return strings.ReplaceAll(UuidV7Formatted(), "-", "")
}

func UuidV7Formatted() string {
	return uuid.Must(uuid.NewV7()).String()
}
