package pyver

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// Version represents a parsed PEP 440 version.
type Version struct {
	Epoch    int
	Release  []int  // e.g. 1.2.3 -> [1,2,3]
	Pre      any    // e.g. ["a", 1] or null
	Post     any    // e.g. 1 or null
	Dev      any    // e.g. 1 or null
	Local    any    // e.g. ["abc", 1] or null
	Original string // original string
}

// BackendPath is the path to the backend binary or script.
var BackendPath = "pyver/pyver_backend.py"

// Parse parses a version string into a Version struct using the backend.
func Parse(s string) (Version, error) {
	v := Version{Original: s}
	cmd := exec.Command("python3", BackendPath, "parse", s)
	out, err := cmd.Output()
	if err != nil {
		return v, fmt.Errorf("pyver backend error: %v", err)
	}
	var resp map[string]any
	if err := json.Unmarshal(out, &resp); err != nil {
		return v, fmt.Errorf("pyver backend JSON error: %v", err)
	}
	if epoch, ok := resp["epoch"].(float64); ok {
		v.Epoch = int(epoch)
	}
	if rel, ok := resp["release"].([]any); ok {
		for _, n := range rel {
			if f, ok := n.(float64); ok {
				v.Release = append(v.Release, int(f))
			}
		}
	}
	v.Pre = resp["pre"]
	v.Post = resp["post"]
	v.Dev = resp["dev"]
	v.Local = resp["local"]
	return v, nil
}

// MustParse parses a version string or panics.
func MustParse(s string) Version {
	v, err := Parse(s)
	if err != nil {
		panic(err)
	}
	return v
}

// Compare returns -1 if v1 < v2, 0 if v1 == v2, 1 if v1 > v2, using the backend.
func Compare(v1, v2 Version) int {
	cmd := exec.Command("python3", BackendPath, "compare", v1.Original, v2.Original)
	out, err := cmd.Output()
	if err != nil {
		panic(fmt.Errorf("pyver backend error: %v", err))
	}
	cmp, err := strconv.Atoi(strings.TrimSpace(string(out)))
	if err != nil {
		panic(fmt.Errorf("pyver backend output error: %v", err))
	}
	return cmp
}

// String returns the original version string.
func (v Version) String() string {
	return v.Original
}
