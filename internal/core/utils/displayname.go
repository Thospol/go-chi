package utils

import (
	"strings"
	"unicode/utf8"

	"saaa-api/internal/core/config"
)

const (
	// MinimumDisplayNameLength minimum display name length
	MinimumDisplayNameLength = 6
	// MaximumDisplayNameLength maximum display name length
	MaximumDisplayNameLength = 50
)

// DisplayNameValidator display name validator interface
type DisplayNameValidator interface {
	ValidDisplayname(dn string) error
}

// New new display name validator
func New() DisplayNameValidator {
	return &displayNameValidator{}
}

type displayNameValidator struct{}

// ValidDisplayname validate display name(dn)
func (dv *displayNameValidator) ValidDisplayname(displayName string) error {
	displayName = strings.TrimSpace(displayName)
	if displayName == "" {
		return config.RR.InvalidName
	}

	if lengthOfName := utf8.RuneCountInString(displayName); lengthOfName > MaximumDisplayNameLength {
		return config.RR.OverMaxSizeOfName
	}

	return nil
}
