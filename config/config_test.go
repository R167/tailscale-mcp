package config

import (
	"context"
	"strings"
	"testing"

	"github.com/R167/tailscale-mcp/internal"
	tailscale "tailscale.com/client/tailscale/v2"
)

func TestLoad_Success_APIKey(t *testing.T) {
	// Setup environment variables
	t.Setenv("TAILSCALE_TAILNET", "test-tailnet")
	t.Setenv("TAILSCALE_API_KEY", "test-api-key")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if cfg.Tailnet != "test-tailnet" {
		t.Errorf("Expected tailnet 'test-tailnet', got '%s'", cfg.Tailnet)
	}

	if cfg.Port != "8080" {
		t.Errorf("Expected default port '8080', got '%s'", cfg.Port)
	}

	if cfg.Client == nil {
		t.Error("Expected client to be initialized")
	}
}

func TestLoad_Success_OAuth(t *testing.T) {
	// Setup environment variables
	t.Setenv("TAILSCALE_TAILNET", "test-tailnet")
	t.Setenv("TAILSCALE_CLIENT_ID", "test-client-id")
	t.Setenv("TAILSCALE_CLIENT_SECRET", "test-client-secret")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if cfg.Tailnet != "test-tailnet" {
		t.Errorf("Expected tailnet 'test-tailnet', got '%s'", cfg.Tailnet)
	}

	if cfg.Client == nil {
		t.Error("Expected client to be initialized")
	}
}

func TestLoad_Success_CustomPort(t *testing.T) {
	// Setup environment variables
	t.Setenv("TAILSCALE_TAILNET", "test-tailnet")
	t.Setenv("TAILSCALE_API_KEY", "test-api-key")
	t.Setenv("PORT", "9000")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if cfg.Port != "9000" {
		t.Errorf("Expected port '9000', got '%s'", cfg.Port)
	}
}

func TestLoad_Error_MissingTailnet(t *testing.T) {
	// No environment variables set - t.Setenv automatically isolates

	_, err := Load()
	if err == nil {
		t.Fatal("Expected error for missing TAILSCALE_TAILNET")
	}

	expected := "TAILSCALE_TAILNET environment variable is required"
	if err.Error() != expected {
		t.Errorf("Expected error '%s', got '%s'", expected, err.Error())
	}
}

func TestLoad_Error_MissingAuth(t *testing.T) {
	// Setup with tailnet but no auth
	t.Setenv("TAILSCALE_TAILNET", "test-tailnet")

	_, err := Load()
	if err == nil {
		t.Fatal("Expected error for missing authentication")
	}

	expected := "either TAILSCALE_API_KEY or TAILSCALE_CLIENT_ID/TAILSCALE_CLIENT_SECRET environment variables are required"
	if !strings.Contains(err.Error(), expected) {
		t.Errorf("Expected error containing '%s', got '%s'", expected, err.Error())
	}
}

func TestLoad_Error_MissingClientSecret(t *testing.T) {
	// Setup with client ID but no secret
	t.Setenv("TAILSCALE_TAILNET", "test-tailnet")
	t.Setenv("TAILSCALE_CLIENT_ID", "test-client-id")

	_, err := Load()
	if err == nil {
		t.Fatal("Expected error for missing client secret")
	}

	expected := "TAILSCALE_CLIENT_SECRET environment variable is required when using OAuth"
	if !strings.Contains(err.Error(), expected) {
		t.Errorf("Expected error containing '%s', got '%s'", expected, err.Error())
	}
}

func TestCreateTailscaleClient_APIKey(t *testing.T) {
	t.Setenv("TAILSCALE_API_KEY", "test-api-key")

	client, err := createTailscaleClient("test-tailnet")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if client == nil {
		t.Error("Expected client to be created")
	}
}

func TestCreateTailscaleClient_OAuth(t *testing.T) {
	t.Setenv("TAILSCALE_CLIENT_ID", "test-client-id")
	t.Setenv("TAILSCALE_CLIENT_SECRET", "test-client-secret")

	client, err := createTailscaleClient("test-tailnet")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if client == nil {
		t.Error("Expected client to be created")
	}
}

func TestCreateTailscaleClient_NoAuth(t *testing.T) {
	// No environment variables set - t.Setenv automatically isolates

	_, err := createTailscaleClient("test-tailnet")
	if err == nil {
		t.Fatal("Expected error for no authentication")
	}

	expected := "either TAILSCALE_API_KEY or TAILSCALE_CLIENT_ID/TAILSCALE_CLIENT_SECRET environment variables are required"
	if err.Error() != expected {
		t.Errorf("Expected error '%s', got '%s'", expected, err.Error())
	}
}

func TestCreateTailscaleClient_MissingSecret(t *testing.T) {
	t.Setenv("TAILSCALE_CLIENT_ID", "test-client-id")

	_, err := createTailscaleClient("test-tailnet")
	if err == nil {
		t.Fatal("Expected error for missing client secret")
	}

	expected := "TAILSCALE_CLIENT_SECRET environment variable is required when using OAuth"
	if err.Error() != expected {
		t.Errorf("Expected error '%s', got '%s'", expected, err.Error())
	}
}

