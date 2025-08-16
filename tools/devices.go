package tools

import (
	"context"
	"encoding/json"

	"github.com/R167/tailscale-mcp/internal"
	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func RegisterDeviceTools(server *mcp.Server, client internal.TailscaleClient) {
	// List devices tool
	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "list_devices",
			Description: "List all devices in the Tailscale network",
			InputSchema: &jsonschema.Schema{
				Type:                 "object",
				AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
			},
		},
		func(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[map[string]any]) (*mcp.CallToolResultFor[any], error) {
			devices, err := client.Devices().ListWithAllFields(ctx)
			if err != nil {
				return toolError("Failed to list devices", err), nil
			}

			output, err := json.MarshalIndent(devices, "", "  ")
			if err != nil {
				return toolError("Failed to serialize device list", err), nil
			}

			return toolSuccess(string(output)), nil
		},
	)

	// Get device details tool
	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "get_device_details",
			Description: "Get detailed information about a specific device",
			InputSchema: &jsonschema.Schema{
				Type: "object",
				Properties: map[string]*jsonschema.Schema{
					"deviceID": {
						Type:        "string",
						Description: "The device ID to get details for",
					},
				},
				Required:             []string{"deviceID"},
				AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
			},
		},
		func(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[map[string]any]) (*mcp.CallToolResultFor[any], error) {
			deviceID, err := getStringParam(params.Arguments, "deviceID")
			if err != nil {
				return toolError("Invalid device ID parameter", err), nil
			}

			if err := validateDeviceID(deviceID); err != nil {
				return toolError("Device ID validation failed", err), nil
			}

			device, err := client.Devices().GetWithAllFields(ctx, deviceID)
			if err != nil {
				return toolError("Failed to get device details", err), nil
			}

			output, err := json.MarshalIndent(device, "", "  ")
			if err != nil {
				return toolError("Failed to serialize device details", err), nil
			}

			return toolSuccess(string(output)), nil
		},
	)

	// Get device routes tool
	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "get_device_routes",
			Description: "Get subnet routes for a specific device",
			InputSchema: &jsonschema.Schema{
				Type: "object",
				Properties: map[string]*jsonschema.Schema{
					"deviceID": {
						Type:        "string",
						Description: "The device ID to get routes for",
					},
				},
				Required:             []string{"deviceID"},
				AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
			},
		},
		func(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[map[string]any]) (*mcp.CallToolResultFor[any], error) {
			deviceID, err := getStringParam(params.Arguments, "deviceID")
			if err != nil {
				return toolError("Invalid device ID parameter", err), nil
			}

			if err := validateDeviceID(deviceID); err != nil {
				return toolError("Device ID validation failed", err), nil
			}

			routes, err := client.Devices().SubnetRoutes(ctx, deviceID)
			if err != nil {
				return toolError("Failed to get device routes", err), nil
			}

			output, err := json.MarshalIndent(routes, "", "  ")
			if err != nil {
				return toolError("Failed to serialize device routes", err), nil
			}

			return toolSuccess(string(output)), nil
		},
	)
}
