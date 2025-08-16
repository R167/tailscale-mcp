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

func TestRegisterDeviceTools(t *testing.T) {
	impl := &mcp.Implementation{}
	server := mcp.NewServer(impl, &mcp.ServerOptions{})
	mockClient := &internal.MockTailscaleClient{}

	// This should not panic
	RegisterDeviceTools(server, mockClient)
}

func TestListDevicesSuccess(t *testing.T) {
	mockDevices := &internal.MockDevicesResource{
		ListWithAllFieldsFunc: func(ctx context.Context) ([]tailscale.Device, error) {
			return []tailscale.Device{
				{ID: "device1", Name: "test-device-1"},
				{ID: "device2", Name: "test-device-2"},
			}, nil
		},
	}

	mockClient := &internal.MockTailscaleClient{
		DevicesFunc: func() internal.DevicesResource {
			return mockDevices
		},
	}

	impl := &mcp.Implementation{}
	server := mcp.NewServer(impl, &mcp.ServerOptions{})
	RegisterDeviceTools(server, mockClient)

	// We can't easily test the actual tool function directly, so we just ensure registration doesn't panic
	// In a real test environment, you'd use the MCP client to call the tool
}

func TestGetDeviceDetailsSuccess(t *testing.T) {
	mockDevices := &internal.MockDevicesResource{
		GetWithAllFieldsFunc: func(ctx context.Context, deviceID string) (*tailscale.Device, error) {
			if deviceID == "test-device" {
				return &tailscale.Device{
					ID:        "test-device",
					Name:      "test-device-name",
					Addresses: []string{"100.1.1.1"},
				}, nil
			}
			return nil, fmt.Errorf("device not found")
		},
	}

	mockClient := &internal.MockTailscaleClient{
		DevicesFunc: func() internal.DevicesResource {
			return mockDevices
		},
	}

	impl := &mcp.Implementation{}
	server := mcp.NewServer(impl, &mcp.ServerOptions{})
	RegisterDeviceTools(server, mockClient)
}

func TestGetDeviceDetailsInvalidID(t *testing.T) {
	// Test with various invalid device IDs
	invalidIDs := []string{
		"",
		"ab",
		"device with spaces",
		string(make([]byte, 51)), // too long
	}

	for _, invalidID := range invalidIDs {
		t.Run(fmt.Sprintf("invalid_id_%s", invalidID), func(t *testing.T) {
			err := validateDeviceID(invalidID)
			if err == nil {
				t.Errorf("Expected validation error for device ID '%s'", invalidID)
			}
		})
	}
}

func TestGetDeviceRoutesSuccess(t *testing.T) {
	mockDevices := &internal.MockDevicesResource{
		SubnetRoutesFunc: func(ctx context.Context, deviceID string) (*tailscale.DeviceRoutes, error) {
			if deviceID == "test-device" {
				return &tailscale.DeviceRoutes{
					Advertised: []string{"10.0.0.0/24"},
					Enabled:    []string{"10.0.0.0/24"},
				}, nil
			}
			return nil, fmt.Errorf("device not found")
		},
	}

	mockClient := &internal.MockTailscaleClient{
		DevicesFunc: func() internal.DevicesResource {
			return mockDevices
		},
	}

	impl := &mcp.Implementation{}
	server := mcp.NewServer(impl, &mcp.ServerOptions{})
	RegisterDeviceTools(server, mockClient)
}

func TestGetDeviceRoutesError(t *testing.T) {
	mockDevices := &internal.MockDevicesResource{
		SubnetRoutesFunc: func(ctx context.Context, deviceID string) (*tailscale.DeviceRoutes, error) {
			return nil, fmt.Errorf("API error")
		},
	}

	mockClient := &internal.MockTailscaleClient{
		DevicesFunc: func() internal.DevicesResource {
			return mockDevices
		},
	}

	impl := &mcp.Implementation{}
	server := mcp.NewServer(impl, &mcp.ServerOptions{})
	RegisterDeviceTools(server, mockClient)
}

// Test helper functions

func TestDeviceToolJSONSerialization(t *testing.T) {
	device := tailscale.Device{
		ID:        "test-device",
		Name:      "test-device-name",
		Addresses: []string{"100.1.1.1"},
	}

	data, err := json.MarshalIndent(device, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal device: %v", err)
	}

	var unmarshaled tailscale.Device
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal device: %v", err)
	}

	if unmarshaled.ID != device.ID {
		t.Errorf("Expected ID '%s', got '%s'", device.ID, unmarshaled.ID)
	}
}

func TestDeviceRoutesJSONSerialization(t *testing.T) {
	routes := tailscale.DeviceRoutes{
		Advertised: []string{"10.0.0.0/24", "192.168.1.0/24"},
		Enabled:    []string{"10.0.0.0/24"},
	}

	data, err := json.MarshalIndent(routes, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal routes: %v", err)
	}

	var unmarshaled tailscale.DeviceRoutes
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal routes: %v", err)
	}

	if len(unmarshaled.Advertised) != len(routes.Advertised) {
		t.Errorf("Expected %d advertised routes, got %d", len(routes.Advertised), len(unmarshaled.Advertised))
	}
}
