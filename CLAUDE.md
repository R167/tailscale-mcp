# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Build and Run
```bash
go build -o tailscale-mcp     # Build the binary
go run .                      # Run directly from source
```

### Testing
```bash
go test ./...                 # Run all tests
go test ./config              # Run specific package tests
go test -v ./internal         # Verbose output for a package
go test -run TestValidatePort # Run specific test
```

### Code Quality
```bash
go fmt ./...                  # Auto-format code (runs automatically via hooks)
go vet ./...                  # Static analysis
go mod tidy                   # Clean up dependencies
```

### Development Tools
```bash
go doc tailscale.com/client/tailscale/v2    # View Tailscale API docs
go list -m all                               # Show all dependencies
```

### CI/CD Commands
```bash
golangci-lint run                            # Run local linting (same as CI)
goreleaser check                             # Validate .goreleaser.yml config
goreleaser build --snapshot --clean          # Test local build
goreleaser release --snapshot --clean        # Test full release process locally
```

## Architecture Overview

This is a **Model Context Protocol (MCP) server** that wraps the Tailscale v2 API for AI assistant integration. The architecture follows a **modular, interface-based design** with clear separation of concerns.

### Core Design Patterns

**1. Interface-Based Dependency Injection**
- `internal/interfaces.go` defines abstraction interfaces for Tailscale client operations
- `internal/mocks.go` provides test implementations
- This pattern enables comprehensive unit testing and loose coupling

**2. Modular Package Structure**
```
config/     - Configuration loading, validation, and client creation
server/     - HTTP server lifecycle and middleware (logging, metrics, context)
tools/      - MCP tool implementations (devices, ACL, keys)
internal/   - Shared interfaces, mocks, and metrics collection
```

**3. Authentication Strategy**
- Supports both **API key** and **OAuth client credentials** authentication
- OAuth provides more granular scoping and better security
- Client creation abstracted in `config/config.go`

### Key Components

**Configuration (`config/`)**
- `Load()` function validates environment variables and creates clients
- Supports both authentication methods with comprehensive validation
- Uses `TailscaleClientAdapter` to wrap the real client with our interfaces

**Server (`server/`)**
- HTTP transport with **streamable MCP protocol** support
- Middleware provides request correlation IDs, structured logging, and metrics
- Graceful shutdown with proper signal handling

**Tools (`tools/`)**
- Each file implements a category of MCP tools (devices, ACL, keys)
- Common error handling patterns in `helpers.go`
- All tools use the abstracted client interfaces for testability

**Internal (`internal/`)**
- Interface definitions that abstract Tailscale client operations
- Metrics collection with thread-safe request/error tracking
- Mock implementations for comprehensive testing

### Environment Configuration

**Required:**
- `TAILSCALE_TAILNET`: Tailnet identifier (validated for format)

**Authentication (one required):**
- `TAILSCALE_API_KEY`: Traditional API key auth
- `TAILSCALE_CLIENT_ID` + `TAILSCALE_CLIENT_SECRET`: OAuth client credentials

**Optional:**
- `PORT`: HTTP server port (defaults to 8080, validated range 0-65535)

### Testing Strategy

**Test Organization:**
- Package-level tests in `*_test.go` files alongside source
- Integration tests in root-level `integration_test.go`
- All tests use `t.Setenv()` for environment isolation
- Mock interfaces enable unit testing without real API calls

**Test Patterns:**
- Configuration validation tests cover edge cases and error conditions
- Tool tests verify JSON schema validation and error handling
- Metrics tests include concurrency safety verification
- Integration tests verify end-to-end configuration loading

### Code Quality Automation

**.claude/settings.json Configuration:**
- Auto-formatting with `go fmt` and `goimports` after every edit via PostToolUse hooks
- Pre-approved Go development commands for productivity
- Hooks ensure consistent code formatting across all files

**Critical Linting Requirements:**
- **Automatic formatting**: Every Edit/Write/MultiEdit triggers `.claude/hooks/format-go.sh`
- **Final newlines**: All files must end with a newline character (enforced by goimports)
- **Import organization**: goimports automatically sorts and groups imports correctly

**CI/CD Pipeline:**
- **GitHub Actions CI**: Comprehensive linting with golangci-lint, multi-version testing, race detection
- **Security Scanning**: Weekly vulnerability checks with govulncheck, Nancy, and Trivy
- **Automated Releases**: GoReleaser builds cross-platform binaries, Docker images, and Homebrew formulas
- **Dependency Updates**: Dependabot automatically updates Go modules and GitHub Actions weekly

**Key Quality Practices:**
- Structured logging with request correlation IDs
- Comprehensive error wrapping with context
- Thread-safe metrics collection with sliding window averages
- Configuration validation at startup prevents runtime failures

### MCP Protocol Implementation

**Transport:** Uses HTTP with streaming (not stdio) for better reliability and debugging

**Tool Registration Pattern:**
```go
mcp.AddTool(server, toolDef, handlerFunc)
```

**Error Handling:** All tools follow consistent error response patterns with descriptive messages for AI assistants

**JSON Schema:** Each tool defines input schemas for validation and documentation

This architecture prioritizes testability, maintainability, and production readiness while providing a clean interface for AI assistants to interact with Tailscale networks.