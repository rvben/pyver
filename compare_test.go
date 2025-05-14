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
		{"pre-release a<b", "1.0.0a1", "1.0.0b1", -1},
		{"pre-release b<rc", "1.0.0b1", "1.0.0rc1", -1},
		{"pre-release rc<final", "1.0.0rc1", "1.0.0", -1},
		// Post-releases
		{"post-release", "1.0.0", "1.0.0.post1", -1},
		// Dev releases
		{"dev<a", "1.0.0.dev1", "1.0.0a1", -1},
		// Epochs
		{"epoch", "1!1.0.0", "1.0.0", 1},
		// Local versions (should be equal for ordering)
		{"local eq", "1.0.0+abc", "1.0.0+xyz", -1},
		// Local version tie-breakers
		{"local tie-breaker numeric", "1.0.0+1", "1.0.0+2", -1},
		{"local tie-breaker lexicographic", "1.0.0+abc", "1.0.0+abd", -1},
		{"local tie-breaker mixed", "1.0.0+1.abc", "1.0.0+1.abd", -1},
		{"local vs no local", "1.0.0", "1.0.0+abc", -1},      // PEP 440: local only matters if public is equal, but local wins
		{"local numeric vs string", "1.0.0+1", "1.0.0+a", 1}, // numeric < string in local, so 1 > a
		// Normalization and segment equivalence
		{"normalization 1.0 vs 1.0.0", "1.0", "1.0.0", 0},
		{"normalization 1.0.0 vs 1.0.0.0", "1.0.0", "1.0.0.0", 0},
		// Real-world messy versions
		{"messy rc dash", "1.0.0-rc1", "1.0.0rc1", 0}, // Should normalize
		{"messy post dash", "1.0.0-post1", "1.0.0.post1", 0},
		{"messy dev dash", "1.0.0-dev1", "1.0.0.dev1", 0},
		// Pre-release vs dev/post
		{"pre vs dev", "1.0.0a1", "1.0.0.dev1", 1},
		{"pre vs post", "1.0.0a1", "1.0.0.post1", -1},
		// Epoch with pre/post/dev
		{"epoch with pre", "1!1.0.0a1", "1.0.0a1", 1},
		{"epoch with post", "1!1.0.0.post1", "1.0.0.post1", 1},
		// Leading zeros (should normalize)
		{"leading zeros", "1.02.3", "1.2.3", 0},
		// Complex real-world
		{"complex dev/post", "1.0.0.post1.dev2", "1.0.0.post1.dev3", -1},
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
