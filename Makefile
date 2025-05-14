# Makefile for pyver: setup, install, and test with uv-managed Python venv

# Name of the virtual environment directory
VENV_DIR=.venv
PYTHON=$(VENV_DIR)/bin/python3
PIP=$(VENV_DIR)/bin/pip

.PHONY: venv install test test-python clean release-patch release-minor release-major

# Create a Python virtual environment using uv
venv:
	uv venv $(VENV_DIR)

# Install required Python dependencies in the venv
install: venv
	uv sync

# Run Go tests using Go-native implementation (default)
test:
	GO_PYTHON=$(PYTHON) USE_GO_NATIVE=1 go test -v -coverprofile=coverage.out

# Run Go tests using Python reference implementation
test-python: install
	GO_PYTHON=$(PYTHON) USE_GO_NATIVE=0 go test -v -coverprofile=coverage.out

# Remove the virtual environment
distclean clean:
	rm -rf $(VENV_DIR)

# Unified release target (git tag only, no version file, simplified)
release-%:
	@type=$*; \
	branch=$$(git rev-parse --abbrev-ref HEAD); \
	if [ "$$branch" != "main" ]; then \
		echo "You must be on the 'main' branch to create a release (current: $$branch)"; \
		exit 1; \
	fi; \
	latest=$$(git tag --list 'v*' --sort=-v:refname | head -n1 | sed 's/^v//'); \
	major=$$(echo $$latest | cut -d. -f1); \
	minor=$$(echo $$latest | cut -d. -f2); \
	patch=$$(echo $$latest | cut -d. -f3); \
	if [ "$$type" = "patch" ]; then \
		patch=$$((patch+1)); \
	elif [ "$$type" = "minor" ]; then \
		minor=$$((minor+1)); patch=0; \
	elif [ "$$type" = "major" ]; then \
		major=$$((major+1)); minor=0; patch=0; \
	else \
		echo "Unknown release type: $$type"; exit 1; \
	fi; \
	new_tag="v$${major}.$${minor}.$${patch}"; \
	echo "Tagging $$new_tag"; \
	git tag -a $$new_tag -m "Release $$new_tag"; \
	git push origin $$new_tag

# Help message
help:
	@echo "make venv             # Create Python venv with uv"
	@echo "make install          # Install Python dependencies in venv"
	@echo "make test             # Run Go tests using Go-native implementation (default)"
	@echo "make test-python      # Run Go tests using Python reference implementation"
	@echo "make clean            # Remove the venv"
	@echo "make release-patch    # Tag and push next patch release (vX.Y.Z+1)"
	@echo "make release-minor    # Tag and push next minor release (vX.Y+1.0)"
	@echo "make release-major    # Tag and push next major release (vX+1.0.0)"