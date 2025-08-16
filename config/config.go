package config

import (
	"fmt"
	"os"

	tailscale "tailscale.com/client/tailscale/v2"
)

type Config struct {
	Tailnet string
	Port    string
	Client  *tailscale.Client
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

	client, err := createTailscaleClient(tailnet)
	if err != nil {
		return nil, fmt.Errorf("failed to create Tailscale client: %w", err)
	}

	return &Config{
		Tailnet: tailnet,
		Port:    port,
		Client:  client,
	}, nil
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
			Scopes:       []string{"devices", "routes"},
		}

		return &tailscale.Client{
			HTTP:    oauthConfig.HTTPClient(),
			Tailnet: tailnet,
		}, nil
	}

	return nil, fmt.Errorf("either TAILSCALE_API_KEY or TAILSCALE_CLIENT_ID/TAILSCALE_CLIENT_SECRET environment variables are required")
}
