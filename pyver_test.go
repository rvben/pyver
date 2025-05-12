package pyver

import (
	"testing"
)

func TestCompare(t *testing.T) {
	tests := []struct {
		name   string
		v1     string
		v2     string
		expect int // -1 if v1 < v2, 0 if v1 == v2, 1 if v1 > v2
	}{
		{"numeric lt", "1.2.3", "1.2.4", -1},
		{"numeric eq", "1.2.3", "1.2.3", 0},
		{"numeric gt", "1.2.4", "1.2.3", 1},
		// Pre-releases
		{"pre-release a<b", "1.0.0a1", "1.0.0b1", -1},     // TODO: will fail until pre-release supported
		{"pre-release b<rc", "1.0.0b1", "1.0.0rc1", -1},   // TODO
		{"pre-release rc<final", "1.0.0rc1", "1.0.0", -1}, // TODO
		// Post-releases
		{"post-release", "1.0.0", "1.0.0.post1", -1}, // TODO
		// Dev releases
		{"dev<a", "1.0.0.dev1", "1.0.0a1", -1}, // TODO
		// Epochs
		{"epoch", "1!1.0.0", "1.0.0", 1}, // TODO
		// Local versions (should be equal for ordering)
		{"local eq", "1.0.0+abc", "1.0.0+xyz", 0}, // TODO
		// Wildcards (treat as higher than any patch in that minor)
		{"wildcard", "1.2.*", "1.1.9", 1}, // TODO
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v1, err1 := Parse(tc.v1)
			v2, err2 := Parse(tc.v2)
			if err1 != nil || err2 != nil {
				t.Fatalf("parse error: %v, %v", err1, err2)
			}
			cmp := Compare(v1, v2)
			if cmp != tc.expect {
				t.Errorf("Compare(%q, %q) = %d, want %d", tc.v1, tc.v2, cmp, tc.expect)
			}
		})
	}
}
