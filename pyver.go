package pyver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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

// BackendPath is the absolute path to the backend script.
var BackendPath string

func init() {
	// Use runtime.Caller to get the directory of this source file
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("could not determine pyver.go location via runtime.Caller")
	}
	backend := filepath.Join(filepath.Dir(filename), "pyver_backend.py")
	if _, err := os.Stat(backend); err == nil {
		BackendPath = backend
		return
	}
	// fallback: look one directory up (for monorepo or test layouts)
	backend = filepath.Join(filepath.Dir(filepath.Dir(filename)), "pyver", "pyver_backend.py")
	if _, err := os.Stat(backend); err == nil {
		BackendPath = backend
		return
	}
	panic("pyver_backend.py not found relative to pyver.go")
}

// pyver.go: Go interface to Python PEP 440 version parsing/comparison.
//
// The Python interpreter used for the backend is determined by the GO_PYTHON environment variable.
// If GO_PYTHON is not set, 'python3' is used.
// GO_PYTHON may be a multi-word command (e.g., 'uv run --with packaging python3').
//
// Example:
//   export GO_PYTHON='uv run --with packaging python3'
//
// The backend must have the 'packaging' library available.

// getPythonArgs returns the Python interpreter and its arguments as a slice.
// If GO_PYTHON is set, uses that (split by spaces).
// Otherwise, if 'uv' is available, uses 'uv run --with packaging python3'.
// Otherwise, falls back to 'python3'.
func getPythonArgs() []string {
	if py := os.Getenv("GO_PYTHON"); py != "" {
		return strings.Fields(py)
	}
	if uvPath, err := exec.LookPath("uv"); err == nil {
		return []string{uvPath, "run", "--with", "packaging", "python3"}
	}
	return []string{"python3"}
}

// Parse parses a version string into a Version struct using the backend.
func Parse(s string) (Version, error) {
	v := Version{Original: s}
	args := append(getPythonArgs(), BackendPath, "parse", s)
	cmd := exec.Command(args[0], args[1:]...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[pyver debug] Parse failed for input: %q\n  Command: %v\n  Stderr: %s\n  Error: %v\n", s, args, stderr.String(), err)
		return v, fmt.Errorf("pyver backend error: %v\nCommand: %v\nStderr: %s", err, args, stderr.String())
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
	args := append(getPythonArgs(), BackendPath, "compare", v1.Original, v2.Original)
	cmd := exec.Command(args[0], args[1:]...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		panic(fmt.Errorf("pyver backend error: %v\nCommand: %v\nStderr: %s", err, args, stderr.String()))
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
