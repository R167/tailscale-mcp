package server

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/R167/tailscale-mcp/config"
	"github.com/R167/tailscale-mcp/tools"
)

func Start() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	impl := &mcp.Implementation{}
	server := mcp.NewServer(impl, &mcp.ServerOptions{})

	// Register all tools
	tools.RegisterDeviceTools(server, cfg.Client)
	tools.RegisterACLTools(server, cfg.Client)
	tools.RegisterKeyTools(server, cfg.Client)

	// Create HTTP handler
	mcpHandler := mcp.NewStreamableHTTPHandler(
		func(req *http.Request) *mcp.Server {
			return server
		},
		&mcp.StreamableHTTPOptions{},
	)

	// Wrap with middleware for request context and logging
	handler := RequestMiddleware(mcpHandler)

	// Setup HTTP server
	httpServer := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: handler,
	}

	// Start server in goroutine
	go func() {
		slog.Info("Starting MCP server", "port", cfg.Port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("HTTP server failed", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	slog.Info("Shutting down server...")

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		slog.Error("Server shutdown error", "error", err)
	}

	slog.Info("Server stopped")
}