func TestValidatePort(t *testing.T) {
	testCases := []struct {
		name      string
		port      string
		expectErr bool
		errMsg    string
	}{
		{
			name:      "ValidPort",
			port:      "8080",
			expectErr: false,
		},
		{
			name:      "ValidPortZero",
			port:      "0",
			expectErr: false,
		},
		{
			name:      "ValidPortMax",
			port:      "65535",
			expectErr: false,
		},
		{
			name:      "EmptyPort",
			port:      "",
			expectErr: true,
			errMsg:    "port cannot be empty",
		},
		{
			name:      "InvalidPortNonNumeric",
			port:      "abc",
			expectErr: true,
			errMsg:    "port must be a valid integer",
		},
		{
			name:      "InvalidPortNegative",
			port:      "-1",
			expectErr: true,
			errMsg:    "port must be between 0 and 65535",
		},
		{
			name:      "InvalidPortTooHigh",
			port:      "65536",
			expectErr: true,
			errMsg:    "port must be between 0 and 65535",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validatePort(tc.port)
			if tc.expectErr {
				if err == nil {
					t.Error("Expected error but got none")
				} else if !strings.Contains(err.Error(), tc.errMsg) {
					t.Errorf("Expected error containing '%s', got '%s'", tc.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestValidateTailnet(t *testing.T) {
	testCases := []struct {
		name      string
		tailnet   string
		expectErr bool
		errMsg    string
	}{
		{
			name:      "ValidTailnet",
			tailnet:   "example.ts.net",
			expectErr: false,
		},
		{
			name:      "ValidTailnetShort",
			tailnet:   "abc",
			expectErr: false,
		},
		{
			name:      "EmptyTailnet",
			tailnet:   "",
			expectErr: true,
			errMsg:    "tailnet cannot be empty",
		},
		{
			name:      "TooShortTailnet",
			tailnet:   "ab",
			expectErr: true,
			errMsg:    "tailnet must be at least 3 characters long",
		},
		{
			name:      "TailnetWithSpace",
			tailnet:   "example ts.net",
			expectErr: true,
			errMsg:    "tailnet contains invalid whitespace characters",
		},
		{
			name:      "TailnetWithTab",
			tailnet:   "example\tts.net",
			expectErr: true,
			errMsg:    "tailnet contains invalid whitespace characters",
		},
		{
			name:      "TailnetWithNewline",
			tailnet:   "example\nts.net",
			expectErr: true,
			errMsg:    "tailnet contains invalid whitespace characters",
		},
		{
			name:      "TailnetWithDotDot",
			tailnet:   "example..ts.net",
			expectErr: true,
			errMsg:    "tailnet contains invalid '..' pattern",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateTailnet(tc.tailnet)
			if tc.expectErr {
				if err == nil {
					t.Error("Expected error but got none")
				} else if !strings.Contains(err.Error(), tc.errMsg) {
					t.Errorf("Expected error containing '%s', got '%s'", tc.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestValidateConfig(t *testing.T) {
	testCases := []struct {
		name      string
		config    *Config
		expectErr bool
		errMsg    string
	}{
		{
			name: "ValidConfig",
			config: &Config{
				Tailnet: "example.ts.net",
				Port:    "8080",
				Client:  &mockClient{},
			},
			expectErr: false,
		},
		{
			name:      "NilConfig",
			config:    nil,
			expectErr: true,
			errMsg:    "configuration is nil",
		},
		{
			name: "NilClient",
			config: &Config{
				Tailnet: "example.ts.net",
				Port:    "8080",
				Client:  nil,
			},
			expectErr: true,
			errMsg:    "Tailscale client is not initialized",
		},
		{
			name: "EmptyTailnet",
			config: &Config{
				Tailnet: "",
				Port:    "8080",
				Client:  &mockClient{},
			},
			expectErr: true,
			errMsg:    "tailnet is empty",
		},
		{
			name: "EmptyPort",
			config: &Config{
				Tailnet: "example.ts.net",
				Port:    "",
				Client:  &mockClient{},
			},
			expectErr: true,
			errMsg:    "port is empty",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateConfig(tc.config)
			if tc.expectErr {
				if err == nil {
					t.Error("Expected error but got none")
				} else if !strings.Contains(err.Error(), tc.errMsg) {
					t.Errorf("Expected error containing '%s', got '%s'", tc.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}

// mockClient is a simple mock for testing validateConfig
type mockClient struct{}

func (m *mockClient) Devices() internal.DevicesResource {
	return &mockDevicesResource{}
}

func (m *mockClient) PolicyFile() internal.PolicyFileResource {
	return &mockPolicyFileResource{}
}

func (m *mockClient) Keys() internal.KeysResource {
	return &mockKeysResource{}
}

// Mock resource implementations
type mockDevicesResource struct{}
type mockPolicyFileResource struct{}
type mockKeysResource struct{}

func (m *mockDevicesResource) ListWithAllFields(ctx context.Context) ([]tailscale.Device, error) {
	return nil, nil
}

func (m *mockDevicesResource) GetWithAllFields(ctx context.Context, deviceID string) (*tailscale.Device, error) {
	return nil, nil
}

func (m *mockDevicesResource) SubnetRoutes(ctx context.Context, deviceID string) (*tailscale.DeviceRoutes, error) {
	return nil, nil
}

func (m *mockPolicyFileResource) Get(ctx context.Context) (*tailscale.ACL, error) {
	return nil, nil
}

func (m *mockKeysResource) List(ctx context.Context, all bool) ([]tailscale.Key, error) {
	return nil, nil
}
