package config

import (
	"fmt"
	"strconv"
	"strings"
)

const MAJOR_VERSION = 0

type VersionNumber struct {
	Major int
	Minor int
}

func (v VersionNumber) IsCompatibleVersion() bool {
	return v.Major == MAJOR_VERSION
}

func (v VersionNumber) NextMinorVersion() VersionNumber {
	return VersionNumber{
		Major: v.Major,
		Minor: v.Minor + 1,
	}
}

func (v VersionNumber) MarshalText() ([]byte, error) {
	return fmt.Appendf(nil, "%d.%d", v.Major, v.Minor), nil
}

func (v *VersionNumber) UnmarshalText(data []byte) error {
	parts := strings.Split(string(data), ".")
	if len(parts) != 2 {
		return fmt.Errorf("invalid version %q", data)
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("invalid major version: %w", err)
	}

	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("invalid minor version: %w", err)
	}

	v.Major = major
	v.Minor = minor
	return nil
}
