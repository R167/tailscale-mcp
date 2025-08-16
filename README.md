# Tailscale MCP Server

A Model Context Protocol (MCP) server that provides tools for interacting with Tailscale networks through the Tailscale API.

## Design

### Overview

This MCP server enables AI assistants to interact with Tailscale networks by providing a set of tools that wrap the Tailscale HTTP API. The server is built using:

- **Tailscale Go Client v2**: Uses `tailscale.com/client/tailscale/v2` for comprehensive API interactions
- **MCP Go SDK**: Uses `github.com/modelcontextprotocol/go-sdk` for MCP protocol implementation
- **Streamable HTTP Transport**: Uses HTTP with streaming for robust communication with MCP clients
- **OAuth2 Support**: Supports both API key and OAuth client credentials authentication

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
- **Description**: List all devices in the Tailscale network with full details
- **Input**: No parameters required
- **Output**: JSON array of comprehensive device information including names, IPs, status, routes, and metadata

#### `get_device_details`
- **Description**: Get detailed information about a specific device
- **Input**: `deviceID` (string) - The device ID to get details for
- **Output**: JSON object with full device information including routing, tags, and connectivity status

#### `get_device_routes`
- **Description**: Get subnet routes for a specific device
- **Input**: `deviceID` (string) - The device ID to get routes for
- **Output**: JSON object with advertised and enabled subnet routes for the device

#### `get_acl`
- **Description**: Get the current Access Control List (ACL) policy file for the tailnet
- **Input**: No parameters required  
- **Output**: JSON representation of the current ACL configuration

#### `list_keys`
- **Description**: List all API keys for the tailnet (both user and tailnet level)
- **Input**: No parameters required
- **Output**: JSON array of API key information including capabilities and expiration

### Error Handling

The server implements robust error handling:

- **API Errors**: Tailscale API errors are captured and returned as tool errors with descriptive messages
- **JSON Marshaling Errors**: Any issues serializing responses are handled gracefully
- **Authentication Errors**: Missing or invalid credentials result in clear error messages

## Configuration

The server requires these environment variables:

**Required:**
- `TAILSCALE_TAILNET`: Your tailnet identifier (e.g., `example.com` or `user@domain.com`)

**Authentication (choose one):**
- `TAILSCALE_API_KEY`: Your Tailscale API key (obtain from Tailscale Admin Console)
- OR `TAILSCALE_CLIENT_ID` + `TAILSCALE_CLIENT_SECRET`: OAuth client credentials

**Optional:**
- `PORT`: HTTP server port (defaults to 8080)

## Usage

### Building

```bash
go build -o tailscale-mcp
```

### Running

**Using API Key:**
```bash
export TAILSCALE_API_KEY="your-api-key-here"
export TAILSCALE_TAILNET="your-tailnet-here"
export PORT="8080"  # optional, defaults to 8080
./tailscale-mcp
```

**Using OAuth Client Credentials:**
```bash
export TAILSCALE_CLIENT_ID="your-client-id"
export TAILSCALE_CLIENT_SECRET="your-client-secret"
export TAILSCALE_TAILNET="your-tailnet-here"
export PORT="8080"  # optional, defaults to 8080
./tailscale-mcp
```

The server will start listening on the specified port (default 8080) and provide logs indicating when it's ready.

### Integration with MCP Clients

The server uses streamable HTTP transport as per the MCP specification. Example configuration for Claude Desktop:

**Using API Key:**
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

**Using OAuth Client Credentials:**
```json
{
  "mcpServers": {
    "tailscale": {
      "url": "http://localhost:8080",
      "env": {
        "TAILSCALE_CLIENT_ID": "your-client-id",
        "TAILSCALE_CLIENT_SECRET": "your-client-secret",
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
- Valid Tailscale API key OR OAuth client credentials

### Project Structure

The project is organized into modular packages:

- `main.go`: Entry point
- `config/`: Configuration and client initialization
- `server/`: HTTP server setup and lifecycle management
- `tools/`: MCP tool implementations organized by functionality
  - `devices.go`: Device management tools
  - `acl.go`: Access control list tools
  - `keys.go`: API key management tools

### Testing

To test the server functionality:

1. **Local Testing**: Run the server with proper environment variables and test the HTTP endpoints
2. **MCP Client Integration**: Integrate with an MCP client for end-to-end testing
3. **Health Check**: The server will log startup status and any errors

```bash
# Test server startup (API Key)
export TAILSCALE_API_KEY="your-key"
export TAILSCALE_TAILNET="your-tailnet"
./tailscale-mcp

# OR test with OAuth credentials
export TAILSCALE_CLIENT_ID="your-client-id"
export TAILSCALE_CLIENT_SECRET="your-client-secret"
export TAILSCALE_TAILNET="your-tailnet"
./tailscale-mcp

# In another terminal, you can test basic connectivity
curl -X POST http://localhost:8080 -H "Content-Type: application/json"
```

### OAuth Client Setup

To use OAuth client credentials instead of API keys:

1. Go to the Tailscale Admin Console
2. Navigate to Settings > OAuth Clients
3. Create a new OAuth client with appropriate scopes:
   - `devices` - for device listing and management
   - `routes` - for subnet route information
   - `dns` - for DNS configuration access
4. Use the generated Client ID and Client Secret with the server

### Future Enhancements

Potential additions to the server:

- Device management operations (enable/disable, rename, set routes)
- ACL modification capabilities  
- DNS configuration management
- User and group management
- Audit log access
- Real-time status monitoring
- Webhook support for notifications

## Security Considerations

- All credentials (API keys, OAuth secrets) are handled securely through environment variables
- No credentials are logged or exposed in error messages
- OAuth client credentials provide more granular access control than API keys
- All API operations respect Tailscale's built-in permissions and access controls
- The server provides both read-only and read-write operations based on the configured scopes
- OAuth tokens are automatically managed and refreshed by the client library

## License

This project follows standard Go module practices and is designed for integration with Tailscale's official Go client library.
