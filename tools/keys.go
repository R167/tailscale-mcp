package tools

import (
	"context"
	"encoding/json"

	"github.com/R167/tailscale-mcp/internal"
	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func RegisterKeyTools(server *mcp.Server, client internal.TailscaleClient) {
	// List keys tool
	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "list_keys",
			Description: "List all API keys for the tailnet",
			InputSchema: &jsonschema.Schema{
				Type:                 "object",
				AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
			},
		},
		func(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[map[string]any]) (*mcp.CallToolResultFor[any], error) {
			keys, err := client.Keys().List(ctx, true)
			if err != nil {
				return toolError("Failed to list API keys", err), nil
			}

			output, err := json.MarshalIndent(keys, "", "  ")
			if err != nil {
				return toolError("Failed to serialize API keys", err), nil
			}

			return toolSuccess(string(output)), nil
		},
	)
}
