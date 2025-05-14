package pyver

import "testing"

func TestInvalidVersions(t *testing.T) {
	cases := []string{
		"",                  // empty
		"1..0.0",            // double dot
		"1.0.",              // trailing dot
		".1.0.0",            // leading dot
		"1.0.0..1",          // double dot in middle
		"1!1!1.0.0",         // multiple epochs
		"1.0.0++abc",        // double plus
		"1.0.0+abc+def",     // multiple local segments
		"1.0.0@abc",         // invalid character
		"1.0.0#meta",        // invalid character
		"1.0.0..dev1",       // double dot before dev
		"1.0.0.dev1.dev2",   // multiple dev segments
		"1.0.0.post1.post2", // multiple post segments
		"1.0.0a1a2",         // multiple pre segments
		"1.0.0 dev1",        // space in version
		"-1.0.0",            // negative release segment
		"1.0.-1",            // negative release segment
		"1.0.0+",            // local with no identifier
		"1.0.0+abc..def",    // double dot in local
	}
	for _, s := range cases {
		t.Run(s, func(t *testing.T) {
			_, err := Parse(s)
			if err == nil {
				t.Errorf("expected error for version %q, got nil", s)
			}
		})
	}
}
