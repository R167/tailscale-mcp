package main

import (
	"testing"
	"time"

	"github.com/R167/tailscale-mcp/config"
)

func TestConfigIntegration(t *testing.T) {
	t.Run("APIKeyAuth", func(t *testing.T) {
		// Setup test environment
		t.Setenv("TAILSCALE_TAILNET", "test-tailnet.ts.net")
		t.Setenv("TAILSCALE_API_KEY", "tskey-test-key")
		t.Setenv("PORT", "8080")

		cfg, err := config.Load()
		if err != nil {
			t.Fatalf("Expected successful config load, got error: %v", err)
		}

		if cfg.Tailnet != "test-tailnet.ts.net" {
			t.Errorf("Expected tailnet 'test-tailnet.ts.net', got '%s'", cfg.Tailnet)
		}

		if cfg.Port != "8080" {
			t.Errorf("Expected port '8080', got '%s'", cfg.Port)
		}

		if cfg.Client == nil {
			t.Error("Expected client to be initialized")
		}
	})

	t.Run("OAuthAuth", func(t *testing.T) {
		// Setup test environment
		t.Setenv("TAILSCALE_TAILNET", "test-tailnet.ts.net")
		t.Setenv("TAILSCALE_CLIENT_ID", "test-client-id")
		t.Setenv("TAILSCALE_CLIENT_SECRET", "test-client-secret")
		t.Setenv("PORT", "9000")

		cfg, err := config.Load()
		if err != nil {
			t.Fatalf("Expected successful config load, got error: %v", err)
		}

		if cfg.Tailnet != "test-tailnet.ts.net" {
			t.Errorf("Expected tailnet 'test-tailnet.ts.net', got '%s'", cfg.Tailnet)
		}

		if cfg.Port != "9000" {
			t.Errorf("Expected port '9000', got '%s'", cfg.Port)
		}

		if cfg.Client == nil {
			t.Error("Expected client to be initialized")
		}
	})
}

func TestServerStartupIntegration(t *testing.T) {
	// This test verifies that the server can be configured properly
	// but doesn't actually start the HTTP server to avoid port conflicts

	// Setup test environment
	t.Setenv("TAILSCALE_TAILNET", "test-tailnet.ts.net")
	t.Setenv("TAILSCALE_API_KEY", "tskey-test-key")
	t.Setenv("PORT", "0") // Use port 0 to get any available port

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Expected successful config load, got error: %v", err)
	}

	// Verify configuration is valid
	if cfg.Client == nil {
		t.Error("Expected client to be configured")
	}

	if cfg.Tailnet == "" {
		t.Error("Expected tailnet to be configured")
	}

	if cfg.Port != "0" {
		t.Errorf("Expected port '0', got '%s'", cfg.Port)
	}
}

func TestConfigValidation(t *testing.T) {
	// Test various invalid configurations
	testCases := []struct {
		name      string
		setup     func(t *testing.T)
		expectErr bool
	}{
		{
			name: "MissingTailnet",
			setup: func(t *testing.T) {
				t.Setenv("TAILSCALE_API_KEY", "test-key")
			},
			expectErr: true,
		},
		{
			name: "MissingAuth",
			setup: func(t *testing.T) {
				t.Setenv("TAILSCALE_TAILNET", "test-tailnet")
			},
			expectErr: true,
		},
		{
			name: "IncompleteOAuth",
			setup: func(t *testing.T) {
				t.Setenv("TAILSCALE_TAILNET", "test-tailnet")
				t.Setenv("TAILSCALE_CLIENT_ID", "test-id")
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup(t)

			_, err := config.Load()
			if tc.expectErr && err == nil {
				t.Error("Expected error but got none")
			}
			if !tc.expectErr && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestPerformanceBasics(t *testing.T) {
	// Basic performance test to ensure config loading is fast
	t.Setenv("TAILSCALE_TAILNET", "test-tailnet")
	t.Setenv("TAILSCALE_API_KEY", "test-key")

	start := time.Now()
	_, err := config.Load()
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("Expected successful config load, got error: %v", err)
	}

	// Config loading should be very fast (under 100ms)
	if duration > 100*time.Millisecond {
		t.Errorf("Config loading took too long: %v", duration)
	}
}
