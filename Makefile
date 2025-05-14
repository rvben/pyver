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

# Release targets (git tag only, no version file, simplified)
release-patch:
	@latest=$$(git tag --list 'v*' --sort=-v:refname | head -n1 | sed 's/^v//'); \
	major=$$(echo $$latest | cut -d. -f1); \
	minor=$$(echo $$latest | cut -d. -f2); \
	patch=$$(echo $$latest | cut -d. -f3); \
	new_patch=$$((patch+1)); \
	new_tag="v$${major}.$${minor}.$${new_patch}"; \
	echo "Tagging $$new_tag"; \
	git tag -a $$new_tag -m "Release $$new_tag"; \
	git push origin $$new_tag

release-minor:
	@latest=$$(git tag --list 'v*' --sort=-v:refname | head -n1 | sed 's/^v//'); \
	major=$$(echo $$latest | cut -d. -f1); \
	minor=$$(echo $$latest | cut -d. -f2); \
	new_minor=$$((minor+1)); \
	new_tag="v$${major}.$${new_minor}.0"; \
	echo "Tagging $$new_tag"; \
	git tag -a $$new_tag -m "Release $$new_tag"; \
	git push origin $$new_tag

release-major:
	@latest=$$(git tag --list 'v*' --sort=-v:refname | head -n1 | sed 's/^v//'); \
	major=$$(echo $$latest | cut -d. -f1); \
	new_major=$$((major+1)); \
	new_tag="v$${new_major}.0.0"; \
	echo "Tagging $$new_tag"; \
	git tag -a $$new_tag -m "Release $$new_tag"; \
	git push origin $$new_tag

# Help message
help:
	@echo "make venv           # Create Python venv with uv"
	@echo "make install        # Install Python dependencies in venv"
	@echo "make test           # Run Go tests using Go-native implementation (default)"
	@echo "make test-python    # Run Go tests using Python reference implementation"
	@echo "make clean          # Remove the venv"
	@echo "make release-patch  # Tag and push next patch release (vX.Y.Z+1)"
	@echo "make release-minor  # Tag and push next minor release (vX.Y+1.0)"
	@echo "make release-major  # Tag and push next major release (vX+1.0.0)"