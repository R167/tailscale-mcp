package tools

import (
	"context"
	"encoding/json"

	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	tailscale "tailscale.com/client/tailscale/v2"
)

func RegisterACLTools(server *mcp.Server, client *tailscale.Client) {
	// Get ACL tool
	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "get_acl",
			Description: "Get the current ACL (Access Control List) for the tailnet",
			InputSchema: &jsonschema.Schema{
				Type:                 "object",
				AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
			},
		},
		func(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[map[string]any]) (*mcp.CallToolResultFor[any], error) {
			acl, err := client.PolicyFile().Get(ctx)
			if err != nil {
				return toolError("Failed to get ACL policy", err), nil
			}

			output, err := json.MarshalIndent(acl, "", "  ")
			if err != nil {
				return toolError("Failed to serialize ACL policy", err), nil
			}

			return toolSuccess(string(output)), nil
		},
	)
}
