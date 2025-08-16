package internal

import (
	"context"
	"fmt"

	tailscale "tailscale.com/client/tailscale/v2"
)

// MockTailscaleClient is a mock implementation for testing
type MockTailscaleClient struct {
	DevicesFunc    func() DevicesResource
	PolicyFileFunc func() PolicyFileResource
	KeysFunc       func() KeysResource
}

func (m *MockTailscaleClient) Devices() DevicesResource {
	if m.DevicesFunc != nil {
		return m.DevicesFunc()
	}
	return &MockDevicesResource{}
}

func (m *MockTailscaleClient) PolicyFile() PolicyFileResource {
	if m.PolicyFileFunc != nil {
		return m.PolicyFileFunc()
	}
	return &MockPolicyFileResource{}
}

func (m *MockTailscaleClient) Keys() KeysResource {
	if m.KeysFunc != nil {
		return m.KeysFunc()
	}
	return &MockKeysResource{}
}

// MockDevicesResource is a mock implementation for testing
type MockDevicesResource struct {
	ListFunc              func(ctx context.Context) ([]tailscale.Device, error)
	ListWithAllFieldsFunc func(ctx context.Context) ([]tailscale.Device, error)
	GetWithAllFieldsFunc  func(ctx context.Context, deviceID string) (*tailscale.Device, error)
	SubnetRoutesFunc      func(ctx context.Context, deviceID string) (*tailscale.DeviceRoutes, error)
}

func (m *MockDevicesResource) List(ctx context.Context) ([]tailscale.Device, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx)
	}
	return []tailscale.Device{
		{
			ID:        "device1",
			Name:      "test-device-1",
			Addresses: []string{"100.1.1.1"},
		},
		{
			ID:        "device2",
			Name:      "test-device-2",
			Addresses: []string{"100.1.1.2"},
		},
	}, nil
}

func (m *MockDevicesResource) ListWithAllFields(ctx context.Context) ([]tailscale.Device, error) {
	if m.ListWithAllFieldsFunc != nil {
		return m.ListWithAllFieldsFunc(ctx)
	}
	return []tailscale.Device{
		{
			ID:        "device1",
			Name:      "test-device-1",
			Addresses: []string{"100.1.1.1"},
		},
		{
			ID:        "device2",
			Name:      "test-device-2",
			Addresses: []string{"100.1.1.2"},
		},
	}, nil
}

func (m *MockDevicesResource) GetWithAllFields(ctx context.Context, deviceID string) (*tailscale.Device, error) {
	if m.GetWithAllFieldsFunc != nil {
		return m.GetWithAllFieldsFunc(ctx, deviceID)
	}
	if deviceID == "invalid" {
		return nil, fmt.Errorf("device not found")
	}
	return &tailscale.Device{
		ID:        deviceID,
		Name:      "test-device",
		Addresses: []string{"100.1.1.1"},
	}, nil
}

func (m *MockDevicesResource) SubnetRoutes(ctx context.Context, deviceID string) (*tailscale.DeviceRoutes, error) {
	if m.SubnetRoutesFunc != nil {
		return m.SubnetRoutesFunc(ctx, deviceID)
	}
	if deviceID == "invalid" {
		return nil, fmt.Errorf("device not found")
	}
	return &tailscale.DeviceRoutes{
		Advertised: []string{"10.0.0.0/24"},
		Enabled:    []string{"10.0.0.0/24"},
	}, nil
}

// MockPolicyFileResource is a mock implementation for testing
type MockPolicyFileResource struct {
	GetFunc func(ctx context.Context) (*tailscale.ACL, error)
}

func (m *MockPolicyFileResource) Get(ctx context.Context) (*tailscale.ACL, error) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx)
	}
	return &tailscale.ACL{
		ACLs: []tailscale.ACLEntry{
			{
				Action:      "accept",
				Source:      []string{"*"},
				Destination: []string{"*:*"},
			},
		},
	}, nil
}

// MockKeysResource is a mock implementation for testing
type MockKeysResource struct {
	ListFunc func(ctx context.Context, all bool) ([]tailscale.Key, error)
}

func (m *MockKeysResource) List(ctx context.Context, all bool) ([]tailscale.Key, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx, all)
	}
	return []tailscale.Key{
		{
			ID:          "key1",
			Description: "Test API Key 1",
		},
		{
			ID:          "key2",
			Description: "Test API Key 2",
		},
	}, nil
}
