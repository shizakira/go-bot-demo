package helpers

import "testing"

func TestGeneratePlaceholders(t *testing.T) {
	tests := []struct {
		name     string
		count    int
		expected string
	}{
		{"zero", 0, ""},
		{"negative", -1, ""},
		{"single", 1, "$1"},
		{"multiple", 3, "$1, $2, $3"},
	}

	for _, tc := range tests {
		if got := GeneratePlaceholders(tc.count); got != tc.expected {
			t.Fatalf("%s: expected %q, got %q", tc.name, tc.expected, got)
		}
	}
}
