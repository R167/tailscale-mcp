package main

import (
	"flag"
	"fmt"
	"os"

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
	)
	flag.Parse()

	if *showVersion {
		fmt.Printf("tailscale-mcp %s, commit %s, built at %s\n", version, commit, date)
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
		fmt.Println("\nUsage:")
		fmt.Println("  tailscale-mcp              Start the server")
		fmt.Println("  tailscale-mcp --version    Show version information")
		fmt.Println("  tailscale-mcp --help       Show this help")
		fmt.Println("\nFor more information, visit: https://github.com/R167/tailscale-mcp")
		os.Exit(0)
	}

	server.Start()
}
