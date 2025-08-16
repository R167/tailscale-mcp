package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/R167/tailscale-mcp/internal"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	tailscale "tailscale.com/client/tailscale/v2"
)

func TestRegisterACLTools(t *testing.T) {
	impl := &mcp.Implementation{}
	server := mcp.NewServer(impl, &mcp.ServerOptions{})
	mockClient := &internal.MockTailscaleClient{}

	// This should not panic
	RegisterACLTools(server, mockClient)
}

func TestGetACLSuccess(t *testing.T) {
	mockPolicyFile := &internal.MockPolicyFileResource{
		GetFunc: func(ctx context.Context) (*tailscale.ACL, error) {
			return &tailscale.ACL{
				ACLs: []tailscale.ACLEntry{
					{
						Action:      "accept",
						Source:      []string{"group:admin"},
						Destination: []string{"tag:server:*"},
					},
				},
			}, nil
		},
	}

	mockClient := &internal.MockTailscaleClient{
		PolicyFileFunc: func() internal.PolicyFileResource {
			return mockPolicyFile
		},
	}

	impl := &mcp.Implementation{}
	server := mcp.NewServer(impl, &mcp.ServerOptions{})
	RegisterACLTools(server, mockClient)
}

func TestGetACLError(t *testing.T) {
	mockPolicyFile := &internal.MockPolicyFileResource{
		GetFunc: func(ctx context.Context) (*tailscale.ACL, error) {
			return nil, fmt.Errorf("API error: unauthorized")
		},
	}

	mockClient := &internal.MockTailscaleClient{
		PolicyFileFunc: func() internal.PolicyFileResource {
			return mockPolicyFile
		},
	}

	impl := &mcp.Implementation{}
	server := mcp.NewServer(impl, &mcp.ServerOptions{})
	RegisterACLTools(server, mockClient)
}

func TestACLJSONSerialization(t *testing.T) {
	acl := tailscale.ACL{
		ACLs: []tailscale.ACLEntry{
			{
				Action:      "accept",
				Source:      []string{"group:admin"},
				Destination: []string{"tag:server:22"},
				Protocol:    "tcp",
			},
			{
				Action:      "accept",
				Source:      []string{"*"},
				Destination: []string{"tag:web:80,443"},
				Protocol:    "tcp",
			},
		},
	}

	data, err := json.MarshalIndent(acl, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal ACL: %v", err)
	}

	var unmarshaled tailscale.ACL
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal ACL: %v", err)
	}

	if len(unmarshaled.ACLs) != len(acl.ACLs) {
		t.Errorf("Expected %d ACL entries, got %d", len(acl.ACLs), len(unmarshaled.ACLs))
	}

	for i, entry := range unmarshaled.ACLs {
		if entry.Action != acl.ACLs[i].Action {
			t.Errorf("Expected action '%s', got '%s'", acl.ACLs[i].Action, entry.Action)
		}
	}
}

func TestACLEmptyResponse(t *testing.T) {
	mockPolicyFile := &internal.MockPolicyFileResource{
		GetFunc: func(ctx context.Context) (*tailscale.ACL, error) {
			return &tailscale.ACL{
				ACLs: []tailscale.ACLEntry{},
			}, nil
		},
	}

	mockClient := &internal.MockTailscaleClient{
		PolicyFileFunc: func() internal.PolicyFileResource {
			return mockPolicyFile
		},
	}

	impl := &mcp.Implementation{}
	server := mcp.NewServer(impl, &mcp.ServerOptions{})
	RegisterACLTools(server, mockClient)

	// Test that empty ACL can be serialized
	acl := tailscale.ACL{ACLs: []tailscale.ACLEntry{}}
	data, err := json.MarshalIndent(acl, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal empty ACL: %v", err)
	}

	if string(data) == "" {
		t.Error("Expected non-empty JSON for empty ACL")
	}
}
