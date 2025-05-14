# Makefile for pyver: setup, install, and test with uv-managed Python venv

# Name of the virtual environment directory
VENV_DIR=.venv
PYTHON=$(VENV_DIR)/bin/python3
PIP=$(VENV_DIR)/bin/pip

.PHONY: venv install test test-native clean

# Create a Python virtual environment using uv
venv:
	uv venv $(VENV_DIR)

# Install required Python dependencies in the venv
install: venv
	uv sync

# Run Go tests (ensure venv is set up and dependencies are installed)
test: install
	GO_PYTHON=$(PYTHON) go test -v -coverprofile=coverage.out

# Run Go tests using Go-native implementation (set USE_GO_NATIVE=1)
test-native: install
	GO_PYTHON=$(PYTHON) USE_GO_NATIVE=1 go test -v -coverprofile=coverage.out

# Remove the virtual environment
distclean clean:
	rm -rf $(VENV_DIR)

# Help message
help:
	@echo "make venv         # Create Python venv with uv"
	@echo "make install      # Install Python dependencies in venv"
	@echo "make test         # Run Go tests using Python backend (default)"
	@echo "make test-native  # Run Go tests using Go-native implementation (in development)"
	@echo "make clean        # Remove the venv"