package pyver

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Version represents a parsed PEP 440 version.
type Version struct {
	Epoch    int
	Release  []int  // e.g. 1.2.3 -> [1,2,3]
	Pre      string // canonical string, e.g. "a1" or ""
	Post     string // canonical string, e.g. "post2" or ""
	Dev      string // canonical string, e.g. "dev3" or ""
	Local    string // canonical string, e.g. "abc.5" or ""
	Original string // original string
	Norm     string // normalized/canonical string
}

// BackendPath is the path to the backend binary or script.
var BackendPath = "pyver_backend.py"

func getPython() string {
	if py := os.Getenv("GO_PYTHON"); py != "" {
		return py
	}
	return "python3"
}

// Parse parses a version string into a Version struct using the backend.
func Parse(s string) (Version, error) {
	v := Version{Original: s}
	cmd := exec.Command(getPython(), BackendPath, "parse", s)
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
	if norm, ok := resp["normalized"].(string); ok {
		v.Norm = norm
	}
	if pre, ok := resp["pre"].(string); ok {
		v.Pre = pre
	}
	if post, ok := resp["post"].(string); ok {
		v.Post = post
	}
	if dev, ok := resp["dev"].(string); ok {
		v.Dev = dev
	}
	if local, ok := resp["local"].(string); ok {
		v.Local = local
	}
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
	cmd := exec.Command(getPython(), BackendPath, "compare", v1.Original, v2.Original)
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

// String returns the normalized version string.
func (v Version) String() string {
	if v.Norm != "" {
		return v.Norm
	}
	return v.Original
}
