package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/modelcontextprotocol/go-sdk/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/tailscale/tailscale-client-go/tailscale"
)

func main() {
	apiKey := os.Getenv("TAILSCALE_API_KEY")
	if apiKey == "" {
		log.Fatal("TAILSCALE_API_KEY environment variable is required")
	}

	tailnet := os.Getenv("TAILSCALE_TAILNET")
	if tailnet == "" {
		log.Fatal("TAILSCALE_TAILNET environment variable is required")
	}

	client, err := tailscale.NewClient(apiKey, tailnet)
	if err != nil {
		log.Fatalf("Failed to create Tailscale client: %v", err)
	}

	impl := &mcp.Implementation{}
	server := mcp.NewServer(impl, &mcp.ServerOptions{})

	// Add tools
	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "list_devices",
			Description: "List all devices in the Tailscale network",
			InputSchema: &jsonschema.Schema{
				Type:                 "object",
				AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
			},
		},
		func(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[map[string]any]) (*mcp.CallToolResultFor[any], error) {
			devices, err := client.Devices(ctx)
			if err != nil {
				return &mcp.CallToolResultFor[any]{
					IsError: true,
					Content: []mcp.Content{
						&mcp.TextContent{
							Text: fmt.Sprintf("Error listing devices: %v", err),
						},
					},
				}, nil
			}

			output, err := json.MarshalIndent(devices, "", "  ")
			if err != nil {
				return &mcp.CallToolResultFor[any]{
					IsError: true,
					Content: []mcp.Content{
						&mcp.TextContent{
							Text: fmt.Sprintf("Error marshaling devices: %v", err),
						},
					},
				}, nil
			}

			return &mcp.CallToolResultFor[any]{
				Content: []mcp.Content{
					&mcp.TextContent{
						Text: string(output),
					},
				},
			}, nil
		},
	)

	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "get_acl",
			Description: "Get the current ACL (Access Control List) for the tailnet",
			InputSchema: &jsonschema.Schema{
				Type:                 "object",
				AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
			},
		},
		func(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[map[string]any]) (*mcp.CallToolResultFor[any], error) {
			acl, err := client.ACL(ctx)
			if err != nil {
				return &mcp.CallToolResultFor[any]{
					IsError: true,
					Content: []mcp.Content{
						&mcp.TextContent{
							Text: fmt.Sprintf("Error getting ACL: %v", err),
						},
					},
				}, nil
			}

			output, err := json.MarshalIndent(acl, "", "  ")
			if err != nil {
				return &mcp.CallToolResultFor[any]{
					IsError: true,
					Content: []mcp.Content{
						&mcp.TextContent{
							Text: fmt.Sprintf("Error marshaling ACL: %v", err),
						},
					},
				}, nil
			}

			return &mcp.CallToolResultFor[any]{
				Content: []mcp.Content{
					&mcp.TextContent{
						Text: string(output),
					},
				},
			}, nil
		},
	)

	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:        "list_keys",
			Description: "List all API keys for the tailnet",
			InputSchema: &jsonschema.Schema{
				Type:                 "object",
				AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
			},
		},
		func(ctx context.Context, session *mcp.ServerSession, params *mcp.CallToolParamsFor[map[string]any]) (*mcp.CallToolResultFor[any], error) {
			keys, err := client.Keys(ctx)
			if err != nil {
				return &mcp.CallToolResultFor[any]{
					IsError: true,
					Content: []mcp.Content{
						&mcp.TextContent{
							Text: fmt.Sprintf("Error listing keys: %v", err),
						},
					},
				}, nil
			}

			output, err := json.MarshalIndent(keys, "", "  ")
			if err != nil {
				return &mcp.CallToolResultFor[any]{
					IsError: true,
					Content: []mcp.Content{
						&mcp.TextContent{
							Text: fmt.Sprintf("Error marshaling keys: %v", err),
						},
					},
				}, nil
			}

			return &mcp.CallToolResultFor[any]{
				Content: []mcp.Content{
					&mcp.TextContent{
						Text: string(output),
					},
				},
			}, nil
		},
	)

	// Setup HTTP server with streamable transport
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Create HTTP handler
	handler := mcp.NewStreamableHTTPHandler(
		func(req *http.Request) *mcp.Server {
			return server
		},
		&mcp.StreamableHTTPOptions{},
	)

	// Setup HTTP server
	httpServer := &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Starting MCP server on port %s", port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down server...")

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	log.Println("Server stopped")
}