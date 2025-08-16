package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	tailscale "tailscale.com/client/tailscale/v2"

	"github.com/R167/tailscale-mcp/internal"
)

type Config struct {
	Tailnet string
	Port    string
	Client  internal.TailscaleClient
}

func Load() (*Config, error) {
	tailnet := os.Getenv("TAILSCALE_TAILNET")
	if tailnet == "" {
		return nil, fmt.Errorf("TAILSCALE_TAILNET environment variable is required")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Validate port
	if err := validatePort(port); err != nil {
		return nil, fmt.Errorf("invalid port configuration: %w", err)
	}

	// Validate tailnet format
	if err := validateTailnet(tailnet); err != nil {
		return nil, fmt.Errorf("invalid tailnet configuration: %w", err)
	}

	client, err := createTailscaleClient(tailnet)
	if err != nil {
		return nil, fmt.Errorf("failed to create Tailscale client: %w", err)
	}

	cfg := &Config{
		Tailnet: tailnet,
		Port:    port,
		Client:  &internal.TailscaleClientAdapter{Client: client},
	}

	// Validate complete configuration
	if err := validateConfig(cfg); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return cfg, nil
}

func createTailscaleClient(tailnet string) (*tailscale.Client, error) {
	if apiKey := os.Getenv("TAILSCALE_API_KEY"); apiKey != "" {
		// Use API Key authentication
		return &tailscale.Client{
			APIKey:  apiKey,
			Tailnet: tailnet,
		}, nil
	}

	if clientID := os.Getenv("TAILSCALE_CLIENT_ID"); clientID != "" {
		// Use OAuth Client credentials
		clientSecret := os.Getenv("TAILSCALE_CLIENT_SECRET")
		if clientSecret == "" {
			return nil, fmt.Errorf("TAILSCALE_CLIENT_SECRET environment variable is required when using OAuth")
		}

		oauthConfig := tailscale.OAuthConfig{
			ClientID:     clientID,
			ClientSecret: clientSecret,
		}

		return &tailscale.Client{
			HTTP:    oauthConfig.HTTPClient(),
			Tailnet: tailnet,
		}, nil
	}

	return nil, fmt.Errorf("either TAILSCALE_API_KEY or TAILSCALE_CLIENT_ID/TAILSCALE_CLIENT_SECRET environment variables are required")
}

// validatePort checks if the port is valid
func validatePort(port string) error {
	if port == "" {
		return fmt.Errorf("port cannot be empty")
	}

	portNum, err := strconv.Atoi(port)
	if err != nil {
		return fmt.Errorf("port must be a valid integer: %w", err)
	}

	if portNum < 0 || portNum > 65535 {
		return fmt.Errorf("port must be between 0 and 65535, got %d", portNum)
	}

	return nil
}

// validateTailnet checks if the tailnet format is valid
func validateTailnet(tailnet string) error {
	if tailnet == "" {
		return fmt.Errorf("tailnet cannot be empty")
	}

	// Basic validation for tailnet format
	if len(tailnet) < 3 {
		return fmt.Errorf("tailnet must be at least 3 characters long")
	}

	// Check for invalid characters
	if strings.ContainsAny(tailnet, " \t\n\r") {
		return fmt.Errorf("tailnet contains invalid whitespace characters")
	}

	// Check for suspicious patterns
	if strings.Contains(tailnet, "..") {
		return fmt.Errorf("tailnet contains invalid '..' pattern")
	}

	return nil
}

// validateConfig performs final validation on the complete configuration
func validateConfig(cfg *Config) error {
	if cfg == nil {
		return fmt.Errorf("configuration is nil")
	}

	if cfg.Client == nil {
		return fmt.Errorf("Tailscale client is not initialized")
	}

	if cfg.Tailnet == "" {
		return fmt.Errorf("tailnet is empty")
	}

	if cfg.Port == "" {
		return fmt.Errorf("port is empty")
	}

	return nil
}
