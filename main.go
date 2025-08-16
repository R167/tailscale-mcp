package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"github.com/R167/tailscale-mcp/server"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	var (
		showVersion = flag.Bool("version", false, "Show version information")
		showHelp    = flag.Bool("help", false, "Show help information")
		loadEnv     = flag.Bool("env", true, "Load .env file if present (default: true)")
	)
	flag.Parse()

	// Load .env file if requested
	if *loadEnv {
		if err := godotenv.Load(); err != nil {
			// Only log in development - production deployments typically use env vars directly
			if os.Getenv("ENVIRONMENT") == "development" {
				fmt.Printf("Note: .env file not found (%v)\n", err)
			}
		}
	}

	if *showVersion {
		fmt.Printf("tailscale-mcp %s, commit %s, built at %s\n", version, commit, date)

		// Debug information
		fmt.Println("\nConfiguration:")
		if tailnet := os.Getenv("TAILSCALE_TAILNET"); tailnet != "" {
			fmt.Printf("  Tailnet: %s\n", tailnet)
		} else {
			fmt.Println("  Tailnet: not configured")
		}

		// Check authentication method
		if apiKey := os.Getenv("TAILSCALE_API_KEY"); apiKey != "" {
			keyPreview := apiKey
			if len(apiKey) > 10 {
				keyPreview = apiKey[:10]
			}
			fmt.Printf("  Auth Method: API Key (key: %s...)\n", keyPreview)
		} else if clientID := os.Getenv("TAILSCALE_CLIENT_ID"); clientID != "" {
			fmt.Printf("  Auth Method: OAuth (client ID: %s)\n", clientID)
		} else {
			fmt.Println("  Auth Method: not configured")
		}

		if port := os.Getenv("PORT"); port != "" {
			fmt.Printf("  Port: %s\n", port)
		} else {
			fmt.Println("  Port: 8080 (default)")
		}

		os.Exit(0)
	}

	if *showHelp {
		fmt.Printf("Tailscale MCP Server %s\n\n", version)
		fmt.Println("A Model Context Protocol server for Tailscale network management.")
		fmt.Println("\nEnvironment Variables:")
		fmt.Println("  TAILSCALE_TAILNET          Your tailnet identifier (required)")
		fmt.Println("  TAILSCALE_API_KEY          Your Tailscale API key")
		fmt.Println("  TAILSCALE_CLIENT_ID        OAuth client ID (alternative to API key)")
		fmt.Println("  TAILSCALE_CLIENT_SECRET    OAuth client secret (required with CLIENT_ID)")
		fmt.Println("  PORT                       HTTP server port (default: 8080)")
		fmt.Println("\nConfiguration:")
		fmt.Println("  Environment variables can be set via .env file or system environment.")
		fmt.Println("\nUsage:")
		fmt.Println("  tailscale-mcp              Start the server")
		fmt.Println("  tailscale-mcp --version    Show version and configuration")
		fmt.Println("  tailscale-mcp --help       Show this help")
		fmt.Println("  tailscale-mcp --env=false  Disable .env file loading")
		fmt.Println("\nFor more information, visit: https://github.com/R167/tailscale-mcp")
		os.Exit(0)
	}

	server.Start()
}
