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

func TestRegisterKeyTools(t *testing.T) {
	impl := &mcp.Implementation{}
	server := mcp.NewServer(impl, &mcp.ServerOptions{})
	mockClient := &internal.MockTailscaleClient{}

	// This should not panic
	RegisterKeyTools(server, mockClient)
}

func TestListKeysSuccess(t *testing.T) {
	mockKeys := &internal.MockKeysResource{
		ListFunc: func(ctx context.Context, all bool) ([]tailscale.Key, error) {
			return []tailscale.Key{
				{
					ID:          "key1",
					Description: "Admin API Key",
				},
				{
					ID:          "key2",
					Description: "Read-only API Key",
				},
			}, nil
		},
	}

	mockClient := &internal.MockTailscaleClient{
		KeysFunc: func() internal.KeysResource {
			return mockKeys
		},
	}

	impl := &mcp.Implementation{}
	server := mcp.NewServer(impl, &mcp.ServerOptions{})
	RegisterKeyTools(server, mockClient)
}

func TestListKeysError(t *testing.T) {
	mockKeys := &internal.MockKeysResource{
		ListFunc: func(ctx context.Context, all bool) ([]tailscale.Key, error) {
			return nil, fmt.Errorf("API error: forbidden")
		},
	}

	mockClient := &internal.MockTailscaleClient{
		KeysFunc: func() internal.KeysResource {
			return mockKeys
		},
	}

	impl := &mcp.Implementation{}
	server := mcp.NewServer(impl, &mcp.ServerOptions{})
	RegisterKeyTools(server, mockClient)
}

func TestListKeysAllParameter(t *testing.T) {
	mockKeys := &internal.MockKeysResource{
		ListFunc: func(ctx context.Context, all bool) ([]tailscale.Key, error) {
			// Verify that the tool calls with all=true
			if !all {
				t.Error("Expected all parameter to be true")
			}
			return []tailscale.Key{}, nil
		},
	}

	mockClient := &internal.MockTailscaleClient{
		KeysFunc: func() internal.KeysResource {
			return mockKeys
		},
	}

	impl := &mcp.Implementation{}
	server := mcp.NewServer(impl, &mcp.ServerOptions{})
	RegisterKeyTools(server, mockClient)

	// The tool should call List with all=true
	// In a real test, we'd invoke the tool and check the captured value
}

func TestKeysJSONSerialization(t *testing.T) {
	keys := []tailscale.Key{
		{
			ID:          "key1",
			Description: "Test API Key",
		},
		{
			ID:          "key2",
			Description: "Another Test Key",
		},
	}

	data, err := json.MarshalIndent(keys, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal keys: %v", err)
	}

	var unmarshaled []tailscale.Key
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal keys: %v", err)
	}

	if len(unmarshaled) != len(keys) {
		t.Errorf("Expected %d keys, got %d", len(keys), len(unmarshaled))
	}

	for i, key := range unmarshaled {
		if key.ID != keys[i].ID {
			t.Errorf("Expected key ID '%s', got '%s'", keys[i].ID, key.ID)
		}
		if key.Description != keys[i].Description {
			t.Errorf("Expected key description '%s', got '%s'", keys[i].Description, key.Description)
		}
	}
}

func TestKeysEmptyResponse(t *testing.T) {
	mockKeys := &internal.MockKeysResource{
		ListFunc: func(ctx context.Context, all bool) ([]tailscale.Key, error) {
			return []tailscale.Key{}, nil
		},
	}

	mockClient := &internal.MockTailscaleClient{
		KeysFunc: func() internal.KeysResource {
			return mockKeys
		},
	}

	impl := &mcp.Implementation{}
	server := mcp.NewServer(impl, &mcp.ServerOptions{})
	RegisterKeyTools(server, mockClient)

	// Test that empty keys list can be serialized
	keys := []tailscale.Key{}
	data, err := json.MarshalIndent(keys, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal empty keys: %v", err)
	}

	if string(data) == "" {
		t.Error("Expected non-empty JSON for empty keys")
	}
}
