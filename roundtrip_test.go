package pyver

import "testing"

func TestRoundtripParseString(t *testing.T) {
	cases := []struct {
		input    string
		expected string // expected output after normalization, or same as input if no normalization
	}{
		{"1.2.3", "1.2.3"},
		{"1.0", "1.0"}, // normalization: packaging normalizes '1.0' to '1.0'
		{"1.0.0", "1.0.0"},
		{"1.0.0a1", "1.0.0a1"},
		{"1.0.0.post1", "1.0.0.post1"},
		{"1.0.0.dev2", "1.0.0.dev2"},
		{"1!1.0.0", "1!1.0.0"},
		{"1.0.0+abc", "1.0.0+abc"},
		{"1.0.0-rc1", "1.0.0rc1"}, // normalization
		{"1.02.3", "1.2.3"},       // leading zero normalization
		{"1.0.0.post1.dev2", "1.0.0.post1.dev2"},
		{"2!3.4.5a1.post2.dev3+meta", "2!3.4.5a1.post2.dev3+meta"},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			v, err := Parse(tc.input)
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}
			s := v.String()
			if s != tc.expected {
				t.Errorf("roundtrip: got %q, want %q", s, tc.expected)
			}
		})
	}
}
