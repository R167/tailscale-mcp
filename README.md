# Tailscale MCP Server

A Model Context Protocol (MCP) server that provides tools for interacting with Tailscale networks through the Tailscale API.

## Design

### Overview

This MCP server enables AI assistants to interact with Tailscale networks by providing a set of tools that wrap the Tailscale HTTP API. The server is built using:

- **Tailscale Go Client**: Uses `github.com/tailscale/tailscale-client-go` for API interactions
- **MCP Go SDK**: Uses `github.com/modelcontextprotocol/go-sdk` for MCP protocol implementation
- **Streamable HTTP Transport**: Uses HTTP with streaming for robust communication with MCP clients

### Architecture

```
AI Assistant <-> MCP Client <-> HTTP <-> Tailscale MCP Server <-> Tailscale API
```

The server follows these key design principles:

1. **Idiomatic Go**: Uses standard Go patterns and error handling
2. **Security-First**: Requires explicit environment variable configuration for credentials
3. **Minimal Dependencies**: Only essential packages for Tailscale API and MCP protocol
4. **Tool-Based Architecture**: Each Tailscale operation is exposed as a separate MCP tool

### Available Tools

The server provides the following tools:

#### `list_devices`
- **Description**: List all devices in the Tailscale network
- **Input**: No parameters required
- **Output**: JSON array of device information including names, IPs, and status

#### `get_acl`
- **Description**: Get the current Access Control List (ACL) for the tailnet
- **Input**: No parameters required  
- **Output**: JSON representation of the current ACL configuration

#### `list_keys`
- **Description**: List all API keys for the tailnet
- **Input**: No parameters required
- **Output**: JSON array of API key information

### Error Handling

The server implements robust error handling:

- **API Errors**: Tailscale API errors are captured and returned as tool errors with descriptive messages
- **JSON Marshaling Errors**: Any issues serializing responses are handled gracefully
- **Authentication Errors**: Missing or invalid credentials result in clear error messages

## Configuration

The server requires these environment variables:

- `TAILSCALE_API_KEY`: Your Tailscale API key (obtain from Tailscale Admin Console)
- `TAILSCALE_TAILNET`: Your tailnet identifier (e.g., `example.com` or `user@domain.com`)
- `PORT`: HTTP server port (optional, defaults to 8080)

## Usage

### Building

```bash
go build -o tailscale-mcp
```

### Running

```bash
export TAILSCALE_API_KEY="your-api-key-here"
export TAILSCALE_TAILNET="your-tailnet-here"
export PORT="8080"  # optional, defaults to 8080
./tailscale-mcp
```

The server will start listening on the specified port (default 8080) and provide logs indicating when it's ready.

### Integration with MCP Clients

The server uses streamable HTTP transport as per the MCP specification. Example configuration for Claude Desktop:

```json
{
  "mcpServers": {
    "tailscale": {
      "url": "http://localhost:8080",
      "env": {
        "TAILSCALE_API_KEY": "your-api-key",
        "TAILSCALE_TAILNET": "your-tailnet"
      }
    }
  }
}
```

For production deployments, configure with HTTPS and proper authentication.

## Development

### Dependencies

- Go 1.25+
- Tailscale account with API access
- Valid Tailscale API key

### Testing

To test the server functionality:

1. **Local Testing**: Run the server with proper environment variables and test the HTTP endpoints
2. **MCP Client Integration**: Integrate with an MCP client for end-to-end testing
3. **Health Check**: The server will log startup status and any errors

```bash
# Test server startup
export TAILSCALE_API_KEY="your-key"
export TAILSCALE_TAILNET="your-tailnet"
./tailscale-mcp

# In another terminal, you can test basic connectivity
curl -X POST http://localhost:8080 -H "Content-Type: application/json"
```

### Future Enhancements

Potential additions to the server:

- Device management operations (enable/disable, rename)
- ACL modification capabilities  
- Route management
- User and group management
- Audit log access
- Real-time status monitoring

## Security Considerations

- API keys are handled securely through environment variables
- No credentials are logged or exposed in error messages
- All API operations respect Tailscale's built-in permissions and access controls
- The server only provides read-only operations by default for safety

## License

This project follows standard Go module practices and is designed for integration with Tailscale's official Go client library.
