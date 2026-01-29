package ui

import "testing"

func TestWordWrapShort(t *testing.T) {
	result := wordWrap("short", 39)
	if len(result) != 1 || result[0] != "short" {
		t.Errorf("expected [short], got %v", result)
	}
}

func TestWordWrapLong(t *testing.T) {
	input := "A SQL query walks into a bar, walks up to two tables and asks, Can I join you?"
	result := wordWrap(input, 39)
	for _, line := range result {
		if len(line) > 39 {
			t.Errorf("line too long: %d chars", len(line))
		}
	}
}

func TestWordWrapPreservesNewlines(t *testing.T) {
	result := wordWrap("line1\nline2", 39)
	if len(result) != 2 {
		t.Errorf("expected 2 lines, got %d", len(result))
	}
}

func TestWordWrapEmpty(t *testing.T) {
	result := wordWrap("", 39)
	if len(result) != 1 || result[0] != "" {
		t.Errorf("expected [\"\"], got %v", result)
	}
}
