package auth

// Package auth provides helpers for working with PromptSentinel API keys.

import (
	"crypto/subtle"
	"errors"
	"fmt"
	"strings"
)

// MinimumAPIKeyLength defines the shortest allowed API key. The limit keeps
// space for entropy while remaining easy to document.
const MinimumAPIKeyLength = 16

// APIKey represents a sanitized API key string. It can safely be shared with
// functions that need to compare keys without exposing the raw, user-provided
// value.
type APIKey struct {
	value string
}

// NewAPIKey validates and normalizes an API key string. Leading and trailing
// whitespace is ignored and the resulting key must satisfy MinimumAPIKeyLength.
//
// A normalized APIKey is returned on success. When the key is invalid, the
// returned error explains which constraint failed.
func NewAPIKey(raw string) (APIKey, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return APIKey{}, errors.New("api key cannot be empty")
	}
	if len(trimmed) < MinimumAPIKeyLength {
		return APIKey{}, fmt.Errorf("api key must be at least %d characters", MinimumAPIKeyLength)
	}

	return APIKey{value: trimmed}, nil
}

// Value exposes the normalized API key. Code interacting with external systems
// (such as the database) can use this accessor instead of storing the key
// manually.
func (k APIKey) Value() string {
	return k.value
}

// Prefix returns the first n characters of the key. It is safe to log or show
// the prefix to users because it avoids leaking the full secret. When the
// requested length exceeds the key length the entire key is returned.
func (k APIKey) Prefix(length int) string {
	if length <= 0 {
		return ""
	}
	if length > len(k.value) {
		length = len(k.value)
	}
	return k.value[:length]
}

// Matches compares the current key to a candidate using constant-time
// comparison. Constant-time checks ensure that different keys take the same time
// to compare, preventing timing attacks.
func (k APIKey) Matches(candidate string) bool {
	trimmed := strings.TrimSpace(candidate)
	if len(k.value) == 0 || len(trimmed) != len(k.value) {
		return false
	}

	return subtle.ConstantTimeCompare([]byte(k.value), []byte(trimmed)) == 1
}
