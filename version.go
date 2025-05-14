// Package pyver provides PEP 440-compliant version parsing and comparison in Go.
package pyver

// Version represents a parsed PEP 440 version.
type Version struct {
	Epoch      int      // e.g. 1!1.2.3 -> 1
	Release    []int    // e.g. 1.2.3 -> [1,2,3]
	PreKind    string   // "a", "b", "rc", or "" if not present
	PreNum     int      // e.g. "a1" -> 1, 0 if not present
	PostNum    int      // e.g. "post2" -> 2, 0 if not present
	DevNum     int      // e.g. "dev3" -> 3, 0 if not present
	Local      []string // e.g. "abc.1" -> ["abc", "1"]
	Original   string   // original version string
	Normalized string   // canonical/normalized version string
}
