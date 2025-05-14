package pyver

import (
	"strconv"
	"strings"
	"testing"
)

func releaseToString(release []int) string {
	parts := make([]string, len(release))
	for i, n := range release {
		parts[i] = strconv.Itoa(n)
	}
	return strings.Join(parts, ".")
}

func preToString(kind string, num int) string {
	if kind == "" {
		return ""
	}
	return kind + strconv.Itoa(num)
}

func postToString(num int) string {
	if num == 0 {
		return ""
	}
	return "post" + strconv.Itoa(num)
}

func devToString(num int) string {
	if num == 0 {
		return ""
	}
	return "dev" + strconv.Itoa(num)
}

func TestParsePEP440Fields(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		epoch   int
		release string
		pre     string
		post    string
		dev     string
		local   string
	}{
		{"epoch", "1!1.0.0", 1, "1.0.0", "", "", "", ""},
		{"release", "2.0.1", 0, "2.0.1", "", "", "", ""},
		{"pre", "1.0.0a2", 0, "1.0.0", "a2", "", "", ""},
		{"post", "1.0.0.post3", 0, "1.0.0", "", "post3", "", ""},
		{"dev", "1.0.0.dev4", 0, "1.0.0", "", "", "dev4", ""},
		{"local", "1.0.0+abc.5", 0, "1.0.0", "", "", "", "abc.5"},
		{"all fields", "2!3.4.5a1.post2.dev3+meta", 2, "3.4.5", "a1", "post2", "dev3", "meta"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v, err := Parse(tc.input)
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}
			if v.Epoch != tc.epoch {
				t.Errorf("epoch: got %d, want %d", v.Epoch, tc.epoch)
			}
			releaseStr := releaseToString(v.Release)
			if releaseStr != tc.release {
				t.Errorf("release: got %q, want %q", releaseStr, tc.release)
			}
			if preToString(v.PreKind, v.PreNum) != tc.pre {
				t.Errorf("pre: got %q, want %q", preToString(v.PreKind, v.PreNum), tc.pre)
			}
			if postToString(v.PostNum) != tc.post {
				t.Errorf("post: got %q, want %q", postToString(v.PostNum), tc.post)
			}
			if devToString(v.DevNum) != tc.dev {
				t.Errorf("dev: got %q, want %q", devToString(v.DevNum), tc.dev)
			}
			if strings.Join(v.Local, ".") != tc.local {
				t.Errorf("local: got %q, want %q", strings.Join(v.Local, "."), tc.local)
			}
		})
	}
}

func TestParsePEP440AllValidCases(t *testing.T) {
	cases := []string{
		// Simple releases
		"0.0.1", "1.0.0", "7.1.0", "2.2.3", "10.20.30",
		// Pre-releases
		"1.2.0rc1", "1.0.0a1", "1.0.0b2", "1.0.0rc3", "1.0.0a0", "1.0.0b0", "1.0.0rc0",
		// Post-releases
		"3.0.0.post1", "1.0.0.post2", "1.0.0-1", "1.0.0post1", "1.0.0rev1", "1.0.0r1", "1.0.0.post0", "1.0.0.post",
		// Dev releases
		"1.0.0.dev1", "1.0.0.dev0", "1.0.0dev2",
		// Epochs
		"1!1.0.0", "2!3.4.5a1.post2.dev3+meta",
		// Local versions
		"1.0.0+abc", "1.0.0+abc.5", "1.0.0+abc-def", "1.0.0+abc_def", "1.0.0+abc.def",
		// Normalization and whitespace
		"1.0.0-rc1", "1.0.0_rc1", "v1.0.0", " 1.0.0 ",
		// Leading zeros
		"01.2.3", "1.02.3", "1.2.03",
		// Complex combos
		"1!2.3.4a5.post6.dev7+abc.def",
	}
	for _, input := range cases {
		t.Run(input, func(t *testing.T) {
			v, err := Parse(input)
			if err != nil {
				t.Errorf("Parse(%q) failed: %v", input, err)
			}
			if v.Normalized == "" {
				t.Errorf("Parse(%q) did not produce a normalized version", input)
			}
		})
	}
}
