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
			if v.Pre != tc.pre {
				t.Errorf("pre: got %q, want %q", v.Pre, tc.pre)
			}
			if v.Post != tc.post {
				t.Errorf("post: got %q, want %q", v.Post, tc.post)
			}
			if v.Dev != tc.dev {
				t.Errorf("dev: got %q, want %q", v.Dev, tc.dev)
			}
			if v.Local != tc.local {
				t.Errorf("local: got %q, want %q", v.Local, tc.local)
			}
		})
	}
}
