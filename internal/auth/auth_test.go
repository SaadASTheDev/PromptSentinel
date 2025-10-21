package auth

import "testing"

func TestNewAPIKeyTrimAndValidate(t *testing.T) {
	key, err := NewAPIKey("   sk-1234567890abcdef   ")
	if err != nil {
		t.Fatalf("expected valid key, got error: %v", err)
	}
	expected := "sk-1234567890abcdef"
	if key.Value() != expected {
		t.Fatalf("expected value %q, got %q", expected, key.Value())
	}
}

func TestNewAPIKeyTooShort(t *testing.T) {
	if _, err := NewAPIKey("short"); err == nil {
		t.Fatal("expected error for key shorter than MinimumAPIKeyLength")
	}
}

func TestAPIKeyMatches(t *testing.T) {
	key, err := NewAPIKey("sk-1234567890abcdef")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !key.Matches("sk-1234567890abcdef") {
		t.Fatal("expected keys to match")
	}
	if key.Matches("sk-1234567890abcdeg") {
		t.Fatal("expected mismatched keys to return false")
	}
}

func TestAPIKeyPrefix(t *testing.T) {
	key, err := NewAPIKey("sk-1234567890abcdef")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if prefix := key.Prefix(4); prefix != "sk-1" {
		t.Fatalf("expected prefix 'sk-1', got %q", prefix)
	}
	if prefix := key.Prefix(0); prefix != "" {
		t.Fatalf("expected empty prefix for zero length, got %q", prefix)
	}
	if prefix := key.Prefix(100); prefix != key.Value() {
		t.Fatalf("expected full value when requested prefix is longer, got %q", prefix)
	}
}
