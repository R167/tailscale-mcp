package tools

import (
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// toolError creates a standardized error response for MCP tools
func toolError(message string, err error) *mcp.CallToolResultFor[any] {
	return &mcp.CallToolResultFor[any]{
		IsError: true,
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("%s: %v", message, err),
			},
		},
	}
}

// toolSuccess creates a standardized success response for MCP tools
func toolSuccess(data string) *mcp.CallToolResultFor[any] {
	return &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: data,
			},
		},
	}
}

// validateDeviceID validates that a device ID is not empty and has a reasonable format
func validateDeviceID(deviceID string) error {
	if deviceID == "" {
		return fmt.Errorf("device ID cannot be empty")
	}

	// Basic validation - device IDs should not contain spaces or special characters
	if strings.ContainsAny(deviceID, " \t\n\r") {
		return fmt.Errorf("device ID contains invalid characters")
	}

	// Device IDs should have a reasonable length
	if len(deviceID) < 3 || len(deviceID) > 50 {
		return fmt.Errorf("device ID length should be between 3 and 50 characters")
	}

	return nil
}

// getStringParam safely extracts a string parameter from MCP tool arguments
func getStringParam(params map[string]any, key string) (string, error) {
	value, exists := params[key]
	if !exists {
		return "", fmt.Errorf("%s parameter is required", key)
	}

	str, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("%s parameter must be a string", key)
	}

	return str, nil
}
