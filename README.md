# pyver

## Running Tests

This project uses Go for the main implementation and Python (via `pyver_backend.py`) for PEP 440 parsing and comparison. All tests are run using Go, but require a Python environment with the `packaging` library installed.

### Setup (using uv)

1. Create a Python virtual environment using [uv](https://github.com/astral-sh/uv):

   ```sh
   make venv
   ```

2. Install Python dependencies:

   ```sh
   make install
   ```

### Running All Tests

To run all Go tests (including those that invoke the Python backend):

```sh
make test
```

Or directly with Go:

```sh
go test -v
```

You can also run a specific test file, e.g.:

```sh
go test -v compare_test.go
```

### Checking Backend Presence

Ensure that `pyver_backend.py` is present in the project root and is executable. You can check this with:

```sh
test -x pyver_backend.py && echo "Backend is executable"
```

Or run a direct backend parse:

```sh
.venv/bin/python3 pyver_backend.py parse 1.2.3
```

## Troubleshooting Backend Errors

- If you see errors like `pyver backend error: exit status 1`, ensure that:
  - The Python virtual environment is created and activated.
  - The `packaging` library is installed in the venv.
  - The `pyver_backend.py` script is present and executable.
  - The `GO_PYTHON` environment variable (if set) points to the correct Python interpreter.
- You can test the backend directly as shown above.
- If the backend is missing or broken, all tests will fail. Fix the backend before proceeding.

## Acceptance Criteria for Go-native Implementation

The current test suite (all `*_test.go` files) is the acceptance criteria for a Go-native implementation. Any new implementation must pass all tests to be considered fully PEP 440 compatible.

## Implementation Modes

By default, `pyver` uses a **Go-native implementation** of PEP 440 parsing and comparison, which is fast, dependency-free, and fully standards-compliant. This is recommended for all users.

For advanced users, maintainers, and debugging, `pyver` also includes a **Python backend** (wrapping Python's `packaging.version.Version`) as a reference implementation. This can be toggled via the `UseGoNative` variable in code, or via an environment variable if supported.

- **Go-native mode:** Default, fast, and recommended for all production use.
- **Python backend mode:** For regression testing, debugging, and verifying compliance with the Python standard. Useful for comparing results or investigating subtle edge cases.

CI may run both implementations to ensure ongoing parity and standards compliance.
