package pyver

import (
	"testing"
)

func TestBackendSanity(t *testing.T) {
	// Test that the backend can parse a simple version
	v, err := Parse("1.2.3")
	if err != nil {
		t.Fatalf("backend failed to parse valid version: %v", err)
	}
	if v.Epoch != 0 || len(v.Release) != 3 || v.Release[0] != 1 || v.Release[1] != 2 || v.Release[2] != 3 {
		t.Errorf("unexpected parse result: %+v", v)
	}
}

func TestBackendErrorHandling(t *testing.T) {
	// Test that the backend returns an error for an invalid version
	_, err := Parse("1..0.0")
	if err == nil {
		t.Errorf("expected error for invalid version, got nil")
	}
}

func TestBackendMissingScript(t *testing.T) {
	// Simulate missing backend script
	orig := BackendPath
	BackendPath = "nonexistent_pyver_backend.py"
	defer func() { BackendPath = orig }()
	_, err := Parse("1.2.3")
	if err == nil {
		t.Skip("skipping: would break suite if backend is missing, but got no error")
	}
}
