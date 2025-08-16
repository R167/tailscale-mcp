#!/bin/bash
# Auto-format Go files after edits

# Change to project directory
cd "$CLAUDE_PROJECT_DIR" || exit 1

# Run go fmt on all Go files
go fmt ./...

# Run goimports if available
if command -v goimports >/dev/null 2>&1; then
    goimports -w .
fi