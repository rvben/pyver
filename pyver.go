package pyver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"unicode"
)

// BackendPath is the absolute path to the backend script.
var BackendPath string

// UseGoNative toggles between the Python backend and Go-native implementation.
// Set to true to use Go-native parsing/comparison (in development).
var UseGoNative = true

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

// Parse parses a version string into a Version struct.
func Parse(s string) (Version, error) {
	if UseGoNative {
		return parseGoNative(s)
	}
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
		v.Normalized = norm
	}
	if pre, ok := resp["pre"].(string); ok && pre != "" {
		// e.g. "a1", "b2", "rc3"
		if len(pre) >= 2 {
			v.PreKind = pre[:len(pre)-1]
			if n, err := strconv.Atoi(pre[len(pre)-1:]); err == nil {
				v.PreNum = n
			}
		} else if len(pre) == 1 {
			v.PreKind = pre
			v.PreNum = 0
		}
	}
	if post, ok := resp["post"].(string); ok && post != "" {
		// e.g. "post2"
		if strings.HasPrefix(post, "post") {
			num := post[4:]
			if n, err := strconv.Atoi(num); err == nil {
				v.PostNum = n
			}
		}
	}
	if dev, ok := resp["dev"].(string); ok && dev != "" {
		// e.g. "dev3"
		if strings.HasPrefix(dev, "dev") {
			num := dev[3:]
			if n, err := strconv.Atoi(num); err == nil {
				v.DevNum = n
			}
		}
	}
	if local, ok := resp["local"].(string); ok && local != "" {
		v.Local = strings.Split(local, ".")
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

// Compare returns -1 if v1 < v2, 0 if v1 == v2, 1 if v1 > v2.
func Compare(v1, v2 Version) int {
	if UseGoNative {
		// TODO: Implement Go-native comparison logic here
		panic("Go-native comparison not yet implemented")
	}
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
	if v.Normalized != "" {
		return v.Normalized
	}
	return v.Original
}

// --- Go-native PEP 440 parser and normalizer ---

// parseGoNative parses and normalizes a PEP 440 version string in pure Go.
func parseGoNative(s string) (Version, error) {
	orig := s
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)
	if strings.HasPrefix(s, "v") {
		s = s[1:]
	}

	// Regex for PEP 440 (per Appendix B, with normalization flexibility)
	var pep440Pattern = regexp.MustCompile(`^((?P<epoch>[0-9]+)!)?(?P<release>[0-9]+(?:\.[0-9]+)*)(?P<pre>([-_\.]?(a|b|rc|alpha|beta|c|pre|preview)[-_\.]?[0-9]*)?)?(?P<post>(-(?P<post_n1>[0-9]+))|(([-_\.]?(post|rev|r)[-_\.]?[0-9]*)))?(?P<dev>([-_\.]?dev[-_\.]?[0-9]*))?(\+(?P<local>[a-z0-9]+(?:[-_\.][a-z0-9]+)*))?$`)

	m := pep440Pattern.FindStringSubmatch(s)
	if m == nil {
		return Version{Original: orig}, fmt.Errorf("invalid version: %q", orig)
	}

	v := Version{Original: orig}
	group := func(name string) string {
		for i, n := range pep440Pattern.SubexpNames() {
			if n == name {
				return m[i]
			}
		}
		return ""
	}

	// Epoch
	if e := group("epoch"); e != "" {
		v.Epoch, _ = strconv.Atoi(e)
	}

	// Release
	rel := group("release")
	for _, part := range strings.Split(rel, ".") {
		if part == "" {
			return v, fmt.Errorf("invalid release segment: %q", rel)
		}
		n, err := strconv.Atoi(part)
		if err != nil || n < 0 {
			return v, fmt.Errorf("invalid release segment: %q", rel)
		}
		v.Release = append(v.Release, n)
	}

	// Pre-release
	pre := group("pre")
	if pre != "" {
		// Normalize spelling and separator
		pre = strings.ReplaceAll(pre, "_", "")
		pre = strings.ReplaceAll(pre, "-", "")
		pre = strings.ReplaceAll(pre, ".", "")
		var kind string
		var num int
		for _, k := range []struct{ alt, norm string }{
			{"alpha", "a"}, {"a", "a"},
			{"beta", "b"}, {"b", "b"},
			{"rc", "rc"}, {"c", "rc"}, {"pre", "rc"}, {"preview", "rc"},
		} {
			if strings.HasPrefix(pre, k.alt) {
				kind = k.norm
				pre = pre[len(k.alt):]
				break
			}
		}
		if kind != "" {
			if pre == "" {
				num = 0
			} else {
				num, _ = strconv.Atoi(pre)
			}
			v.PreKind = kind
			v.PreNum = num
		}
	}

	// Post-release
	post := group("post")
	if post != "" {
		// Normalize spelling and separator
		post = strings.ReplaceAll(post, "_", "")
		post = strings.ReplaceAll(post, "-", "")
		post = strings.ReplaceAll(post, ".", "")
		if strings.HasPrefix(post, "post") {
			post = post[4:]
		} else if strings.HasPrefix(post, "rev") {
			post = post[3:]
		} else if strings.HasPrefix(post, "r") {
			post = post[1:]
		}
		if post == "" {
			v.PostNum = 0
		} else {
			v.PostNum, _ = strconv.Atoi(post)
		}
	}
	if n := group("post_n1"); n != "" {
		v.PostNum, _ = strconv.Atoi(n)
	}

	// Dev-release
	dev := group("dev")
	if dev != "" {
		dev = strings.ReplaceAll(dev, "_", "")
		dev = strings.ReplaceAll(dev, "-", "")
		dev = strings.ReplaceAll(dev, ".", "")
		if strings.HasPrefix(dev, "dev") {
			dev = dev[3:]
		}
		if dev == "" {
			v.DevNum = 0
		} else {
			v.DevNum, _ = strconv.Atoi(dev)
		}
	}

	// Local version
	local := group("local")
	if local != "" {
		// Normalize separators to '.'
		for _, sep := range []string{"-", "_"} {
			local = strings.ReplaceAll(local, sep, ".")
		}
		parts := strings.Split(local, ".")
		for _, part := range parts {
			if part == "" {
				return v, fmt.Errorf("invalid local segment: %q", local)
			}
			// Must be alphanumeric
			for _, r := range part {
				if !unicode.IsDigit(r) && !unicode.IsLetter(r) {
					return v, fmt.Errorf("invalid local segment: %q", local)
				}
			}
			v.Local = append(v.Local, part)
		}
		// Must start and end with alphanumeric
		if len(v.Local) == 0 || !isAlnum(v.Local[0]) || !isAlnum(v.Local[len(v.Local)-1]) {
			return v, fmt.Errorf("invalid local segment: %q", local)
		}
	}

	// Normalized string
	v.Normalized = versionToString(v)
	return v, nil
}

func isAlnum(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return true
		}
	}
	return false
}

// versionToString returns the canonical normalized version string.
func versionToString(v Version) string {
	var b strings.Builder
	if v.Epoch > 0 {
		b.WriteString(strconv.Itoa(v.Epoch))
		b.WriteString("!")
	}
	for i, n := range v.Release {
		if i > 0 {
			b.WriteString(".")
		}
		b.WriteString(strconv.Itoa(n))
	}
	if v.PreKind != "" {
		b.WriteString(v.PreKind)
		b.WriteString(strconv.Itoa(v.PreNum))
	}
	if v.PostNum > 0 || (v.PostNum == 0 && (v.PreKind != "" || v.DevNum != 0)) {
		b.WriteString(".post")
		b.WriteString(strconv.Itoa(v.PostNum))
	}
	if v.DevNum > 0 || (v.DevNum == 0 && (v.PreKind != "" || v.PostNum != 0)) {
		b.WriteString(".dev")
		b.WriteString(strconv.Itoa(v.DevNum))
	}
	if len(v.Local) > 0 {
		b.WriteString("+")
		b.WriteString(strings.Join(v.Local, "."))
	}
	return b.String()
}
