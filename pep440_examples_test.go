package pyver

import "testing"

func TestPEP440CanonicalExamples(t *testing.T) {
	tests := []struct {
		name   string
		v1     string
		v2     string
		expect int // -1 if v1 < v2, 0 if v1 == v2, 1 if v1 > v2
	}{
		{"release eq", "1.0", "1.0.0", 0},
		{"release lt", "1.0", "2.0", -1},
		{"pre-release lt", "1.0a1", "1.0b1", -1},
		{"pre-release eq", "1.0rc1", "1.0rc1", 0},
		{"pre-release gt", "1.0rc2", "1.0rc1", 1},
		{"pre-release vs release", "1.0rc1", "1.0", -1},
		{"post-release gt", "1.0.post1", "1.0", 1},
		{"dev-release lt", "1.0.dev1", "1.0a1", -1},
		{"dev-release eq", "1.0.dev1", "1.0.dev1", 0},
		{"epoch gt", "1!1.0", "1.0", 1},
		{"epoch lt", "1.0", "1!1.0", -1},
		{"local gt", "1.0+abc", "1.0+aaa", 1},
		{"local eq", "1.0+abc", "1.0+abc", 0},
		{"local lt", "1.0+abc", "1.0+xyz", -1},
		{"normalize eq", "1.0.0.0", "1.0", 0},
		{"complex eq", "1!1.0.0.post1.dev2+abc", "1!1.0.0.post1.dev2+abc", 0},
		{"pre-release rc/c/preview eq", "1.0rc1", "1.0c1", 0},
		{"pre-release rc/preview eq", "1.0rc1", "1.0preview1", 0},
		{"pre-release a0 eq a", "1.0a0", "1.0a", 0},
		{"pre-release b0 eq b", "1.0b0", "1.0b", 0},
		{"pre-release rc0 eq rc", "1.0rc0", "1.0rc", 0},
		{"post-release post/rev/r eq", "1.0.post1", "1.0-1", 0},
		{"post-release post/rev/r eq2", "1.0.post1", "1.0post1", 0},
		{"post-release post/rev/r eq3", "1.0.post1", "1.0rev1", 0},
		{"post-release post/rev/r eq4", "1.0.post1", "1.0r1", 0},
		{"post-release post0 eq post", "1.0.post0", "1.0.post", 0},
		{"dev-release dev0 eq dev", "1.0.dev0", "1.0.dev", 0},
		{"dash/underscore normalization", "1.0.0-rc1", "1.0.0_rc1", 0},
		{"v prefix normalization", "v1.0", "1.0", 0},
		{"whitespace normalization", " 1.0.0 ", "1.0.0", 0},
		{"local version dash/underscore/period normalization", "1.0.0+abc-def", "1.0.0+abc_def", 0},
		{"local version period/underscore normalization", "1.0.0+abc.def", "1.0.0+abc_def", 0},
		{"case normalization rc/RC", "1.0RC1", "1.0rc1", 0},
		{"leading zeros in release", "1.01.0", "1.1.0", 0},
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
